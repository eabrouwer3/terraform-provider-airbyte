package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
)

// OperationModel describes the data source data model.
type OperationModel struct {
	Id                  types.String  `tfsdk:"id"`
	WorkspaceId         types.String  `tfsdk:"workspace_id"`
	Name                types.String  `tfsdk:"name"`
	OperatorType        types.String  `tfsdk:"operator_type"`
	NormalizationOption types.String  `tfsdk:"normalization_option"`
	Dbt                 *dbtModel     `tfsdk:"dbt"`
	Webhook             *webhookModel `tfsdk:"webhook"`
}

type dbtModel struct {
	GitRepoUrl    types.String `tfsdk:"git_repo_url"`
	GitRepoBranch types.String `tfsdk:"git_repo_branch"`
	DockerImage   types.String `tfsdk:"docker_image"`
	DbtArguments  types.String `tfsdk:"dbt_arguments"`
}

type webhookModel struct {
	ExecutionUrl    types.String `tfsdk:"execution_url"`
	ExecutionBody   types.String `tfsdk:"execution_body"`
	WebhookConfigId types.String `tfsdk:"webhook_config_id"`
}

func FlattenOperation(operation *apiclient.Operation) OperationModel {
	var data OperationModel

	data.Id = types.StringValue(operation.OperationId)
	data.WorkspaceId = types.StringValue(operation.WorkspaceId)
	data.Name = types.StringValue(operation.Name)
	data.OperatorType = types.StringValue(operation.OperatorConfiguration.OperatorType)

	if v := operation.OperatorConfiguration.Normalization; v != nil {
		if v.Option != "" {
			data.NormalizationOption = types.StringValue(v.Option)
		} else {
			data.NormalizationOption = types.StringNull()
		}
	}

	if dbt := operation.OperatorConfiguration.Dbt; dbt != nil {
		model := dbtModel{}
		if v := dbt.GitRepoUrl; v != "" {
			model.GitRepoUrl = types.StringValue(v)
		} else {
			model.GitRepoUrl = types.StringNull()
		}
		if v := dbt.GitRepoBranch; v != "" {
			model.GitRepoBranch = types.StringValue(v)
		} else {
			model.GitRepoBranch = types.StringNull()
		}
		if v := dbt.DockerImage; v != "" {
			model.DockerImage = types.StringValue(v)
		} else {
			model.DockerImage = types.StringNull()
		}
		if v := dbt.DbtArguments; v != "" {
			model.DbtArguments = types.StringValue(v)
		} else {
			model.DbtArguments = types.StringNull()
		}
		data.Dbt = &model
	}

	if dbt := operation.OperatorConfiguration.Webhook; dbt != nil {
		model := webhookModel{}
		if v := dbt.ExecutionUrl; v != "" {
			model.ExecutionUrl = types.StringValue(v)
		} else {
			model.ExecutionUrl = types.StringNull()
		}
		if v := dbt.ExecutionBody; v != "" {
			model.ExecutionBody = types.StringValue(v)
		} else {
			model.ExecutionBody = types.StringNull()
		}
		if v := dbt.WebhookConfigId; v != "" {
			model.WebhookConfigId = types.StringValue(v)
		} else {
			model.WebhookConfigId = types.StringNull()
		}
		data.Webhook = &model
	}

	return data
}

func GetCommonOperationFields(data OperationModel) apiclient.CommonOperationFields {
	fields := apiclient.CommonOperationFields{
		Name: data.Name.ValueString(),
		OperatorConfiguration: apiclient.OperationConfig{
			OperatorType: data.OperatorType.ValueString(),
		},
	}

	if v := data.NormalizationOption; !v.IsUnknown() {
		fields.OperatorConfiguration.Normalization = &apiclient.NormalizationOption{
			Option: v.ValueString(),
		}
	}
	if data.Dbt != nil {
		dbt := apiclient.DbtConfig{}
		if v := data.Dbt.GitRepoUrl; !v.IsUnknown() {
			dbt.GitRepoUrl = v.ValueString()
		}
		if v := data.Dbt.GitRepoBranch; !v.IsUnknown() {
			dbt.GitRepoBranch = v.ValueString()
		}
		if v := data.Dbt.DockerImage; !v.IsUnknown() {
			dbt.DockerImage = v.ValueString()
		}
		if v := data.Dbt.DbtArguments; !v.IsUnknown() {
			dbt.DbtArguments = v.ValueString()
		}
		fields.OperatorConfiguration.Dbt = &dbt
	}
	if data.Webhook != nil {
		webhook := apiclient.WebhookConfig{}
		if v := data.Webhook.ExecutionUrl; !v.IsUnknown() {
			webhook.ExecutionUrl = v.ValueString()
		}
		if v := data.Webhook.ExecutionBody; !v.IsUnknown() {
			webhook.ExecutionBody = v.ValueString()
		}
		if v := data.Webhook.WebhookConfigId; !v.IsUnknown() {
			webhook.WebhookConfigId = v.ValueString()
		}
		fields.OperatorConfiguration.Webhook = &webhook
	}

	return fields
}
