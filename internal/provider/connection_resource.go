package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
var _ resource.Resource = &ConnectionResource{}
var _ resource.ResourceWithImportState = &ConnectionResource{}

func NewConnectionResource() resource.Resource {
	return &ConnectionResource{}
}

// ConnectionResource defines the resource implementation.
type ConnectionResource struct {
	client *apiclient.ApiClient
}

func (r *ConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection"
}

func (r *ConnectionResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Connection resource",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Connection ID",
				Type:        types.StringType,
				Computed:    true,
			},
			"source_id": {
				Description: "Source ID",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"destination_id": {
				Description: "Destination ID",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"status": {
				Description: "Active means that data is flowing through the connection. Inactive means it is not." +
					"Deprecated means the connection is off and cannot be re-activated. The schema field describes " +
					"the elements of the schema that will be synced. Allowed Values: 'active' | 'inactive' | 'deprecated'.",
				Type:     types.StringType,
				Required: true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("active", "inactive", "deprecated"),
				},
			},
			"name": {
				Description: "Optional name of the connection",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"namespace_definition": {
				Description: "Method used for computing final namespace in destination. " +
					"Allowed Values: 'source' | 'destination' | 'customformat'",
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("source", "destination", "customformat"),
				},
			},
			"namespace_format": {
				Description: "Used when namespaceDefinition is 'customformat'. If blank then behaves like " +
					"namespaceDefinition = 'destination'. If \"${SOURCE_NAMESPACE}\" then behaves like " +
					"namespaceDefinition = 'source'.",
				Type:     types.StringType,
				Optional: true,
			},
			"prefix": {
				Description: "Prefix that will be prepended to the name of each stream when it is written to the destination. Example: \"airbyte_\"",
				Type:        types.StringType,
				Optional:    true,
			},
			"operation_ids": {
				Description: "Operation IDs",
				Type:        types.ListType{ElemType: types.StringType},
				Optional:    true,
				Computed:    true,
			},
			"sync_catalog": {
				Description: "Describes the available schema (catalog). Each stream is split in two parts; the " +
					"immutable schema from source and mutable configuration for destination.",
				Required: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"source_schema": {
						Description: "The immutable schema defined by the source",
						Required:    true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"name": {
								Description: "Stream's name",
								Type:        types.StringType,
								Required:    true,
							},
							"json_schema": {
								Description: "Stream schema using json Schema specs",
								Type:        types.StringType,
								Optional:    true,
							},
							"supported_sync_modes": {
								Description: "Allowed Values: 'full_refresh' | 'incremental'",
								Type:        types.ListType{ElemType: types.StringType},
								Optional:    true,
								Validators: []tfsdk.AttributeValidator{
									listvalidator.ValuesAre(stringvalidator.OneOf("source", "destination", "customformat")),
								},
							},
							"source_defined_cursor": {
								Description: "If the source defines the cursor field, then any other cursor field " +
									"inputs will be ignored. If it does not, either the user_provided one is used, " +
									"or the default one is used as a backup.",
								Type:     types.BoolType,
								Optional: true,
							},
							"default_cursor_field": {
								Description: "Path to the field that will be used to determine if a record is new or " +
									"modified since the last sync. If not provided by the source, the end user will " +
									"have to specify the comparable themselves.",
								Type:     types.ListType{ElemType: types.StringType},
								Optional: true,
							},
							"source_defined_primary_key": {
								Description: "If the source defines the primary key, paths to the fields that will be " +
									"used as a primary key. If not provided by the source, the end user will have to " +
									"specify the primary key themselves.",
								Type:     types.ListType{ElemType: types.ListType{ElemType: types.StringType}},
								Optional: true,
							},
							"namespace": {
								Description: "Optional Source-defined namespace. Airbyte streams from the same sources " +
									"should have the same namespace. Currently only used by JDBC destinations to " +
									"determine what schema to write to.",
								Type:     types.StringType,
								Optional: true,
							},
						}),
					},
					"destination_config": {
						Description: "The mutable part of the stream to configure the destination",
						Required:    true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"sync_mode": {
								Description: "Allowed Values: 'full_refresh' | 'incremental'",
								Type:        types.StringType,
								Required:    true,
								Validators: []tfsdk.AttributeValidator{
									stringvalidator.OneOf("full_refresh", "incremental"),
									utils.ValueBasedAlsoRequires("incremental", path.MatchRelative().AtName("cursor_field")),
								},
							},
							"cursor_field": {
								MarkdownDescription: "Path to the field that will be used to determine if a record is " +
									"new or modified since the last sync. This field is REQUIRED if `sync_mode` is " +
									"`incremental`. Otherwise it is ignored.",
								Type:     types.ListType{ElemType: types.StringType},
								Optional: true,
							},
							"destination_sync_mode": {
								Description: "Allowed Values: 'append' | 'overwrite' | 'append_dedup'",
								Type:        types.StringType,
								Required:    true,
								Validators: []tfsdk.AttributeValidator{
									stringvalidator.OneOf("append", "overwrite", "append_dedup"),
									utils.ValueBasedAlsoRequires("append_dedup", path.MatchRelative().AtName("primary_key")),
								},
							},
							"primary_key": {
								MarkdownDescription: "Paths to the fields that will be used as primary key. This field " +
									"is REQUIRED if `destination_sync_mode` is `*_dedup`. Otherwise it is ignored.",
								Type:     types.ListType{ElemType: types.ListType{ElemType: types.StringType}},
								Optional: true,
								Validators: []tfsdk.AttributeValidator{
									schemavalidator.AlsoRequires(),
								},
							},
							"alias_name": {
								Description: "Alias name to the stream to be used in the destination",
								Type:        types.StringType,
								Optional:    true,
							},
							"selected": {
								Description: "Whether this config is selected i.e. should be synced",
								Type:        types.BoolType,
								Required:    true,
							},
						}),
					},
				}),
			},
			"schedule_type": {
				Description: "Determine how the schedule data should be interpreted. Allowed: 'manual' | 'basic' | 'cron'",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("manual", "basic", "cron"),
					utils.ValueBasedAlsoRequires("basic", path.MatchRelative().AtName("basic_schedule")),
					utils.ValueBasedAlsoRequires("cron", path.MatchRelative().AtName("cron_schedule")),
				},
			},
			"basic_schedule": {
				Description: "Basic time schedule - \"Run sync every ...\"",
				Optional:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"time_unit": {
						Description: "Allowed: minutes | hours | days | weeks | months",
						Type:        types.StringType,
						Required:    true,
						Validators: []tfsdk.AttributeValidator{
							stringvalidator.OneOf("minutes", "hours", "days", "weeks", "months"),
						},
					},
					"units": {
						MarkdownDescription: "Count of `time_unit`",
						Type:                types.Int64Type,
						Required:            true,
					},
				}),
			},
			"cron_schedule": {
				Description: "Flexible Cron Schedule",
				Optional:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"cron_expression": {
						MarkdownDescription: "[Cron Expression](http://www.quartz-scheduler.org/documentation/quartz-2.3.0/tutorials/crontrigger.html). " +
							"Example: `0 0 12 * * ?`.",
						Type:     types.StringType,
						Required: true,
					},
					"cron_time_zone": {
						MarkdownDescription: "Time Zone to honor cron expression according to. Examples: `UTC`, `US/Denver`, etc." +
							"See the 'TZ database name' column [here](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) for all options.",
						Type:     types.Int64Type,
						Required: true,
					},
				}),
			},
			"resource_requirements": {
				Description: "Optional resource requirements to run workers (blank for unbounded allocations)",
				Optional:    true,
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
			"source_catalog_id": {
				Description: "Source Catalog ID",
				Type:        types.StringType,
				Optional:    true,
			},
			"geography": {
				Description: "Allowed Values: 'auto' | 'us' | 'eu'",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

func getCommonConnectionFields(data ConnectionModel) apiclient.CommonConnectionFields {
	fields := apiclient.CommonConnectionFields{
		Status: data.Status.ValueString(),
	}

	if v := data.Name; !v.IsUnknown() {
		fields.Name = v.ValueString()
	}
	if v := data.NamespaceDefinition; !v.IsUnknown() {
		fields.NamespaceDefinition = v.ValueString()
	}
	if v := data.NamespaceFormat; !v.IsUnknown() {
		fields.NamespaceFormat = v.ValueString()
	}
	if v := data.Prefix; !v.IsUnknown() {
		fields.Prefix = v.ValueString()
	}
	if v := data.OperationIds; !v.IsUnknown() {
		for _, elem := range v.Elements() {
			fields.OperationIds = append(fields.OperationIds, elem.(types.String).ValueString())
		}
	}
	if data.SyncCatalog != nil {
		var streams []apiclient.Stream
		for _, cfg := range *data.SyncCatalog {
			stream := apiclient.Stream{
				Stream: apiclient.SourceStreamSchema{
					Name: cfg.SourceSchema.Name.ValueString(),
				},
				Config: apiclient.DestinationStreamConfig{
					SyncMode: cfg.DestinationConfig.SyncMode.ValueString(),
				},
			}

			// Source Schema Fields
			if v := cfg.SourceSchema.JsonSchema; !v.IsUnknown() {
				stream.Stream.JsonSchema = json.RawMessage(v.ValueString())
			}
			if v := cfg.SourceSchema.SupportedSyncModes; !v.IsUnknown() {
				for _, elem := range v.Elements() {
					stream.Stream.SupportedSyncModes = append(stream.Stream.SupportedSyncModes, elem.(types.String).ValueString())
				}
			}
			if v := cfg.SourceSchema.SourceDefinedCursor; !v.IsUnknown() {
				b := v.ValueBool()
				stream.Stream.SourceDefinedCursor = &b
			}
			if v := cfg.SourceSchema.DefaultCursorField; !v.IsUnknown() {
				for _, elem := range v.Elements() {
					stream.Stream.DefaultCursorField = append(stream.Stream.DefaultCursorField, elem.(types.String).ValueString())
				}
			}
			if v := cfg.SourceSchema.SourceDefinedPrimaryKey; !v.IsUnknown() {
				for _, elem := range v.Elements() {
					var arr []string
					for _, inner := range elem.(types.List).Elements() {
						arr = append(arr, inner.(types.String).ValueString())
					}
					stream.Stream.SourceDefinedPrimaryKey = append(stream.Stream.SourceDefinedPrimaryKey, arr)
				}
			}
			if v := cfg.SourceSchema.Namespace; !v.IsUnknown() {
				stream.Stream.Namespace = v.ValueString()
			}

			// Destination Config Fields
			if v := cfg.DestinationConfig.DestinationSyncMode; !v.IsUnknown() {
				stream.Config.DestinationSyncMode = v.ValueString()
			}
			if v := cfg.DestinationConfig.CursorField; !v.IsUnknown() {
				for _, elem := range v.Elements() {
					stream.Config.CursorField = append(stream.Config.CursorField, elem.(types.String).ValueString())
				}
			}
			if v := cfg.DestinationConfig.PrimaryKey; !v.IsUnknown() {
				for _, elem := range v.Elements() {
					var arr []string
					for _, inner := range elem.(types.List).Elements() {
						arr = append(arr, inner.(types.String).ValueString())
					}
					stream.Config.PrimaryKey = append(stream.Config.PrimaryKey, arr)
				}
			}
			if v := cfg.DestinationConfig.AliasName; !v.IsUnknown() {
				stream.Config.AliasName = v.ValueString()
			}
			if v := cfg.DestinationConfig.Selected; !v.IsUnknown() {
				b := v.ValueBool()
				stream.Config.Selected = &b
			}
		}
		fields.SyncCatalog.Streams = streams
	}
	if v := data.ScheduleType; !v.IsUnknown() {
		fields.ScheduleType = v.ValueString()
	}
	if data.BasicSchedule != nil || data.CronSchedule != nil {
		fields.ScheduleData = &apiclient.ScheduleData{}
		if data.BasicSchedule != nil {
			fields.ScheduleData.BasicSchedule = &apiclient.ScheduleSpec{
				Units:    data.BasicSchedule.Units.ValueInt64(),
				TimeUnit: data.BasicSchedule.TimeUnit.ValueString(),
			}
		}
		if data.CronSchedule != nil {
			fields.ScheduleData.Cron = &apiclient.CronScheduleSpec{
				CronExpression: data.CronSchedule.CronExpression.ValueString(),
				CronTimeZone:   data.CronSchedule.CronTimeZone.ValueString(),
			}
		}
	}
	if reqs := data.ResourceRequirements; reqs != nil {
		fields.ResourceRequirements = getResourceRequirementOptions(reqs)
	}
	if v := data.SourceCatalogId; !v.IsUnknown() {
		fields.SourceCatalogId = v.ValueString()
	}
	if v := data.BreakingChange; !v.IsUnknown() {
		b := v.ValueBool()
		fields.BreakingChange = &b
	}

	return fields
}

func (r *ConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConnectionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newConnection := apiclient.NewConnection{
		CommonConnectionFields: getCommonConnectionFields(plan),
		SourceIdBody: apiclient.SourceIdBody{
			SourceId: plan.SourceId.ValueString(),
		},
		DestinationIdBody: apiclient.DestinationIdBody{
			DestinationId: plan.DestinationId.ValueString(),
		},
	}

	connection, err := r.client.CreateConnection(newConnection)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating connection",
			"Could not create connection, unexpected error: "+err.Error(),
		)
		return
	}

	state, diags := FlattenConnection(connection)
	resp.Diagnostics.Append(diags...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ConnectionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	connectionId := state.Id.ValueString()

	connection, err := r.client.GetConnectionById(connectionId)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read connection, got error: %s", err))
		return
	}

	state, diags := FlattenConnection(connection)
	resp.Diagnostics.Append(diags...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConnectionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updatedConnection := apiclient.UpdatedConnection{
		ConnectionIdBody: apiclient.ConnectionIdBody{
			ConnectionId: plan.Id.ValueString(),
		},
		CommonConnectionFields: getCommonConnectionFields(plan),
	}

	connection, err := r.client.UpdateConnection(updatedConnection)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating connection",
			"Could not update connection, unexpected error: "+err.Error(),
		)
		return
	}

	state, diags := FlattenConnection(connection)
	resp.Diagnostics.Append(diags...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ConnectionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	connectionId := state.Id.ValueString()
	err := r.client.DeleteConnection(connectionId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating connection",
			"Could not update connection, unexpected error: "+err.Error(),
		)
	}
}

func (r *ConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
