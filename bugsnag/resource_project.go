package bugsnag

import (
	"github.com/codehex/terraform-provider-bugsnag/api"
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
			"type":       {Type: schema.TypeString, Required: true, ValidateFunc: validateValueFunc(validProjectTypes())},
			"api_key":    {Type: schema.TypeString, Computed: true},
			"slug":       {Type: schema.TypeString, Computed: true},
			"html_url":   {Type: schema.TypeString, Computed: true},
			"created_at": {Type: schema.TypeString, Computed: true},
			"updated_at": {Type: schema.TypeString, Computed: true},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*api.Client)
	project := &api.Project{
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
	c := m.(*api.Client)
	project, err := c.GetProject(d.Id())
	if err != nil {
		return err
	}

	fields := map[string]string{
		"name":       project.Name,
		"type":       project.Type,
		"api_key":    project.APIKey,
		"slug":       project.Slug,
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
	c := m.(*api.Client)
	project := &api.Project{
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
	c := m.(*api.Client)
	return c.DeleteProject(d.Id())
}
