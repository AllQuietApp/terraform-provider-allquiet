// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIntegrationResourceConfig("Integration One"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_integration.test", "display_name", "Integration One"),
					resource.TestCheckResourceAttrSet("allquiet_integration.test", "webhook_url"),

					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.selected_days.0", "mon"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.selected_days.1", "tue"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.selected_days.2", "wed"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.selected_days.3", "thu"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.selected_days.4", "fri"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.from", "22:00"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.until", "07:00"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.snooze_until_absolute", "07:00"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.1.selected_days.0", "sat"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.1.selected_days.1", "sun"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.1.snooze_window_in_minutes", "10"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.url", "https://example.com"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.method", "GET"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.timeout_in_milliseconds", "1000"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.interval_in_seconds", "60"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.authentication_type", "Bearer"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.bearer_authentication_token", "my-token"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.headers.Content-Type", "application/json"),
					resource.TestCheckResourceAttr("allquiet_integration.heartbeat_monitor", "integration_settings.heartbeat_monitor.interval_in_sec", "60"),
					resource.TestCheckResourceAttr("allquiet_integration.heartbeat_monitor", "integration_settings.heartbeat_monitor.grace_period_in_sec", "10"),
					resource.TestCheckResourceAttr("allquiet_integration.heartbeat_monitor", "integration_settings.heartbeat_monitor.severity", "Warning"),
					resource.TestCheckResourceAttr("allquiet_integration.cronjob_monitor", "integration_settings.cronjob_monitor.cron_expression", "0 0 * * *"),
					resource.TestCheckResourceAttr("allquiet_integration.cronjob_monitor", "integration_settings.cronjob_monitor.grace_period_in_sec", "10"),
					resource.TestCheckResourceAttr("allquiet_integration.cronjob_monitor", "integration_settings.cronjob_monitor.severity", "Critical"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIntegrationResourceConfig("Integration Two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_integration.test", "display_name", "Integration Two"),
					resource.TestCheckResourceAttrSet("allquiet_integration.test", "webhook_url"),

					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.selected_days.0", "mon"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.selected_days.1", "tue"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.selected_days.2", "wed"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.selected_days.3", "thu"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.selected_days.4", "fri"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.from", "22:00"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.until", "07:00"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.0.snooze_until_absolute", "07:00"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.1.selected_days.0", "sat"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.1.selected_days.1", "sun"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "snooze_settings.filters.1.snooze_window_in_minutes", "10"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "webhook_authentication.type", "bearer"),
					resource.TestCheckResourceAttr("allquiet_integration.webhook_snooze_absolute", "webhook_authentication.bearer.token", "my-token"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.url", "https://example.com"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.method", "GET"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.timeout_in_milliseconds", "1000"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.interval_in_seconds", "60"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.authentication_type", "Bearer"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.bearer_authentication_token", "my-token"),
					resource.TestCheckResourceAttr("allquiet_integration.heartbeat_monitor", "integration_settings.heartbeat_monitor.interval_in_sec", "60"),
					resource.TestCheckResourceAttr("allquiet_integration.heartbeat_monitor", "integration_settings.heartbeat_monitor.grace_period_in_sec", "10"),
					resource.TestCheckResourceAttr("allquiet_integration.heartbeat_monitor", "integration_settings.heartbeat_monitor.severity", "Warning"),
					resource.TestCheckResourceAttr("allquiet_integration.cronjob_monitor", "integration_settings.cronjob_monitor.cron_expression", "0 0 * * *"),
					resource.TestCheckResourceAttr("allquiet_integration.cronjob_monitor", "integration_settings.cronjob_monitor.grace_period_in_sec", "10"),
					resource.TestCheckResourceAttr("allquiet_integration.cronjob_monitor", "integration_settings.cronjob_monitor.severity", "Critical"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIntegrationResourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIntegrationResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_integration.datadog", "display_name", "My Datadog Integration"),
					resource.TestCheckResourceAttr("allquiet_integration.amazon_cloudwatch", "display_name", "My Amazon CloudWatch Integration"),
					resource.TestCheckResourceAttrSet("allquiet_integration.amazon_cloudwatch", "webhook_url"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_integration.datadog",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIntegrationResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_integration.datadog", "display_name", "My Datadog Integration"),
					resource.TestCheckResourceAttr("allquiet_integration.amazon_cloudwatch", "display_name", "My Amazon CloudWatch Integration"),
					resource.TestCheckResourceAttrSet("allquiet_integration.amazon_cloudwatch", "webhook_url"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIntegrationResourceConfig(display_name string) string {
	return fmt.Sprintf(`
resource "allquiet_team" "test" {
  display_name = "Root"
}

resource "allquiet_integration" "test" {
  display_name = %[1]q
  team_id = allquiet_team.test.id
  type = "Datadog"
}

resource "allquiet_integration" "webhook_snooze_absolute" {
	display_name = "My Webhook Integration"
	team_id = allquiet_team.test.id
	type = "Webhook"
	snooze_settings = {
		filters = [
			{
				selected_days = ["mon", "tue", "wed", "thu", "fri"]
				from = "22:00"
				until = "07:00"
				snooze_until_absolute = "07:00"
			},
			{
				selected_days = ["sat", "sun"]
				snooze_window_in_minutes = 10
			}
		]
	}
	webhook_authentication = {
		type = "bearer"
		bearer = {
			token = "my-token"
		}
	}
}

	
resource "allquiet_integration" "http_monitoring" {
	display_name = "My HTTP Monitoring Integration"
	team_id = allquiet_team.test.id
	type = "HttpMonitoring"
	integration_settings = {
		http_monitoring = {
			url = "https://example.com"
			method = "GET"
			timeout_in_milliseconds = 1000
			interval_in_seconds = 60
			authentication_type = "Bearer"
			bearer_authentication_token = "my-token"
			headers = {
				"Content-Type" = "application/json"
			}
			body = "{\"message\": \"Hello, world!\"}"
			is_paused = false
			content_test = "Hello, world!"
			ssl_certificate_max_age_in_days_degraded = 30
			ssl_certificate_max_age_in_days_down = 10
			severity_degraded = "Warning"
			severity_down = "Critical"
		}
	}
}
	
resource "allquiet_integration" "heartbeat_monitor" {
	display_name = "My Heartbeat Monitoring Integration"
	team_id = allquiet_team.test.id
	type = "HeartbeatMonitor"
	integration_settings = {
		heartbeat_monitor = {
			interval_in_sec = 60
			grace_period_in_sec = 10
			severity = "Warning"
		}
	}
}
	
resource "allquiet_integration" "cronjob_monitor" {
	display_name = "My Cronjob Monitoring Integration"
	team_id = allquiet_team.test.id
	type = "CronJobMonitor"
	integration_settings = {
		cronjob_monitor = {
			cron_expression = "0 0 * * *"
			grace_period_in_sec = 10
			severity = "Critical"
			time_zone_id = "Europe/Amsterdam"
		}
	}
}


`, display_name)

}

func testAccIntegrationResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_integration/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return RandomizeExample(string(dat))
}
