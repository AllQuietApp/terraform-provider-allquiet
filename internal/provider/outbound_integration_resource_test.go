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
