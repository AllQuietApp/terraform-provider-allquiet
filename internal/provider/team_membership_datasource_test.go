// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTeamMembershipDataSource(t *testing.T) {
	uid := uuid.New().String()
	email := fmt.Sprintf("acceptance-tests+millie+%s@allquiet.app", uid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamMembershipDataSourceConfig("Millie Bobby Brown", email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.allquiet_team_membership.test", "user_id"),
					resource.TestCheckResourceAttrSet("data.allquiet_team_membership.test", "team_id"),
					resource.TestCheckResourceAttr("data.allquiet_team_membership.test", "role", "Member"),
				),
			},
		},
	})
}

func TestAccTeamMembershipDataSourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamMembershipDataSourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.allquiet_team_membership.test", "user_id"),
					resource.TestCheckResourceAttrSet("data.allquiet_team_membership.test", "team_id"),
					resource.TestCheckResourceAttr("data.allquiet_team_membership.test", "role", "Member"),
				),
			},
		},
	})
}

func testAccTeamMembershipDataSourceConfig(displayName, email string) string {
	return fmt.Sprintf(`

		resource "allquiet_user" "test" {
			display_name = %[1]q
			email        = %[2]q
		}

		resource "allquiet_team" "test" {
			display_name = "TF Data Source Team"
		}

		resource "allquiet_team_membership" "test" {
			user_id = allquiet_user.test.id
			team_id = allquiet_team.test.id
			role = "Member"
		}

		data "allquiet_team_membership" "test" {
			team_id = allquiet_team.test.id
			user_id = allquiet_user.test.id
			role = "Member"

			depends_on = [allquiet_team_membership.test]
		}
	`, displayName, email)
}

func testAccTeamMembershipDataSourceExample() string {
	absPath, _ := filepath.Abs("../../examples/data-sources/allquiet_team_membership/data-source.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return string(dat)
}
