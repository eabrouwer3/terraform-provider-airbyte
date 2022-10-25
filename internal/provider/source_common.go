package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
)

// SourceModel describes the data source data model.
type SourceModel struct {
	Id                 types.String `tfsdk:"id"`
	SourceDefinitionId types.String `tfsdk:"source_definition_id"`
	WorkspaceId        types.String `tfsdk:"workspace_id"`
	Name               types.String `tfsdk:"name"`
	SourceName         types.String `tfsdk:"source_name"`
	Icon               types.String `tfsdk:"icon"`
	// TODO: Implement something like this: https://www.youtube.com/watch?v=N2oGuPb_CIM
	//       See this ticket: https://github.com/hashicorp/terraform-plugin-framework/issues/147
	ConnectionConfiguration types.String `tfsdk:"connection_configuration"`
}

func FlattenSource(source *apiclient.Source) (*SourceModel, error) {
	var data SourceModel

	data.Id = types.String{Value: source.SourceId}
	data.SourceDefinitionId = types.String{Value: source.SourceDefinitionId}
	data.WorkspaceId = types.String{Value: source.WorkspaceId}
	data.Name = types.String{Value: source.Name}
	data.SourceName = types.String{Value: source.SourceName}
	data.Icon = types.String{Value: source.Icon}

	config, err := source.ConnectionConfiguration.MarshalJSON()
	if err != nil {
		return nil, err
	}
	data.ConnectionConfiguration = types.String{Value: string(config)}

	return &data, nil
}
