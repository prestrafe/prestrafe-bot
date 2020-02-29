package gsi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"prestrafe-bot/config"
)

// Defines the public API for the Game State Integration server. The server acts as a rely between the CSGO GSI API,
// which sends game state data to a configured web-hook and potential clients, which may wish to consume this data as a
// service, without providing their own HTTP server. The GSI server supports multiple tenants, which are identified by
// their authentication token, that is send with each GSI web-hook call.
type Server interface {
	// Starts the server in the current thread and blocks until an error occurs.
	Start() error
}

type server struct {
	port       int
	store      Store
	httpServer *http.Server
	upgrader   *websocket.Upgrader
}

// Creates a new GSI server.
func NewServer(config *config.GsiConfig) Server {
	return &server{
		config.Port,
		NewStore(time.Duration(config.TTL) * time.Second),
		nil,
		nil,
	}
}

func (server *server) Start() error {
	router := mux.NewRouter()
	router.Path("/get").Methods("GET").HandlerFunc(server.handleGsiGet)
	router.Path("/update").Methods("POST").HandlerFunc(server.handleGsiUpdate)
	router.Path("/websocket/{authToken}").Methods("GET").HandlerFunc(server.handleGsiWebsocket)
	router.PathPrefix("/").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("Unhandled route: %s %s\n", request.Method, request.URL)
	})

	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", server.port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	server.upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(request *http.Request) bool {
			return true
		},
	}

	return server.httpServer.ListenAndServe()
}

func (server *server) Stop() error {
	return server.httpServer.Shutdown(context.Background())
}

func (server *server) handleGsiGet(writer http.ResponseWriter, request *http.Request) {
	if !strings.HasPrefix(request.Header.Get("Authorization"), "GSI ") {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	authToken := request.Header.Get("Authorization")[4:]
	gameState, hasGameState := server.store.Get(authToken)
	if !hasGameState {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(gameState)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if _, err = writer.Write(response); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (server *server) handleGsiUpdate(writer http.ResponseWriter, request *http.Request) {
	body, ioError := ioutil.ReadAll(request.Body)
	if ioError != nil || body == nil || len(body) <= 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	gameState := new(GameState)
	if jsonError := json.Unmarshal(body, gameState); jsonError != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	authToken := gameState.Auth.Token
	gameState.Auth = nil

	if gameState.Provider != nil {
		server.store.Put(authToken, gameState)
	} else {
		server.store.Remove(authToken)
	}

	writer.WriteHeader(http.StatusOK)
}

func (server *server) handleGsiWebsocket(writer http.ResponseWriter, request *http.Request) {
	authToken, authTokenPresent := mux.Vars(request)["authToken"]
	if !authTokenPresent {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := server.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		_ = conn.Close()
		return
	}

	channel := server.store.Channel(authToken)

	for {
		select {
		case gameState, more := <-channel:
			if err := conn.WriteJSON(gameState); err != nil || !more {
				_ = conn.Close()
				return
			}
		}
	}
}
