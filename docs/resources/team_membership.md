---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "allquiet_team_membership Resource - allquiet"
subcategory: ""
description: |-
  The user resource represents a user in All Quiet. Users can be members of users and receive notifications for incidents.
---

# allquiet_team_membership (Resource)

The user resource represents a user in All Quiet. Users can be members of users and receive notifications for incidents.

## Example Usage

```terraform
resource "allquiet_team" "my_team" {
  display_name = "My Team"
  time_zone_id = "America/Los_Angeles"
}

resource "allquiet_user" "millie_brown" {
  display_name = "Millie Bobby Brown"
  email        = "acceptance-tests+millie@allquiet.app"
}

resource "allquiet_user" "taylor" {
  display_name = "Taylor Swift"
  email        = "acceptance-tests+taylor@allquiet.app"
}

resource "allquiet_team_membership" "my_team_millie_brown" {
  team_id = allquiet_team.my_team.id
  user_id = allquiet_user.millie_brown.id
  role    = "Administrator"
}

resource "allquiet_team_membership" "my_team_taylor" {
  team_id = allquiet_team.my_team.id
  user_id = allquiet_user.taylor.id
  role    = "Member"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `role` (String) Role of the member. Possible values are: Member, Administrator
- `team_id` (String) The team id that the user is a member of
- `user_id` (String) The user id of the user

### Read-Only

- `id` (String) Id
