resource "allquiet_team" "root" {
  display_name = "Root"
}

resource "allquiet_outbound_integration" "slack" {
  display_name = "My Slack Integration"
  team_id      = allquiet_team.root.id
  type         = "Slack"
}