package apiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BaseUrl = "api/v1"

type ApiClient struct {
	HostURL    string // http://localhost:8000
	HTTPClient *http.Client
}

type HealthCheckResponse struct {
	available bool
}

type CommonErrorResponseFields struct {
	Message            string `json:"message"`
	exceptionClassName string
	exceptionStack     []string
}

type Response422 struct {
	CommonErrorResponseFields
	ValidationErrors []ValidationError `json:"validationErrors"`
}

type ValidationError struct {
	PropertyPath string `json:"propertyPath"`
	InvalidValue string `json:"invalidValue"`
	Message      string `json:"message"`
}

type Response404 struct {
	Id string `json:"id"`
	CommonErrorResponseFields
	rootCauseExceptionClassName string
	rootCauseExceptionStack     []string
}

type Response500 struct {
	CommonErrorResponseFields
}

func (c *ApiClient) check() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/health", c.HostURL, BaseUrl), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	hcr := HealthCheckResponse{}
	err = json.Unmarshal(body, &hcr)
	if err != nil {
		return err
	}

	return nil
}

func (c *ApiClient) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		if res.StatusCode == http.StatusUnprocessableEntity {
			r := Response422{}
			err := json.Unmarshal(body, &r)
			if err == nil {
				body, _ = json.Marshal(r)
			}
		} else if res.StatusCode == http.StatusNotFound {
			r := Response404{}
			err := json.Unmarshal(body, &r)
			if err == nil {
				body, _ = json.Marshal(r)
			}
		} else if res.StatusCode == http.StatusInternalServerError {
			r := Response500{}
			err := json.Unmarshal(body, &r)
			if err == nil {
				body, _ = json.Marshal(r)
			}
		}
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
