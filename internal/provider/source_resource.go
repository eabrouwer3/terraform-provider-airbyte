package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &SourceResource{}
var _ resource.ResourceWithImportState = &SourceResource{}

func NewSourceResource() resource.Resource {
	return &SourceResource{}
}

// SourceResource defines the resource implementation.
type SourceResource struct {
	client *apiclient.ApiClient
}

func (r *SourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (r *SourceResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Source resource",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Source ID",
				Type:        types.StringType,
				Computed:    true,
			},
			"definition_id": {
				Description: "Source Definition ID",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"workspace_id": {
				Description: "Workspace ID",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"name": {
				Description: "Source Name",
				Type:        types.StringType,
				Required:    true,
			},
			"connection_configuration": {
				Description: "Connection Configuration",
				Type:        types.StringType,
				Required:    true,
				Sensitive:   true,
			},
			"definition_name": {
				Description: "Source Definition Name",
				Type:        types.StringType,
				Computed:    true,
			},
			"icon": {
				Description: "Icon SVG/URL",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

func (r *SourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(apiclient.ApiClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *apiclient.ApiClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = &client
}

func (r *SourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConnectorModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newSource := apiclient.NewConnector{
		SourceDefinitionIdBody: apiclient.SourceDefinitionIdBody{
			SourceDefinitionId: plan.DefinitionId.ValueString(),
		},
		WorkspaceIdBody: apiclient.WorkspaceIdBody{
			WorkspaceId: plan.WorkspaceId.ValueString(),
		},
		CommonConnectorFields: GetCommonConnectorFields(plan),
	}

	source, err := r.client.CreateConnector(newSource, apiclient.SourceType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Source",
			"Could not create Source, unexpected error: "+err.Error(),
		)
		return
	}

	state, err := FlattenConnector(source)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Source",
			"Could not create Source, unexpected error: "+err.Error(),
		)
		return
	}
	state.ConnectionConfiguration = plan.ConnectionConfiguration

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan ConnectorModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sourceId := plan.Id.ValueString()

	source, err := r.client.GetConnectorById(sourceId, apiclient.SourceType)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Source, got error: %s", err))
		return
	}

	state, err := FlattenConnector(source)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read Source, got error: %s", err),
		)
		return
	}
	state.ConnectionConfiguration = plan.ConnectionConfiguration

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConnectorModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updatedSource := apiclient.UpdatedConnector{
		SourceIdBody:          apiclient.SourceIdBody{SourceId: plan.Id.ValueString()},
		CommonConnectorFields: GetCommonConnectorFields(plan),
	}

	source, err := r.client.UpdateConnector(updatedSource, apiclient.SourceType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating source",
			"Could not update Source, unexpected error: "+err.Error(),
		)
		return
	}

	state, err := FlattenConnector(source)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating source",
			"Could not update Source, unexpected error: "+err.Error(),
		)
		return
	}
	state.ConnectionConfiguration = plan.ConnectionConfiguration

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ConnectorModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sourceId := state.Id.ValueString()
	err := r.client.DeleteConnector(sourceId, apiclient.SourceType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating source",
			"Could not update Source, unexpected error: "+err.Error(),
		)
	}
}

func (r *SourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
