package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type OperationIdBody struct {
	OperationId string `json:"operationId"`
}

type Operation struct {
	WorkspaceIdBody
	OperationIdBody
	CommonOperationFields
}

type NewOperation struct {
	WorkspaceIdBody
	CommonOperationFields
}

type UpdatedOperation struct {
	OperationIdBody
	CommonOperationFields
}

type CommonOperationFields struct {
	Name                  string          `json:"name"`
	OperatorConfiguration OperationConfig `json:"operatorConfiguration"`
}

type OperationConfig struct {
	OperatorType  string               `json:"operatorType"`
	Normalization *NormalizationOption `json:"normalization,omitempty"`
	Dbt           *DbtConfig           `json:"dbt,omitempty"`
	Webhook       *WebhookConfig       `json:"webhook,omitempty"`
}

type NormalizationOption struct {
	Option string `json:"option,omitempty"`
}

type DbtConfig struct {
	GitRepoUrl    string `json:"gitRepoUrl"`
	GitRepoBranch string `json:"gitRepoBranch,omitempty"`
	DockerImage   string `json:"dockerImage,omitempty"`
	DbtArguments  string `json:"dbtArguments,omitempty"`
}

type WebhookConfig struct {
	ExecutionUrl    string `json:"executionUrl"`
	ExecutionBody   string `json:"executionBody,omitempty"`
	WebhookConfigId string `json:"webhookConfigId,omitempty"`
}

type OperationCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (c *ApiClient) GetOperationById(operationId string) (*Operation, error) {
	rb, err := json.Marshal(OperationIdBody{
		OperationId: operationId,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/operations/get", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	operation := Operation{}
	err = json.Unmarshal(body, &operation)
	if err != nil {
		return nil, err
	}

	return &operation, nil
}

func (c *ApiClient) CheckOperation(opCfg OperationConfig) (*OperationCheckResponse, error) {
	rb, err := json.Marshal(opCfg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/operations/check", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	check := OperationCheckResponse{}
	err = json.Unmarshal(body, &check)
	if err != nil {
		return nil, err
	}

	return &check, nil
}

func (c *ApiClient) CreateOperation(newOperation NewOperation) (*Operation, error) {
	rb, err := json.Marshal(newOperation)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/operations/create", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	operation := Operation{}
	err = json.Unmarshal(body, &operation)
	if err != nil {
		return nil, err
	}

	return &operation, nil
}

func (c *ApiClient) UpdateOperation(updatedOperation UpdatedOperation) (*Operation, error) {
	rb, err := json.Marshal(updatedOperation)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/operations/update", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	operation := Operation{}
	err = json.Unmarshal(body, &operation)
	if err != nil {
		return nil, err
	}

	return &operation, nil
}

func (c *ApiClient) DeleteOperation(operationId string) error {
	rb, err := json.Marshal(OperationIdBody{OperationId: operationId})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/operations/delete", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
