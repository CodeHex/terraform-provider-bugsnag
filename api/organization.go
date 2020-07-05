package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// OrgCreator defines details about the creator of the organization
type OrgCreator struct {
	Email string `json:"email,omitempty"`
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
}

// Organization defines the details of a Bugsnag organization as provided by the Bugsnag data access API
type Organization struct {
	ID            string      `json:"id,omitempty"`
	Name          string      `json:"name,omitempty"`
	BillingEmails []string    `json:"billing_emails,omitempty"`
	AutoUpgrade   *bool       `json:"auto_upgrade,omitempty"`
	Slug          string      `json:"slug,omitempty"`
	Creator       *OrgCreator `json:"creator,omitempty"`
	CreatedAt     string      `json:"created_at,omitempty"`
	UpdatedAt     string      `json:"updated_at,omitempty"`
}

const orgsPath = "user/organizations"
const orgIDPath = "organizations/%s"

// GetCurrentOrganization returns the ID of the organization associated with the configured data access token
func (c *Client) GetCurrentOrganization() (string, error) {
	var orgs []Organization
	err := c.callAPI(http.MethodGet, orgsPath, nil, &orgs, http.StatusOK)
	if err != nil {
		return "", err
	}

	// We only expect one org as we only support data access API tokens
	if len(orgs) != 1 {
		return "", fmt.Errorf("unexpected number of orgs for token (%d), expecting 1", len(orgs))
	}
	return orgs[0].ID, nil
}

// GetOrganization returns the organization associated with the configured data access token
func (c *Client) GetOrganization() (*Organization, error) {
	var org Organization
	err := c.callAPI(http.MethodGet, fmt.Sprintf(orgIDPath, c.OrgID), nil, &org, http.StatusOK)
	return &org, err
}

// UpdateOrganization updates organization settings associated with the configured data access token
func (c *Client) UpdateOrganization(org *Organization) (*Organization, error) {
	body, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}

	var updatedOrg Organization
	err = c.callAPI(http.MethodPatch, fmt.Sprintf(orgIDPath, org.ID), body, &updatedOrg, http.StatusOK)
	return &updatedOrg, err
}
