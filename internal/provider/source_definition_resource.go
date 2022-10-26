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
var _ resource.Resource = &SourceDefinitionResource{}
var _ resource.ResourceWithImportState = &SourceDefinitionResource{}

func NewSourceDefinitionResource() resource.Resource {
	return &SourceDefinitionResource{}
}

// SourceDefinitionResource defines the resource implementation.
type SourceDefinitionResource struct {
	client *apiclient.ApiClient
}

func (r *SourceDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_definition"
}

func (r *SourceDefinitionResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "SourceDefinition resource",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Source Definition ID",
				Type:        types.StringType,
				Computed:    true,
			},
			"name": {
				Description: "Source Definition Name",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"docker_repository": {
				Description: "Docker Repository URL (e.g. 112233445566.dkr.ecr.us-east-1.amazonaws.com/source-custom) or DockerHub identifier (e.g. airbyte/source-postgres)",
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

func (r *SourceDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SourceDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConnectorDefinitionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newSourceDefinition := GetCommonConnectorDefinitionFields(plan)

	sourceDefinition, err := r.client.CreateConnectorDefinition(newSourceDefinition, apiclient.SourceType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Source Definition",
			"Could not create Source Definition, unexpected error: "+err.Error(),
		)
		return
	}

	state, err := FlattenConnectorDefinition(sourceDefinition)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Source Definition",
			"Could not create Source Definition, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SourceDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan ConnectorDefinitionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sourceDefinitionId := plan.Id.ValueString()

	sourceDefinition, err := r.client.GetConnectorDefinitionById(sourceDefinitionId, apiclient.SourceType)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Source Definition, got error: %s", err))
		return
	}

	state, err := FlattenConnectorDefinition(sourceDefinition)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Source Definition, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SourceDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConnectorDefinitionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updatedSourceDefinition := apiclient.UpdatedConnectorDefinition{
		SourceDefinitionIdBody: apiclient.SourceDefinitionIdBody{
			SourceDefinitionId: plan.Id.ValueString(),
		},
		DockerImageTag:       plan.DockerImageTag.ValueString(),
		ResourceRequirements: getResourceRequirementFields(plan),
	}

	sourceDefinition, err := r.client.UpdateConnectorDefinition(updatedSourceDefinition, apiclient.SourceType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Source Definition",
			"Could not update Source Definition, unexpected error: "+err.Error(),
		)
		return
	}

	state, err := FlattenConnectorDefinition(sourceDefinition)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SourceDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ConnectorDefinitionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sourceDefinitionId := state.Id.ValueString()
	err := r.client.DeleteConnectorDefinition(sourceDefinitionId, apiclient.SourceType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Source Definition",
			"Could not update Source Definition, unexpected error: "+err.Error(),
		)
	}
}

func (r *SourceDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
