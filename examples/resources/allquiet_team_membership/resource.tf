
resource "allquiet_team" "my_team" {
  display_name = "My Team"
  time_zone_id = "America/Los_Angeles"
}

resource "allquiet_user" "millie_brown" {
  display_name = "Millie Bobby Brown"
  email        = "millie@acme.com"
}

resource "allquiet_user" "taylor" {
  display_name = "Taylor Swift"
  email        = "taylor@acme.com"
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