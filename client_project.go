package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type APIProject struct {
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

func (c *Client) GetProject(id string) (*APIProject, error) {
	var project APIProject
	err := c.callAPI(http.MethodGet, fmt.Sprintf(getProjectPath, id), nil, &project, http.StatusOK)
	return &project, err
}

func (c *Client) CreateProject(project *APIProject) (*APIProject, error) {
	body, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	var createdProject APIProject
	err = c.callAPI(http.MethodPost, fmt.Sprintf(createProjectPath, c.orgID), body, &createdProject, http.StatusOK)
	return &createdProject, err
}

func (c *Client) DeleteProject(id string) error {
	return c.callAPI(http.MethodDelete, fmt.Sprintf(getProjectPath, id), nil, nil, http.StatusNoContent)
}

func (c *Client) UpdateProject(project *APIProject) (*APIProject, error) {
	body, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	var updatedProject APIProject
	err = c.callAPI(http.MethodPatch, fmt.Sprintf(getProjectPath, project.ID), body, &updatedProject, http.StatusOK)
	return &updatedProject, err
}

func (c *Client) ListProjects(query string) ([]APIProject, error) {
	projects := make([]APIProject, 0)
	uri, err := url.Parse(fmt.Sprintf(listProjects, c.orgID))
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("per_page", "100")
	if query != "" {
		v.Set("q", query)
	}
	uri.RawQuery = v.Encode()
	for {
		projectsPage := make([]APIProject, 0)
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
