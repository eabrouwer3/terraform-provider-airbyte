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
var _ resource.Resource = &DestinationResource{}
var _ resource.ResourceWithImportState = &DestinationResource{}

func NewDestinationResource() resource.Resource {
	return &DestinationResource{}
}

// DestinationResource defines the resource implementation.
type DestinationResource struct {
	client *apiclient.ApiClient
}

func (r *DestinationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination"
}

func (r *DestinationResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Destination resource",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Destination ID",
				Type:        types.StringType,
				Computed:    true,
			},
			"definition_id": {
				Description: "Destination Definition ID",
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
				Description: "Destination Name",
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
				Description: "Destination Definition Name",
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

func (r *DestinationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(apiclient.ApiClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Destination Configure Type",
			fmt.Sprintf("Expected *apiclient.ApiClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = &client
}

func (r *DestinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConnectorModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newDestination := apiclient.NewConnector{
		DestinationDefinitionIdBody: apiclient.DestinationDefinitionIdBody{
			DestinationDefinitionId: plan.DefinitionId.ValueString(),
		},
		WorkspaceIdBody: apiclient.WorkspaceIdBody{
			WorkspaceId: plan.WorkspaceId.ValueString(),
		},
		CommonConnectorFields: GetCommonConnectorFields(plan),
	}

	destination, err := r.client.CreateConnector(newDestination, apiclient.DestinationType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Destination",
			"Could not create Destination, unexpected error: "+err.Error(),
		)
		return
	}

	state, err := FlattenConnector(destination)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Destination",
			"Could not create Destination, unexpected error: "+err.Error(),
		)
		return
	}
	state.ConnectionConfiguration = plan.ConnectionConfiguration

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DestinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan ConnectorModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	destinationId := plan.Id.ValueString()

	destination, err := r.client.GetConnectorById(destinationId, apiclient.DestinationType)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Destination, got error: %s", err))
		return
	}

	state, err := FlattenConnector(destination)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read Destination, got error: %s", err),
		)
		return
	}
	state.ConnectionConfiguration = plan.ConnectionConfiguration

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DestinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConnectorModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updatedDestination := apiclient.UpdatedConnector{
		DestinationIdBody:     apiclient.DestinationIdBody{DestinationId: plan.Id.ValueString()},
		CommonConnectorFields: GetCommonConnectorFields(plan),
	}

	destination, err := r.client.UpdateConnector(updatedDestination, apiclient.DestinationType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating destination",
			"Could not update Destination, unexpected error: "+err.Error(),
		)
		return
	}

	state, err := FlattenConnector(destination)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating destination",
			"Could not update Destination, unexpected error: "+err.Error(),
		)
		return
	}
	state.ConnectionConfiguration = plan.ConnectionConfiguration

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DestinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ConnectorModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	destinationId := state.Id.ValueString()
	err := r.client.DeleteConnector(destinationId, apiclient.DestinationType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating destination",
			"Could not update Destination, unexpected error: "+err.Error(),
		)
	}
}

func (r *DestinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
