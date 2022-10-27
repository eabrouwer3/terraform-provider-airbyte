package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
	"strconv"
	"time"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ datasource.DataSource              = &SourceSchemaCatalogDataSource{}
	_ datasource.DataSourceWithConfigure = &SourceSchemaCatalogDataSource{}
)

func NewSourceSchemaCatalogDataSource() datasource.DataSource {
	return &SourceSchemaCatalogDataSource{}
}

// SourceSchemaCatalogDataSource defines the data source implementation.
type SourceSchemaCatalogDataSource struct {
	client *apiclient.ApiClient
}

func (d *SourceSchemaCatalogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_schema_catalog"
}

func (d *SourceSchemaCatalogDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Get an Airbyte Source Schema Catalog by Source id",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Unique ID - Force to always get a new version",
				Type:        types.StringType,
				Computed:    true,
			},
			"source_id": {
				Description: "Source to get the schema for",
				Type:        types.StringType,
				Required:    true,
			},
			"sync_catalog": {
				Description: "Describes the available schema (catalog). Each stream is split in two parts; the " +
					"immutable schema from source and mutable configuration for destination.",
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"source_schema": {
						Description: "The immutable schema defined by the source",
						Computed:    true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"name": {
								Description: "Stream's name",
								Type:        types.StringType,
								Computed:    true,
							},
							"json_schema": {
								Description: "Stream schema using json Schema specs",
								Type:        types.StringType,
								Computed:    true,
							},
							"supported_sync_modes": {
								Description: "Allowed Values: 'full_refresh' | 'incremental'",
								Type:        types.ListType{ElemType: types.StringType},
								Computed:    true,
							},
							"source_defined_cursor": {
								Description: "If the source defines the cursor field, then any other cursor field " +
									"inputs will be ignored. If it does not, either the user_provided one is used, " +
									"or the default one is used as a backup.",
								Type:     types.BoolType,
								Computed: true,
							},
							"default_cursor_field": {
								Description: "Path to the field that will be used to determine if a record is new or " +
									"modified since the last sync. If not provided by the source, the end user will " +
									"have to specify the comparable themselves.",
								Type:     types.ListType{ElemType: types.StringType},
								Computed: true,
							},
							"source_defined_primary_key": {
								Description: "If the source defines the primary key, paths to the fields that will be " +
									"used as a primary key. If not provided by the source, the end user will have to " +
									"specify the primary key themselves.",
								Type:     types.ListType{ElemType: types.ListType{ElemType: types.StringType}},
								Computed: true,
							},
							"namespace": {
								Description: "Optional Source-defined namespace. Airbyte streams from the same sources " +
									"should have the same namespace. Currently only used by JDBC destinations to " +
									"determine what schema to write to.",
								Type:     types.StringType,
								Computed: true,
							},
						}),
					},
					"destination_config": {
						Description: "The mutable part of the stream to configure the destination",
						Computed:    true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"sync_mode": {
								Description: "Allowed Values: 'full_refresh' | 'incremental'",
								Type:        types.StringType,
								Computed:    true,
							},
							"cursor_field": {
								MarkdownDescription: "Path to the field that will be used to determine if a record is " +
									"new or modified since the last sync. This field is REQUIRED if `sync_mode` is " +
									"`incremental`. Otherwise it is ignored.",
								Type:     types.ListType{ElemType: types.StringType},
								Computed: true,
							},
							"destination_sync_mode": {
								Description: "Allowed Values: 'append' | 'overwrite' | 'append_dedup'",
								Type:        types.StringType,
								Computed:    true,
							},
							"primary_key": {
								MarkdownDescription: "Paths to the fields that will be used as primary key. This field " +
									"is REQUIRED if `destination_sync_mode` is `*_dedup`. Otherwise it is ignored.",
								Type:     types.ListType{ElemType: types.ListType{ElemType: types.StringType}},
								Computed: true,
							},
							"alias_name": {
								Description: "Alias name to the stream to be used in the destination",
								Type:        types.StringType,
								Computed:    true,
							},
							"selected": {
								Description: "Whether this config is selected i.e. should be synced",
								Type:        types.BoolType,
								Computed:    true,
							},
						}),
					},
				}),
			},
		},
	}, nil
}

func (d *SourceSchemaCatalogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SourceSchemaCatalogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config SourceSchemaCatalogModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sourceId := config.SourceId.ValueString()

	sourceSchemaCatalog, err := d.client.GetSourceSchemaCatalogById(sourceId)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Source Schema Catalog, got error: %s", err))
		return
	}

	syncCatalog, diags := FlattenSyncCatalog(&sourceSchemaCatalog.Catalog)
	resp.Diagnostics.Append(diags...)
	state := SourceSchemaCatalogModel{
		Id:          types.StringValue(strconv.FormatInt(time.Now().Unix(), 10)),
		SourceId:    types.StringValue(sourceId),
		SyncCatalog: syncCatalog,
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
