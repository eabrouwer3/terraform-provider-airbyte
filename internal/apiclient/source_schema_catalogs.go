package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SourceSchemaCatalog struct {
	Catalog   SyncCatalog `json:"catalog"`
	JobInfo   JobInfo     `json:"jobInfo"`
	catalogId string
}

type SyncCatalog struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	Stream SourceStreamSchema      `json:"stream"`
	Config DestinationStreamConfig `json:"config"`
}

type SourceStreamSchema struct {
	// Required Field
	Name string `json:"name"`
	// Optional Fields
	JsonSchema              json.RawMessage `json:"jsonSchema,omitempty"`
	SupportedSyncModes      []string        `json:"supportedSyncModes,omitempty"`
	SourceDefinedCursor     *bool           `json:"sourceDefinedCursor,omitempty"`
	DefaultCursorField      []string        `json:"defaultCursorField,omitempty"`
	SourceDefinedPrimaryKey [][]string      `json:"sourceDefinedPrimaryKey,omitempty"`
	Namespace               string          `json:"namespace,omitempty"`
}

type DestinationStreamConfig struct {
	// Required Fields
	SyncMode            string `json:"syncMode"`
	DestinationSyncMode string `json:"destinationSyncMode"`
	// Optional Fields
	CursorField []string   `json:"cursorField,omitempty"`
	PrimaryKey  [][]string `json:"primaryKey,omitempty"`
	AliasName   string     `json:"aliasName,omitempty"`
	Selected    *bool      `json:"selected,omitempty"`
}

func (c *ApiClient) GetSourceSchemaCatalogById(sourceId string) (*SourceSchemaCatalog, error) {
	rb, err := json.Marshal(SourceIdBody{SourceId: sourceId})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/sources/discover_schema", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	connection := SourceSchemaCatalog{}
	err = json.Unmarshal(body, &connection)
	if err != nil {
		return nil, err
	}

	if !connection.JobInfo.Succeeded {
		return nil, fmt.Errorf("job to get source schema catalog failed")
	}

	return &connection, nil
}
