package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ConnectionIdBody struct {
	ConnectionId string `json:"connectionId"`
}

type Connection struct {
	ConnectionIdBody
	CommonConnectionFields
	SourceIdBody
	DestinationIdBody
	Geography string `json:"geography,omitempty"`
}

type NewConnection struct {
	CommonConnectionFields
	SourceIdBody
	DestinationIdBody
}

type UpdatedConnection struct {
	ConnectionIdBody
	CommonConnectionFields
}

type CommonConnectionFields struct {
	// Required Fields
	Status string `json:"status"`
	// Optional Fields
	Name                 string                       `json:"name,omitempty"`
	NamespaceDefinition  string                       `json:"namespaceDefinition,omitempty"`
	NamespaceFormat      string                       `json:"namespaceFormat,omitempty"`
	Prefix               string                       `json:"prefix,omitempty"`
	OperationIds         []string                     `json:"operationIds,omitempty"`
	SyncCatalog          *SyncCatalog                 `json:"syncCatalog,omitempty"`
	ScheduleType         string                       `json:"scheduleType,omitempty"`
	ScheduleData         *ScheduleData                `json:"scheduleData,omitempty"`
	ResourceRequirements *ResourceRequirementsOptions `json:"resourceRequirements,omitempty"`
	SourceCatalogId      string                       `json:"sourceCatalogId,omitempty"`
	BreakingChange       *bool                        `json:"breakingChange,omitempty"`
}

type ScheduleSpec struct {
	// Both are required
	Units    int64  `json:"units"`
	TimeUnit string `json:"timeUnit"`
}

type ScheduleData struct {
	BasicSchedule *ScheduleSpec     `json:"basicSchedule,omitempty"`
	Cron          *CronScheduleSpec `json:"cron,omitempty"`
}

type CronScheduleSpec struct {
	// Both are required
	CronExpression string `json:"cronExpression"`
	CronTimeZone   string `json:"cronTimeZone"`
}

func (c *ApiClient) GetConnectionById(connectionId string) (*Connection, error) {
	rb, err := json.Marshal(ConnectionIdBody{
		ConnectionId: connectionId,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/connections/get", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	connection := Connection{}
	err = json.Unmarshal(body, &connection)
	if err != nil {
		return nil, err
	}

	return &connection, nil
}

func (c *ApiClient) CreateConnection(newConnection NewConnection) (*Connection, error) {
	rb, err := json.Marshal(newConnection)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/connections/create", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	connection := Connection{}
	err = json.Unmarshal(body, &connection)
	if err != nil {
		return nil, err
	}

	return &connection, nil
}

func (c *ApiClient) UpdateConnection(updatedConnection UpdatedConnection) (*Connection, error) {
	rb, err := json.Marshal(updatedConnection)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/connections/update", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	connection := Connection{}
	err = json.Unmarshal(body, &connection)
	if err != nil {
		return nil, err
	}

	return &connection, nil
}

func (c *ApiClient) DeleteConnection(connectionId string) error {
	rb, err := json.Marshal(ConnectionIdBody{ConnectionId: connectionId})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/connections/delete", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
