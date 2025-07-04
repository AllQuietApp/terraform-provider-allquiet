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

func TestAccStatusPageResource(t *testing.T) {
	var slug = "public-status-page-test" + uuid.New().String()
	var host = "status-page-test-" + uuid.New().String() + ".allquiet.com"
	var host2 = "status-page-test-" + uuid.New().String() + ".allquiet.com"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccStatusPageResourceConfig("Status Page One", slug, host),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_status_page.test", "display_name", "Status Page One"),
					resource.TestCheckResourceAttr("allquiet_status_page.test", "public_title", "Status Page One"),
					resource.TestCheckResourceAttr("allquiet_status_page.test", "public_description", "Payment APIs and integrations"),
					resource.TestCheckResourceAttr("allquiet_status_page.test", "history_in_days", "30"),
					resource.TestCheckResourceAttr("allquiet_status_page.test", "disable_public_subscription", "false"),
					resource.TestCheckResourceAttr("allquiet_status_page.test_custom_host_settings", "custom_host_settings.host", host),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_status_page.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "allquiet_status_page.test_custom_host_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccStatusPageResourceConfig("Status Page Two", slug, host2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_status_page.test", "display_name", "Status Page Two"),
					resource.TestCheckResourceAttr("allquiet_status_page.test", "public_title", "Status Page Two"),
					resource.TestCheckResourceAttr("allquiet_status_page.test", "public_description", "Payment APIs and integrations"),
					resource.TestCheckResourceAttr("allquiet_status_page.test", "history_in_days", "30"),
					resource.TestCheckResourceAttr("allquiet_status_page.test", "disable_public_subscription", "false"),
					resource.TestCheckResourceAttr("allquiet_status_page.test_custom_host_settings", "custom_host_settings.host", host2),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccStatusPageResourceExample(t *testing.T) {
	var config = testAccStatusPageResourceExample()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_status_page.public_status_page", "display_name", "Public Status Page"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_status_page.public_status_page",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_status_page.public_status_page", "display_name", "Public Status Page"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccStatusPageResourceConfig(display_name string, slug string, host string) string {
	return fmt.Sprintf(`
resource "allquiet_status_page" "test" {
  display_name = %[1]q
  public_title = %[1]q
  public_description = "Payment APIs and integrations"  
  history_in_days = 30
  disable_public_subscription = false
  banner_background_color = "#000000"
  banner_background_color_dark_mode = "#447788"
  banner_text_color = "#ffffff"
  banner_text_color_dark_mode = "#ffffff"
  slug = %[2]q
}
resource "allquiet_status_page" "test_custom_host_settings" {
  display_name = %[1]q
  public_title = %[1]q
  public_description = "Payment APIs and integrations"  
  history_in_days = 30
  disable_public_subscription = false
  banner_background_color = "#000000"
  banner_background_color_dark_mode = "#447788"
  banner_text_color = "#ffffff"
  banner_text_color_dark_mode = "#ffffff"
  custom_host_settings = {
    host = %[3]q
  }
}

	`, display_name, slug, host)
}

func testAccStatusPageResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_status_page/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return strings.Replace(RandomizeExample(string(dat)), "public-status-page-test", "public-status-page-test"+uuid.New().String(), -1)
}
