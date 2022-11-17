package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type JobInfo struct {
	Succeeded  bool   `json:"succeeded"`
	Id         string `json:"id"`
	configType string
	configId   string
	createdAt  int
	endedAt    int
	logs       map[string][]string
}

type CheckConnectionResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	JobInfo JobInfo `json:"jobInfo"`
}

func (c *ApiClient) CheckNewConnector(connector NewConnector, t ConnectorType) (*CheckConnectionResponse, error) {
	// This API endpoint takes everything except the name, which it'll yell about, so we do this
	connectorToCheck := connector
	connectorToCheck.Name = ""

	rb, err := json.Marshal(connectorToCheck)
	var urlPath string
	if t == SourceType {
		urlPath = "sources"
	} else if t == DestinationType {
		urlPath = "destinations"
	} else {
		err = fmt.Errorf("invalid ConnectorType: %d", t)
	}
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/scheduler/%s/check_connection", c.HostURL, BaseUrl, urlPath), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	res := CheckConnectionResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *ApiClient) CheckUpdatedConnector(connector UpdatedConnector, t ConnectorType) (*CheckConnectionResponse, error) {
	rb, err := json.Marshal(connector)
	var urlPath string
	if t == SourceType {
		urlPath = "sources"
	} else if t == DestinationType {
		urlPath = "destinations"
	} else {
		err = fmt.Errorf("invalid ConnectorType: %d", t)
	}
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/%s/check_connection_for_update", c.HostURL, BaseUrl, urlPath), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	res := CheckConnectionResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
