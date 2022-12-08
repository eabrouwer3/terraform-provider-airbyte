package apiclient

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

const BaseUrl = "api/v1"

type ApiClient struct {
	HostURL           string // http://localhost:8000
	Username          string
	Password          string
	HTTPClient        *retryablehttp.Client
	AdditionalHeaders map[string]string
}

type HealthCheckResponse struct {
	Available bool `json:"available"`
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

func (c *ApiClient) Check() error {
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

	if !hcr.Available {
		return fmt.Errorf("url: %s, available: %t, body: %s", req.URL, hcr.Available, body)
	}

	return nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *ApiClient) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	if c.Username != "" && c.Password != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth(c.Username, c.Password)))
	}
	for k, v := range c.AdditionalHeaders {
		if strings.ToLower(k) == "host" {
			req.Host = v
		} else {
			req.Header.Set(k, v)
		}
	}

	retryableReq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(retryableReq)
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
			fmt.Println(err)
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
		return nil, fmt.Errorf(
			"url: %s, status: %d, body: %s",
			req.URL,
			res.StatusCode,
			body,
		)
	}

	return body, err
}
