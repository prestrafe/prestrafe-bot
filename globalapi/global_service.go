package globalapi

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type QueryParameters map[string]string

func globalApiGet(path string, result interface{}, queryParams QueryParameters) error {
	response, restErr := resty.New().
		R().
		SetQueryParams(queryParams).
		Get("https://kztimerglobal.com/api/v1.0/" + path)
	if restErr != nil {
		return restErr
	}

	if response.StatusCode() != 200 {
		return errors.New(fmt.Sprintf("Expected status '%d' but got '%d', with response: %s", 200, response.StatusCode(), response.Body()))
	}
	if jsonErr := json.Unmarshal(response.Body(), result); jsonErr != nil {
		return jsonErr
	}

	return nil
}
