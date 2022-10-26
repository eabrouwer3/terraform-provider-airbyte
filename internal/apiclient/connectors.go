package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SourceIdBody struct {
	SourceId string `json:"sourceId,omitempty"`
}
type DestinationIdBody struct {
	DestinationId string `json:"destinationId,omitempty"`
}

type CommonConnectorFields struct {
	Name                    string          `json:"name"`
	ConnectionConfiguration json.RawMessage `json:"connectionConfiguration"`
}

type Connector struct {
	SourceIdBody
	DestinationIdBody
	SourceDefinitionIdBody
	DestinationDefinitionIdBody
	WorkspaceIdBody
	CommonConnectorFields
	SourceName      string `json:"sourceName"`
	DestinationName string `json:"destinationName"`
	Icon            string `json:"icon"`
}

type NewConnector struct {
	SourceDefinitionIdBody
	DestinationDefinitionIdBody
	WorkspaceIdBody
	CommonConnectorFields
}

type UpdatedConnector struct {
	SourceIdBody
	DestinationIdBody
	CommonConnectorFields
}

func (c *ApiClient) GetConnectorById(connectorId string, t ConnectorType) (*Connector, error) {
	var (
		rb      []byte
		err     error
		urlPath string
	)
	if t == SourceType {
		rb, err = json.Marshal(SourceIdBody{SourceId: connectorId})
		urlPath = "sources"
	} else if t == DestinationType {
		rb, err = json.Marshal(DestinationIdBody{DestinationId: connectorId})
		urlPath = "destinations"
	} else {
		err = fmt.Errorf("invalid ConnectorType: %d", t)
	}
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/%s/get", c.HostURL, BaseUrl, urlPath), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	s := Connector{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (c *ApiClient) CreateConnector(newConnector NewConnector, t ConnectorType) (*Connector, error) {
	rb, err := json.Marshal(newConnector)
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

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/%s/create", c.HostURL, BaseUrl, urlPath), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	s := Connector{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (c *ApiClient) UpdateConnector(updatedConnector UpdatedConnector, t ConnectorType) (*Connector, error) {
	rb, err := json.Marshal(updatedConnector)
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

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/%s/update", c.HostURL, BaseUrl, urlPath), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	s := Connector{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (c *ApiClient) DeleteConnector(connectorId string, t ConnectorType) error {
	var (
		rb      []byte
		err     error
		urlPath string
	)
	if t == SourceType {
		rb, err = json.Marshal(SourceIdBody{SourceId: connectorId})
		urlPath = "sources"
	} else if t == DestinationType {
		rb, err = json.Marshal(DestinationIdBody{DestinationId: connectorId})
		urlPath = "destinations"
	} else {
		err = fmt.Errorf("invalid ConnectorType: %d", t)
	}
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/%s/delete", c.HostURL, BaseUrl, urlPath), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
