# Terraform Provider For Bugsnag

![build](https://github.com/CodeHex/terraform-provider-bugsnag/workflows/build/badge.svg?branch=master)
[![Go Report](https://goreportcard.com/badge/codehex/terraform-provider-bugsnag)](https://goreportcard.com/badge/codehex/terraform-provider-bugsnag)

Provides management of organizations, projects and collaborators using personal access tokens.

Requires [Terraform 0.12.x](https://www.terraform.io/downloads.html)

## Installing

Coming Soon...

## Documentation

- [Bugsnag Provider](#bugsnag-provider)
- [Data Sources](#data-sources)
  - [bugsnag_projects](#bugsnag_projects)
- [Resources](#resources)
  - [bugsnag_current_org](#bugsnag_current_org)
  - [bugsnag_project](#bugsnag_project)
  - [bugsnag_collaborator](#bugsnag_collaborator)

For examples... (comming soon)

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
This resource allows the organization associated with the access token can be updated. There should only be **one** resource of this type.

**NOTE that this resource cannot be created or deleted.** 

To manage this resource, it should be imported
```bash
terraform import bugsnag_current_org.<resource_name> ""
```
and removed from the state if its no longer being managed using
```bash
terraform state rm bugsnag_current_org.<resource_name>
```

#### Example Usage
 ```HCL
resource "bugsnag_current_org" "example" {
  name = "Example Org"
}
 ```
#### Arguments Reference
- `name` (Optional) - The name of the organization
- `billing_emails` (Optional) - A list of email addresses to send invoice information to. Setting this as empty indicates the only billing email should be the creators email.
- `auto_upgrade` (Optional) - Whether the organizations plan is upgraded in response to the organization reaching its plan limit of events. If this value is false your events will be throttled when you reach your plan limit.

#### Attributes Reference
- `slug` - The organizations slug used in the URLs (e.g. `example-org`)
- `creator` - Contains details about the creator of the organization
  - `email` - The email address of the organizations creator (e.g `richard.jones@example.com`)
  - `id` - - The ID of the organizations creator  (e.g. `86925ec410def315d6d2bffd91f51da1`)
  - `name` - The name of the organizations creator  (e.g. `Richard Jones`)
- `created_at` - The date and time the organization was created (e.g. `2017-04-24T22:17:13.000Z`)
- `updated_at` - The date and time the organization was last updated at (e.g. `2017-04-24T22:17:13.000Z`)

---

### bugsnag_project
This resource allows you to create and manage projects within your Bugsnag organizations.

#### Example Usage
 ```HCL
resource "bugsnag_project" "example" {
  name = "Example iOS App"
  type = "ios"
}
 ```

#### Arguments Reference
- `name` (Required) - The name of the project
- `type` (Required) - The language and/or framework of the project. See [Creating a Project](https://bugsnagapiv2.docs.apiary.io/#reference/projects/projects/create-a-project-in-an-organization) in the Bugsnag data access API under `type` for possible valid values.

#### Attributes Reference
- `api_key` - The API key for the project (e.g. `86925ec410def315d6d2bffd91f51da1`)
- `slug` - The projects slug used in the URLs (e.g. `example-ios-app`)
- `html_url` - The URL to the Bugsnag dashboard for this project (e.g. `https://app.bugsnag.com/example-org/example-ios-app`)
- `created_at` - The date and time the project was created (e.g. `2017-04-24T22:17:13.000Z`)
- `updated_at` - The date and time the project was last updated at (e.g. `2017-04-24T22:17:13.000Z`)

---

### bugsnag_collaborator
This resource allows you to invite and manage collaborators.
#### Example Usage
 ```HCL
resource "bugsnag_project" "example_ios_app" {
  name = "Example iOS App"
  type = "ios"
}

resource "bugsnag_project" "example_android_app" {
  name = "Example Android App"
  type = "android"
}

resource "bugsnag_collaborator" "example_non_admin" {
   email = "richard.jones@example.com"
   project_ids = [
     bugsnag_project.example_ios_app.id,
     bugsnag_project.example_android_app.id
   ]
}

resource "bugsnag_collaborator" "example_admin" {
   email = "lucy.simmons@example.com"
   admin = true
 ```
#### Arguments Reference
- `email` (Required) - The email address of the collaborator
- `admin` (Optional) - Set to `true` if the collaborator should have admin access to the organization. Defaults to `false`. The `project_ids` field must be blank if admin is set to `true`
- `project_ids` (Optional) - A list of project IDs the collaborator should have access to. Should only be provided when admin is set to `false`
#### Attributes Reference
- `name` - The name of the collaborator
- `two_factor_enabled` - Set to `true` if 2FA has is enabled for this user
- `two_factor_enabled_on` - The date and time 2FA was enabled for this user, if it is active (e.g. `2017-04-24T22:17:13.000Z`)
- `recovery_codes_remaining` - The number of two factor recovery codes the User has left. (e.g. `2`)
- `password_updated_on` - The date and time the users password was last changed (e.g. `2017-04-24T22:17:13.000Z`)
- `show_time_in_utc` - Set to `true` if the user has opted to have times displayed in UTC instead of local time
- `created_at` - The date and time the project was created (e.g. `2017-04-24T22:17:13.000Z`)
- `pending_invitation` - Set to `true` if the user has not yet accepted their Bugsnag invitation
- `last_request_at`- The last time the user interacted with the Bugsnag dashboard or related APIs. This is not set if the user has not interacted with the dashboard before (e.g. `2017-04-24T22:17:13.000Z)
- `paid_for` - Set to `true` if the user has Bugsnag dashboard access under the organization's current plan. If this is `false` for a collaborator, they will see a "locked out" message when they attempt to log in to the Bugsnag dashboard

---
