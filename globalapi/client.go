package globalapi

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

const (
	apiEndpointBase = "https://kztimerglobal.com/api"
	apiVersion      = "2.0"
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
	baseUrl string
	rest    *resty.Client
	logger  *log.Logger
}

func NewClient(apiToken string) Client {
	logger := log.New(os.Stdout, "Global API > ", log.LstdFlags)

	restClient := resty.New().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json")
	if apiToken != "" {
		restClient.SetHeader("X-ApiKey", apiToken)
	}

	return &client{
		fmt.Sprintf("%s/v%s", apiEndpointBase, apiVersion),
		restClient,
		logger,
	}
}

func (c *client) Get(path string, result interface{}) error {
	return c.GetWithParameters(path, nil, result)
}

func (c *client) GetWithParameters(path string, queryParams QueryParameters, result interface{}) error {
	response, restErr := resty.New().
		R().
		SetQueryParams(queryParams).
		Get(fmt.Sprintf("%s/%s", c.baseUrl, path))

	if restErr != nil {
		return restErr
	}

	c.logger.Printf("%s -> Status: %d, Body: %s\n", response.Request.URL, response.StatusCode(), response.Body())
	if response.StatusCode() != 200 {
		return fmt.Errorf("Expected status '%d' but got '%d'", 200, response.StatusCode())
	}

	if jsonErr := json.Unmarshal(response.Body(), result); jsonErr != nil {
		return jsonErr
	}

	return nil
}
