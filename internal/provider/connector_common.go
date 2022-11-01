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
		data.Id = types.StringValue(connector.SourceId)
		data.DefinitionId = types.StringValue(connector.SourceDefinitionId)
		data.DefinitionName = types.StringValue(connector.SourceName)
	} else if connector.DestinationId != "" {
		data.Id = types.StringValue(connector.DestinationId)
		data.DefinitionId = types.StringValue(connector.DestinationDefinitionId)
		data.DefinitionName = types.StringValue(connector.DestinationName)
	} else {
		return nil, fmt.Errorf("either SourceId or DestinationId must be set, not empty")
	}

	data.WorkspaceId = types.StringValue(connector.WorkspaceId)
	data.Name = types.StringValue(connector.Name)
	data.Icon = types.StringValue(connector.Icon)

	return &data, nil
}

func GetCommonConnectorFields(data ConnectorModel) apiclient.CommonConnectorFields {
	return apiclient.CommonConnectorFields{
		Name:                    data.Name.ValueString(),
		ConnectionConfiguration: json.RawMessage(data.ConnectionConfiguration.ValueString()),
	}
}
