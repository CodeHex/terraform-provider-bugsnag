package bugsnag

import (
	"github.com/codehex/terraform-provider-bugsnag/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectsRead,

		Schema: map[string]*schema.Schema{
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":      {Type: schema.TypeString, Computed: true},
						"name":    {Type: schema.TypeString, Computed: true},
						"api_key": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceProjectsRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*api.Client)

	apiProjects, err := c.ListProjects(d.Get("query").(string))
	if err != nil {
		return err
	}
	d.SetId(c.OrgID)
	projects := make([]interface{}, 0)
	for _, apiProject := range apiProjects {
		project := make(map[string]interface{})
		project["id"] = apiProject.ID
		project["name"] = apiProject.Name
		project["api_key"] = apiProject.APIKey
		projects = append(projects, project)
	}

	d.Set("project", projects)
	return nil
}
