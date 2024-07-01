// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationMappingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIntegrationMappingResourceConfig(),
				Check:  resource.ComposeAggregateTestCheckFunc(
				//resource.TestCheckResourceAttr("allquiet_integration.test", "display_name", "Integration One"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_integration_mapping.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIntegrationMappingResourceConfig(),
				Check:  resource.ComposeAggregateTestCheckFunc(
				//resource.TestCheckResourceAttr("allquiet_integration.test", "display_name", "Integration Two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIntegrationMappingResourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIntegrationMappingResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_integration_mapping.datadog_custom_mapping", "integration_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_integration_mapping.datadog_custom_mapping",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIntegrationMappingResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_integration_mapping.datadog_custom_mapping", "integration_id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIntegrationMappingResourceConfig() string {
	return `
resource "allquiet_team" "test" {
  display_name = "Root"
}

resource "allquiet_integration" "test" {
  display_name = "My Datadog Integration"
  team_id = allquiet_team.test.id
  type = "Datadog"
}

resource "allquiet_integration_mapping" "test" {
   integration_id = allquiet_integration.test.id
	attributes_mapping = {
		attributes = [
		{
			name = "Severity",
			mappings = [
			{ json_path = "$.jsonBody.title" },
			{ map = "A->Critical,->Warning" }
			]	
		},
		{
			name = "Status",
			mappings = [
			{ static = "Open" }
			]	
		} ,
		{
			name = "Title",
			mappings = [
			{ xpath = "//json" },
			{ json_path = "$.jsonBody.status" },
			{ regex = "\\d+", replace = "$1" },
			{ map = "->Open" },
			]	
		} 		
		]
	}
}
`

}

func testAccIntegrationMappingResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/integration_mapping/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return string(dat)
}
