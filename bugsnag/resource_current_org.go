package bugsnag

import (
	"errors"
	"fmt"

	"github.com/codehex/terraform-provider-bugsnag/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCurrentOrg() *schema.Resource {
	return &schema.Resource{
		Create: resourceCurrentOrgCreate,
		Read:   resourceCurrentOrgRead,
		Update: resourceCurrentOrgUpdate,
		Delete: resourceCurrentOrgDelete,

		Schema: map[string]*schema.Schema{
			"name": {Type: schema.TypeString, Optional: true, Computed: true},
			"billing_emails": {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{
				Type: schema.TypeString,
			}},
			"auto_upgrade": {Type: schema.TypeBool, Optional: true, Default: false},

			"slug": {Type: schema.TypeString, Computed: true},
			"creator": {Type: schema.TypeList, MaxItems: 1, Computed: true, Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"email": {Type: schema.TypeString, Computed: true},
					"id":    {Type: schema.TypeString, Computed: true},
					"name":  {Type: schema.TypeString, Computed: true},
				},
			}},
			"created_at": {Type: schema.TypeString, Computed: true},
			"updated_at": {Type: schema.TypeString, Computed: true},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Note this is currently restricted to user authentication only
func resourceCurrentOrgCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*api.Client)
	return fmt.Errorf("current organization must be imported, please run `terraform import bugsnag_current_org.<resource_name> %s`", c.OrgID)
}

func resourceCurrentOrgRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*api.Client)
	org, err := c.GetOrganization()
	if err != nil {
		return err
	}

	fields := map[string]interface{}{
		"name":           org.Name,
		"billing_emails": org.BillingEmails,
		"auto_upgrade":   *org.AutoUpgrade,
		"slug":           org.Slug,
		"created_at":     org.CreatedAt,
		"updated_at":     org.UpdatedAt,
	}

	for field, val := range fields {
		err = d.Set(field, val)
		if err != nil {
			return err
		}
	}

	creator := map[string]string{
		"email": org.Creator.Email,
		"id":    org.Creator.ID,
		"name":  org.Creator.Name,
	}
	err = d.Set("creator", []interface{}{creator})
	if err != nil {
		return err
	}

	creatorBillingOnly := len(org.BillingEmails) == 1 && org.BillingEmails[0] == org.Creator.Email
	if d.State().Attributes["billing_email"] == "" && creatorBillingOnly {
		err = d.Set("billing_emails", nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceCurrentOrgUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*api.Client)
	org := &api.Organization{
		ID: d.Id(),
	}

	if d.HasChange("name") {
		org.Name = d.Get("name").(string)
	}
	if d.HasChange("billing_emails") {
		org.BillingEmails = readBillingEmails(d)
	}
	if d.HasChange("auto_upgrade") {
		autoUpgrade := d.Get("auto_upgrade").(bool)
		org.AutoUpgrade = &autoUpgrade
	}
	_, err := c.UpdateOrganization(org)
	if err != nil {
		return err
	}
	return resourceCurrentOrgRead(d, m)
}

func resourceCurrentOrgDelete(d *schema.ResourceData, m interface{}) error {
	return errors.New("current organization can not be deleted\n. Please remove from state by running `terraform state rm bugsnag_current_org.<resource_name>`")
}

func readBillingEmails(d *schema.ResourceData) []string {
	emailsSet := d.Get("billing_emails").(*schema.Set)
	emails := make([]string, 0)
	for _, val := range emailsSet.List() {
		emails = append(emails, val.(string))
	}
	return emails
}
