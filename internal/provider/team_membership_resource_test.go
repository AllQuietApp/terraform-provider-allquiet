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

func TestAccTeamMembershipResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamMembershipResourceConfig("taylor_swift"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_team_membership.test", "user_id", "allquiet_user.taylor_swift.id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_team_membership.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccTeamMembershipResourceConfig("millie_brown"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_team_membership.test", "user_id", "allquiet_user.millie_brown.id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTeamMembershipResourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamMembershipResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_team_membership.my_team_millie_brown", "user_id", "allquiet_user.millie_brown.id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_team_membership.millie_brown",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccTeamMembershipResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_team_membership.my_team_millie_brown", "user_id", "allquiet_user.millie_brown.id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTeamMembershipResourceConfig(user_name string) string {
	return fmt.Sprintf(`
resource "allquiet_team" "team" {
  display_name = "Root"
}

resource "allquiet_user" "millie_brown" {
  display_name =  "Millie Bobby Brown"
  email = "millie@acme.com"
}

resource "allquiet_user" "taylor_swift" {
  display_name =  "Taylor Swift"
  email = "taylor@acme.com"
}

resource "allquiet_team_membership" "test" {
	user_id = allquiet_user.%[1]q.id
  	team_id = allquiet_team.team.id
	role = "Member"
}

`, user_name)

}

func testAccTeamMembershipResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_team_membership/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return string(dat)
}
