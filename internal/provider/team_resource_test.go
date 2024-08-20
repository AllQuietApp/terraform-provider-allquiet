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

func TestAccTeamResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamResourceConfig("Team One", "mon"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_team.test", "display_name", "Team One"),
					resource.TestCheckResourceAttr("allquiet_team.test", "time_zone_id", "Europe/Berlin"),
					resource.TestCheckResourceAttr("allquiet_team.test", "incident_engagement_report_settings.day_of_week", "mon"),
					resource.TestCheckResourceAttr("allquiet_team.test", "incident_engagement_report_settings.time", "09:00"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_team.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccTeamResourceConfig("Team Two", "tue"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_team.test", "display_name", "Team Two"),
					resource.TestCheckResourceAttr("allquiet_team.test", "time_zone_id", "Europe/Berlin"),
					resource.TestCheckResourceAttr("allquiet_team.test", "incident_engagement_report_settings.day_of_week", "tue"),
					resource.TestCheckResourceAttr("allquiet_team.test", "incident_engagement_report_settings.time", "09:00"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTeamResourceConfig(display_name string, day_of_week string) string {
	return fmt.Sprintf(`
resource "allquiet_team" "test" {
  display_name = %[1]q
  time_zone_id = "Europe/Berlin"
  incident_engagement_report_settings = {
    	day_of_week = %[2]q
		time = "09:00"
	}
}
`, display_name, day_of_week)
}

func TestAccTeamExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_team.my_team", "display_name", "My Team"),
					resource.TestCheckResourceAttr("allquiet_team.my_team", "time_zone_id", "America/Los Angeles"),
					resource.TestCheckResourceAttr("allquiet_team.my_team", "incident_engagement_report_settings.day_of_week", "tue"),
					resource.TestCheckResourceAttr("allquiet_team.my_team", "incident_engagement_report_settings.time", "09:00"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_team.my_team",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccTeamResourceExample(),
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTeamResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_team/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return string(dat)

}
