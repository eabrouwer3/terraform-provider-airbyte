package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type WorkspaceIdBody struct {
	WorkspaceId string `json:"workspaceId"`
}

type WorkspaceNameBody struct {
	Name string `json:"name,omitempty"`
}

type CommonWorkspaceFields struct {
	Email                   string         `json:"email,omitempty"`
	AnonymousDataCollection *bool          `json:"anonymousDataCollection,omitempty"`
	News                    *bool          `json:"news,omitempty"`
	SecurityUpdates         *bool          `json:"securityUpdates,omitempty"`
	Notifications           []Notification `json:"notifications,omitempty"`
	DisplaySetupWizard      *bool          `json:"displaySetupWizard,omitempty"`
}

type Workspace struct {
	WorkspaceIdBody
	WorkspaceNameBody
	CommonWorkspaceFields
	CustomerId           string         `json:"customerId"`
	Slug                 string         `json:"slug"`
	InitialSetupComplete *bool          `json:"initialSetupComplete,omitempty"`
	SecurityUpdates      *bool          `json:"securityUpdates,omitempty"`
	Notifications        []Notification `json:"notifications"`
	FirstCompletedSync   *bool          `json:"firstCompletedSync,omitempty"`
	FeedbackDone         *bool          `json:"feedbackDone,omitempty"`
	DefaultGeography     string         `json:"defaultGeography"`
}

type NewWorkspace struct {
	WorkspaceNameBody
	CommonWorkspaceFields
}

type UpdatedWorkspace struct {
	WorkspaceIdBody
	CommonWorkspaceFields
}

type Notification struct {
	NotificationType   string             `json:"notificationType"`
	SendOnSuccess      *bool              `json:"sendOnSuccess,omitempty"`
	SendOnFailure      *bool              `json:"sendOnFailure,omitempty"`
	SlackConfiguration SlackConfiguration `json:"slackConfiguration"`
}

type SlackConfiguration struct {
	Webhook string `json:"webhook"`
}

type WorkspaceList struct {
	Workspaces []*Workspace `json:"workspaces"`
}

func (c *ApiClient) GetWorkspaceById(workspaceId string) (*Workspace, error) {
	rb, err := json.Marshal(WorkspaceIdBody{
		WorkspaceId: workspaceId,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/workspaces/get", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workspace := Workspace{}
	err = json.Unmarshal(body, &workspace)
	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (c *ApiClient) GetWorkspaceBySlug(slug string) (*Workspace, error) {
	rb, err := json.Marshal(struct {
		Slug string `json:"slug"`
	}{
		Slug: slug,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/workspaces/get_by_slug", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workspace := Workspace{}
	err = json.Unmarshal(body, &workspace)
	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (c *ApiClient) GetWorkspaces() ([]*Workspace, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/workspaces/list", c.HostURL, BaseUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	wl := WorkspaceList{}
	err = json.Unmarshal(body, &wl)
	if err != nil {
		return nil, err
	}

	return wl.Workspaces, nil
}

func (c *ApiClient) CreateWorkspace(newWorkspace NewWorkspace) (*Workspace, error) {
	rb, err := json.Marshal(newWorkspace)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/workspaces/create", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workspace := Workspace{}
	err = json.Unmarshal(body, &workspace)
	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (c *ApiClient) UpdateWorkspace(updatedWorkspace UpdatedWorkspace) (*Workspace, error) {
	rb, err := json.Marshal(updatedWorkspace)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/workspaces/update", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workspace := Workspace{}
	err = json.Unmarshal(body, &workspace)
	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (c *ApiClient) DeleteWorkspace(workspaceId string) error {
	rb, err := json.Marshal(WorkspaceIdBody{WorkspaceId: workspaceId})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/workspaces/delete", c.HostURL, BaseUrl), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
