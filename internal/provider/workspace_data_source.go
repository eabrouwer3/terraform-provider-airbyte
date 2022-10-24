package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ datasource.DataSource              = &WorkspaceDataSource{}
	_ datasource.DataSourceWithConfigure = &WorkspaceDataSource{}
)

func NewWorkspaceDataSource() datasource.DataSource {
	return &WorkspaceDataSource{}
}

// WorkspaceDataSource defines the data source implementation.
type WorkspaceDataSource struct {
	client *apiclient.ApiClient
}

func (d *WorkspaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (d *WorkspaceDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Get an Airbyte Workspace by id or slug",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Workspace ID",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("slug"),
					}...),
				},
			},
			"slug": {
				Description: "Workspace Slug",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("id"),
					}...),
				},
			},
			"customer_id": {
				Description: "Customer ID",
				Type:        types.StringType,
				Computed:    true,
			},
			"email": {
				Description: "Customer Email",
				Type:        types.StringType,
				Computed:    true,
			},
			"name": {
				Description: "Workspace Name",
				Type:        types.StringType,
				Computed:    true,
			},
			"initial_setup_complete": {
				Description: "Is the initial setup complete",
				Type:        types.BoolType,
				Computed:    true,
			},
			"display_setup_wizard": {
				Description: "Should the UI display the setup wizard",
				Type:        types.BoolType,
				Computed:    true,
			},
			"anonymous_data_collection": {
				Description: "Is anonymous data collection turned on",
				Type:        types.BoolType,
				Computed:    true,
			},
			"news": {
				Description: "Should the UI show news updates",
				Type:        types.BoolType,
				Computed:    true,
			},
			"security_updates": {
				Description: "Should the UI show security updates",
				Type:        types.BoolType,
				Computed:    true,
			},
			"notification_config": {
				Description: "Notification systems set up",
				Computed:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"notification_type": {
						Description: "Possible value: slack",
						Type:        types.StringType,
						Computed:    true,
					},
					"send_on_success": {
						Description: "Should the notification be sent for successes",
						Type:        types.BoolType,
						Computed:    true,
					},
					"send_on_failure": {
						Description: "Should the notification be sent for failures",
						Type:        types.BoolType,
						Computed:    true,
					},
					"slack_webhook": {
						Description: "Configuration for Slack notifications - See https://slack.com/help/articles/115005265063-Incoming-webhooks-for-Slack",
						Type:        types.StringType,
						Computed:    true,
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

func (d *WorkspaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = &client
}

func (d *WorkspaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config WorkspaceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId := config.Id.Value
	slug := config.Slug.Value

	if workspaceId != "" && slug != "" {
		resp.Diagnostics.AddError(
			"Error getting workspace",
			"Only one of `id` and `slug` can be set",
		)
		return
	}

	var workspace *apiclient.Workspace
	var err error
	if workspaceId != "" {
		workspace, err = d.client.GetWorkspaceById(workspaceId)
	} else if slug != "" {
		workspace, err = d.client.GetWorkspaceBySlug(slug)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace, got error: %s", err))
		return
	}

	state := FlattenWorkspace(workspace)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
