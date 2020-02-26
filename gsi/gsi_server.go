package gsi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	verificationToken string
	ttl               time.Duration
	port              int

	gameState  *GameState
	lastUpdate time.Time
}

func CreateServer(verificationToken string, ttl time.Duration) *Server {
	checkTTL := ttl
	if ttl == 0 {
		checkTTL = time.Duration(999)
	}

	return &Server{
		verificationToken: verificationToken,
		ttl:               checkTTL,
	}
}

func (server *Server) ListenAndServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/update", server.handleGsiUpdate)
	mux.HandleFunc("/get", server.handleGsiGet)
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("Unhandled route: %s %s\n", request.Method, request.URL)
	})

	return http.ListenAndServe(":8337", mux)
}

func (server *Server) handleGsiUpdate(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("GSI: Method not allowed from %s\n", request.Host)
		return
	}
	if request.Body == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Printf("GSI: No body from %s\n", request.Host)
		return
	}

	body, ioError := ioutil.ReadAll(request.Body)
	if ioError != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Printf("GSI: Empty body from %s\n", request.Host)
		return
	}

	gameState := new(GameState)
	if jsonError := json.Unmarshal(body, gameState); jsonError != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Printf("GSI: Bad body from %s\n", request.Host)
		return
	}

	if gameState.Auth.Token != server.verificationToken {
		writer.WriteHeader(http.StatusForbidden)
		log.Printf("GSI: Invalid toke from %s\n", request.Host)
		return
	}

	if isValidGameState(gameState) {
		server.gameState = gameState
		server.gameState.Auth = nil
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
	if gameState.Player == nil || gameState.Map == nil {
		return false
	}

	return startsWithAny(gameState.Map.Name, []string{"kz", "kzpro", "skz", "vnl", "xc"})
}

func startsWithAny(s string, prefixes []string) bool {
	if len(prefixes) == 0 {
		return false
	}

	return strings.HasPrefix(s, prefixes[0]) || startsWithAny(s, prefixes[1:])
}
