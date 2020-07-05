package bugsnag

import (
	"github.com/codehex/terraform-provider-bugsnag/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider creates the Bugsnag Provider for terraform
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
			"api_endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BUGSNAG_API_ENDPOINT", "https://api.bugsnag.com"),
			},
		},
		ConfigureFunc: func(data *schema.ResourceData) (interface{}, error) {
			return api.New(
				data.Get("auth_token").(string),
				data.Get("api_endpoint").(string),
			)
		},
	}
}
