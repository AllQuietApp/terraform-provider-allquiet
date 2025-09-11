resource "allquiet_team" "root" {
  display_name = "Root"
}

resource "allquiet_team" "engineering" {
  display_name = "Engineering"
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

resource "allquiet_outbound_integration" "slack_selected_teams" {
  display_name               = "My Slack Integration (Selected Teams)"
  team_id                    = allquiet_team.root.id
  type                       = "Slack"
  triggers_only_on_forwarded = false
  team_connection_settings = {
    team_connection_mode = "SelectedTeams"
    team_ids             = [allquiet_team.root.id, allquiet_team.engineering.id]
  }
}

resource "allquiet_outbound_integration" "webhook" {
  display_name                   = "My Webhook Integration"
  team_id                        = allquiet_team.root.id
  type                           = "Webhook"
  triggers_only_on_forwarded     = true
  skip_updating_after_forwarding = true
}