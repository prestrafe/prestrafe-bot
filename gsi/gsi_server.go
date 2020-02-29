package gsi

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"

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
	port        int
	ttl         time.Duration
	gameStates  map[string]*GameState
	lastUpdates map[string]time.Time
	httpServer  *http.Server
}

// Creates a new GSI server.
func NewServer(config *config.GsiConfig) Server {
	return &server{
		config.Port,
		time.Duration(config.TTL) * time.Second,
		make(map[string]*GameState),
		make(map[string]time.Time),
		nil,
	}
}

func (server *server) Start() error {
	router := mux.NewRouter()
	router.Path("/gsi").Methods("GET").HandlerFunc(server.handleGsiGet)
	router.Path("/gsi").Methods("POST").HandlerFunc(server.handleGsiUpdate)
	router.PathPrefix("/").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("Unhandled route: %s %s\n", request.Method, request.URL)
	})

	server.httpServer = &http.Server{
		Addr:         "0.0.0.0",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return server.httpServer.ListenAndServe()
}

func (server *server) Stop() error {
	return server.httpServer.Shutdown(context.Background())
}

func (server *server) handleGsiGet(writer http.ResponseWriter, request *http.Request) {
	if !strings.HasPrefix(request.Header.Get("Authorization"), "GSI ") {
		writer.WriteHeader(http.StatusUnauthorized)
	}

	authToken := request.Header.Get("Authorization")[4:]

	if lastUpdate, hasLastUpdate := server.lastUpdates[authToken]; hasLastUpdate {
		if lastUpdate.Before(time.Now().Add(-server.ttl)) {
			delete(server.gameStates, authToken)
		}
	}

	gameState, hasGameState := server.gameStates[authToken]
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

	if isValidGameState(gameState) {
		gameState.Map.Name = cleanupMapName(gameState.Map.Name)

		server.gameStates[authToken] = gameState
		server.lastUpdates[authToken] = time.Now()
	} else {
		delete(server.gameStates, authToken)
		delete(server.lastUpdates, authToken)
	}

	writer.WriteHeader(http.StatusOK)
}

func isValidGameState(gameState *GameState) bool {
	if gameState.Player == nil || gameState.Map == nil {
		return false
	}

	matchString, err := regexp.MatchString("^(workshop/[0-9]+/)?(kz|kzpro|skz|vnl|xc)_.+$", gameState.Map.Name)
	return matchString && err == nil
}

func cleanupMapName(mapName string) string {
	if strings.HasPrefix(mapName, "workshop") {
		return mapName[strings.LastIndex(mapName, "/")+1:]
	}

	return mapName
}
