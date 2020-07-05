package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Project defines the details of a Bugsnag project as provided by the Bugsnag data access API.
type Project struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	APIKey    string `json:"api_key,omitempty"`
	Slug      string `json:"slug,omitempty"`
	URL       string `json:"url,omitempty"`
	HTMLURL   string `json:"html_url,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

const getProjectPath = "projects/%s"
const createProjectPath = "organizations/%s/projects"
const listProjects = "organizations/%s/projects"

// GetProject returns a specific project using the provided project ID
func (c *Client) GetProject(id string) (*Project, error) {
	var project Project
	err := c.callAPI(http.MethodGet, fmt.Sprintf(getProjectPath, id), nil, &project, http.StatusOK)
	return &project, err
}

// CreateProject creates a new project in organization
func (c *Client) CreateProject(project *Project) (*Project, error) {
	body, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	var createdProject Project
	err = c.callAPI(http.MethodPost, fmt.Sprintf(createProjectPath, c.OrgID), body, &createdProject, http.StatusOK)
	return &createdProject, err
}

// DeleteProject deletes a project in the organization
func (c *Client) DeleteProject(id string) error {
	return c.callAPI(http.MethodDelete, fmt.Sprintf(getProjectPath, id), nil, nil, http.StatusNoContent)
}

// UpdateProject updates a project in the organization
func (c *Client) UpdateProject(project *Project) (*Project, error) {
	body, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	var updatedProject Project
	err = c.callAPI(http.MethodPatch, fmt.Sprintf(getProjectPath, project.ID), body, &updatedProject, http.StatusOK)
	return &updatedProject, err
}

// ListProjects provides a list of projects in the organization. A query string can be provided to narrow down the
// list of projects based on project name. The list can be sorted.
func (c *Client) ListProjects(query string, sort string, direction string) ([]Project, error) {
	projects := make([]Project, 0)
	uri, err := url.Parse(fmt.Sprintf(listProjects, c.OrgID))
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("per_page", "100")
	if query != "" {
		v.Set("q", query)
	}
	if sort != "" {
		v.Set("sort", sort)
	}
	if direction != "" {
		v.Set("direction", direction)
	}
	uri.RawQuery = v.Encode()
	for {
		projectsPage := make([]Project, 0)
		uri, err = c.callPagedAPI(http.MethodGet, uri.String(), nil, &projectsPage, http.StatusOK)
		if err != nil {
			return nil, err
		}
		projects = append(projects, projectsPage...)
		if uri == nil {
			break
		}
	}
	return projects, err
}
