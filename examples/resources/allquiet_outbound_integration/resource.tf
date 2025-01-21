resource "allquiet_team" "root" {
  display_name = "Root"
}

resource "allquiet_outbound_integration" "slack" {
  display_name = "My Slack Integration"
  team_id      = allquiet_team.root.id
  type         = "Slack"
}

resource "allquiet_outbound_integration" "notion" {
  display_name = "My Notion Integration"
  team_id      = allquiet_team.root.id
  type         = "Notion"
}

resource "allquiet_outbound_integration" "notion_directly" {
  display_name               = "My Notion Integration (Directly)"
  team_id                    = allquiet_team.root.id
  type                       = "Notion"
  triggers_only_on_forwarded = false
}