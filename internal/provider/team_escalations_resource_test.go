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

func TestAccTeamEscalationsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamEscalationsResourceConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.schedules.0.rotations.0.members.0.team_membership_id"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.auto_escalation_severities.0"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.auto_escalation_severities.1"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.auto_escalation_time_filters.0.selected_days.0"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.auto_escalation_time_filters.0.from"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.auto_escalation_time_filters.0.until"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.auto_assign_to_teams.0"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.auto_assign_to_teams_severities.0"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.auto_assign_to_teams_time_filters.0.selected_days.0"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.auto_assign_to_teams_time_filters.0.from"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.auto_assign_to_teams_time_filters.0.until"),
					resource.TestCheckResourceAttr("allquiet_team_escalations.my_team_with_round_robin", "escalation_tiers.0.schedules.0.round_robin_settings.round_robin_size", "3"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_team_escalations.my_team",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccTeamEscalationsResourceConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.schedules.0.rotations.0.members.0.team_membership_id"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.auto_escalation_severities.0"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.auto_escalation_severities.1"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.repeats"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team", "escalation_tiers.0.repeats_after_minutes"),
					resource.TestCheckResourceAttr("allquiet_team_escalations.my_team_with_round_robin", "escalation_tiers.0.schedules.0.round_robin_settings.round_robin_size", "5"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTeamEscalationsResourceConfigCreate() string {
	return fmt.Sprintf(`
	  resource "allquiet_user" "galois" {
		display_name = "Galois"
		email        = "acceptance-tests+galois+%s@allquiet.app"
	  }
	  
	  resource "allquiet_user" "kolmogorov" {
		display_name = "Kolmogorov"
		email        = "acceptance-tests+kolmogorov+%s@allquiet.app"
	  }

	  resource "allquiet_team" "my_team" {
		display_name = "My team with weekend rotation"
		time_zone_id = "America/Los_Angeles"
	  }
	  resource "allquiet_team" "my_team2" {
		display_name = "My team with weekend rotation"
		time_zone_id = "America/Los_Angeles"
	  }

	  resource "allquiet_team" "my_team_round_robin" {
		display_name = "My team with round robin"
		time_zone_id = "America/Los_Angeles"
	  }

	  resource "allquiet_team" "engineering" {
		display_name = "Engineering"
		time_zone_id = "America/Los_Angeles"
	  }
	  
	  resource "allquiet_team_membership" "my_team_galois" {
		team_id = allquiet_team.my_team.id
		user_id = allquiet_user.galois.id
		role    = "Administrator"
	  }
	  
	  resource "allquiet_team_membership" "my_team_kolmogorov" {
		team_id = allquiet_team.my_team.id
		user_id = allquiet_user.kolmogorov.id
		role    = "Member"
	  }

	  resource "allquiet_team_membership" "my_team_round_robin_galois" {
		team_id = allquiet_team.my_team_round_robin.id
		user_id = allquiet_user.galois.id
		role    = "Administrator"
	  }
	  
	  resource "allquiet_team_membership" "my_team_round_robin_kolmogorov" {
		team_id = allquiet_team.my_team_round_robin.id
		user_id = allquiet_user.kolmogorov.id
		role    = "Member"
	  }

	  resource "allquiet_team" "my_team_with_empty_weekly_schedules" {
		display_name = "My team with empty weekly schedules"
		time_zone_id = "America/Los_Angeles"
	  }

	  resource "allquiet_team_membership" "my_team_with_empty_weekly_schedules_kolmogorov" {
		team_id = allquiet_team.my_team_with_empty_weekly_schedules.id
		user_id = allquiet_user.kolmogorov.id
		role    = "Member"
	  }
	  
	  resource "allquiet_team_membership" "my_team_with_empty_weekly_schedules_galois" {
		team_id = allquiet_team.my_team_with_empty_weekly_schedules.id
		user_id = allquiet_user.galois.id
		role    = "Member"
	  }

	  resource "allquiet_team_escalations" "my_team_with_empty_weekly_schedules" {
		escalation_tiers = [
			{
			auto_assign_to_teams_severities   = ["Critical", "Warning"]
			auto_assign_to_teams_time_filters = null
			auto_escalation_after_minutes     = 10
			auto_escalation_enabled           = false
			auto_escalation_severities        = null
			auto_escalation_stop_mode         = "acknowledged"
			auto_escalation_time_filters      = null
			repeats                           = 1
			repeats_after_minutes             = 10
			repeats_stop_mode                 = "acknowledged"
			schedules = [
				{
				rotation_settings = {
					auto_rotation_size      = null
					custom_repeat_unit      = null
					custom_repeat_value     = null
					effective_from          = "2025-06-11"
					repeats                 = "weekly"
					rotation_mode           = "explicit"
					starts_on_date_of_month = null
					starts_on_day_of_week   = "wed"
					starts_on_time          = "14:00"
				}
				rotations = [
					{
					members = [
						{
						team_membership_id = allquiet_team_membership.my_team_with_empty_weekly_schedules_galois.id
						},
					]
					},
					{
					members = [
						{
						team_membership_id = allquiet_team_membership.my_team_with_empty_weekly_schedules_kolmogorov.id
						},
					]
					},
				]
				schedule_settings = {
					weekly_schedules = null
				}
				},
			]
			},
		]
		team_id = allquiet_team.my_team_with_empty_weekly_schedules.id
		tier_settings = null
		}


	  resource "allquiet_team_escalations" "my_team" {
		team_id = allquiet_team.my_team.id
		escalation_tiers = [
		  {
			auto_escalation_enabled = true
			auto_escalation_after_minutes = 5
			auto_escalation_severities = ["Critical", "Warning"]
			auto_escalation_time_filters = [
				{
					selected_days = ["mon", "tue", "wed", "thu", "fri"],
					from = "06:00",
					until = "18:00"
				},
				{
					selected_days = ["sat", "sun"],
					from = "10:00",
					until = "16:00"
				}
			],
			auto_assign_to_teams = [allquiet_team.engineering.id]
			auto_assign_to_teams_severities = ["Critical", "Warning"]
			auto_assign_to_teams_time_filters = [
				{
					selected_days = ["mon", "tue", "wed", "thu", "fri"],
					from = "06:00",
					until = "18:00"
				}
			]
			repeats = 1
			repeats_after_minutes = 5
			schedules = [
			  {
				schedule_settings = {
				  selected_days = ["mon", "tue", "wed", "thu", "fri"]
				}
				rotations = [
				  {
					members = [
					  {
						team_membership_id = allquiet_team_membership.my_team_galois.id
					  },
					  {
						team_membership_id = allquiet_team_membership.my_team_kolmogorov.id
					  }
					]
				  }
				]
			  }
			]
		  }
		]
	  }

	  resource "allquiet_team_escalations" "my_team_with_round_robin" {
		team_id = allquiet_team.my_team_round_robin.id
		escalation_tiers = [
		  {
			auto_escalation_enabled = true
			auto_escalation_after_minutes = 5
			auto_escalation_severities = ["Critical", "Warning"]
			schedules = [
			  {
				schedule_settings = {
				  selected_days = ["mon", "tue", "wed", "thu", "fri"]
				}
				round_robin_settings = {
				  round_robin_size = 3
				}
				rotations = [
				  {
					members = [
					  {
						team_membership_id = allquiet_team_membership.my_team_round_robin_galois.id
					  },
					  {
						team_membership_id = allquiet_team_membership.my_team_round_robin_kolmogorov.id
					  }
					]
				  }
				]
			  }
			]
		  }
		]
	  }


	  resource "allquiet_team_escalations" "my_team_with_empty_members" {
		team_id = allquiet_team.my_team2.id
		escalation_tiers = [
		  {
			auto_escalation_after_minutes = 5
			auto_escalation_severities = ["Critical", "Warning"]
			auto_escalation_stop_mode = "acknowledged"
			repeats = 1
			repeats_after_minutes = 5
			repeats_stop_mode = "resolved"
			schedules = [
			  {
				schedule_settings = {
				  selected_days = ["mon", "tue", "wed", "thu", "fri"]
				}
				rotations = [
				  {
					members = []
				  }
				]
			  },
			  {
				schedule_settings = {
				  weekly_schedules = [
					{
						selected_days = ["mon", "tue", "wed", "thu", "fri"]
						from = "06:00"
						until = "18:00"
					},
					{
						selected_days = ["sat", "sun"]
						from = "10:00"
						until = "16:00"
					}
				  ]
				}
				rotations = [
				  {
					members = []
				  }
				]
			  }
			]
		  }
		]
	  }
	  
`, uuid.New().String(), uuid.New().String())
}

func testAccTeamEscalationsResourceConfigUpdate() string {
	return fmt.Sprintf(`
	  
	  resource "allquiet_user" "galois" {
		display_name = "Galois"
		email        = "acceptance-tests+galois+%s@allquiet.app"
	  }

      resource "allquiet_user" "gauss" {
		display_name = "Gauss"
		email        = "acceptance-tests+gauss+%s@allquiet.app"
	  }
	  
	  resource "allquiet_team" "my_team" {
		display_name = "My team with weekend rotation"
		time_zone_id = "America/Los_Angeles"
	  }
	  resource "allquiet_team" "my_team2" {
		display_name = "My team with weekend rotation"
		time_zone_id = "America/Los_Angeles"
	  }

	  resource "allquiet_team" "my_team_round_robin" {
		display_name = "My team with round robin"
		time_zone_id = "America/Los_Angeles"
	  }
	  
	  resource "allquiet_team_membership" "my_team_galois" {
		team_id = allquiet_team.my_team.id
		user_id = allquiet_user.galois.id
		role    = "Administrator"
	  }
	  
	  resource "allquiet_team_membership" "my_team_gauss" {
		team_id = allquiet_team.my_team.id
		user_id = allquiet_user.gauss.id
		role    = "Member"
	  }
	  
	  resource "allquiet_team_membership" "my_team2_gauss" {
		team_id = allquiet_team.my_team2.id
		user_id = allquiet_user.gauss.id
		role    = "Member"
	  }

	  resource "allquiet_team_membership" "my_team_round_robin_galois" {
		team_id = allquiet_team.my_team_round_robin.id
		user_id = allquiet_user.galois.id
		role    = "Administrator"
	  }
	  
	  resource "allquiet_team_membership" "my_team_round_robin_gauss" {
		team_id = allquiet_team.my_team_round_robin.id
		user_id = allquiet_user.gauss.id
		role    = "Member"
	  }

	  resource "allquiet_team_escalations" "my_team" {
		team_id = allquiet_team.my_team.id
		escalation_tiers = [
		  {
			auto_escalation_after_minutes = 5
			auto_escalation_severities = ["Critical", "Warning"]
			auto_escalation_stop_mode = "acknowledged"
			repeats = 1
			repeats_after_minutes = 5
			repeats_stop_mode = "resolved"
			schedules = [
			  {
				schedule_settings = {
				  selected_days = ["mon", "tue", "wed", "thu", "fri"]
				}
				rotations = [
				  {
					members = [
					  {
						team_membership_id = allquiet_team_membership.my_team_galois.id
					  },
					  {
						team_membership_id = allquiet_team_membership.my_team_gauss.id
					  }
					]
				  }
				]
			  },
			  {
				schedule_settings = {
				  weekly_schedules = [
					{
						selected_days = ["mon", "tue", "wed", "thu", "fri"]
						from = "06:00"
						until = "18:00"
					},
					{
						selected_days = ["sat", "sun"]
						from = "10:00"
						until = "16:00"
					}
				  ]
				}
				rotations = [
				  {
					members = [
					  {
						team_membership_id = allquiet_team_membership.my_team_galois.id
					  }
					]
				  }
				]
			  }
			]
		  }
		]
	  }

	  resource "allquiet_team_escalations" "my_team_with_round_robin" {
		team_id = allquiet_team.my_team_round_robin.id
		escalation_tiers = [
		  {
			auto_escalation_enabled = true
			auto_escalation_after_minutes = 5
			auto_escalation_severities = ["Critical", "Warning"]
			schedules = [
			  {
				schedule_settings = {
				  selected_days = ["mon", "tue", "wed", "thu", "fri"]
				}
				round_robin_settings = {
				  round_robin_size = 5
				}
				rotations = [
				  {
					members = [
					  {
						team_membership_id = allquiet_team_membership.my_team_round_robin_galois.id
					  },
					  {
						team_membership_id = allquiet_team_membership.my_team_round_robin_gauss.id
					  }
					]
				  }
				]
			  }
			]
		  }
		]
	  }

	  resource "allquiet_team_escalations" "my_team_with_empty_members" {
		team_id = allquiet_team.my_team2.id
		escalation_tiers = [
		  {
			auto_escalation_after_minutes = 5
			auto_escalation_severities = ["Critical", "Warning"]
			auto_escalation_stop_mode = "acknowledged"
			repeats = 1
			repeats_after_minutes = 5
			repeats_stop_mode = "resolved"
			schedules = [
			  {
				display_name = "Working weekdays schedule"
				schedule_settings = {
				  selected_days = ["mon", "tue", "wed", "thu", "fri"]
				}
				rotations = [
				  {
					members = [
					  {
						team_membership_id = allquiet_team_membership.my_team2_gauss.id
					  }
					]
				  }
				]
			  },
			  {
				schedule_settings = {
				  weekly_schedules = [
					{
						selected_days = ["mon", "tue", "wed", "thu", "fri"]
						from = "06:00"
						until = "18:00"
					},
					{
						selected_days = ["sat", "sun"]
						from = "10:00"
						until = "16:00"
					}
				  ]
				}
				rotations = [
				  {
					members = [
					  {
						team_membership_id = allquiet_team_membership.my_team2_gauss.id
					  }
					]
				  }
				]
			  }
			]
		  }
		]
	  }
	  
`, uuid.New().String(), uuid.New().String())
}

func TestAccTeamEscalationsExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamEscalationsResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team_escalations_with_day_and_night_rotation", "team_id"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team_escalations_with_hourly_rotation", "team_id"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team_escalations_with_weekend_rotation", "team_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_team_escalations.my_team_escalations_with_day_and_night_rotation",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccTeamEscalationsResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team_escalations_with_day_and_night_rotation", "team_id"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team_escalations_with_hourly_rotation", "team_id"),
					resource.TestCheckResourceAttrSet("allquiet_team_escalations.my_team_escalations_with_weekend_rotation", "team_id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTeamEscalationsResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_team_escalations/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return RandomizeExample(string(dat))

}
