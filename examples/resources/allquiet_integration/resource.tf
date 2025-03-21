resource "allquiet_team" "root" {
  display_name = "Root"
}

resource "allquiet_integration" "datadog" {
  display_name = "My Datadog Integration"
  team_id      = allquiet_team.root.id
  type         = "Datadog"
}

resource "allquiet_integration" "amazon_cloudwatch" {
  display_name = "My Amazon CloudWatch Integration"
  team_id      = allquiet_team.root.id
  type         = "AmazonCloudWatch"
}

resource "allquiet_integration" "webhook" {
  display_name = "My Webhook Integration"
  team_id      = allquiet_team.root.id
  type         = "Webhook"
  snooze_settings = {
    snooze_window_in_minutes = 1440
  }
}

resource "allquiet_integration" "webhook_snooze_absolute" {
  display_name = "My Webhook Integration"
  team_id      = allquiet_team.root.id
  type         = "Webhook"
  snooze_settings = {
    filters = [
      {
        selected_days         = ["mon", "tue", "wed", "thu", "fri"]
        from                  = "22:00"
        until                 = "07:00"
        snooze_until_absolute = "07:00"
      }
    ]
  }
}


locals {
  computed_amazon_cloudwatch_webhook_url = allquiet_integration.amazon_cloudwatch.webhook_url
}