---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "allquiet_user Data Source - allquiet"
subcategory: ""
description: |-
  User data source
---

# allquiet_user (Data Source)

User data source

## Example Usage

```terraform
# Read a user by email
data "allquiet_user" "test_by_email" {
  email = "acceptance-tests+millie@allquiet.app"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `display_name` (String) Display name of the user to look up
- `email` (String) Email address of the user to look up
- `scim_external_id` (String) If the user was provisioned by SCIM, this is the SCIM external ID of the user to look up

### Read-Only

- `id` (String) User ID
