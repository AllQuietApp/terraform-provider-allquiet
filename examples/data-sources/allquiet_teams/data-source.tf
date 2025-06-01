data "allquiet_teams" "teams_by_display_name" {
  display_name = "(TF Acceptance Test) Team"
}

data "allquiet_teams" "all_teams" {
}

output "team_ids" {
  value = data.allquiet_teams.teams_by_display_name.teams[*].id
}

output "team_display_names" {
  value = data.allquiet_teams.teams_by_display_name.teams[*].display_name
}

output "team_time_zone_ids" {
  value = data.allquiet_teams.teams_by_display_name.teams[*].time_zone_id
}