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

func TestAccIntegrationMaintenanceWindowResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIntegrationMaintenanceWindowResourceConfig("test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_integration_maintenance_window.test", "integration_id"),
					resource.TestCheckResourceAttrSet("allquiet_integration_maintenance_window.test2", "integration_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_integration_maintenance_window.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "allquiet_integration_maintenance_window.test2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIntegrationMaintenanceWindowResourceConfig("other"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_integration_maintenance_window.test", "integration_id"),
					resource.TestCheckResourceAttrSet("allquiet_integration_maintenance_window.test2", "integration_id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIntegrationMaintenanceWindowResourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIntegrationMaintenanceWindowResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_integration_maintenance_window.test1", "integration_id"),
					resource.TestCheckResourceAttrSet("allquiet_integration_maintenance_window.test2", "integration_id"),
					resource.TestCheckResourceAttrSet("allquiet_integration_maintenance_window.test3", "integration_id"),
					resource.TestCheckResourceAttrSet("allquiet_integration_maintenance_window.test4", "integration_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_integration_maintenance_window.test1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "allquiet_integration_maintenance_window.test2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "allquiet_integration_maintenance_window.test3",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "allquiet_integration_maintenance_window.test4",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIntegrationMaintenanceWindowResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_integration_maintenance_window.test1", "integration_id"),
					resource.TestCheckResourceAttrSet("allquiet_integration_maintenance_window.test2", "integration_id"),
					resource.TestCheckResourceAttrSet("allquiet_integration_maintenance_window.test3", "integration_id"),
					resource.TestCheckResourceAttrSet("allquiet_integration_maintenance_window.test4", "integration_id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIntegrationMaintenanceWindowResourceConfig(integrationIdStr string) string {
	return fmt.Sprintf(`
resource "allquiet_team" "test" {
  display_name = "Root"
}

resource "allquiet_integration" "test" {
  display_name = "My Datadog Integration"
  team_id = allquiet_team.test.id
  type = "Datadog"
}

resource "allquiet_integration" "other" {
  display_name = "My Other Integration"
  team_id = allquiet_team.test.id
  type = "Datadog"
}

resource "allquiet_integration_maintenance_window" "test" {
   integration_id = allquiet_integration.test.id
	start = "2025-01-01T00:00:00Z"
	end = "2025-01-01T00:00:00Z"
	description = "My Maintenance Window"
	type = "maintenance"
}

resource "allquiet_integration_maintenance_window" "test2" {
   integration_id = allquiet_integration.%s.id
	start = "2025-01-01T00:00:00Z"
	end = "2025-01-01T00:00:00Z"
	description = "My Muted Window"
	type = "muted"
}
`, integrationIdStr)
}

func testAccIntegrationMaintenanceWindowResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_integration_maintenance_window/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return RandomizeExample(string(dat))
}
