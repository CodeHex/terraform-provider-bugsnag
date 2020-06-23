package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name":       {Type: schema.TypeString, Required: true},
			"type":       {Type: schema.TypeString, Required: true, ValidateFunc: validateType},
			"api_key":    {Type: schema.TypeString, Computed: true},
			"slug":       {Type: schema.TypeString, Computed: true},
			"url":        {Type: schema.TypeString, Computed: true},
			"html_url":   {Type: schema.TypeString, Computed: true},
			"created_at": {Type: schema.TypeString, Computed: true},
			"updated_at": {Type: schema.TypeString, Computed: true},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func validateType(val interface{}, key string) (warns []string, err []error) {
	v := val.(string)
	isValid := false
	for _, t := range ValidProjectTypes() {
		if t == v {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, []error{fmt.Errorf("unrecognized project type '%s'", v)}
	}
	return nil, nil
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	project := &APIProject{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}
	project, err := c.CreateProject(project)
	if err != nil {
		return err
	}
	d.SetId(project.ID)
	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	project, err := c.GetProject(d.Id())
	if err != nil {
		return err
	}

	fields := map[string]string{
		"name":       project.Name,
		"type":       project.Type,
		"api_key":    project.APIKey,
		"slug":       project.Slug,
		"url":        project.URL,
		"html_url":   project.HTMLURL,
		"created_at": project.CreatedAt,
		"updated_at": project.UpdatedAt,
	}

	for field, val := range fields {
		err = d.Set(field, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	project := &APIProject{
		ID: d.Id(),
	}

	if d.HasChange("name") {
		project.Name = d.Get("name").(string)
	}
	if d.HasChange("type") {
		project.Type = d.Get("type").(string)
	}
	_, err := c.UpdateProject(project)
	if err != nil {
		return err
	}
	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	return c.DeleteProject(d.Id())
}
