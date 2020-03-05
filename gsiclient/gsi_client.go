package gsiclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"

	"gitlab.com/prestrafe/prestrafe-gsi"
)

// This interfaces defines the public API of the GSI client. The client can be used to retrieve information about the
// current game state of a player, by connecting to a running GSI server. It handles authentication automatically.
type Client interface {
	// Retrieves the game state for the player that this client connects to.
	GetGameState() (*gsi.GameState, error)
}

type client struct {
	host      string
	port      int
	authToken string
}

func New(host string, port int, authToken string) Client {
	return &client{host, port, authToken}
}

func (c *client) GetGameState() (*gsi.GameState, error) {
	response, restErr := resty.New().
		R().
		SetHeader("Authorization", fmt.Sprintf("GSI %s", c.authToken)).
		Get(fmt.Sprintf("http://%s:%d/", c.host, c.port))
	if restErr != nil {
		log.Println(restErr)
		return nil, restErr
	}

	if response.StatusCode() != 200 {
		errorMessage := fmt.Sprintf("Expected status '%d' but got '%d', with response: %s", 200, response.StatusCode(), response.Body())
		log.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	result := new(gsi.GameState)
	if jsonErr := json.Unmarshal(response.Body(), result); jsonErr != nil {
		log.Println(jsonErr)
		return nil, jsonErr
	}

	return result, nil
}
