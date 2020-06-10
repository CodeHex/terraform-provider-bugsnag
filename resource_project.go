package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {Type: schema.TypeString, Required: true},
			"type": {Type: schema.TypeString, Required: true, ForceNew: true},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
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
	err = d.Set("name", project.Name)
	if err != nil {
		return err
	}

	// The type of project is not available on the API, so assume the state is correct
	err = d.Set("type", d.State().Attributes["type"])
	if err != nil {
		return err
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
