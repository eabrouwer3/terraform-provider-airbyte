package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/utils"
	"strconv"
	"time"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ datasource.DataSource              = &WorkspaceIdsDataSource{}
	_ datasource.DataSourceWithConfigure = &WorkspaceIdsDataSource{}
)

func NewWorkspaceIdsDataSource() datasource.DataSource {
	return &WorkspaceIdsDataSource{}
}

// WorkspaceIdsDataSource defines the data source implementation.
type WorkspaceIdsDataSource struct {
	client *apiclient.ApiClient
}

type WorkspaceIdsModel struct {
	Id  types.String `tfsdk:"id"`
	Ids types.List   `tfsdk:"ids"`
}

func (d *WorkspaceIdsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_ids"
}

func (d *WorkspaceIdsDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Get all Airbyte Workspace ids (first will always be the default one created by Airbyte on launch)",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"ids": {
				Description: "Workspace Id List",
				Type:        types.ListType{ElemType: types.StringType},
				Computed:    true,
			},
		},
	}, nil
}

func (d *WorkspaceIdsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspaceIdsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config WorkspaceIdsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	wl, err := d.client.GetWorkspaces()

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace, got error: %s", err))
		return
	}

	ids, diags := types.ListValue(
		types.StringType,
		utils.Map(wl, func(w *apiclient.Workspace) attr.Value {
			return types.StringValue(w.WorkspaceId)
		}),
	)
	if diags.HasError() {
		return
	}

	state := WorkspaceIdsModel{
		Id:  types.StringValue(strconv.FormatInt(time.Now().Unix(), 10)),
		Ids: ids,
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
