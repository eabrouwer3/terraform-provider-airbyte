package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OperationResource{}
var _ resource.ResourceWithImportState = &OperationResource{}

func NewOperationResource() resource.Resource {
	return &OperationResource{}
}

// OperationResource defines the resource implementation.
type OperationResource struct {
	client *apiclient.ApiClient
}

func (r *OperationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_operation"
}

func (r *OperationResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Operation resource",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Operation ID",
				Type:        types.StringType,
				Computed:    true,
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
				Description: "Operation Name",
				Type:        types.StringType,
				Required:    true,
			},
			"operator_type": {
				Description: "Operation Name",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("normalization", "dbt", "webhook"),
					utils.ValueBasedAlsoRequires("normalization", path.MatchRelative().AtParent().AtName("normalization_option")),
					utils.ValueBasedAlsoRequires("dbt", path.MatchRelative().AtParent().AtName("dbt")),
					utils.ValueBasedAlsoRequires("webhook", path.MatchRelative().AtParent().AtName("webhook")),
				},
			},
			"normalization_option": {
				Description: "Normalization Option",
				Type:        types.StringType,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("basic"),
					schemavalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("normalization_option"),
						path.MatchRelative().AtParent().AtName("dbt"),
						path.MatchRelative().AtParent().AtName("webhook"),
					),
				},
			},
			"dbt": {
				Description: "DBT Configuration",
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("normalization_option"),
						path.MatchRelative().AtParent().AtName("dbt"),
						path.MatchRelative().AtParent().AtName("webhook"),
					),
				},
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"git_repo_url": {
						Description: "Git repo where DBT Transforms are",
						Type:        types.StringType,
						Required:    true,
					},
					"git_repo_branch": {
						Description: "Branch of above repo that should be used",
						Type:        types.StringType,
						Optional:    true,
					},
					"docker_image": {
						Description: "DBT Docker Image",
						Type:        types.StringType,
						Optional:    true,
					},
					"dbt_arguments": {
						Description: "Arguments to pass to DBT on a run",
						Type:        types.StringType,
						Optional:    true,
					},
				}),
			},
			"webhook": {
				Description: "Webhook Configuration",
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("normalization_option"),
						path.MatchRelative().AtParent().AtName("dbt"),
						path.MatchRelative().AtParent().AtName("webhook"),
					),
				},
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"execution_url": {
						Description: "The URL to call to execute the webhook operation via POST request.",
						Type:        types.StringType,
						Required:    true,
					},
					"execution_body": {
						Description: "If populated, this will be sent with the POST request.",
						Type:        types.StringType,
						Optional:    true,
					},
					"webhook_config_id": {
						Description: "The id of the webhook configs to use from the workspace.",
						Type:        types.StringType,
						Optional:    true,
					},
				}),
			},
		},
	}, nil
}

func (r *OperationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OperationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OperationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newOperation := apiclient.NewOperation{
		WorkspaceIdBody: apiclient.WorkspaceIdBody{
			WorkspaceId: plan.WorkspaceId.ValueString(),
		},
		CommonOperationFields: GetCommonOperationFields(plan),
	}

	checkResponse, err := r.client.CheckOperation(newOperation.OperatorConfiguration)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error checking operation",
			"Could not check operation, unexpected error: "+err.Error(),
		)
		return
	}
	if checkResponse.Status == "failed" {
		resp.Diagnostics.AddError(
			"Error creating operation",
			"Operation check failed, unexpected error: "+checkResponse.Message,
		)
		return
	}
	operation, err := r.client.CreateOperation(newOperation)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating operation",
			"Could not create operation, unexpected error: "+err.Error(),
		)
		return
	}

	state := FlattenOperation(operation)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OperationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OperationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	operationId := state.Id.ValueString()

	operation, err := r.client.GetOperationById(operationId)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read operation, got error: %s", err))
		return
	}

	state = FlattenOperation(operation)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OperationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OperationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updatedOperation := apiclient.UpdatedOperation{
		OperationIdBody: apiclient.OperationIdBody{
			OperationId: plan.Id.ValueString(),
		},
		CommonOperationFields: GetCommonOperationFields(plan),
	}

	operation, err := r.client.UpdateOperation(updatedOperation)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating operation",
			"Could not update operation, unexpected error: "+err.Error(),
		)
		return
	}

	state := FlattenOperation(operation)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OperationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OperationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	operationId := state.Id.ValueString()
	err := r.client.DeleteOperation(operationId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating operation",
			"Could not update operation, unexpected error: "+err.Error(),
		)
	}
}

func (r *OperationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
