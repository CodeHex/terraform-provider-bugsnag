package main

import (
	"github.com/codehex/terraform-provider-bugsnag/bugsnag"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return bugsnag.Provider()
		},
	})
}
