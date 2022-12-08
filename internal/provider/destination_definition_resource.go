package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &DestinationDefinitionResource{}
var _ resource.ResourceWithImportState = &DestinationDefinitionResource{}

func NewDestinationDefinitionResource() resource.Resource {
	return &DestinationDefinitionResource{}
}

// DestinationDefinitionResource defines the resource implementation.
type DestinationDefinitionResource struct {
	client *apiclient.ApiClient
}

func (r *DestinationDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination_definition"
}

func (r *DestinationDefinitionResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DestinationDefinition resource",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Destination Definition ID",
				Type:        types.StringType,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.UseStateForUnknown(),
				},
			},
			"workspace_id": {
				Description: "Workspace ID",
				Type:        types.StringType,
				Required:    true,
			},
			"name": {
				Description: "Destination Definition Name",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"docker_repository": {
				Description: "Docker Repository URL (e.g. 112233445566.dkr.ecr.us-east-1.amazonaws.com/destination-custom) or DockerHub identifier (e.g. airbyte/destination-postgres)",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"docker_image_tag": {
				Description: "Docker image tag",
				Type:        types.StringType,
				Required:    true,
			},
			"documentation_url": {
				Description: "Documentation URL",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"protocol_version": {
				Description: "The Airbyte Protocol version supported by the connector",
				Type:        types.StringType,
				Computed:    true,
			},
			"release_stage": {
				Description: "Allowed: alpha | beta | generally_available | custom",
				Type:        types.StringType,
				Computed:    true,
			},
			"release_date": {
				Description: "The date when this connector was first released, in yyyy-mm-dd format",
				Type:        types.StringType,
				Computed:    true,
			},
			"default_resource_requirements": {
				Description: "Actor definition specific resource requirements. If default is set, these are the requirements " +
					"that should be set for ALL jobs run for this actor definition. It is overridden by the job type specific " +
					"configurations. If not set, the platform will use defaults. These values will be overridden by configuration " +
					"at the connection level.",
				Optional: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"cpu_request": {
						Description: "CPU Requested",
						Type:        types.StringType,
						Optional:    true,
					},
					"cpu_limit": {
						Description: "CPU Limit",
						Type:        types.StringType,
						Optional:    true,
					},
					"memory_request": {
						Description: "Memory Requested",
						Type:        types.StringType,
						Optional:    true,
					},
					"memory_limit": {
						Description: "Memory Limit",
						Type:        types.StringType,
						Optional:    true,
					},
				}),
			},
			"job_specific_resource_requirements": {
				Description: "Sets resource requirements for a specific job type for an actor definition. These values override " +
					"the default, if both are set. These values will be overridden by configuration at the connection level.",
				Optional: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"job_type": {
						Description: "Allowed: get_spec | check_connection | discover_schema | sync | reset_connection | connection_updater | replicate",
						Type:        types.StringType,
						Required:    true,
						Validators: []tfsdk.AttributeValidator{
							stringvalidator.OneOf("get_spec", "check_connection", "discover_schema", "sync", "reset_connection", "connection_updater", "replicate"),
						},
					},
					"cpu_request": {
						Description: "CPU Requested",
						Type:        types.StringType,
						Optional:    true,
					},
					"cpu_limit": {
						Description: "CPU Limit",
						Type:        types.StringType,
						Optional:    true,
					},
					"memory_request": {
						Description: "Memory Requested",
						Type:        types.StringType,
						Optional:    true,
					},
					"memory_limit": {
						Description: "Memory Limit",
						Type:        types.StringType,
						Optional:    true,
					},
				}),
			},
		},
	}, nil
}

func (r *DestinationDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DestinationDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConnectorDefinitionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	commonFields := GetCommonConnectorDefinitionFields(plan)
	newDestinationDefinition := apiclient.NewConnectorDefinition{
		WorkspaceIdBody: apiclient.WorkspaceIdBody{
			WorkspaceId: plan.WorkspaceId.ValueString(),
		},
		DestinationDefinition: &commonFields,
	}

	destinationDefinition, err := r.client.CreateConnectorDefinition(newDestinationDefinition, apiclient.DestinationType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Destination Definition",
			"Could not create Destination Definition, unexpected error: "+err.Error(),
		)
		return
	}

	state, err := FlattenConnectorDefinition(destinationDefinition)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Destination Definition",
			"Could not create Destination Definition, unexpected error: "+err.Error(),
		)
		return
	}
	state.WorkspaceId = plan.WorkspaceId

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DestinationDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan ConnectorDefinitionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	destinationDefinitionId := plan.Id.ValueString()

	destinationDefinition, err := r.client.GetConnectorDefinitionById(destinationDefinitionId, apiclient.DestinationType)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Destination Definition, got error: %s", err))
		return
	}

	state, err := FlattenConnectorDefinition(destinationDefinition)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Destination Definition, got error: %s", err))
		return
	}
	state.WorkspaceId = plan.WorkspaceId

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DestinationDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConnectorDefinitionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updatedDestinationDefinition := apiclient.UpdatedConnectorDefinition{
		DestinationDefinitionIdBody: apiclient.DestinationDefinitionIdBody{
			DestinationDefinitionId: plan.Id.ValueString(),
		},
		DockerImageTag:       plan.DockerImageTag.ValueString(),
		ResourceRequirements: getResourceRequirementFields(plan),
	}

	destinationDefinition, err := r.client.UpdateConnectorDefinition(updatedDestinationDefinition, apiclient.DestinationType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Destination Definition",
			"Could not update Destination Definition, unexpected error: "+err.Error(),
		)
		return
	}

	state, err := FlattenConnectorDefinition(destinationDefinition)
	state.WorkspaceId = plan.WorkspaceId

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DestinationDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ConnectorDefinitionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	destinationDefinitionId := state.Id.ValueString()
	err := r.client.DeleteConnectorDefinition(destinationDefinitionId, apiclient.DestinationType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Destination Definition",
			"Could not update Destination Definition, unexpected error: "+err.Error(),
		)
	}
}

func (r *DestinationDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
