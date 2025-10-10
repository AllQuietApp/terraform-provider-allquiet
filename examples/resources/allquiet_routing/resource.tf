resource "allquiet_team" "root" {
  display_name = "Root"
}

resource "allquiet_team" "infrastructure" {
  display_name = "Infrastructure"
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

resource "allquiet_outbound_integration" "notion" {
  team_id      = allquiet_team.root.id
  display_name = "Notion"
  type         = "Notion"
}

resource "allquiet_service" "pre_sales" {
  display_name = "Pre Sales"
  public_title = "Pre Sales"
}

resource "allquiet_routing" "example_1" {
  team_id      = allquiet_team.root.id
  display_name = "Route to specific team based on attribute"
  team_connection_settings = {
    team_connection_mode = "SelectedTeams"
    team_ids             = [allquiet_team.root.id, allquiet_team.infrastructure.id]
  }
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
        assign_to_teams = [allquiet_team.pres_sales.id]
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
        assign_to_teams = [allquiet_team.after_sales.id]
      }
    }
  ]
}

resource "allquiet_routing" "example_2" {
  team_id      = allquiet_team.root.id
  display_name = "Mute Slack Outbound Integration when Minor"
  team_connection_settings = {
    team_connection_mode = "OrganizationTeams"
  }
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

resource "allquiet_routing" "example_4" {
  team_id      = allquiet_team.root.id
  display_name = "Auto Resolve non Critical incidents from Custom Webhook after 10 minutes"
  rules = [
    {
      conditions = {
        severities   = ["Minor", "Warning"]
        integrations = [allquiet_integration.custom_webhook.id]
      },
      actions = {
        add_interaction          = "Resolved"
        delay_actions_in_minutes = 10
      }
    },
  ]
}

resource "allquiet_routing" "example_5" {
  team_id      = allquiet_team.root.id
  display_name = "Auto Resolve non Critical incidents from Custom Webhook but only in the period from 2025-01-16T22:07:03Z to 2025-01-19T00:00:00Z"
  rules = [
    {
      conditions = {
        severities = ["Minor", "Warning"],
        date_restriction = {
          from  = "2025-01-16T22:07:03Z",
          until = "2025-01-19T00:00:00Z"
        }
      },
      actions = {
        add_interaction = "Resolved"
      }
    },
  ]
}

resource "allquiet_routing" "example_6" {
  team_id      = allquiet_team.root.id
  display_name = "Auto Archive incidents Monday Mornings from 08 until 10"
  rules = [
    {
      conditions = {
        schedule = {
          days_of_week = ["mon"],
          after        = "08:00",
          before       = "10:00"
        }
      },
      actions = {
        add_interaction = "Archived"
      }
    },
  ]
}

resource "allquiet_routing" "example_7" {
  team_id      = allquiet_team.root.id
  display_name = "Auto Forward Critical incidents to Notion"
  rules = [
    {
      conditions = {
        severities = ["Critical"]
      },
      actions = {
        add_interaction                  = "Forwarded"
        forward_to_outbound_integrations = [allquiet_outbound_integration.notion.id]
      }
    }
  ]
}

resource "allquiet_routing" "example_8" {
  team_id      = allquiet_team.root.id
  display_name = "Auto Affect incidents with 'Project' attribute 'Pre Sales' to service 'Pre Sales'"
  rules = [
    {
      conditions = {
        attributes = [
          {
            name     = "Project",
            operator = "=",
            value    = "Pre Sales"
          }
        ]
      },
      actions = {
        add_interaction  = "Affects"
        affects_services = [allquiet_service.pre_sales.id]
      }
    }
  ]
}

resource "allquiet_routing" "example_9" {
  team_id      = allquiet_team.root.id
  display_name = "Add Attributes to incidents"
  rules = [
    {
      conditions = {
        attributes_match_type = "any"
        attributes = [
          {
            name     = "Project",
            operator = "=",
            value    = "Pre Sales"
          },
          {
            name     = "Project",
            operator = "=",
            value    = "After Sales"
          }
        ]
      },
      actions = {
        set_attributes = [
          {
            name             = "Team",
            value            = "Leads",
            hide_in_previews = true
          }
        ]
      }
    }
  ]
}

resource "allquiet_routing" "example_10" {
  team_id      = allquiet_team.root.id
  display_name = "Snooze Minor incidents for 10 minutes"
  team_connection_settings = {
    team_connection_mode = "OrganizationTeams"
  }
  rules = [
    {
      conditions = {
        severities = ["Minor", "Warning"]
      },
      actions = {
        add_interaction                = "Snoozed"
        snooze_for_relative_in_minutes = 10
      }
    },
  ]
}



resource "allquiet_routing" "example_11" {
  team_id      = allquiet_team.root.id
  display_name = "Snooze Minor incidents until Monday 07:00"
  team_connection_settings = {
    team_connection_mode = "OrganizationTeams"
  }
  rules = [
    {
      conditions = {
        severities = ["Minor", "Warning"]
      },
      actions = {
        add_interaction               = "Snoozed"
        snooze_until_absolute         = "07:00"
        snooze_until_weekday_absolute = "mon"
      }
    },
  ]
}