package gsi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	ttl  time.Duration
	port int

	gameState  *GameState
	lastUpdate time.Time
}

func CreateServer(ttl time.Duration) *Server {
	checkTTL := ttl
	if ttl == 0 {
		checkTTL = time.Duration(999)
	}

	return &Server{ttl: checkTTL}
}

func (server *Server) ListenAndServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/update", server.handleGsiUpdate)
	mux.HandleFunc("/get", server.handleGsiGet)

	return http.ListenAndServe(":8337", mux)
}

func (server *Server) handleGsiUpdate(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if request.Body == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	body, ioError := ioutil.ReadAll(request.Body)
	if ioError != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	gameState := new(GameState)
	if jsonError := json.Unmarshal(body, gameState); jsonError != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if isValidGameState(gameState) {
		server.gameState = gameState
	} else {
		server.gameState = nil
	}

	server.lastUpdate = time.Now()
	writer.WriteHeader(http.StatusOK)
}

func (server *Server) handleGsiGet(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if server.lastUpdate.Before(time.Now().Add(-server.ttl)) {
		server.gameState = nil
	}

	if server.gameState == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(server.gameState)
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

func isValidGameState(gameState *GameState) bool {
	hasPlayer := gameState.Player != nil
	hasMap := gameState.Map != nil
	hasKzMap := startsWithAny(gameState.Map.Name, []string{"kz", "kzpro", "skz", "vnl", "xc"})
	return hasPlayer && hasMap && hasKzMap
}

func startsWithAny(s string, prefixes []string) bool {
	if len(prefixes) == 0 {
		return false
	}

	return strings.HasPrefix(s, prefixes[0]) || startsWithAny(s, prefixes[1:])
}