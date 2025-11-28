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
	var host = "spt-" + uuid.New().String() + ".allquiet.com"
	var host2 = "spt-" + uuid.New().String() + ".allquiet.com"
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
					resource.TestCheckResourceAttr("allquiet_status_page.test_custom_host_settings", "custom_host_settings.cloudflare_create_custom_hostname_response.success", "true"),
					resource.TestCheckResourceAttrSet("allquiet_status_page.test_custom_host_settings", "custom_host_settings.cloudflare_create_custom_hostname_response.result.hostname"),
					resource.TestCheckResourceAttrSet("allquiet_status_page.test_custom_host_settings", "custom_host_settings.cloudflare_create_custom_hostname_response.result.id"),
					resource.TestCheckResourceAttrSet("allquiet_status_page.test_custom_host_settings", "custom_host_settings.cloudflare_create_custom_hostname_response.result.ownership_verification.name"),
					resource.TestCheckResourceAttrSet("allquiet_status_page.test_custom_host_settings", "custom_host_settings.cloudflare_create_custom_hostname_response.result.ownership_verification.type"),
					resource.TestCheckResourceAttrSet("allquiet_status_page.test_custom_host_settings", "custom_host_settings.cloudflare_create_custom_hostname_response.result.ownership_verification.value"),
					resource.TestCheckResourceAttrSet("allquiet_status_page.test_custom_host_settings", "custom_host_settings.cloudflare_create_custom_hostname_response.result.ssl.id"),
					resource.TestCheckResourceAttrSet("allquiet_status_page.test_custom_host_settings", "custom_host_settings.cloudflare_create_custom_hostname_response.result.ssl.method"),
					resource.TestCheckResourceAttrSet("allquiet_status_page.test_custom_host_settings", "custom_host_settings.cloudflare_create_custom_hostname_response.result.ssl.status"),
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

func TestAccStatusPageResourceWithGroups(t *testing.T) {
	var slug = "public-status-page-test" + uuid.New().String()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccStatusPageResourceWithGroupsConfig("Status Page One", slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_status_page.test_with_groups", "display_name", "Status Page One"),
					resource.TestCheckResourceAttr("allquiet_status_page.test_with_groups", "public_title", "Status Page One"),
					resource.TestCheckResourceAttr("allquiet_status_page.test_with_groups", "public_description", "Payment APIs and integrations"),
					resource.TestCheckResourceAttr("allquiet_status_page.test_with_groups", "history_in_days", "30"),
					resource.TestCheckResourceAttr("allquiet_status_page.test_with_groups", "disable_public_subscription", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_status_page.test_with_groups",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccStatusPageResourceWithGroupsConfig("Status Page Two", slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_status_page.test_with_groups", "display_name", "Status Page Two"),
					resource.TestCheckResourceAttr("allquiet_status_page.test_with_groups", "public_title", "Status Page Two"),
					resource.TestCheckResourceAttr("allquiet_status_page.test_with_groups", "public_description", "Payment APIs and integrations"),
					resource.TestCheckResourceAttr("allquiet_status_page.test_with_groups", "history_in_days", "30"),
					resource.TestCheckResourceAttr("allquiet_status_page.test_with_groups", "disable_public_subscription", "false"),
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

func testAccStatusPageResourceWithGroupsConfig(display_name string, slug string) string {
	return fmt.Sprintf(`
resource "allquiet_service" "test_service_1" {
  display_name = "Payment API 1"
  public_title = "Payment API 1"
}
resource "allquiet_service" "test_service_2" {
  display_name = "Payment API 2"
  public_title = "Payment API 2"
}
resource "allquiet_service" "test_service_3" {
  display_name = "AI API"
  public_title = "AI API"
}
resource "allquiet_status_page" "test_with_groups" {
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
  service_groups = [
    {
      public_display_name = "Payment APIs"
	  public_description = "Payment APIs and integrations"
      services = [
        allquiet_service.test_service_1.id,
        allquiet_service.test_service_2.id,
      ]
    },
	{
		public_display_name = "Chat Bot"
		public_description = "AI APIs and integrations"
		services = [
			allquiet_service.test_service_3.id,
		]
	}
  ]
}
	`, display_name, slug)
}

func testAccStatusPageResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_status_page/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	var result = RandomizeExample(string(dat))
	result = strings.Replace(result, "public-status-page-test", "public-status-page-test"+uuid.New().String(), -1)
	result = strings.Replace(result, "private-status-page-test", "private-status-page-test"+uuid.New().String(), -1)
	result = strings.Replace(result, "status-page-test-resource.allquiet.com", "status-page-test-resource-"+uuid.New().String()+".allquiet.com", -1)
	return result
}
