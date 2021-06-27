package smclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

type Client interface {
	// Retrieves the game state for the player that this client connects to.
	GetPlayerInfo() (*FullPlayerInfo, error)
}

type client struct {
	host        string
	port        int
	serverToken string
}

func New(host string, port int, serverToken string) Client {
	return &client{host, port, serverToken}
}

func (c *client) GetPlayerInfo() (*FullPlayerInfo, error) {
	response, restErr := resty.New().
		R().
		SetHeader("Authorization", fmt.Sprintf("SM %s", c.serverToken)).
		Get(fmt.Sprintf("http://%s:%d/sm/get", c.host, c.port))
	if restErr != nil {
		log.Println(restErr)
		return nil, restErr
	}

	if response.StatusCode() != 200 {
		errorMessage := fmt.Sprintf("Expected status '%d' but got '%d', with response: %s", 200, response.StatusCode(), response.Body())
		log.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	result := new(FullPlayerInfo)
	if jsonErr := json.Unmarshal(response.Body(), result); jsonErr != nil {
		log.Println(jsonErr)
		return nil, jsonErr
	}

	return result, nil
}