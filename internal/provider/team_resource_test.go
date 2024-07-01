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

func TestAccTeamWithMembersResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamResourceWithMembersConfig("Team One With Members", "billie@allquiet.app", "miley@allquiet.app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_team.test_with_members", "display_name", "Team One With Members"),
					resource.TestCheckResourceAttr("allquiet_team.test_with_members", "time_zone_id", "Europe/Berlin"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"allquiet_team.test_with_members",
						"members.*",
						map[string]string{
							"role":  "Member",
							"email": "billie@allquiet.app",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"allquiet_team.test_with_members",
						"members.*",
						map[string]string{
							"role":  "Administrator",
							"email": "miley@allquiet.app",
						},
					),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_team.test_with_members",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccTeamResourceWithMembersConfig("Team Two With Updated Members", "taylor@allquiet.app", "billie@allquiet.app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("allquiet_team.test_with_members", "display_name", "Team Two With Updated Members"),
					resource.TestCheckResourceAttr("allquiet_team.test_with_members", "time_zone_id", "Europe/Berlin"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"allquiet_team.test_with_members",
						"members.*",
						map[string]string{
							"role":  "Member",
							"email": "taylor@allquiet.app",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"allquiet_team.test_with_members",
						"members.*",
						map[string]string{
							"role":  "Administrator",
							"email": "billie@allquiet.app",
						},
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTeamExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamResourceExample(),
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_team.my_team_with_weekend_rotation",
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

func testAccTeamResourceWithMembersConfig(display_name string, member1_email string, member2_email string) string {
	return fmt.Sprintf(`
resource "allquiet_team" "test_with_members" {
  	display_name = %[1]q
	time_zone_id = "Europe/Berlin"
	members = [
		{
			email = %[2]q
			role = "Member"
		},
		{
			email = %[3]q
			role = "Administrator"
		}
	]
	tiers = [
		{
			auto_escalation_after_minutes = 5
			schedules = [
				{
					schedule_settings = {
						start = "09:00"
						end = "17:00"
						selected_days = ["mon", "tue", "wed", "thu", "fri"]
					},
					rotations = [
						{
							members = [
								{
									email = %[2]q
								},
								{
									email = %[3]q
								}
							]
						}
					]
				},
				{
					schedule_settings = {
						selected_days = ["sat", "sun"]
					}
					rotation_settings = {
						repeats = "weekly"
						starts_on_day_of_week = "mon"
						starts_on_date_of_month = 1
					},
					rotations = [
						{
							members = [
								{
									email = %[2]q
								}
							]
						},
						{
							members = [
								{
									email = %[3]q
								}
							]
						}
					]
				}
			]
		},
		{
			schedules = [
				{
					rotations = [
						{
							members = [
								{
									email = %[2]q
								},
								{
									email = %[3]q
								}
							]
						}
					]
				}
			]
		}
	]
}
`, display_name, member1_email, member2_email)
}

func testAccTeamResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_team/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return string(dat)

}
