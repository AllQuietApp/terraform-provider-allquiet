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
					resource.TestCheckResourceAttr("allquiet_routing.test", "rules.1.actions.set_attributes.0.name", "Team"),
					resource.TestCheckResourceAttr("allquiet_routing.test", "rules.1.actions.set_attributes.0.value", "Sales"),
					resource.TestCheckResourceAttr("allquiet_routing.test", "rules.1.actions.set_attributes.0.hide_in_previews", "true"),
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
					resource.TestCheckResourceAttr("allquiet_routing.test", "rules.1.actions.set_attributes.0.name", "Team"),
					resource.TestCheckResourceAttr("allquiet_routing.test", "rules.1.actions.set_attributes.0.value", "Sales"),
					resource.TestCheckResourceAttr("allquiet_routing.test", "rules.1.actions.set_attributes.0.hide_in_previews", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRoutingResourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRoutingResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_routing.example_1", "display_name", "Route to specific team based on attribute"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_routing.example_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccRoutingResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_routing.example_1", "display_name", "Route to specific team based on attribute"),
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
	   assign_to_teams = [allquiet_team.test.id]
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
	  	set_attributes = [
			{
				name = "Team"
				value = "Sales"
				hide_in_previews = true
			}
		]
      }
	}
  ]
}

resource "allquiet_routing" "test_with_team_connection_settings" {
	display_name = %[1]q
	team_id = allquiet_team.root.id
	team_connection_settings = {	
		team_connection_mode = "SelectedTeams"
		team_ids = [allquiet_team.root.id, allquiet_team.test.id]
	}
	rules = [
	  {
		conditions = {
		  statuses = ["Open"]
		  severities = ["Critical", "Warning"]
	   },
	   channels = {
	   },
	   actions = {
		 assign_to_teams = [allquiet_team.test.id]
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
			set_attributes = [
			  {
				  name = "Team"
				  value = "Sales"
				  hide_in_previews = true
			  }
		  ]
		}
	  }
	]
  }


resource "allquiet_routing" "test_with_team_connection_settings_organization_teams" {
	display_name = %[1]q
	team_id = allquiet_team.root.id
	team_connection_settings = {	
		team_connection_mode = "OrganizationTeams"
	}
	rules = [
	  {
		conditions = {
		  statuses = ["Open"]
		  severities = ["Critical", "Warning"]
	   },
	   channels = {
	   },
	   actions = {
		 assign_to_teams = [allquiet_team.test.id]
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
			set_attributes = [
			  {
				  name = "Team"
				  value = "Sales"
				  hide_in_previews = true
			  }
		  ]
		}
	  }
	]
  }
`, display_name)

}

func testAccRoutingResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_routing/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return RandomizeExample(string(dat))
}
