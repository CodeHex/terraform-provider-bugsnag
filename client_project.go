package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIProject struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

const getProjectPath = "projects/%s"
const createProjectPath = "organizations/%s/projects"

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
	err = c.callAPI(http.MethodPost, fmt.Sprintf(createProjectPath, c.OrgID), body, &createdProject, http.StatusOK)
	return &createdProject, err
}

func (c *Client) DeleteProject(id string) error {
	return c.callAPI(http.MethodDelete, fmt.Sprintf(getProjectPath, id), nil, nil, http.StatusNoContent)
}
