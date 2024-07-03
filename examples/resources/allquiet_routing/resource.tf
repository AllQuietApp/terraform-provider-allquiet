resource "allquiet_team" "root" {
  display_name = "Root"
}

resource "allquiet_team" "pres_sales" {
  display_name = "Pre Sales"
}

resource "allquiet_team" "after_sales" {
  display_name = "After Sales"
}

resource "allquiet_integration" "datadog" {
  team_id      = allquiet_team.root.id
  display_name = "Datadog"
  type         = "Datadog"
}

resource "allquiet_integration" "custom_webhook" {
  team_id      = allquiet_team.root.id
  display_name = "Custom Webhook"
  type         = "Webhook"
}


resource "allquiet_outbound_integration" "slack" {
  team_id      = allquiet_team.root.id
  display_name = "Slack"
  type         = "Slack"
}

resource "allquiet_routing" "example_1" {
  team_id      = allquiet_team.root.id
  display_name = "Route to specific team based on attribute"
  rules = [
    {
      conditions = {
        attributes = [
          {
            name     = "Service",
            operator = "=",
            value    = "Pre Sales"
          }
        ]
      },
      actions = {
        route_to_teams = [allquiet_team.pres_sales.id]
      }
    },

    {
      conditions = {
        attributes = [
          {
            name     = "Service",
            operator = "=",
            value    = "After Sales"
          }
        ]
      },
      actions = {
        route_to_teams = [allquiet_team.after_sales.id]
      }
    }
  ]
}

resource "allquiet_routing" "example_2" {
  team_id      = allquiet_team.root.id
  display_name = "Mute Slack Outbound Integration when Minor"
  rules = [
    {
      conditions = {
        severities = ["Minor"]
      },
      channels = {
        outbound_integrations       = [allquiet_outbound_integration.slack.id]
        outbound_integrations_muted = true
      },
    },
  ]
}

resource "allquiet_routing" "example_3" {
  team_id      = allquiet_team.root.id
  display_name = "Auto Resolve non Critical incidents from Custom Webhook"
  rules = [
    {
      conditions = {
        severities   = ["Minor", "Warning"]
        integrations = [allquiet_integration.custom_webhook.id]
      },
      actions = {
        add_interaction = "Resolved"
      }
    },
  ]
}