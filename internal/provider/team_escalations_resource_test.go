// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"path/filepath"
	"testing"

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
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTeamEscalationsResourceConfigCreate() string {
	return `
	  resource "allquiet_user" "galois" {
		display_name = "Galois"
		email        = "galois@acme.com"
	  }
	  
	  resource "allquiet_user" "kolmogorov" {
		display_name = "Kolmogorov"
		email        = "kolmogorov@acme.com"
	  }

	  resource "allquiet_team" "my_team" {
		display_name = "My team with weekend rotation"
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
	  
	  resource "allquiet_team_escalations" "my_team" {
		team_id = allquiet_team.my_team.id
		escalation_tiers = [
		  {
			auto_escalation_after_minutes = 5
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
	  
`
}

func testAccTeamEscalationsResourceConfigUpdate() string {
	return `
	  
	  resource "allquiet_user" "galois" {
		display_name = "Galois"
		email        = "galois@acme.com"
	  }

      resource "allquiet_user" "gauss" {
		display_name = "Gauss"
		email        = "gauss@acme.com"
	  }
	  
	  resource "allquiet_team" "my_team" {
		display_name = "My team with weekend rotation"
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
	  
	  resource "allquiet_team_escalations" "my_team" {
		team_id = allquiet_team.my_team.id
		escalation_tiers = [
		  {
			auto_escalation_after_minutes = 5
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
			  }
			]
		  }
		]
	  }
	  
`
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

	return string(dat)

}
