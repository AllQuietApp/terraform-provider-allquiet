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
					resource.TestCheckResourceAttr("allquiet_service.test_with_integrations", "integrations.#", "1"),
					resource.TestCheckResourceAttrSet("allquiet_service.test_with_integrations", "integrations.0.id"),
					resource.TestCheckResourceAttrPair("allquiet_service.test_with_integrations", "integrations.0.integration_id", "allquiet_integration.service_integration", "id"),
					resource.TestCheckResourceAttr("allquiet_service.test_with_integrations", "integrations.0.severities.#", "2"),
					resource.TestCheckResourceAttr("allquiet_service.test_with_integrations", "integrations.0.severities.0", "Critical"),
					resource.TestCheckResourceAttr("allquiet_service.test_with_integrations", "integrations.0.severities.1", "Warning"),
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
					resource.TestCheckResourceAttr("allquiet_service.test_with_integrations", "integrations.#", "1"),
					resource.TestCheckResourceAttr("allquiet_service.test_with_integrations", "integrations.0.severities.#", "2"),
					resource.TestCheckResourceAttr("allquiet_service.test_with_integrations", "integrations.0.severities.0", "Critical"),
					resource.TestCheckResourceAttr("allquiet_service.test_with_integrations", "integrations.0.severities.1", "Warning"),
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
resource "allquiet_team" "service_team" {
  display_name = "Service Test Team"
}

resource "allquiet_integration" "service_integration" {
  display_name = "Integration for Service Link"
  team_id      = allquiet_team.service_team.id
  type         = "Webhook"
}

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

resource "allquiet_service" "test_with_integrations" {
  display_name = %[1]q
  public_title = %[1]q
  public_description = "Service with linked integrations"
  team_connection_settings = {
    team_connection_mode = "SelectedTeams"
    team_ids             = [allquiet_team.service_team.id]
  }
  integrations = [
    {
      integration_id = allquiet_integration.service_integration.id
      severities      = ["Critical", "Warning"]
    }
  ]
  templates = [
    {
      display_name = "Refunds delayed"
      message       = "Refunds are currently delayed."
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
