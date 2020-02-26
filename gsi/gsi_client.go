package gsi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func GetGameState() (*GameState, error) {
	response, restErr := resty.New().
		R().
		Get("http://localhost:8337/get")
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
