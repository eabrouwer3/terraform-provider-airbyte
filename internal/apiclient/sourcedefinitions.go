package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SourceDefinitionIdBody struct {
	SourceDefinitionId string `json:"sourceDefinitionId"`
}

type SourceDefinition struct {
	SourceDefinitionIdBody
	CommonSourceDefinitionFields
	ProtocolVersion string `json:"protocolVersion,omitempty"`
	ReleaseStage    string `json:"releaseStage,omitempty"`
	ReleaseDate     string `json:"releaseDate,omitempty"`
	SourceType      string `json:"sourceType,omitempty"`
}

type CommonSourceDefinitionFields struct {
	Name                 string                `json:"name"`
	DockerRepository     string                `json:"dockerRepository"`
	DockerImageTag       string                `json:"dockerImageTag"`
	DocumentationUrl     string                `json:"documentationUrl"`
	Icon                 string                `json:"icon,omitempty"`
	ResourceRequirements *ResourceRequirements `json:"resourceRequirements,omitempty"`
}

type NewSourceDefinition = CommonSourceDefinitionFields

type UpdatedSourceDefinition struct {
	SourceDefinitionIdBody
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

func (c *ApiClient) GetSourceDefinitionById(sourceDefinitionId string) (*SourceDefinition, error) {
	rb, err := json.Marshal(struct {
		SourceDefinitionId string `json:"sourceDefinitionId"`
	}{
		SourceDefinitionId: sourceDefinitionId,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/source_definitions/get", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sd := SourceDefinition{}
	err = json.Unmarshal(body, &sd)
	if err != nil {
		return nil, err
	}

	return &sd, nil
}

// TODO: Implement this

func (c *ApiClient) GetSourceDefinitionSpec(sourceDefinitionId string) (*SourceDefinition, error) {
	rb, err := json.Marshal(struct {
		SourceDefinitionId string `json:"sourceDefinitionId"`
	}{
		SourceDefinitionId: sourceDefinitionId,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/source_definitions/get", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sd := SourceDefinition{}
	err = json.Unmarshal(body, &sd)
	if err != nil {
		return nil, err
	}

	return &sd, nil
}

func (c *ApiClient) CreateSourceDefinition(newSourceDefinition NewSourceDefinition) (*SourceDefinition, error) {
	rb, err := json.Marshal(newSourceDefinition)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/source_definitions/create", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sd := SourceDefinition{}
	err = json.Unmarshal(body, &sd)
	if err != nil {
		return nil, err
	}

	return &sd, nil
}

func (c *ApiClient) UpdateSourceDefinition(updatedSourceDefinition UpdatedSourceDefinition) (*SourceDefinition, error) {
	rb, err := json.Marshal(updatedSourceDefinition)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/source_definitions/update", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sd := SourceDefinition{}
	err = json.Unmarshal(body, &sd)
	if err != nil {
		return nil, err
	}

	return &sd, nil
}

func (c *ApiClient) DeleteSourceDefinition(sourceDefinitionId string) error {
	rb, err := json.Marshal(SourceDefinitionIdBody{SourceDefinitionId: sourceDefinitionId})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/source_definitions/delete", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
