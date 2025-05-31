# Read a team by display name

data "allquiet_team" "test" {
  display_name = "(TF Acceptance Test) Team"
}

data "allquiet_user" "test" {
  email = "acceptance-tests+millie@allquiet.app"
}

data "allquiet_team_membership" "test" {
  team_id = data.allquiet_team.test.id
  user_id = data.allquiet_user.test.id
}