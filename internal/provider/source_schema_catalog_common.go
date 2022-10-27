package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
)

// SourceSchemaCatalogModel describes the data source data model.
type SourceSchemaCatalogModel struct {
	Id          types.String        `tfsdk:"id"`
	SourceId    types.String        `tfsdk:"source_id"`
	SyncCatalog *[]SyncCatalogModel `tfsdk:"sync_catalog"`
}

type SyncCatalogModel struct {
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

func FlattenSyncCatalog(ssc *apiclient.SyncCatalog) (*[]SyncCatalogModel, diag.Diagnostics) {
	var data *[]SyncCatalogModel
	var diags diag.Diagnostics

	if ssc.Streams != nil {
		var streams []SyncCatalogModel
		for _, val := range ssc.Streams {
			stream := SyncCatalogModel{
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
		data = &streams
	}

	return data, diags
}
