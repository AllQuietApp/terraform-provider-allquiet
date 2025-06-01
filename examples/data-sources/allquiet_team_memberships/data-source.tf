# Setup
data "allquiet_team" "test" {
  display_name = "(TF Acceptance Test) Team"
}

data "allquiet_user" "test" {
  email = "acceptance-tests+millie@allquiet.app"
}

# Data sources
data "allquiet_team_memberships" "memberships_by_team" {
  team_id = data.allquiet_team.test.id
}

output "team_memberships_ids" {
  value = data.allquiet_team_memberships.memberships_by_team.team_memberships[*].id
}
