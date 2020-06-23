package main

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCollaborator() *schema.Resource {
	return &schema.Resource{
		Create: resourceCollaboratorCreate,
		Read:   resourceCollaboratorRead,
		Update: resourceCollaboratorUpdate,
		Delete: resourceCollaboratorDelete,

		Schema: map[string]*schema.Schema{
			"name":     {Type: schema.TypeString, Required: true, ForceNew: true},
			"email":    {Type: schema.TypeString, Required: true, ForceNew: true},
			"password": {Type: schema.TypeString, Optional: true, ForceNew: true, Sensitive: true},
			"admin":    {Type: schema.TypeBool, Optional: true, Default: false},
			"project_ids": {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{
				Type: schema.TypeString,
			}},

			"two_factor_enabled":       {Type: schema.TypeBool, Computed: true},
			"two_factor_enabled_on":    {Type: schema.TypeString, Computed: true},
			"recovery_codes_remaining": {Type: schema.TypeInt, Computed: true},
			"password_updated_on":      {Type: schema.TypeString, Computed: true},
			"show_time_in_utc":         {Type: schema.TypeBool, Computed: true},
			"heroku":                   {Type: schema.TypeBool, Computed: true},
			"created_at":               {Type: schema.TypeString, Computed: true},
			"pending_invitation":       {Type: schema.TypeBool, Computed: true},
			"last_request_at":          {Type: schema.TypeString, Computed: true},
			"paid_for":                 {Type: schema.TypeBool, Computed: true},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCollaboratorCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	projIDs := readProjectIDs(d)
	name := d.Get("name").(string)
	password := d.Get("password").(string)
	admin := d.Get("admin").(bool)
	if name == "" && password != "" {
		return errors.New("unable to create collaborator, password is not supported without user name")
	}
	if admin && len(projIDs) != 0 {
		return errors.New("unable to create collaborator, project IDs are not supported when user is admin")
	}
	collab := &APICollaborator{
		Name:       name,
		Email:      d.Get("email").(string),
		Password:   password,
		Admin:      admin,
		ProjectIDs: projIDs,
	}
	collab, err := c.CreateCollaborator(collab)
	if err != nil {
		return err
	}
	d.SetId(collab.ID)
	return resourceCollaboratorRead(d, m)
}

func resourceCollaboratorRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	collaborator, err := c.GetCollaborator(d.Id())
	if err != nil {
		return err
	}

	fields := map[string]interface{}{
		"name":                     collaborator.Name,
		"email":                    collaborator.Email,
		"admin":                    collaborator.IsAdmin,
		"project_ids":              collaborator.ProjectIDs,
		"two_factor_enabled":       collaborator.TwoFAEnabled,
		"two_factor_enabled_on":    collaborator.TwoFAEnabledOn,
		"recovery_codes_remaining": collaborator.RecoveryCodesRemaining,
		"password_updated_on":      collaborator.PasswordUpdatedOn,
		"heroku":                   collaborator.Heroku,
		"show_time_in_utc":         collaborator.ShowTimeInUTC,
		"created_at":               collaborator.CreatedAt,
		"pending_invitation":       collaborator.PendingInvitation,
		"last_request_at":          collaborator.LastRequestAt,
		"paid_for":                 collaborator.PaidFor,
	}

	for field, val := range fields {
		err = d.Set(field, val)
		if err != nil {
			return err
		}
	}

	// Ignore project IDs field if they are admin as they will belong to every one
	if collaborator.IsAdmin {
		err = d.Set("project_ids", nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceCollaboratorUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	collaborator := &APICollaborator{
		ID:         d.Id(),
		Admin:      d.Get("admin").(bool),
		ProjectIDs: readProjectIDs(d),
	}

	if collaborator.Admin && len(collaborator.ProjectIDs) != 0 {
		return errors.New("unable to update collaborator, project IDs provided on admin user")
	}
	_, err := c.UpdateCollaborator(collaborator)
	if err != nil {
		return err
	}
	return resourceCollaboratorRead(d, m)
}

func resourceCollaboratorDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	return c.DeleteCollaborator(d.Id())
}

func readProjectIDs(d *schema.ResourceData) []string {
	projIDSet := d.Get("project_ids").(*schema.Set)
	var projIDs []string
	for _, val := range projIDSet.List() {
		projIDs = append(projIDs, val.(string))
	}
	return projIDs
}
