package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DestinationDefinitionIdBody struct {
	DestinationDefinitionId string `json:"destinationDefinitionId,omitempty"`
}
type SourceDefinitionIdBody struct {
	SourceDefinitionId string `json:"sourceDefinitionId,omitempty"`
}

type ConnectorDefinition struct {
	SourceDefinitionIdBody
	DestinationDefinitionIdBody
	CommonConnectorDefinitionFields
	ProtocolVersion string `json:"protocolVersion,omitempty"`
	ReleaseStage    string `json:"releaseStage,omitempty"`
	ReleaseDate     string `json:"releaseDate,omitempty"`
}

type CommonConnectorDefinitionFields struct {
	Name                 string                `json:"name"`
	DockerRepository     string                `json:"dockerRepository"`
	DockerImageTag       string                `json:"dockerImageTag"`
	DocumentationUrl     string                `json:"documentationUrl"`
	Icon                 string                `json:"icon,omitempty"`
	ResourceRequirements *ResourceRequirements `json:"resourceRequirements,omitempty"`
}

type NewConnectorDefinition = CommonConnectorDefinitionFields

type UpdatedConnectorDefinition struct {
	SourceDefinitionIdBody
	DestinationDefinitionIdBody
	DockerImageTag       string                `json:"dockerImageTag,omitempty"`
	ResourceRequirements *ResourceRequirements `json:"resourceRequirements,omitempty"`
}

type ResourceRequirements struct {
	Default     *ResourceRequirementsOptions       `json:"default,omitempty"`
	JobSpecific *[]JobSpecificResourceRequirements `json:"jobSpecific,omitempty"`
}

type ResourceRequirementsOptions struct {
	CPURequest    string `json:"cpu_request,omitempty"`
	CPULimit      string `json:"cpu_limit,omitempty"`
	MemoryRequest string `json:"memory_request,omitempty"`
	MemoryLimit   string `json:"memory_limit,omitempty"`
}

type JobSpecificResourceRequirements struct {
	JobType              string                      `json:"jobType"`
	ResourceRequirements ResourceRequirementsOptions `json:"resourceRequirements"`
}

type ConnectorType int

const (
	SourceType ConnectorType = iota
	DestinationType
)

func (c *ApiClient) GetConnectorDefinitionById(connectorDefinitionId string, t ConnectorType) (*ConnectorDefinition, error) {
	var (
		rb      []byte
		err     error
		urlPath string
	)
	if t == SourceType {
		rb, err = json.Marshal(SourceDefinitionIdBody{SourceDefinitionId: connectorDefinitionId})
		urlPath = "source_definitions"
	} else if t == DestinationType {
		rb, err = json.Marshal(DestinationDefinitionIdBody{DestinationDefinitionId: connectorDefinitionId})
		urlPath = "destination_definitions"
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

	dd := ConnectorDefinition{}
	err = json.Unmarshal(body, &dd)
	if err != nil {
		return nil, err
	}

	return &dd, nil
}

func (c *ApiClient) CreateConnectorDefinition(newDefinition NewConnectorDefinition, t ConnectorType) (*ConnectorDefinition, error) {
	rb, err := json.Marshal(newDefinition)
	var urlPath string
	if t == SourceType {
		urlPath = "source_definitions"
	} else if t == DestinationType {
		urlPath = "destination_definitions"
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

	dd := ConnectorDefinition{}
	err = json.Unmarshal(body, &dd)
	if err != nil {
		return nil, err
	}

	return &dd, nil
}

func (c *ApiClient) UpdateConnectorDefinition(updatedConnectorDefinition UpdatedConnectorDefinition, t ConnectorType) (*ConnectorDefinition, error) {
	rb, err := json.Marshal(updatedConnectorDefinition)
	var urlPath string
	if t == SourceType {
		urlPath = "source_definitions"
	} else if t == DestinationType {
		urlPath = "destination_definitions"
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

	dd := ConnectorDefinition{}
	err = json.Unmarshal(body, &dd)
	if err != nil {
		return nil, err
	}

	return &dd, nil
}

func (c *ApiClient) DeleteConnectorDefinition(connectorDefinitionId string, t ConnectorType) error {
	var (
		rb      []byte
		err     error
		urlPath string
	)
	if t == SourceType {
		rb, err = json.Marshal(SourceDefinitionIdBody{SourceDefinitionId: connectorDefinitionId})
		urlPath = "source_definitions"
	} else if t == DestinationType {
		rb, err = json.Marshal(DestinationDefinitionIdBody{DestinationDefinitionId: connectorDefinitionId})
		urlPath = "destination_definitions"
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
