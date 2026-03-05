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

func TestAccOutboundIntegrationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOutboundIntegrationResourceConfig("Outbound One"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_outbound_integration.test", "display_name", "Outbound One"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_outbound_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccOutboundIntegrationResourceConfig("Outbound Two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_outbound_integration.test", "display_name", "Outbound Two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccOutboundIntegrationResourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOutboundIntegrationResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack", "display_name", "My Slack Integration"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_outbound_integration.slack",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccOutboundIntegrationResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack", "display_name", "My Slack Integration"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccOutboundIntegrationResourceConfig(display_name string) string {

	return fmt.Sprintf(`
resource "allquiet_team" "test" {
  display_name = "Root"
}

resource "allquiet_outbound_integration" "test" {
  display_name = %[1]q
  team_id = allquiet_team.test.id
  type = "Slack"
  triggers_only_on_forwarded = true
  team_connection_settings = {
    team_connection_mode = "SelectedTeams"
    team_ids = [allquiet_team.test.id]
  }
}

resource "allquiet_outbound_integration" "test_empty_team_ids" {
  display_name = "%[1]s (Empty Team IDs)"	
  team_id = allquiet_team.test.id
  type = "Slack"
  triggers_only_on_forwarded = true
  team_connection_settings = {
    team_connection_mode = "OrganizationTeams"
    team_ids = []
  }
}

resource "allquiet_outbound_integration" "test_null_team_ids" {
  display_name = "%[1]s (Null Team IDs)"	
  team_id = allquiet_team.test.id
  type = "Slack"
  triggers_only_on_forwarded = true
  team_connection_settings = {
    team_connection_mode = "OrganizationTeams"
    team_ids = null
  }
}



resource "allquiet_outbound_integration" "test_webhook" {
	display_name = "%[1]s (Webhook)"	
	team_id = allquiet_team.test.id
	type = "Webhook"
	triggers_only_on_forwarded = true
	skip_updating_after_forwarding = true
	team_connection_settings = {
	  team_connection_mode = "OrganizationTeams"
	  team_ids = null
	}
  }
`, display_name)

}

func testAccOutboundIntegrationResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_outbound_integration/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return RandomizeExample(string(dat))
}

func TestAccOutboundIntegrationResourceSlackSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Slack integration without settings (should remain null)
			{
				Config: testAccOutboundIntegrationSlackSettingsConfig("basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_basic", "display_name", "Slack Basic"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_basic", "type", "Slack"),
					resource.TestCheckNoResourceAttr("allquiet_outbound_integration.slack_basic", "slack_settings"),
				),
			},
			// Create Slack integration with selected_channel_ids
			{
				Config: testAccOutboundIntegrationSlackSettingsConfig("with_channels"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_channels", "display_name", "Slack With Channels"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_channels", "type", "Slack"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_channels", "slack_settings.selected_channel_ids.#", "2"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_channels", "slack_settings.selected_channel_ids.0", "channel1"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_channels", "slack_settings.selected_channel_ids.1", "channel2"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_channels", "slack_settings.tag_on_call_members", "true"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_channels", "slack_settings.is_slack_message_payload_read_only", "false"),
				),
			},
			// Create Slack integration with severity_based_channel_settings
			{
				Config: testAccOutboundIntegrationSlackSettingsConfig("with_severity"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_severity", "display_name", "Slack With Severity"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_severity", "type", "Slack"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_severity", "slack_settings.severity_based_channel_settings.selected_channel_ids_minor.#", "1"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_severity", "slack_settings.severity_based_channel_settings.selected_channel_ids_minor.0", "minor_channel"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_severity", "slack_settings.severity_based_channel_settings.selected_channel_ids_warning.#", "1"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_severity", "slack_settings.severity_based_channel_settings.selected_channel_ids_warning.0", "warning_channel"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_severity", "slack_settings.severity_based_channel_settings.selected_channel_ids_critical.#", "1"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_severity", "slack_settings.severity_based_channel_settings.selected_channel_ids_critical.0", "critical_channel"),
				),
			},
			// Create Slack integration with on_call_reminder settings
			{
				Config: testAccOutboundIntegrationSlackSettingsConfig("with_reminders"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_reminders", "display_name", "Slack With Reminders"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_reminders", "type", "Slack"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_reminders", "slack_settings.on_call_reminder_schedule_settings.run_time", "09:00"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_reminders", "slack_settings.on_call_reminder_schedule_settings.days_of_week.#", "2"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_reminders", "slack_settings.on_call_reminder_schedule_settings.days_of_week.0", "mon"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_reminders", "slack_settings.on_call_reminder_schedule_settings.days_of_week.1", "wed"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_reminders", "slack_settings.on_call_reminder_channel_ids.#", "2"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_reminders", "slack_settings.on_call_reminder_channel_ids.0", "reminder1"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_reminders", "slack_settings.on_call_reminder_channel_ids.1", "reminder2"),
				),
			},
			// Update: Change from selected_channel_ids to severity_based_channel_settings
			{
				Config: testAccOutboundIntegrationSlackSettingsConfig("update_severity"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_channels", "display_name", "Slack With Channels"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_channels", "slack_settings.severity_based_channel_settings.selected_channel_ids_minor.#", "1"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.slack_with_channels", "slack_settings.severity_based_channel_settings.selected_channel_ids_minor.0", "minor_channel"),
					resource.TestCheckNoResourceAttr("allquiet_outbound_integration.slack_with_channels", "slack_settings.selected_channel_ids"),
				),
			},
		},
	})
}

func testAccOutboundIntegrationSlackSettingsConfig(testType string) string {
	baseConfig := `
resource "allquiet_team" "test" {
  display_name = "Root"
}
`

	switch testType {
	case "basic":
		return baseConfig + `
resource "allquiet_outbound_integration" "slack_basic" {
  display_name = "Slack Basic"
  team_id      = allquiet_team.test.id
  type         = "Slack"
}
`
	case "with_channels":
		return baseConfig + `
resource "allquiet_outbound_integration" "slack_with_channels" {
  display_name = "Slack With Channels"
  team_id      = allquiet_team.test.id
  type         = "Slack"
  
  slack_settings = {
    selected_channel_ids              = ["channel1", "channel2"]
    tag_on_call_members               = true
    is_slack_message_payload_read_only = false
  }
}
`
	case "with_severity":
		return baseConfig + `
resource "allquiet_outbound_integration" "slack_with_severity" {
  display_name = "Slack With Severity"
  team_id      = allquiet_team.test.id
  type         = "Slack"
  
  slack_settings = {
    severity_based_channel_settings = {
      selected_channel_ids_minor    = ["minor_channel"]
      selected_channel_ids_warning   = ["warning_channel"]
      selected_channel_ids_critical  = ["critical_channel"]
    }
  }
}
`
	case "with_reminders":
		return baseConfig + `
resource "allquiet_outbound_integration" "slack_with_reminders" {
  display_name = "Slack With Reminders"
  team_id      = allquiet_team.test.id
  type         = "Slack"
  
  slack_settings = {
    on_call_reminder_schedule_settings = {
      run_time    = "09:00"
      days_of_week = ["mon", "wed"]
    }
    on_call_reminder_channel_ids = ["reminder1", "reminder2"]
  }
}
`
	case "update_severity":
		return baseConfig + `
resource "allquiet_outbound_integration" "slack_with_channels" {
  display_name = "Slack With Channels"
  team_id      = allquiet_team.test.id
  type         = "Slack"
  
  slack_settings = {
    severity_based_channel_settings = {
      selected_channel_ids_minor    = ["minor_channel"]
    }
  }
}
`
	default:
		return baseConfig
	}
}

func TestAccOutboundIntegrationResourceMattermostSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Mattermost integration without settings (should remain null)
			{
				Config: testAccOutboundIntegrationMattermostSettingsConfig("basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_basic", "display_name", "Mattermost Basic"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_basic", "type", "Mattermost"),
					resource.TestCheckNoResourceAttr("allquiet_outbound_integration.mattermost_basic", "mattermost_settings"),
				),
			},
			// Create Mattermost integration with selected_channel_ids
			{
				Config: testAccOutboundIntegrationMattermostSettingsConfig("with_channels"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_channels", "display_name", "Mattermost With Channels"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_channels", "type", "Mattermost"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_channels", "mattermost_settings.send_incidents_to_mattermost", "true"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_channels", "mattermost_settings.create_incidents_from_mattermost", "false"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_channels", "mattermost_settings.selected_team_id", "team-id-1"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_channels", "mattermost_settings.selected_channel_ids.#", "2"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_channels", "mattermost_settings.selected_channel_ids.0", "channel1"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_channels", "mattermost_settings.selected_channel_ids.1", "channel2"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_channels", "mattermost_settings.is_message_read_only", "false"),
				),
			},
			// Create Mattermost integration with severity_based_channel_settings
			{
				Config: testAccOutboundIntegrationMattermostSettingsConfig("with_severity"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_severity", "display_name", "Mattermost With Severity"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_severity", "type", "Mattermost"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_severity", "mattermost_settings.send_incidents_to_mattermost", "true"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_severity", "mattermost_settings.selected_team_id", "team-id-2"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_severity", "mattermost_settings.severity_based_channel_settings.selected_channel_ids_minor.#", "1"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_severity", "mattermost_settings.severity_based_channel_settings.selected_channel_ids_minor.0", "minor_channel"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_severity", "mattermost_settings.severity_based_channel_settings.selected_channel_ids_warning.#", "1"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_severity", "mattermost_settings.severity_based_channel_settings.selected_channel_ids_warning.0", "warning_channel"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_severity", "mattermost_settings.severity_based_channel_settings.selected_channel_ids_critical.#", "1"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_severity", "mattermost_settings.severity_based_channel_settings.selected_channel_ids_critical.0", "critical_channel"),
				),
			},
			// Update: Change from selected_channel_ids to severity_based_channel_settings
			{
				Config: testAccOutboundIntegrationMattermostSettingsConfig("update_severity"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_channels", "display_name", "Mattermost With Channels"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_channels", "mattermost_settings.severity_based_channel_settings.selected_channel_ids_minor.#", "1"),
					resource.TestCheckResourceAttr("allquiet_outbound_integration.mattermost_with_channels", "mattermost_settings.severity_based_channel_settings.selected_channel_ids_minor.0", "minor_channel"),
				),
			},
		},
	})
}

func testAccOutboundIntegrationMattermostSettingsConfig(testType string) string {
	baseConfig := `
resource "allquiet_team" "test" {
  display_name = "Root"
}
`

	switch testType {
	case "basic":
		return baseConfig + `
resource "allquiet_outbound_integration" "mattermost_basic" {
  display_name = "Mattermost Basic"
  team_id      = allquiet_team.test.id
  type         = "Mattermost"
}
`
	case "with_channels":
		return baseConfig + `
resource "allquiet_outbound_integration" "mattermost_with_channels" {
  display_name = "Mattermost With Channels"
  team_id      = allquiet_team.test.id
  type         = "Mattermost"

  mattermost_settings = {
    send_incidents_to_mattermost     = true
    create_incidents_from_mattermost = false
    base_url                         = "https://mattermost.com"
    selected_team_id                 = "team-id-1"
    selected_channel_ids            = ["channel1", "channel2"]
    is_message_read_only            = false
  }
}
`
	case "with_severity":
		return baseConfig + `
resource "allquiet_outbound_integration" "mattermost_with_severity" {
  display_name = "Mattermost With Severity"
  team_id      = allquiet_team.test.id
  type         = "Mattermost"

  mattermost_settings = {
    send_incidents_to_mattermost = true
	create_incidents_from_mattermost = false
    base_url                     = "https://mattermost.com"
    selected_team_id             = "team-id-2"
    severity_based_channel_settings = {
      selected_channel_ids_minor    = ["minor_channel"]
      selected_channel_ids_warning  = ["warning_channel"]
      selected_channel_ids_critical = ["critical_channel"]
    }
  }
}
`
	case "update_severity":
		return baseConfig + `
resource "allquiet_outbound_integration" "mattermost_with_channels" {
  display_name = "Mattermost With Channels"
  team_id      = allquiet_team.test.id
  type         = "Mattermost"

  mattermost_settings = {
    send_incidents_to_mattermost = true
	create_incidents_from_mattermost = false
    base_url                     = "https://mattermost.com"
    selected_team_id             = "team-id-1"
    severity_based_channel_settings = {
      selected_channel_ids_minor = ["minor_channel"]
    }
  }
}
`
	default:
		return baseConfig
	}
}
