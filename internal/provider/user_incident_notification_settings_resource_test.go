// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserIncidentNotificationSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserIncidentNotificationSettingsResourceConfig(true, 5, []string{"Critical"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_user_incident_notification_settings.test", "user_id"),
					resource.TestCheckResourceAttr("allquiet_user_incident_notification_settings.test", "should_send_sms", "true"),
					resource.TestCheckResourceAttr("allquiet_user_incident_notification_settings.test", "delay_in_min_sms", "5"),
					resource.TestCheckResourceAttr("allquiet_user_incident_notification_settings.test", "severities_sms.#", "1"),
					resource.TestCheckResourceAttr("allquiet_user_incident_notification_settings.test", "severities_sms.0", "Critical"),
				),
			},
			{
				ResourceName:      "allquiet_user_incident_notification_settings.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccUserIncidentNotificationSettingsResourceConfig(false, 10, []string{"Critical", "Warning"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_user_incident_notification_settings.test", "should_send_sms", "false"),
					resource.TestCheckResourceAttr("allquiet_user_incident_notification_settings.test", "delay_in_min_sms", "10"),
					resource.TestCheckResourceAttr("allquiet_user_incident_notification_settings.test", "severities_sms.#", "2"),
				),
			},
		},
	})
}

func TestAccUserIncidentNotificationSettingsResourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserIncidentNotificationSettingsResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_user_incident_notification_settings.taylor", "user_id"),
					resource.TestCheckResourceAttr("allquiet_user_incident_notification_settings.billie_eilish", "disabled_intents_sms.#", "1"),
					resource.TestCheckResourceAttr("allquiet_user_incident_notification_settings.billie_eilish", "disabled_intents_sms.0", "Resolved"),
					resource.TestCheckResourceAttr("allquiet_user_incident_notification_settings.billie_eilish", "disabled_intents_voice.#", "1"),
					resource.TestCheckResourceAttr("allquiet_user_incident_notification_settings.billie_eilish", "disabled_intents_voice.0", "Resolved"),
					resource.TestCheckResourceAttr("allquiet_user_incident_notification_settings.billie_eilish", "disabled_intents_push.#", "0"),
					resource.TestCheckResourceAttr("allquiet_user_incident_notification_settings.billie_eilish", "disabled_intents_email.#", "0"),
				),
			},
			{
				ResourceName:      "allquiet_user_incident_notification_settings.taylor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccUserIncidentNotificationSettingsResourceConfig(shouldSendSMS bool, delayInMinSMS int, severitiesSMS []string) string {
	severitiesList := ""
	for i, s := range severitiesSMS {
		if i > 0 {
			severitiesList += ", "
		}
		severitiesList += fmt.Sprintf("%q", s)
	}

	return fmt.Sprintf(`
resource "allquiet_user" "test" {
  display_name = "Acceptance Tests"
  email        = "acceptance-tests+notification-settings+%s@allquiet.app"
}

resource "allquiet_user_incident_notification_settings" "test" {
  user_id = allquiet_user.test.id

  should_send_sms  = %t
  delay_in_min_sms = %d
  severities_sms   = [%s]

  should_call_voice  = false
  delay_in_min_voice = 0
  severities_voice   = []

  should_send_push  = true
  delay_in_min_push = 0
  severities_push   = ["Critical", "Warning"]

  should_send_email  = true
  delay_in_min_email = 0
  severities_email   = ["Critical", "Warning", "Minor"]
}
`, uuid.New().String(), shouldSendSMS, delayInMinSMS, severitiesList)
}

func testAccUserIncidentNotificationSettingsResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_user_incident_notification_settings/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return RandomizeExample(string(dat))
}
