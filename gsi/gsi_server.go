package gsi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

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
}

// Creates a new GSI server.
func NewServer(config *config.GsiConfig) Server {
	return &server{
		config.Port,
		time.Duration(config.TTL) * time.Second,
		make(map[string]*GameState),
		make(map[string]time.Time),
	}
}

func (server *server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/update", server.handleGsiUpdate)
	mux.HandleFunc("/get", server.handleGsiGet)
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("Unhandled route: %s %s\n", request.Method, request.URL)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", server.port), mux)
}

func (server *server) handleGsiUpdate(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("GSI-UPDATE: Method not allowed from %s\n", request.Host)
		return
	}
	if request.Body == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Printf("GSI-UPDATE: No body from %s\n", request.Host)
		return
	}

	body, ioError := ioutil.ReadAll(request.Body)
	if ioError != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Printf("GSI-UPDATE: Empty body from %s\n", request.Host)
		return
	}

	gameState := new(GameState)
	if jsonError := json.Unmarshal(body, gameState); jsonError != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Printf("GSI-UPDATE: Bad body from %s\n", request.Host)
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

func (server *server) handleGsiGet(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("GSI-GET: Method not allowed from %s\n", request.Host)
		return
	}

	if !strings.HasPrefix(request.Header.Get("Authorization"), "GSI ") {
		writer.WriteHeader(http.StatusUnauthorized)
		log.Printf("GSI-GET: No GSI token provided %s\n", request.Host)
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
		log.Printf("GSI-GET: Could not serialize game state %s\n", request.Host)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if _, err = writer.Write(response); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("GSI-GET: Could not write game state %s\n", request.Host)
		return
	}
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
