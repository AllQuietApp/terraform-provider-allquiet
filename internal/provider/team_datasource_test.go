// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTeamDataSource(t *testing.T) {
	uid := uuid.New().String()
	teamName := fmt.Sprintf("team+%s", uid)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamDataSourceConfig(teamName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.allquiet_team.test_by_display_name", "display_name", teamName),
				),
			},
		},
	})
}

func testAccTeamDataSourceConfig(teamName string) string {
	return fmt.Sprintf(`

		resource "allquiet_team" "test" {
			display_name = %[1]q
		}

		data "allquiet_team" "test_by_display_name" {
			display_name = %[1]q
			depends_on = [allquiet_team.test]
		}

	`, teamName)
}
