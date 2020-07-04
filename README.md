# Terraform Provider For Bugsnag

![build](https://github.com/CodeHex/terraform-provider-bugsnag/workflows/build/badge.svg?branch=master)
[![Go Report](https://goreportcard.com/badge/codehex/terraform-provider-bugsnag)](https://goreportcard.com/badge/codehex/terraform-provider-bugsnag)

Provides management of organizations, projects and collaborators using personal access tokens.

Requires [Terraform 0.12.x](https://www.terraform.io/downloads.html)

## Installing

---
## Bugsnag Provider
The Bugsnag provider is used to interact with Bugsnag organization related resources.

The provider needs to be configured with a [Personal Auth Token](https://bugsnagapiv2.docs.apiary.io/introduction/authentication/personal-auth-tokens-(recommended)), which will be restricted to a single organization. Only settings, projects and collaborators that the token can access will be accessible.

[Multiple provider instances](https://www.terraform.io/docs/configuration/providers.html#alias-multiple-provider-instances) should be used for managing multiple organizations.

#### Example Usage
 ```HCL
provider "bugsnag" {
  auth_token = "${var.bugsnag_auth_token}"
}
 ```
#### Argument Reference
- `auth_token` - (Required) This is the Bugsnag personal auth token. It can also be provided from the `BUGSNAG_DATA_ACCESS_TOKEN` environment variable.
---

## Data Sources
### bugsnag_projects
Use this data source to provide a list of projects that can be accessed, using an optional search query

#### Example Usage
 ```HCL
data "bugsnag_projects" "example" {
  query     = "Website"
  sort      = "name"
  direction = "asc"
}
 ```
#### Arguments Reference
- `query` - (Optional) Filter the projects with names matching this query value. Defaults to no filtering.
- `sort` - (Optional) Sort the projects by either `created_at`, `name` or `favorite`. Defaults to `created_at`
- `direction` - (Optional) Which direction to sort the projects, either `asc` or `desc`. Defaults to `asc` when sorting by favorite, `desc` otherwise.

#### Attributes Reference
- `id` - The ID of the project (e.g. `5ed73feb00002b509b040000`)
- `name` - The name of the project (e.g. `Website Frontend`)
- `api_key` - The API key for the project (e.g. `86925ec410def315d6d2bffd91f51da1`)

---

## Resources
### bugsnag_current_org
#### Example Usage
#### Arguments Reference
#### Attributes Reference

---

### bugsnag_project
#### Example Usage
#### Arguments Reference
#### Attributes Reference

---

### bugsnag_collaborator
#### Example Usage
#### Arguments Reference
#### Attributes Reference

---

## Developing the provider

The provider requires [Go +1.14](https://golang.org/doc/install) to build
