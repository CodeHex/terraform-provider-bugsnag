package bugsnag

import (
	"errors"

	"github.com/codehex/terraform-provider-bugsnag/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCollaborator() *schema.Resource {
	return &schema.Resource{
		Create: resourceCollaboratorCreate,
		Read:   resourceCollaboratorRead,
		Update: resourceCollaboratorUpdate,
		Delete: resourceCollaboratorDelete,

		Schema: map[string]*schema.Schema{

			"email": {Type: schema.TypeString, Required: true, ForceNew: true},
			"admin": {Type: schema.TypeBool, Optional: true, Default: false},
			"project_ids": {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{
				Type: schema.TypeString,
			}},

			"name":                     {Type: schema.TypeString, Computed: true},
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
	c := m.(*api.Client)
	projIDs := readProjectIDs(d)
	admin := d.Get("admin").(bool)
	if admin && len(projIDs) != 0 {
		return errors.New("unable to create collaborator, project IDs are not supported when user is admin")
	}
	collab := &api.Collaborator{
		Email:      d.Get("email").(string),
		Admin:      &admin,
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
	c := m.(*api.Client)
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
	c := m.(*api.Client)
	admin := d.Get("admin").(bool)
	collaborator := &api.Collaborator{
		ID:         d.Id(),
		Admin:      &admin,
		ProjectIDs: readProjectIDs(d),
	}

	if admin && len(collaborator.ProjectIDs) != 0 {
		return errors.New("unable to update collaborator, project IDs provided on admin user")
	}
	_, err := c.UpdateCollaborator(collaborator)
	if err != nil {
		return err
	}
	return resourceCollaboratorRead(d, m)
}

func resourceCollaboratorDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*api.Client)
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
