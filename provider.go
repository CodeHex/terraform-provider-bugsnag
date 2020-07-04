package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"bugsnag_current_org":  resourceCurrentOrg(),
			"bugsnag_project":      resourceProject(),
			"bugsnag_collaborator": resourceCollaborator(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bugsnag_projects": dataSourceProjects(),
		},
		Schema: map[string]*schema.Schema{
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BUGSNAG_DATA_ACCESS_TOKEN", nil),
			},
		},
		ConfigureFunc: func(data *schema.ResourceData) (interface{}, error) {
			return NewClient(data.Get("auth_token").(string))
		},
	}
}
