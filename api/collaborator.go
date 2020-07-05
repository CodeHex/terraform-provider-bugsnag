package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Collaborator defines the details of a Bugsnag collaborator as provided by the Bugsnag data access API.
type Collaborator struct {
	ID                     string   `json:"id,omitempty"`
	Name                   string   `json:"name,omitempty"`
	Email                  string   `json:"email,omitempty"`
	Admin                  *bool    `json:"admin,omitempty"`
	Password               string   `json:"password,omitempty"`
	IsAdmin                bool     `json:"is_admin,omitempty"`
	ProjectIDs             []string `json:"project_ids,omitempty"`
	TwoFAEnabled           bool     `json:"two_factor_enabled,omitempty"`
	TwoFAEnabledOn         string   `json:"two_factor_enabled_on,omitempty"`
	RecoveryCodesRemaining int      `json:"recovery_codes_remaining,omitempty"`
	PasswordUpdatedOn      string   `json:"password_updated_on,omitempty"`
	ShowTimeInUTC          bool     `json:"show_time_in_utc,omitempty"`
	CreatedAt              string   `json:"created_at,omitempty"`
	PendingInvitation      bool     `json:"pending_invitation,omitempty"`
	LastRequestAt          string   `json:"last_request_at,omitempty"`
	PaidFor                bool     `json:"paid_for,omitempty"`
}

const collaboratorPath = "organizations/%s/collaborators/%s"
const createCollaboratorPath = "organizations/%s/collaborators"

// GetCollaborator returns a specific collaborator using the provided user ID
func (c *Client) GetCollaborator(id string) (*Collaborator, error) {
	var collab Collaborator
	err := c.callAPI(http.MethodGet, fmt.Sprintf(collaboratorPath, c.OrgID, id), nil, &collab, http.StatusOK)
	return &collab, err
}

// CreateCollaborator creates a new collaborator. If the user does not exist, a user will be created and invited to join the organization
func (c *Client) CreateCollaborator(collab *Collaborator) (*Collaborator, error) {
	body, err := json.Marshal(collab)
	if err != nil {
		return nil, err
	}

	var createdCollab Collaborator
	err = c.callAPI(http.MethodPost, fmt.Sprintf(createCollaboratorPath, c.OrgID), body, &createdCollab, http.StatusOK)
	return &createdCollab, err
}

// DeleteCollaborator will remove the collaborator from the organization
func (c *Client) DeleteCollaborator(id string) error {
	return c.callAPI(http.MethodDelete, fmt.Sprintf(collaboratorPath, c.OrgID, id), nil, nil, http.StatusNoContent)
}

// UpdateCollaborator will update access permissions for the specified collaborator
func (c *Client) UpdateCollaborator(collab *Collaborator) (*Collaborator, error) {
	body, err := json.Marshal(collab)
	if err != nil {
		return nil, err
	}

	var updatedCollaborator Collaborator
	err = c.callAPI(http.MethodPatch, fmt.Sprintf(collaboratorPath, c.OrgID, collab.ID), body, &updatedCollaborator, http.StatusOK)
	return &updatedCollaborator, err
}
