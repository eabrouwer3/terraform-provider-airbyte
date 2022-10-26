package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
)

// ConnectionModel describes the data source data model.
type ConnectionModel struct {
	Id                   types.String               `tfsdk:"id"`
	SourceId             types.String               `tfsdk:"source_id"`
	DestinationId        types.String               `tfsdk:"destination_id"`
	Status               types.String               `tfsdk:"status"`
	Name                 types.String               `tfsdk:"name"`
	NamespaceDefinition  types.String               `tfsdk:"namespace_definition"`
	NamespaceFormat      types.String               `tfsdk:"namespace_format"`
	Prefix               types.String               `tfsdk:"prefix"`
	OperationIds         types.List                 `tfsdk:"operation_ids"`
	SyncCatalog          *[]syncCatalogModel        `tfsdk:"sync_catalog"`
	ScheduleType         types.String               `tfsdk:"schedule_type"`
	BasicSchedule        *basicScheduleModule       `tfsdk:"basic_schedule"`
	CronSchedule         *cronScheduleModel         `tfsdk:"cron_schedule"`
	ResourceRequirements *ResourceRequirementsModel `tfsdk:"resource_requirements"`
	SourceCatalogId      types.String               `tfsdk:"source_catalog_id"`
	Geography            types.String               `tfsdk:"geography"`
	BreakingChange       types.Bool                 `tfsdk:"breaking_change"`
}

type syncCatalogModel struct {
	SourceSchema      sourceStreamSchemaModel      `tfsdk:"source_schema"`
	DestinationConfig destinationStreamConfigModel `tfsdk:"destination_config"`
}

type sourceStreamSchemaModel struct {
	Name                    types.String `tfsdk:"name"`
	JsonSchema              types.String `tfsdk:"json_schema"`
	SupportedSyncModes      types.List   `tfsdk:"supported_sync_modes"`
	SourceDefinedCursor     types.Bool   `tfsdk:"source_defined_cursor"`
	DefaultCursorField      types.List   `tfsdk:"default_cursor_field"`
	SourceDefinedPrimaryKey types.List   `tfsdk:"source_defined_primary_key"`
	Namespace               types.String `tfsdk:"namespace"`
}

type destinationStreamConfigModel struct {
	SyncMode            types.String `tfsdk:"sync_mode"`
	DestinationSyncMode types.String `tfsdk:"destination_sync_mode"`
	CursorField         types.List   `tfsdk:"cursor_field"`
	PrimaryKey          types.List   `tfsdk:"primary_key"`
	AliasName           types.String `tfsdk:"alias_name"`
	Selected            types.Bool   `tfsdk:"selected"`
}

type basicScheduleModule struct {
	Units    types.Int64  `tfsdk:"units"`
	TimeUnit types.String `tfsdk:"time_unit"`
}

type cronScheduleModel struct {
	CronExpression types.String `tfsdk:"cron_expression"`
	CronTimeZone   types.String `tfsdk:"cron_time_zone"`
}

func FlattenConnection(connection *apiclient.Connection) (ConnectionModel, diag.Diagnostics) {
	var data ConnectionModel
	var diags diag.Diagnostics

	data.Id = types.StringValue(connection.ConnectionId)
	data.SourceId = types.StringValue(connection.SourceId)
	data.DestinationId = types.StringValue(connection.DestinationId)
	data.Status = types.StringValue(connection.Status)
	data.Name = types.StringValue(connection.Name)
	data.NamespaceDefinition = types.StringValue(connection.NamespaceDefinition)
	if connection.NamespaceFormat != "" {
		data.NamespaceFormat = types.StringValue(connection.NamespaceFormat)
	} else {
		data.NamespaceFormat = types.StringNull()
	}
	if connection.Prefix != "" {
		data.Prefix = types.StringValue(connection.Prefix)
	} else {
		data.Prefix = types.StringNull()
	}
	if connection.ScheduleType != "" {
		data.ScheduleType = types.StringValue(connection.ScheduleType)
	} else {
		data.ScheduleType = types.StringNull()
	}
	if connection.SourceCatalogId != "" {
		data.SourceCatalogId = types.StringValue(connection.SourceCatalogId)
	} else {
		data.SourceCatalogId = types.StringNull()
	}
	if connection.Geography != "" {
		data.Geography = types.StringValue(connection.Geography)
	} else {
		data.Geography = types.StringNull()
	}
	if connection.BreakingChange != nil {
		data.BreakingChange = types.BoolValue(*connection.BreakingChange)
	} else {
		data.BreakingChange = types.BoolNull()
	}

	if connection.OperationIds != nil {
		var operationIds []attr.Value
		for _, op := range connection.OperationIds {
			operationIds = append(operationIds, types.StringValue(op))
		}
		data.OperationIds, diags = types.ListValue(types.StringType, operationIds)
		if diags.HasError() {
			return data, diags
		}
	} else {
		data.OperationIds, diags = types.ListValue(types.StringType, make([]attr.Value, 0))
	}

	if connection.SyncCatalog != nil {
		if connection.SyncCatalog.Streams != nil {
			var streams []syncCatalogModel
			for _, val := range connection.SyncCatalog.Streams {
				stream := syncCatalogModel{
					SourceSchema: sourceStreamSchemaModel{
						Name: types.StringValue(val.Stream.Name),
					},
					DestinationConfig: destinationStreamConfigModel{
						SyncMode:            types.StringValue(val.Config.SyncMode),
						DestinationSyncMode: types.StringValue(val.Config.DestinationSyncMode),
					},
				}

				config, err := val.Stream.JsonSchema.MarshalJSON()
				if err != nil {
					diags.AddError("Client Error", fmt.Sprintf("Unable to read connection, got error: %s", err))
					return data, diags
				}
				stream.SourceSchema.JsonSchema = types.StringValue(string(config))

				if val.Stream.SupportedSyncModes != nil {
					var modes []attr.Value
					for _, mode := range val.Stream.SupportedSyncModes {
						modes = append(modes, types.StringValue(mode))
					}
					stream.SourceSchema.SupportedSyncModes, diags = types.ListValue(types.StringType, modes)
					if diags.HasError() {
						return data, diags
					}
				} else {
					stream.SourceSchema.SupportedSyncModes, diags = types.ListValue(types.StringType, make([]attr.Value, 0))
				}

				if val.Stream.SourceDefinedCursor != nil {
					stream.SourceSchema.SourceDefinedCursor = types.BoolValue(*val.Stream.SourceDefinedCursor)
				} else {
					stream.SourceSchema.SourceDefinedCursor = types.BoolNull()
				}

				if val.Stream.DefaultCursorField != nil {
					var cursors []attr.Value
					for _, cursor := range val.Stream.DefaultCursorField {
						cursors = append(cursors, types.StringValue(cursor))
					}
					stream.SourceSchema.DefaultCursorField, diags = types.ListValue(types.StringType, cursors)
					if diags.HasError() {
						return data, diags
					}
				} else {
					stream.SourceSchema.DefaultCursorField, diags = types.ListValue(types.StringType, make([]attr.Value, 0))
				}

				if val.Stream.SourceDefinedPrimaryKey != nil {
					var keys []attr.Value
					for _, key := range val.Stream.SourceDefinedPrimaryKey {
						var keyParts []attr.Value
						for _, keyPart := range key {
							keyParts = append(keyParts, types.StringValue(keyPart))
						}
						keyPartsVal, diags := types.ListValue(types.StringType, keyParts)
						if diags.HasError() {
							return data, diags
						}
						keys = append(keys, keyPartsVal)
					}
					stream.SourceSchema.SourceDefinedPrimaryKey, diags = types.ListValue(types.ListType{ElemType: types.StringType}, keys)
					if diags.HasError() {
						return data, diags
					}
				} else {
					stream.SourceSchema.SourceDefinedPrimaryKey, diags = types.ListValue(types.ListType{ElemType: types.StringType}, make([]attr.Value, 0))
				}

				if val.Stream.Namespace != "" {
					stream.SourceSchema.Namespace = types.StringValue(val.Stream.Namespace)
				} else {
					stream.SourceSchema.Namespace = types.StringNull()
				}

				if val.Config.CursorField != nil {
					var cursors []attr.Value
					for _, cursor := range val.Config.CursorField {
						cursors = append(cursors, types.StringValue(cursor))
					}
					stream.DestinationConfig.CursorField, diags = types.ListValue(types.StringType, cursors)
					if diags.HasError() {
						return data, diags
					}
				} else {
					stream.DestinationConfig.CursorField, diags = types.ListValue(types.StringType, make([]attr.Value, 0))
				}

				if val.Config.PrimaryKey != nil {
					var keys []attr.Value
					for _, key := range val.Config.PrimaryKey {
						var keyParts []attr.Value
						for _, keyPart := range key {
							keyParts = append(keyParts, types.StringValue(keyPart))
						}
						keyPartsVal, diags := types.ListValue(types.StringType, keyParts)
						if diags.HasError() {
							return data, diags
						}
						keys = append(keys, keyPartsVal)
					}
					stream.DestinationConfig.PrimaryKey, diags = types.ListValue(types.ListType{ElemType: types.StringType}, keys)
					if diags.HasError() {
						return data, diags
					}
				} else {
					stream.DestinationConfig.PrimaryKey, diags = types.ListValue(types.ListType{ElemType: types.StringType}, make([]attr.Value, 0))
				}

				if val.Config.AliasName != "" {
					stream.DestinationConfig.AliasName = types.StringValue(val.Config.AliasName)
				} else {
					stream.DestinationConfig.AliasName = types.StringNull()
				}

				if val.Config.Selected != nil {
					stream.DestinationConfig.Selected = types.BoolValue(*val.Config.Selected)
				} else {
					stream.DestinationConfig.Selected = types.BoolNull()
				}

				streams = append(streams, stream)
			}
			data.SyncCatalog = &streams
		}
	}

	if connection.ScheduleData != nil {
		if connection.ScheduleData.BasicSchedule != nil {
			data.BasicSchedule = &basicScheduleModule{
				Units:    types.Int64Value(connection.ScheduleData.BasicSchedule.Units),
				TimeUnit: types.StringValue(connection.ScheduleData.BasicSchedule.TimeUnit),
			}
		}
		if connection.ScheduleData.Cron != nil {
			data.CronSchedule = &cronScheduleModel{
				CronExpression: types.StringValue(connection.ScheduleData.Cron.CronExpression),
				CronTimeZone:   types.StringValue(connection.ScheduleData.Cron.CronTimeZone),
			}
		}
	}

	if connection.ResourceRequirements != nil {
		data.ResourceRequirements = FlattenResourceRequirements(connection.ResourceRequirements)
	}

	return data, diags
}
