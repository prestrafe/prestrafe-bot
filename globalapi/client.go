package globalapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

const (
	apiEndpointBase = "https://kztimerglobal.com/api"
	apiVersion      = "1.0"
)

// A type definition for query parameter maps. These can be passed to many methods of the global API to specify which
// data should be retrieved.
type QueryParameters map[string]string

// Defines the API for clients of the global KZ API. It handles some common logic like authentication, parameter
// injection and the actual HTTP communication.
type Client interface {
	// Gets the result of an API call to the given path.
	Get(path string, result interface{}) error
	// Gets the result of an API call to the given path, passing the given parameters.
	GetWithParameters(path string, queryParams QueryParameters, result interface{}) error
}

type client struct {
	baseUrl  string
	version  string
	apiToken string
	logger   *log.Logger
}

func NewClient(apiToken string) Client {
	logger := log.New(os.Stdout, "Global API > ", log.LstdFlags)
	return &client{apiEndpointBase, apiVersion, apiToken, logger}
}

func (c *client) Get(path string, result interface{}) error {
	return c.GetWithParameters(path, nil, result)
}

func (c *client) GetWithParameters(path string, queryParams QueryParameters, result interface{}) error {
	apiEndpoint := fmt.Sprintf("%s/v%s/%s", c.baseUrl, c.version, path)

	response, restErr := resty.New().
		R().
		SetQueryParams(queryParams).
		Get(apiEndpoint)

	if restErr != nil {
		return restErr
	}

	if response.StatusCode() != 200 {
		c.logger.Printf("%s -> Status: %d, Body: %s\n", apiEndpoint, response.StatusCode(), response.Body())
		return errors.New(fmt.Sprintf("Expected status '%d' but got '%d'", 200, response.StatusCode()))
	}

	if jsonErr := json.Unmarshal(response.Body(), result); jsonErr != nil {
		return jsonErr
	}

	return nil
}
