package gsi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

// This interfaces defines the public API of the GSI client. The client can be used to retrieve information about the
// current game state of a player, by connecting to a running GSI server. It handles authentication automatically.
type Client interface {
	// Retrieves the game state for the player that this client connects to.
	GetGameState(name string) (*GameState, error)
}

type client struct {
	Host       string
	Port       int
	AuthTokens map[string]string
}

func NewClient(host string, port int, authTokens map[string]string) Client {
	return &client{
		Host:       host,
		Port:       port,
		AuthTokens: authTokens,
	}
}

func (c *client) GetGameState(name string) (*GameState, error) {
	response, restErr := resty.New().
		R().
		SetHeader("Authorization", fmt.Sprintf("GSI %s", c.AuthTokens[name])).
		Get(fmt.Sprintf("http://%s:%d/get", c.Host, c.Port))
	if restErr != nil {
		log.Println(restErr)
		return nil, restErr
	}

	if response.StatusCode() != 200 {
		errorMessage := fmt.Sprintf("Expected status '%d' but got '%d', with response: %s", 200, response.StatusCode(), response.Body())
		log.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	result := new(GameState)
	if jsonErr := json.Unmarshal(response.Body(), result); jsonErr != nil {
		log.Println(jsonErr)
		return nil, jsonErr
	}

	return result, nil
}
