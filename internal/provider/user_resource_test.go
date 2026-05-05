// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

func TestAccUserResourceRejectsInlineNotificationSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "allquiet_user" "test" {
  display_name = "Millie"
  email        = "acceptance-tests+millie+%s@allquiet.app"
  incident_notification_settings = {
    should_send_sms    = true
    delay_in_min_sms   = 5
    severities_sms     = ["Critical"]
    should_call_voice  = false
    delay_in_min_voice = 0
    severities_voice   = []
    should_send_push   = true
    delay_in_min_push  = 0
    severities_push    = ["Critical"]
    should_send_email  = true
    delay_in_min_email = 0
    severities_email   = ["Critical"]
  }
}
`, uuid.New().String()),
				ExpectError: regexp.MustCompile("incident_notification_settings on allquiet_user has been removed"),
			},
		},
	})
}

func TestAccUserResourceRejectsInlinePhoneNumber(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "allquiet_user" "test" {
  display_name = "Millie"
  email        = "acceptance-tests+millie+%s@allquiet.app"
  phone_number = "+12035479055"
}
`, uuid.New().String()),
				ExpectError: regexp.MustCompile("phone_number on allquiet_user has been removed"),
			},
		},
	})
}
