package provider

import (
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
	SyncCatalog          *[]SyncCatalogModel        `tfsdk:"sync_catalog"`
	ScheduleType         types.String               `tfsdk:"schedule_type"`
	BasicSchedule        *basicScheduleModule       `tfsdk:"basic_schedule"`
	CronSchedule         *cronScheduleModel         `tfsdk:"cron_schedule"`
	ResourceRequirements *ResourceRequirementsModel `tfsdk:"resource_requirements"`
	SourceCatalogId      types.String               `tfsdk:"source_catalog_id"`
	Geography            types.String               `tfsdk:"geography"`
	BreakingChange       types.Bool                 `tfsdk:"breaking_change"`
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
			data.SyncCatalog, diags = FlattenSyncCatalog(connection.SyncCatalog)
			if diags.HasError() {
				return data, diags
			}
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
