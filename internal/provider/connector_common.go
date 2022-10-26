package provider

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
)

// ConnectorModel describes the data connector data model.
type ConnectorModel struct {
	Id             types.String `tfsdk:"id"`
	DefinitionId   types.String `tfsdk:"definition_id"`
	DefinitionName types.String `tfsdk:"definition_name"`
	WorkspaceId    types.String `tfsdk:"workspace_id"`
	Name           types.String `tfsdk:"name"`
	Icon           types.String `tfsdk:"icon"`
	// TODO: Implement something like this: https://www.youtube.com/watch?v=N2oGuPb_CIM
	//       See this ticket: https://github.com/hashicorp/terraform-plugin-framework/issues/147
	ConnectionConfiguration types.String `tfsdk:"connection_configuration"`
}

func FlattenConnector(connector *apiclient.Connector) (*ConnectorModel, error) {
	var data ConnectorModel

	if connector.SourceId != "" {
		data.Id = types.String{Value: connector.SourceId}
		data.DefinitionId = types.String{Value: connector.SourceDefinitionId}
		data.DefinitionName = types.String{Value: connector.SourceName}
	} else if connector.DestinationId != "" {
		data.Id = types.String{Value: connector.DestinationId}
		data.DefinitionId = types.String{Value: connector.DestinationDefinitionId}
		data.DefinitionName = types.String{Value: connector.DestinationName}
	} else {
		return nil, fmt.Errorf("either SourceId or DestinationId must be set, not empty")
	}

	data.WorkspaceId = types.String{Value: connector.WorkspaceId}
	data.Name = types.String{Value: connector.Name}
	data.Icon = types.String{Value: connector.Icon}

	config, err := connector.ConnectionConfiguration.MarshalJSON()
	if err != nil {
		return nil, err
	}
	data.ConnectionConfiguration = types.String{Value: string(config)}

	return &data, nil
}

func GetCommonConnectorFields(data ConnectorModel) apiclient.CommonConnectorFields {
	return apiclient.CommonConnectorFields{
		Name:                    data.Name.Value,
		ConnectionConfiguration: json.RawMessage(data.ConnectionConfiguration.Value),
	}
}
