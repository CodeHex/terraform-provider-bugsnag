package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APICollaborator struct {
	ID                     string   `json:"id,omitempty"`
	Name                   string   `json:"name,omitempty"`
	Email                  string   `json:"email,omitempty"`
	Admin                  bool     `json:"admin,omitempty"`
	Password               string   `json:"password,omitempty"`
	IsAdmin                bool     `json:"is_admin,omitempty"`
	ProjectIDs             []string `json:"project_ids,omitempty"`
	TwoFAEnabled           bool     `json:"two_factor_enabled,omitempty"`
	TwoFAEnabledOn         string   `json:"two_factor_enabled_on,omitempty"`
	RecoveryCodesRemaining int      `json:"recovery_codes_remaining,omitempty"`
	PasswordUpdatedOn      string   `json:"password_updated_on,omitempty"`
	ShowTimeInUTC          bool     `json:"show_time_in_utc,omitempty"`
	Heroku                 bool     `json:"heroku,omitempty"`
	CreatedAt              string   `json:"created_at,omitempty"`
	PendingInvitation      bool     `json:"pending_invitation,omitempty"`
	LastRequestAt          string   `json:"last_request_at,omitempty"`
	PaidFor                bool     `json:"paid_for,omitempty"`
}

const collaboratorPath = "organizations/%s/collaborators/%s"
const createCollaboratorPath = "organizations/%s/collaborators"

func (c *Client) GetCollaborator(id string) (*APICollaborator, error) {
	var collab APICollaborator
	err := c.callAPI(http.MethodGet, fmt.Sprintf(collaboratorPath, c.OrgID, id), nil, &collab, http.StatusOK)
	return &collab, err
}

func (c *Client) CreateCollaborator(collab *APICollaborator) (*APICollaborator, error) {
	body, err := json.Marshal(collab)
	if err != nil {
		return nil, err
	}

	var createdCollab APICollaborator
	err = c.callAPI(http.MethodPost, fmt.Sprintf(createCollaboratorPath, c.OrgID), body, &createdCollab, http.StatusOK)
	return &createdCollab, err
}

func (c *Client) DeleteCollaborator(id string) error {
	return c.callAPI(http.MethodDelete, fmt.Sprintf(collaboratorPath, c.OrgID, id), nil, nil, http.StatusNoContent)
}

func (c *Client) UpdateCollaborator(collab *APICollaborator) (*APICollaborator, error) {
	body, err := json.Marshal(collab)
	if err != nil {
		return nil, err
	}

	var updatedCollaborator APICollaborator
	err = c.callAPI(http.MethodPatch, fmt.Sprintf(collaboratorPath, c.OrgID, collab.ID), body, &updatedCollaborator, http.StatusOK)
	return &updatedCollaborator, err
}
