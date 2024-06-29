// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoutingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRoutingResourceConfig("Routing One"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_routing.test", "display_name", "Routing One"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_routing.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccRoutingResourceConfig("Routing Two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_routing.test", "display_name", "Routing Two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRoutingResourceConfig(display_name string) string {
	return fmt.Sprintf(`
resource "allquiet_team" "root" {
  display_name = "Root"
}

resource "allquiet_team" "test" {
	display_name = "Test"
  }

resource "allquiet_routing" "test" {
  display_name = %[1]q
  team_id = allquiet_team.root.id
  rules = [
    {
	  conditions = {
	    statuses = ["Open"]
		severities = ["Critical", "Warning"]
     },
	 channels = {
	 },
	 actions = {
	   route_to_teams = [allquiet_team.test.id]
     }
    },
	{
	  conditions = {
	    attributes = [
		  {
		    name = "source"
		    operator = "=" 
		    value = "web"
		  }
		]
	  },
	  channels = {
	    notification_channels = ["VoiceCall"]
	  },
	  actions = {
      }
	}
  ]
}
`, display_name)

}
