// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
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
					resource.TestCheckResourceAttr("allquiet_integration.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("allquiet_integration.test", "labels.0", "prod"),
					resource.TestCheckResourceAttr("allquiet_integration.test", "labels.1", "monitoring"),
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
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.ignore_non_http_errors", "true"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.content_test", "Hello, world!"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.content_test_mode", "Contains"),
					resource.TestCheckResourceAttr("allquiet_integration.heartbeat_monitor", "integration_settings.heartbeat_monitor.interval_in_sec", "60"),
					resource.TestCheckResourceAttr("allquiet_integration.heartbeat_monitor", "integration_settings.heartbeat_monitor.grace_period_in_sec", "10"),
					resource.TestCheckResourceAttr("allquiet_integration.heartbeat_monitor", "integration_settings.heartbeat_monitor.severity", "Warning"),
					resource.TestCheckResourceAttr("allquiet_integration.cronjob_monitor", "integration_settings.cronjob_monitor.cron_expression", "0 0 * * *"),
					resource.TestCheckResourceAttr("allquiet_integration.cronjob_monitor", "integration_settings.cronjob_monitor.grace_period_in_sec", "10"),
					resource.TestCheckResourceAttr("allquiet_integration.cronjob_monitor", "integration_settings.cronjob_monitor.severity", "Critical"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring_with_max_retries", "integration_settings.http_monitoring.max_retries", "0"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring_with_accepted_status_codes", "integration_settings.http_monitoring.url", "https://example.com/health"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring_with_accepted_status_codes", "integration_settings.http_monitoring.override_accepted_status_codes.#", "4"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring_with_accepted_status_codes", "integration_settings.http_monitoring.override_accepted_status_codes.0", "200"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring_with_accepted_status_codes", "integration_settings.http_monitoring.override_accepted_status_codes.1", "201"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring_with_accepted_status_codes", "integration_settings.http_monitoring.override_accepted_status_codes.2", "204"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring_with_accepted_status_codes", "integration_settings.http_monitoring.override_accepted_status_codes.3", "302"),
					resource.TestCheckResourceAttr("allquiet_integration.ping_monitor_with_max_retries", "integration_settings.ping_monitor.max_retries", "5"),
					resource.TestCheckResourceAttr("allquiet_integration.email", "display_name", "My Email Integration"),
					resource.TestCheckResourceAttr("allquiet_integration.email", "type", "Email"),
					resource.TestCheckResourceAttrSet("allquiet_integration.email", "integration_settings.email.email_address"),
					resource.TestCheckResourceAttr("allquiet_integration.email_with_aliases", "display_name", "My Email Integration"),
					resource.TestCheckResourceAttr("allquiet_integration.email_with_aliases", "type", "Email"),
					resource.TestCheckResourceAttr("allquiet_integration.email_with_aliases", "integration_settings.email.aliases.#", "2"),
					resource.TestCheckResourceAttrSet("allquiet_integration.email_with_aliases", "integration_settings.email.email_address"),
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
					resource.TestCheckResourceAttr("allquiet_integration.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("allquiet_integration.test", "labels.0", "prod"),
					resource.TestCheckResourceAttr("allquiet_integration.test", "labels.1", "monitoring"),
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
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.ignore_non_http_errors", "true"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.content_test", "Hello, world!"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring", "integration_settings.http_monitoring.content_test_mode", "Contains"),
					resource.TestCheckResourceAttr("allquiet_integration.heartbeat_monitor", "integration_settings.heartbeat_monitor.interval_in_sec", "60"),
					resource.TestCheckResourceAttr("allquiet_integration.heartbeat_monitor", "integration_settings.heartbeat_monitor.grace_period_in_sec", "10"),
					resource.TestCheckResourceAttr("allquiet_integration.heartbeat_monitor", "integration_settings.heartbeat_monitor.severity", "Warning"),
					resource.TestCheckResourceAttr("allquiet_integration.cronjob_monitor", "integration_settings.cronjob_monitor.cron_expression", "0 0 * * *"),
					resource.TestCheckResourceAttr("allquiet_integration.cronjob_monitor", "integration_settings.cronjob_monitor.grace_period_in_sec", "10"),
					resource.TestCheckResourceAttr("allquiet_integration.cronjob_monitor", "integration_settings.cronjob_monitor.severity", "Critical"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring_with_max_retries", "integration_settings.http_monitoring.max_retries", "0"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring_with_accepted_status_codes", "integration_settings.http_monitoring.override_accepted_status_codes.#", "4"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring_with_accepted_status_codes", "integration_settings.http_monitoring.override_accepted_status_codes.0", "200"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring_with_accepted_status_codes", "integration_settings.http_monitoring.override_accepted_status_codes.1", "201"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring_with_accepted_status_codes", "integration_settings.http_monitoring.override_accepted_status_codes.2", "204"),
					resource.TestCheckResourceAttr("allquiet_integration.http_monitoring_with_accepted_status_codes", "integration_settings.http_monitoring.override_accepted_status_codes.3", "302"),
					resource.TestCheckResourceAttr("allquiet_integration.ping_monitor_with_max_retries", "integration_settings.ping_monitor.max_retries", "5"),
					resource.TestCheckResourceAttr("allquiet_integration.email", "display_name", "My Email Integration"),
					resource.TestCheckResourceAttr("allquiet_integration.email", "type", "Email"),
					resource.TestCheckResourceAttrSet("allquiet_integration.email", "integration_settings.email.email_address"),
					resource.TestCheckResourceAttr("allquiet_integration.email_with_aliases", "display_name", "My Email Integration"),
					resource.TestCheckResourceAttr("allquiet_integration.email_with_aliases", "type", "Email"),
					resource.TestCheckResourceAttr("allquiet_integration.email_with_aliases", "integration_settings.email.aliases.#", "2"),
					resource.TestCheckResourceAttrSet("allquiet_integration.email_with_aliases", "integration_settings.email.email_address"),
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
	result := fmt.Sprintf(`
resource "allquiet_team" "test" {
  display_name = "Root"
}

resource "allquiet_integration" "test" {
  display_name = %[1]q
  team_id = allquiet_team.test.id
  type = "Datadog"
  labels = ["prod", "monitoring"]
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

resource "allquiet_integration" "webhook_snooze_absolute_with_weekday" {
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
				snooze_until_weekday_absolute = "tue"
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
			content_test_mode = "Contains"
			ssl_certificate_max_age_in_days_degraded = 30
			ssl_certificate_max_age_in_days_down = 10
			severity_degraded = "Warning"
			severity_down = "Critical"
			ignore_non_http_errors = true
		}
	}
}

resource "allquiet_integration" "http_monitoring_with_max_retries" {
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
			max_retries = 0
		}
	}
}

resource "allquiet_integration" "http_monitoring_with_accepted_status_codes" {
	display_name = "My HTTP Monitoring Integration With Custom Status Codes"
	team_id = allquiet_team.test.id
	type = "HttpMonitoring"
	integration_settings = {
		http_monitoring = {
			url = "https://example.com/health"
			method = "GET"
			timeout_in_milliseconds = 1000
			interval_in_seconds = 60
			authentication_type = "None"
			override_accepted_status_codes = [200, 201, 204, 302]
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

resource "allquiet_integration" "ping_monitor" {
	display_name = "My Ping Monitoring Integration"
	team_id = allquiet_team.test.id
	type = "PingMonitor"
	integration_settings = {
		ping_monitor = {
			host = "google.com"
			timeout_in_milliseconds = 1000
			interval_in_seconds = 60
			is_paused = false
			severity_degraded = "Warning"
			severity_down = "Critical"
		}
	}
}

resource "allquiet_integration" "ping_monitor_with_max_retries" {
	display_name = "My Ping Monitoring Integration"
	team_id = allquiet_team.test.id
	type = "PingMonitor"
	integration_settings = {
		ping_monitor = {
			host = "google.com"
			timeout_in_milliseconds = 1000
			interval_in_seconds = 60
			is_paused = false
			severity_degraded = "Warning"
			severity_down = "Critical"
			max_retries = 5
		}
	}
}

resource "allquiet_integration" "email" {
	display_name = "My Email Integration"
	team_id      = allquiet_team.test.id
	type         = "Email"
	integration_settings = {
		email = {}
	}
}

resource "allquiet_integration" "email_with_aliases" {
	display_name = "My Email Integration"
	team_id      = allquiet_team.test.id
	type         = "Email"
	integration_settings = {
	  email = {
		aliases = ["email-alias-acc-test@integrations.allquiet.app", "email-alias-acc-test2@integrations.allquiet.app"]
	  }
	}
  }
  

`, display_name)

	return replaceEmailAliases(result)
}

func testAccIntegrationResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_integration/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return RandomizeExample(replaceEmailAliases(string(dat)))
}

func replaceEmailAliases(str string) string {
	env := GetAccTestEnv()
	if env == "local" {
		return strings.Replace(str, "@integrations.allquiet.app", "+"+uuid.New().String()+"+"+env+"@integrations-dev.allquiet-test.app", -1)
	} else {
		return strings.Replace(str, "@integrations.allquiet.app", "+"+uuid.New().String()+"+"+env+"@integrations.allquiet.app", -1)
	}

}

func TestAccIntegrationResourceTwilio(t *testing.T) {
	accountSid := os.Getenv("ALLQUIET_TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("ALLQUIET_TWILIO_AUTH_TOKEN")
	phoneNumber := os.Getenv("ALLQUIET_TWILIO_PHONE_NUMBER")
	welcomeAudioPath := os.Getenv("ALLQUIET_TWILIO_WELCOME_AUDIO_PATH")
	if accountSid == "" || authToken == "" || phoneNumber == "" {
		t.Skip("Skipping Twilio acceptance test: set ALLQUIET_TWILIO_ACCOUNT_SID, ALLQUIET_TWILIO_AUTH_TOKEN, and ALLQUIET_TWILIO_PHONE_NUMBER")
	}
	if welcomeAudioPath == "" {
		welcomeAudioPath = "/Users/madsquist/Downloads/Max Welcome.mp3"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationResourceTwilioConfig(accountSid, authToken, phoneNumber, welcomeAudioPath),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_integration.twilio", "type", "Twilio"),
					resource.TestCheckResourceAttrSet("allquiet_integration.twilio", "integration_settings.twilio.access_settings.api_base_domain"),
					resource.TestCheckResourceAttr("allquiet_integration.twilio", "integration_settings.twilio.call_flow_config.welcome_message.type", "Audio"),
					resource.TestCheckResourceAttrSet("allquiet_integration.twilio", "integration_settings.twilio.call_flow_config.welcome_message.audio_file.object_key"),
				),
			},
			{
				ResourceName:      "allquiet_integration.twilio",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"integration_settings.twilio.access_settings.auth_token",
					"integration_settings.twilio.call_flow_config.welcome_message.audio_file.source",
				},
			},
		},
	})
}

func testAccIntegrationResourceTwilioConfig(accountSid, authToken, phoneNumber, welcomeAudioPath string) string {
	return fmt.Sprintf(`
resource "allquiet_team" "twilio" {
  display_name = "Twilio Call Routing"
}

resource "allquiet_integration" "twilio" {
  display_name = "Twilio inbound number"
  team_id      = allquiet_team.twilio.id
  type         = "Twilio"

  integration_settings = {
    twilio = {
      access_settings = {
        account_sid  = %[1]q
        auth_token   = %[2]q
        phone_number = %[3]q
      }
      call_flow_config = {
        welcome_message = {
          type = "Audio"
          audio_file = {
            source = %[4]q
          }
        }
        route = {
          team_id = allquiet_team.twilio.id
          voicemail_message = {
            type = "Text"
            text = "Please leave a message after the tone."
            locale = "en-US"
          }
        }
      }
    }
  }
}
`, accountSid, authToken, phoneNumber, welcomeAudioPath)
}
