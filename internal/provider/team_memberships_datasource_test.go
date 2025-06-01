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

func TestAccTeamMembershipsDataSource(t *testing.T) {
	uid := uuid.New().String()
	email := fmt.Sprintf("acceptance-tests+millie+%s@allquiet.app", uid)
	uid2 := uuid.New().String()
	email2 := fmt.Sprintf("acceptance-tests+millie+%s@allquiet.app", uid2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamMembershipsDataSourceConfig("Millie Bobby Brown", email, "Miley Cyrus", email2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.allquiet_team_memberships.test", "team_memberships.#", "2"),
					resource.TestCheckResourceAttrSet("data.allquiet_team_memberships.test", "team_memberships.0.id"),
					resource.TestCheckResourceAttrSet("data.allquiet_team_memberships.test", "team_memberships.0.user_id"),
					resource.TestCheckResourceAttrSet("data.allquiet_team_memberships.test", "team_memberships.0.team_id"),
					resource.TestCheckResourceAttr("data.allquiet_team_memberships.test", "team_memberships.0.role", "Member"),
					resource.TestCheckResourceAttrSet("data.allquiet_team_memberships.test", "team_memberships.1.id"),
					resource.TestCheckResourceAttrSet("data.allquiet_team_memberships.test", "team_memberships.1.user_id"),
					resource.TestCheckResourceAttrSet("data.allquiet_team_memberships.test", "team_memberships.1.team_id"),
					resource.TestCheckResourceAttr("data.allquiet_team_memberships.test", "team_memberships.1.role", "Member"),
				),
			},
		},
	})
}

func testAccTeamMembershipsDataSourceConfig(displayName, email, displayName2, email2 string) string {
	return fmt.Sprintf(`

		resource "allquiet_user" "test" {
			display_name = %[1]q
			email        = %[2]q
		}

		resource "allquiet_user" "test2" {
			display_name = %[3]q
			email        = %[4]q
		}


		resource "allquiet_team" "test" {
			display_name = "TF Data Source Team"
		}

		resource "allquiet_team_membership" "test" {
			user_id = allquiet_user.test.id
			team_id = allquiet_team.test.id
			role = "Member"
		}

		resource "allquiet_team_membership" "test2" {
			user_id = allquiet_user.test2.id
			team_id = allquiet_team.test.id
			role = "Member"
		}

		data "allquiet_team_memberships" "test" {
			team_id = allquiet_team.test.id
			depends_on = [allquiet_team_membership.test, allquiet_team_membership.test2]
		}
	`, displayName, email, displayName2, email2)
}

func TestAccTeamMembershipsDataSourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamMembershipsDataSourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.allquiet_team_memberships.memberships_by_team", "team_memberships.#"),
				),
			},
		},
	})
}

func testAccTeamMembershipsDataSourceExample() string {
	absPath, _ := filepath.Abs("../../examples/data-sources/allquiet_team_memberships/data-source.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return string(dat)
}
