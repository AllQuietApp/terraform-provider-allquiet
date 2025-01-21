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

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig("Millie Brown"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_user.test", "display_name", "Millie Brown"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_user.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccUserResourceConfig("Millie Bobby Brown"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_user.test", "display_name", "Millie Bobby Brown"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccUserResourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_user.millie_brown", "display_name", "Millie Bobby Brown"),
					resource.TestCheckResourceAttr("allquiet_user.billie_eilish", "incident_notification_settings.disabled_intents_sms.0", "Resolved"),
					resource.TestCheckResourceAttr("allquiet_user.billie_eilish", "incident_notification_settings.disabled_intents_sms.#", "1"),
					resource.TestCheckResourceAttr("allquiet_user.billie_eilish", "incident_notification_settings.disabled_intents_voice.0", "Resolved"),
					resource.TestCheckResourceAttr("allquiet_user.billie_eilish", "incident_notification_settings.disabled_intents_voice.#", "1"),
					resource.TestCheckResourceAttr("allquiet_user.billie_eilish", "incident_notification_settings.disabled_intents_push.#", "0"),
					resource.TestCheckResourceAttr("allquiet_user.billie_eilish", "incident_notification_settings.disabled_intents_email.#", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_user.millie_brown",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccUserResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_user.millie_brown", "display_name", "Millie Bobby Brown"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccUserResourceConfig(display_name string) string {
	return fmt.Sprintf(`
resource "allquiet_user" "test" {
  display_name =  %[1]q
  email = "acceptance-tests+millie+%s@allquiet.app"
  phone_number = "+12035479055"
}

`, display_name, uuid.New().String())

}

func testAccUserResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_user/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	result := RandomizeExample(string(dat))

	return result
}
