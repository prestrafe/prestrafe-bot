package gsi

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

func GetGameState() (*GameState, error) {
	response, restErr := resty.New().
		R().
		Get("http://localhost:8337/get")
	if restErr != nil {
		return nil, restErr
	}

	if response.StatusCode() != 200 {
		return nil, errors.New(fmt.Sprintf("Expected status '%d' but got '%d', with response: %s", 200, response.StatusCode(), response.Body()))
	}

	result := new(GameState)
	if jsonErr := json.Unmarshal(response.Body(), result); jsonErr != nil {
		return nil, jsonErr
	}

	return result, nil
}
