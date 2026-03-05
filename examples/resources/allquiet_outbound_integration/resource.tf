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

resource "allquiet_outbound_integration" "slack_with_channels" {
  display_name = "My Slack Integration (With Channels)"
  team_id      = allquiet_team.root.id
  type         = "Slack"

  slack_settings = {
    selected_channel_ids               = ["C1234567890", "C0987654321"]
    tag_on_call_members                = true
    is_slack_message_payload_read_only = false
  }
}

resource "allquiet_outbound_integration" "slack_with_severity_channels" {
  display_name = "My Slack Integration (Severity-Based Channels)"
  team_id      = allquiet_team.root.id
  type         = "Slack"

  slack_settings = {
    severity_based_channel_settings = {
      selected_channel_ids_minor    = ["C1111111111"]
      selected_channel_ids_warning  = ["C2222222222"]
      selected_channel_ids_critical = ["C3333333333"]
    }
  }
}

resource "allquiet_outbound_integration" "slack_with_reminders" {
  display_name = "My Slack Integration (With Reminders)"
  team_id      = allquiet_team.root.id
  type         = "Slack"

  slack_settings = {
    on_call_reminder_schedule_settings = {
      run_time     = "09:00"
      days_of_week = ["mon", "wed", "fri"]
    }
    on_call_reminder_channel_ids = ["C4444444444"]
    tag_on_call_members          = true
  }
}

resource "allquiet_outbound_integration" "slack_full_config" {
  display_name = "My Slack Integration (Full Config)"
  team_id      = allquiet_team.root.id
  type         = "Slack"

  slack_settings = {
    selected_channel_ids = ["C5555555555"]
    on_call_reminder_schedule_settings = {
      run_time     = "08:00"
      days_of_week = ["mon", "tue", "wed", "thu", "fri"]
    }
    on_call_reminder_channel_ids       = ["C6666666666"]
    tag_on_call_members                = true
    is_slack_message_payload_read_only = false
  }
}

resource "allquiet_outbound_integration" "mattermost" {
  display_name = "My Mattermost Integration"
  team_id      = allquiet_team.root.id
  type         = "Mattermost"

  mattermost_settings = {
    send_incidents_to_mattermost     = true
    create_incidents_from_mattermost = false
    base_url                         = "https://mattermost.com"
    selected_team_id                 = "your-team-id"
    selected_channel_ids             = ["channel-id-1", "channel-id-2"]
    is_message_read_only             = false
  }
}

resource "allquiet_outbound_integration" "mattermost_severity_channels" {
  display_name = "My Mattermost Integration (Severity-Based Channels)"
  team_id      = allquiet_team.root.id
  type         = "Mattermost"

  mattermost_settings = {
    send_incidents_to_mattermost     = true
    create_incidents_from_mattermost = false
    base_url                         = "https://mattermost.com"
    selected_team_id                 = "your-team-id"
    severity_based_channel_settings = {
      selected_channel_ids_minor    = ["minor-channel-id"]
      selected_channel_ids_warning  = ["warning-channel-id"]
      selected_channel_ids_critical = ["critical-channel-id"]
    }
  }
}