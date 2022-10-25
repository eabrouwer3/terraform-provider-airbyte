package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SourceIdBody struct {
	SourceId string `json:"sourceId"`
}

type CommonSourceFields struct {
	Name                    string          `json:"name"`
	ConnectionConfiguration json.RawMessage `json:"connectionConfiguration"`
}

type Source struct {
	SourceIdBody
	SourceDefinitionIdBody
	WorkspaceIdBody
	CommonSourceFields
	SourceName string `json:"sourceName"`
	Icon       string `json:"icon"`
}

type NewSource struct {
	SourceDefinitionIdBody
	WorkspaceIdBody
	CommonSourceFields
}

type UpdatedSource struct {
	SourceIdBody
	CommonSourceFields
}

func (c *ApiClient) GetSourceById(sourceId string) (*Source, error) {
	rb, err := json.Marshal(SourceIdBody{SourceId: sourceId})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/sources/get", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	s := Source{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (c *ApiClient) CreateSource(newSource NewSource) (*Source, error) {
	rb, err := json.Marshal(newSource)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/sources/create", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	s := Source{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (c *ApiClient) UpdateSource(updatedSource UpdatedSource) (*Source, error) {
	rb, err := json.Marshal(updatedSource)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/sources/update", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	s := Source{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (c *ApiClient) DeleteSource(sourceId string) error {
	rb, err := json.Marshal(SourceIdBody{SourceId: sourceId})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/sources/delete", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
