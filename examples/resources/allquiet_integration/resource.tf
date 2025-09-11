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
  webhook_authentication = {
    type = "bearer"
    bearer = {
      token = "your_secret_token"
    }
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

resource "allquiet_integration" "webhook_snooze_absolute_with_weekday" {
  display_name = "My Webhook Integration"
  team_id      = allquiet_team.root.id
  type         = "Webhook"
  snooze_settings = {
    filters = [
      {
        selected_days                 = ["sat", "sun"]
        snooze_until_absolute         = "07:00"
        snooze_until_weekday_absolute = "mon"
      }
    ]
  }
}

resource "allquiet_integration" "heartbeat_monitor" {
  display_name = "My Heartbeat Monitoring Integration"
  team_id      = allquiet_team.root.id
  type         = "HeartbeatMonitor"
  integration_settings = {
    heartbeat_monitor = {
      interval_in_sec     = 60
      grace_period_in_sec = 10
      severity            = "Warning"
    }
  }
}

resource "allquiet_integration" "cronjob_monitor" {
  display_name = "My Cronjob Monitoring Integration"
  team_id      = allquiet_team.root.id
  type         = "CronJobMonitor"
  integration_settings = {
    cronjob_monitor = {
      cron_expression     = "0 0 * * *"
      grace_period_in_sec = 25
      severity            = "Critical"
      time_zone_id        = "Europe/Amsterdam"
    }
  }
}

resource "allquiet_integration" "http_monitoring_post" {
  display_name = "My HTTP Monitoring POST Integration"
  team_id      = allquiet_team.root.id
  type         = "HttpMonitoring"
  integration_settings = {
    http_monitoring = {
      url                         = "https://allquiet.com"
      method                      = "POST"
      timeout_in_milliseconds     = 1000 # 1 second
      interval_in_seconds         = 300  # 5 minutes
      authentication_type         = "Bearer"
      bearer_authentication_token = "your_secret_token"
      headers = {
        "Content-Type" = "application/json"
      }
      body                                     = "{\"message\": \"Hello, world!\"}"
      is_paused                                = false
      content_test                             = "Expected response text"
      ssl_certificate_max_age_in_days_degraded = 30
      ssl_certificate_max_age_in_days_down     = 10
      severity_degraded                        = "Warning"
      severity_down                            = "Critical"
    }
  }
}

resource "allquiet_integration" "http_monitoring_head" {
  display_name = "My HTTP Monitoring HEAD Integration"
  team_id      = allquiet_team.root.id
  type         = "HttpMonitoring"
  integration_settings = {
    http_monitoring = {
      url = "https://allquiet.com"

      method                                   = "HEAD"
      timeout_in_milliseconds                  = 4000 # 4 seconds
      interval_in_seconds                      = 300  # 5 minutes
      authentication_type                      = "Bearer"
      bearer_authentication_token              = "your_secret_token"
      is_paused                                = false
      ssl_certificate_max_age_in_days_degraded = 30
      ssl_certificate_max_age_in_days_down     = 10
      severity_degraded                        = "Warning"
      severity_down                            = "Critical"
    }
  }
}

resource "allquiet_integration" "ping_monitor" {
  display_name = "My Ping Monitoring Integration"
  team_id      = allquiet_team.root.id
  type         = "PingMonitor"
  integration_settings = {
    ping_monitor = {
      host = "google.com"

      timeout_in_milliseconds = 1000
      interval_in_seconds     = 300
      is_paused               = false
      severity_degraded       = "Warning"
      severity_down           = "Critical"
    }
  }
}


locals {
  computed_amazon_cloudwatch_webhook_url = allquiet_integration.amazon_cloudwatch.webhook_url
}