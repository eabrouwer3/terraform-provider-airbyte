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

	data.Id = types.String{Value: workspace.WorkspaceId}
	data.Slug = types.String{Value: workspace.Slug}
	data.CustomerId = types.String{Value: workspace.CustomerId}
	data.Name = types.String{Value: workspace.Name}
	if len(workspace.Notifications) > 0 {
		data.NotificationConfig = []workspaceNotificationConfigModel{}
		for _, notifConfig := range workspace.Notifications {
			data.NotificationConfig = append(data.NotificationConfig, workspaceNotificationConfigModel{
				NotificationType: types.String{Value: notifConfig.NotificationType},
				SendOnSuccess:    types.Bool{Value: *notifConfig.SendOnSuccess},
				SendOnFailure:    types.Bool{Value: *notifConfig.SendOnFailure},
				SlackWebhook:     types.String{Value: notifConfig.SlackConfiguration.Webhook},
			})
		}
	}
	data.Email = types.String{Value: workspace.Email}
	data.InitialSetupComplete = types.Bool{Value: *workspace.InitialSetupComplete}
	data.DisplaySetupWizard = types.Bool{Value: *workspace.DisplaySetupWizard}
	data.AnonymousDataCollection = types.Bool{Value: *workspace.AnonymousDataCollection}
	data.News = types.Bool{Value: *workspace.News}
	data.SecurityUpdates = types.Bool{Value: *workspace.SecurityUpdates}
	if v := workspace.FirstCompletedSync; v != nil {
		data.FirstCompletedSync = types.Bool{Value: *v}
	} else {
		data.FirstCompletedSync = types.Bool{Null: true}
	}
	if v := workspace.FeedbackDone; v != nil {
		data.FeedbackDone = types.Bool{Value: *v}
	} else {
		data.FeedbackDone = types.Bool{Null: true}
	}
	data.DefaultGeography = types.String{Value: workspace.DefaultGeography}

	return data
}
