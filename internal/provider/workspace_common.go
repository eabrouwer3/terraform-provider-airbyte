package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
)

// WorkspaceModel describes the data source data model.
type WorkspaceModel struct {
	Id                      types.String                       `tfsdk:"id"`
	CustomerId              types.String                       `tfsdk:"customer_id"`
	Email                   types.String                       `tfsdk:"email"`
	Name                    types.String                       `tfsdk:"name"`
	Slug                    types.String                       `tfsdk:"slug"`
	InitialSetupComplete    types.Bool                         `tfsdk:"initial_setup_complete"`
	DisplaySetupWizard      types.Bool                         `tfsdk:"display_setup_wizard"`
	AnonymousDataCollection types.Bool                         `tfsdk:"anonymous_data_collection"`
	News                    types.Bool                         `tfsdk:"news"`
	SecurityUpdates         types.Bool                         `tfsdk:"security_updates"`
	NotificationConfig      []workspaceNotificationConfigModel `tfsdk:"notification_config"`
	FirstCompletedSync      types.Bool                         `tfsdk:"first_completed_sync"`
	FeedbackDone            types.Bool                         `tfsdk:"feedback_done"`
	DefaultGeography        types.String                       `tfsdk:"default_geography"`
}

type workspaceNotificationConfigModel struct {
	NotificationType types.String `tfsdk:"notification_type"`
	SendOnSuccess    types.Bool   `tfsdk:"send_on_success"`
	SendOnFailure    types.Bool   `tfsdk:"send_on_failure"`
	SlackWebhook     types.String `tfsdk:"slack_webhook"`
}

func FlattenWorkspace(workspace *apiclient.Workspace) WorkspaceModel {
	var data WorkspaceModel

	data.Id = types.StringValue(workspace.WorkspaceId)
	data.Slug = types.StringValue(workspace.Slug)
	data.CustomerId = types.StringValue(workspace.CustomerId)
	data.Name = types.StringValue(workspace.Name)
	if len(workspace.Notifications) > 0 {
		data.NotificationConfig = []workspaceNotificationConfigModel{}
		for _, notifConfig := range workspace.Notifications {
			data.NotificationConfig = append(data.NotificationConfig, workspaceNotificationConfigModel{
				NotificationType: types.StringValue(notifConfig.NotificationType),
				SendOnSuccess:    types.BoolValue(*notifConfig.SendOnSuccess),
				SendOnFailure:    types.BoolValue(*notifConfig.SendOnFailure),
				SlackWebhook:     types.StringValue(notifConfig.SlackConfiguration.Webhook),
			})
		}
	}
	data.Email = types.StringValue(workspace.Email)
	data.InitialSetupComplete = types.BoolValue(*workspace.InitialSetupComplete)
	data.DisplaySetupWizard = types.BoolValue(*workspace.DisplaySetupWizard)
	data.AnonymousDataCollection = types.BoolValue(*workspace.AnonymousDataCollection)
	data.News = types.BoolValue(*workspace.News)
	data.SecurityUpdates = types.BoolValue(*workspace.SecurityUpdates)
	if v := workspace.FirstCompletedSync; v != nil {
		data.FirstCompletedSync = types.BoolValue(*v)
	} else {
		data.FirstCompletedSync = types.BoolNull()
	}
	if v := workspace.FeedbackDone; v != nil {
		data.FeedbackDone = types.BoolValue(*v)
	} else {
		data.FeedbackDone = types.BoolNull()
	}
	data.DefaultGeography = types.StringValue(workspace.DefaultGeography)

	return data
}
