resource "allquiet_team" "root" {
  display_name = "Root"
}

resource "allquiet_integration" "datadog" {
  display_name = "My Datadog Integration"
  team_id      = allquiet_team.root.id
  type         = "Datadog"
}