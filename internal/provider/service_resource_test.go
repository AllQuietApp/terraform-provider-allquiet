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

func TestAccServiceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccServiceResourceConfig("Service One"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_service.test", "display_name", "Service One"),
					resource.TestCheckResourceAttr("allquiet_service.test", "public_title", "Service One"),
					resource.TestCheckResourceAttr("allquiet_service.test", "public_description", "Payment APIs and integrations"),
					resource.TestCheckResourceAttr("allquiet_service.test", "templates.#", "2"),
					resource.TestCheckResourceAttr("allquiet_service.test", "templates.0.display_name", "Refunds delayed"),
					resource.TestCheckResourceAttr("allquiet_service.test", "templates.0.message", "Refunds are currently delayed. All refunds will be processed but can currently take longer than usual to complete."),
					resource.TestCheckResourceAttr("allquiet_service.test", "templates.1.display_name", "Payment gateway down"),
					resource.TestCheckResourceAttr("allquiet_service.test", "templates.1.message", "Our payment gateway is currently down. We are working to resolve the issue as soon as possible."),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_service.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccServiceResourceConfig("Service Two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_service.test", "display_name", "Service Two"),
					resource.TestCheckResourceAttr("allquiet_service.test", "public_title", "Service Two"),
					resource.TestCheckResourceAttr("allquiet_service.test", "public_description", "Payment APIs and integrations"),
					resource.TestCheckResourceAttr("allquiet_service.test", "templates.#", "2"),
					resource.TestCheckResourceAttr("allquiet_service.test", "templates.0.display_name", "Refunds delayed"),
					resource.TestCheckResourceAttr("allquiet_service.test", "templates.0.message", "Refunds are currently delayed. All refunds will be processed but can currently take longer than usual to complete."),
					resource.TestCheckResourceAttr("allquiet_service.test", "templates.1.display_name", "Payment gateway down"),
					resource.TestCheckResourceAttr("allquiet_service.test", "templates.1.message", "Our payment gateway is currently down. We are working to resolve the issue as soon as possible."),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccServiceResourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccServiceResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_service.payment_api", "display_name", "Payment API"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_service.payment_api",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccServiceResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_service.payment_api", "display_name", "Payment API"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccServiceResourceConfig(display_name string) string {
	return fmt.Sprintf(`
resource "allquiet_service" "test" {
  display_name = %[1]q
  public_title = %[1]q
  public_description = "Payment APIs and integrations"
  templates = [
		{
			display_name = "Refunds delayed"
			message = "Refunds are currently delayed. All refunds will be processed but can currently take longer than usual to complete."
		},
		{
			display_name = "Payment gateway down"
			message = "Our payment gateway is currently down. We are working to resolve the issue as soon as possible."
		}
	]
}

resource "allquiet_service" "test_with_team_connection_settings" {
  display_name = %[1]q
  public_title = %[1]q
  public_description = "Payment APIs and integrations"
  team_connection_settings = {
    team_connection_mode = "OrganizationTeams"
  }
  templates = [
		{
			display_name = "Refunds delayed"
			message = "Refunds are currently delayed. All refunds will be processed but can currently take longer than usual to complete."
		}
	]
}
`, display_name)

}

func testAccServiceResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_service/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return RandomizeExample(string(dat))
}
