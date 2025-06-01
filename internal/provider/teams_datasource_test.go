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

func TestAccTeamsDataSource(t *testing.T) {
	uid := uuid.New().String()
	displayName := fmt.Sprintf("TF Acceptance Test %s", uid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamsDataSourceConfig(displayName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.allquiet_teams.test_by_display_name", "teams.#", "3"),
				),
			},
		},
	})
}

func testAccTeamsDataSourceConfig(displayName string) string {
	return fmt.Sprintf(`

		resource "allquiet_team" "test1" {
			display_name = "%[1]s 1"
		}

		resource "allquiet_team" "test2" {
			display_name = "%[1]s 2"
		}

		resource "allquiet_team" "test3" {
			display_name = "%[1]s 3"
		}

		data "allquiet_teams" "test_by_display_name" {
			display_name = "%[1]s"
			depends_on = [allquiet_team.test1, allquiet_team.test2, allquiet_team.test3]
		}
	`, displayName)
}

func TestAccTeamsDataSourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamsDataSourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.allquiet_teams.teams_by_display_name", "teams.#", "1"),
					resource.TestCheckResourceAttrSet("data.allquiet_teams.all_teams", "teams.#"),
				),
			},
		},
	})
}
func testAccTeamsDataSourceExample() string {
	absPath, _ := filepath.Abs("../../examples/data-sources/allquiet_teams/data-source.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return string(dat)
}
