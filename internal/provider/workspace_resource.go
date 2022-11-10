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
var _ resource.Resource = &WorkspaceResource{}
var _ resource.ResourceWithImportState = &WorkspaceResource{}

func NewWorkspaceResource() resource.Resource {
	return &WorkspaceResource{}
}

// WorkspaceResource defines the resource implementation.
type WorkspaceResource struct {
	client *apiclient.ApiClient
}

func (r *WorkspaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (r *WorkspaceResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Workspace resource",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Workspace ID",
				Type:        types.StringType,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.UseStateForUnknown(),
				},
			},
			"slug": {
				Description: "Workspace Slug",
				Type:        types.StringType,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.UseStateForUnknown(),
				},
			},
			"customer_id": {
				Description: "Customer ID",
				Type:        types.StringType,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.UseStateForUnknown(),
				},
			},
			"email": {
				Description: "Customer Email",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Description: "Workspace Name",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"initial_setup_complete": {
				Description: "Is the initial setup complete",
				Type:        types.BoolType,
				Computed:    true,
			},
			"display_setup_wizard": {
				Description: "Should the UI display the setup wizard",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"anonymous_data_collection": {
				Description: "Is anonymous data collection turned on",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"news": {
				Description: "Should the UI show news updates",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"security_updates": {
				Description: "Should the UI show security updates",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"notification_config": {
				Description: "Notification systems set up",
				Optional:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"notification_type": {
						Description: "Possible value: slack",
						Type:        types.StringType,
						Required:    true,
						Validators: []tfsdk.AttributeValidator{
							stringvalidator.OneOf("slack"),
						},
					},
					"send_on_success": {
						Description: "Should the notification be sent for successes",
						Type:        types.BoolType,
						Optional:    true,
						Computed:    true,
					},
					"send_on_failure": {
						Description: "Should the notification be sent for failures",
						Type:        types.BoolType,
						Optional:    true,
						Computed:    true,
					},
					"slack_webhook": {
						Description: "Configuration for Slack notifications - See https://slack.com/help/articles/115005265063-Incoming-webhooks-for-Slack",
						Type:        types.StringType,
						Required:    true,
					},
				}),
			},
			"first_completed_sync": {
				Description: "Has a first sync completed",
				Type:        types.BoolType,
				Computed:    true,
			},
			"feedback_done": {
				Description: "Is Feedback done",
				Type:        types.BoolType,
				Computed:    true,
			},
			"default_geography": {
				Description: "Possible values: auto | us | eu",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

func getCommonWorkspaceFields(data WorkspaceModel) apiclient.CommonWorkspaceFields {
	fields := apiclient.CommonWorkspaceFields{
		Email: data.Email.ValueString(),
	}

	if v := data.AnonymousDataCollection; !v.IsUnknown() {
		b := v.ValueBool()
		fields.AnonymousDataCollection = &b
	}
	if v := data.News; !v.IsUnknown() {
		b := v.ValueBool()
		fields.News = &b
	}
	if v := data.SecurityUpdates; !v.IsUnknown() {
		b := v.ValueBool()
		fields.SecurityUpdates = &b
	}
	if v := data.DisplaySetupWizard; !v.IsUnknown() {
		b := v.ValueBool()
		fields.DisplaySetupWizard = &b
	}
	for _, notif := range data.NotificationConfig {
		n := apiclient.Notification{
			NotificationType: notif.NotificationType.ValueString(),
			SlackConfiguration: apiclient.SlackConfiguration{
				Webhook: notif.SlackWebhook.ValueString(),
			},
		}
		if v := notif.SendOnSuccess; !v.IsUnknown() {
			b := v.ValueBool()
			n.SendOnSuccess = &b
		}
		if v := notif.SendOnFailure; !v.IsUnknown() {
			b := v.ValueBool()
			n.SendOnFailure = &b
		}
		fields.Notifications = append(fields.Notifications, n)
	}

	return fields
}

func (r *WorkspaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WorkspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WorkspaceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newWorkspace := apiclient.NewWorkspace{
		WorkspaceNameBody: apiclient.WorkspaceNameBody{
			Name: plan.Name.ValueString(),
		},
		CommonWorkspaceFields: getCommonWorkspaceFields(plan),
	}

	workspace, err := r.client.CreateWorkspace(newWorkspace)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating workspace",
			"Could not create workspace, unexpected error: "+err.Error(),
		)
		return
	}

	state := FlattenWorkspace(workspace)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkspaceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId := state.Id.ValueString()

	workspace, err := r.client.GetWorkspaceById(workspaceId)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace, got error: %s", err))
		return
	}

	state = FlattenWorkspace(workspace)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WorkspaceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updatedWorkspace := apiclient.UpdatedWorkspace{
		WorkspaceIdBody: apiclient.WorkspaceIdBody{
			WorkspaceId: plan.Id.ValueString(),
		},
		CommonWorkspaceFields: getCommonWorkspaceFields(plan),
	}

	workspace, err := r.client.UpdateWorkspace(updatedWorkspace)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating workspace",
			"Could not update workspace, unexpected error: "+err.Error(),
		)
		return
	}

	state := FlattenWorkspace(workspace)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkspaceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId := state.Id.ValueString()
	err := r.client.DeleteWorkspace(workspaceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating workspace",
			"Could not update workspace, unexpected error: "+err.Error(),
		)
	}
}

func (r *WorkspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
