package provider

import (
	"context"
	"encoding/json"
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
			"source_definition_id": {
				Description: "Source ID",
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
			},
			"source_name": {
				Description: "Source Name",
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

func getCommonSourceFields(data SourceModel) apiclient.CommonSourceFields {
	return apiclient.CommonSourceFields{
		Name:                    data.Name.Value,
		ConnectionConfiguration: json.RawMessage(data.ConnectionConfiguration.Value),
	}
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
	var plan SourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newSource := apiclient.NewSource{
		SourceDefinitionIdBody: apiclient.SourceDefinitionIdBody{
			SourceDefinitionId: plan.SourceDefinitionId.Value,
		},
		WorkspaceIdBody: apiclient.WorkspaceIdBody{
			WorkspaceId: plan.WorkspaceId.Value,
		},
		CommonSourceFields: getCommonSourceFields(plan),
	}

	source, err := r.client.CreateSource(newSource)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Source",
			"Could not create Source, unexpected error: "+err.Error(),
		)
		return
	}

	state, err := FlattenSource(source)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Source",
			"Could not create Source, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan SourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sourceId := plan.Id.Value

	source, err := r.client.GetSourceById(sourceId)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Source, got error: %s", err))
		return
	}

	state, err := FlattenSource(source)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read Source, got error: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updatedSource := apiclient.UpdatedSource{
		SourceIdBody:       apiclient.SourceIdBody{SourceId: plan.Id.Value},
		CommonSourceFields: getCommonSourceFields(plan),
	}

	source, err := r.client.UpdateSource(updatedSource)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating source",
			"Could not update Source, unexpected error: "+err.Error(),
		)
		return
	}

	state, err := FlattenSource(source)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating source",
			"Could not update Source, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sourceId := state.Id.Value
	err := r.client.DeleteSource(sourceId)
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
