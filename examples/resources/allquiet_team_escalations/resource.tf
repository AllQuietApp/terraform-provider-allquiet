

################################################################################
# Create users
################################################################################

resource "allquiet_user" "riemann" {
  display_name = "Riemann"
  email        = "acceptance-tests+riemann@allquiet.app"
}

resource "allquiet_user" "galois" {
  display_name = "Galois"
  email        = "acceptance-tests+galois@allquiet.app"
}

resource "allquiet_user" "gauss" {
  display_name = "Gauss"
  email        = "acceptance-tests+gauss@allquiet.app"
}

resource "allquiet_user" "kolmogorov" {
  display_name = "Kolmogorov"
  email        = "acceptance-tests+kolmogorov@allquiet.app"
}

################################################################################
# Example 1: My team with weekend rotation
################################################################################

resource "allquiet_team" "my_team_with_weekend_rotation" {
  display_name = "My team with weekend rotation"
  time_zone_id = "America/Los_Angeles"
}

resource "allquiet_team_membership" "my_team_with_weekend_rotation_riemann" {
  team_id = allquiet_team.my_team_with_weekend_rotation.id
  user_id = allquiet_user.riemann.id
  role    = "Administrator"
}

resource "allquiet_team_membership" "my_team_with_weekend_rotation_galois" {
  team_id = allquiet_team.my_team_with_weekend_rotation.id
  user_id = allquiet_user.galois.id
  role    = "Member"
}

resource "allquiet_team_membership" "my_team_with_weekend_rotation_gauss" {
  team_id = allquiet_team.my_team_with_weekend_rotation.id
  user_id = allquiet_user.gauss.id
  role    = "Member"
}

resource "allquiet_team_membership" "my_team_with_weekend_rotation_kolmogorov" {
  team_id = allquiet_team.my_team_with_weekend_rotation.id
  user_id = allquiet_user.kolmogorov.id
  role    = "Member"
}

resource "allquiet_team_escalations" "my_team_escalations_with_weekend_rotation" {
  team_id = allquiet_team.my_team_with_weekend_rotation.id
  escalation_tiers = [
    {
      auto_escalation_enabled       = true
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
                  team_membership_id = allquiet_team_membership.my_team_with_weekend_rotation_riemann.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_weekend_rotation_galois.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_weekend_rotation_gauss.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_weekend_rotation_kolmogorov.id
                }
              ]
            }
          ]
        },

        {
          display_name = "Weekend schedule"
          schedule_settings = {
            selected_days = ["sat", "sun"]
          }
          rotation_settings = {
            repeats               = "weekly"
            starts_on_day_of_week = "sat"
          }
          rotations = [
            {
              members = [
                {
                  team_membership_id = allquiet_team_membership.my_team_with_weekend_rotation_riemann.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_weekend_rotation_galois.id
                }
              ]
            },
            {
              members = [
                {
                  team_membership_id = allquiet_team_membership.my_team_with_weekend_rotation_gauss.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_weekend_rotation_kolmogorov.id
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}

################################################################################
# Example 2: My team with day and night rotation
################################################################################

resource "allquiet_team" "my_team_with_day_and_night_rotation" {
  display_name = "My team with day and night rotation"
  time_zone_id = "America/Los_Angeles"
}


resource "allquiet_team_membership" "my_team_with_day_and_night_rotation_riemann" {
  team_id = allquiet_team.my_team_with_day_and_night_rotation.id
  user_id = allquiet_user.riemann.id
  role    = "Administrator"
}

resource "allquiet_team_membership" "my_team_with_day_and_night_rotation_galois" {
  team_id = allquiet_team.my_team_with_day_and_night_rotation.id
  user_id = allquiet_user.galois.id
  role    = "Member"
}

resource "allquiet_team_membership" "my_team_with_day_and_night_rotation_gauss" {
  team_id = allquiet_team.my_team_with_day_and_night_rotation.id
  user_id = allquiet_user.gauss.id
  role    = "Member"
}

resource "allquiet_team_membership" "my_team_with_day_and_night_rotation_kolmogorov" {
  team_id = allquiet_team.my_team_with_day_and_night_rotation.id
  user_id = allquiet_user.kolmogorov.id
  role    = "Member"
}

resource "allquiet_team_escalations" "my_team_escalations_with_day_and_night_rotation" {
  team_id = allquiet_team.my_team_with_day_and_night_rotation.id
  escalation_tiers = [
    {
      auto_escalation_enabled       = true
      auto_escalation_after_minutes = 5
      schedules = [
        {
          display_name = "Working hours"
          schedule_settings = {
            start : "08:00"
            end : "20:00"
          }
          rotations = [
            {
              members = [
                {
                  team_membership_id = allquiet_team_membership.my_team_with_day_and_night_rotation_riemann.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_day_and_night_rotation_galois.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_day_and_night_rotation_gauss.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_day_and_night_rotation_kolmogorov.id
                }
              ]
            }
          ]
        },
        {
          display_name = "Night hours"
          schedule_settings = {
            start : "20:00"
            end : "08:00"
          }
          rotation_settings = {
            repeats               = "weekly"
            starts_on_day_of_week = "sat"
          }
          rotations = [
            {
              members = [
                {
                  team_membership_id = allquiet_team_membership.my_team_with_day_and_night_rotation_riemann.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_day_and_night_rotation_galois.id
                }
              ]
            },
            {
              members = [
                {
                  team_membership_id = allquiet_team_membership.my_team_with_day_and_night_rotation_gauss.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_day_and_night_rotation_kolmogorov.id
                }
              ]
            }
          ]
        }
      ]
    }
  ]

}

################################################################################
# Example 3: My team with hourly rotation
################################################################################

resource "allquiet_team" "my_team_with_hourly_rotation" {
  display_name = "My team with hourly rotation"
  time_zone_id = "America/Los_Angeles"
}


resource "allquiet_team_membership" "my_team_with_hourly_rotation_riemann" {
  team_id = allquiet_team.my_team_with_hourly_rotation.id
  user_id = allquiet_user.riemann.id
  role    = "Administrator"
}

resource "allquiet_team_membership" "my_team_with_hourly_rotation_galois" {
  team_id = allquiet_team.my_team_with_hourly_rotation.id
  user_id = allquiet_user.galois.id
  role    = "Member"
}

resource "allquiet_team_escalations" "my_team_escalations_with_hourly_rotation" {
  team_id = allquiet_team.my_team_with_hourly_rotation.id
  escalation_tiers = [
    {
      auto_escalation_enabled       = true
      auto_escalation_after_minutes = 5
      schedules = [
        {
          rotation_settings = {
            repeats             = "custom"
            custom_repeat_unit  = "hours"
            custom_repeat_value = 6
            starts_on_time      = "09:00"
            effective_from      = "2024-07-10"
          }
          rotations = [
            {
              members = [
                {
                  team_membership_id = allquiet_team_membership.my_team_with_hourly_rotation_riemann.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_hourly_rotation_galois.id
                },
              ]
            }
          ]
        },
      ]
    }
  ]
}

################################################################################
# Example 4: My team with auto rotation
################################################################################

resource "allquiet_team" "my_team_with_auto_rotation" {
  display_name = "My team with auto rotation"
  time_zone_id = "America/Los_Angeles"
}


resource "allquiet_team_membership" "my_team_with_auto_rotation_riemann" {
  team_id = allquiet_team.my_team_with_auto_rotation.id
  user_id = allquiet_user.riemann.id
  role    = "Administrator"
}

resource "allquiet_team_membership" "my_team_with_auto_rotation_galois" {
  team_id = allquiet_team.my_team_with_auto_rotation.id
  user_id = allquiet_user.galois.id
  role    = "Member"
}

resource "allquiet_team_membership" "my_team_with_auto_rotation_gauss" {
  team_id = allquiet_team.my_team_with_auto_rotation.id
  user_id = allquiet_user.gauss.id
  role    = "Member"
}

resource "allquiet_team_membership" "my_team_with_auto_rotation_kolmogorov" {
  team_id = allquiet_team.my_team_with_auto_rotation.id
  user_id = allquiet_user.kolmogorov.id
  role    = "Member"
}

resource "allquiet_team_escalations" "my_team_escalations_with_auto_rotation" {
  team_id = allquiet_team.my_team_with_auto_rotation.id
  escalation_tiers = [
    {
      schedules = [
        {
          rotation_settings = {
            repeats               = "weekly"
            starts_on_day_of_week = "sat"
            rotation_mode         = "auto"
            auto_rotation_size    = 3
          }
          rotations = [
            {
              members = [
                {
                  team_membership_id = allquiet_team_membership.my_team_with_auto_rotation_riemann.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_auto_rotation_galois.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_auto_rotation_gauss.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_auto_rotation_kolmogorov.id
                },
              ]
            }
          ]
        },
      ]
    }
  ]
}

################################################################################
# Example 5: My team with repeating tier
################################################################################

resource "allquiet_team" "my_team_with_repeating_tier" {
  display_name = "My team with repeating tier"
  time_zone_id = "America/Los_Angeles"
}

resource "allquiet_team_membership" "my_team_with_repeating_tier_riemann" {
  team_id = allquiet_team.my_team_with_repeating_tier.id
  user_id = allquiet_user.riemann.id
  role    = "Administrator"
}

resource "allquiet_team_membership" "my_team_with_repeating_tier_galois" {
  team_id = allquiet_team.my_team_with_repeating_tier.id
  user_id = allquiet_user.galois.id
  role    = "Member"
}

resource "allquiet_team_membership" "my_team_with_repeating_tier_gauss" {
  team_id = allquiet_team.my_team_with_repeating_tier.id
  user_id = allquiet_user.gauss.id
  role    = "Member"
}

resource "allquiet_team_membership" "my_team_with_repeating_tier_kolmogorov" {
  team_id = allquiet_team.my_team_with_repeating_tier.id
  user_id = allquiet_user.kolmogorov.id
  role    = "Member"
}

resource "allquiet_team_escalations" "my_team_escalations_with_repeating_tier" {
  team_id = allquiet_team.my_team_with_repeating_tier.id
  escalation_tiers = [
    {
      auto_escalation_enabled       = true
      auto_escalation_after_minutes = 5
      repeats                       = 2
      repeats_after_minutes         = null
      auto_escalation_severities    = ["Critical"]
      schedules = [
        {
          rotations = [
            {
              members = [
                {
                  team_membership_id = allquiet_team_membership.my_team_with_repeating_tier_riemann.id
                },
                {
                  team_membership_id = allquiet_team_membership.my_team_with_repeating_tier_galois.id
                },
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
                  team_membership_id = allquiet_team_membership.my_team_with_repeating_tier_gauss.id
                },
              ]
            }
          ]
        }
      ]
    }
  ]
}

################################################################################
# Example 6: Root Team with Auto Assign To Teams
################################################################################

resource "allquiet_team" "root_team" {
  display_name = "Root team"
  time_zone_id = "America/Los_Angeles"
}

resource "allquiet_team" "second_level_support_team" {
  display_name = "2nd level support team"
  time_zone_id = "America/Los_Angeles"
}

resource "allquiet_team_membership" "root_team_riemann" {
  team_id = allquiet_team.root_team.id
  user_id = allquiet_user.riemann.id
  role    = "Member"
}

resource "allquiet_team_escalations" "root_team_escalations" {
  team_id = allquiet_team.root_team.id
  escalation_tiers = [
    {
      auto_escalation_enabled       = true
      auto_escalation_after_minutes = 5
      auto_escalation_severities    = ["Critical", "Warning"]
      auto_escalation_time_filters = [
        {
          selected_days = ["mon", "tue", "wed", "thu", "fri"]
          from          = "08:00"
          until         = "17:00"
        }
      ]
      auto_assign_to_teams            = [allquiet_team.second_level_support_team.id]
      auto_assign_to_teams_severities = ["Critical"]
      auto_assign_to_teams_time_filters = [
        {
          selected_days = ["mon", "tue", "wed", "thu", "fri"]
          from          = "08:00"
          until         = "17:00"
        }
      ]
      schedules = [
        {
          rotations = [
            {
              members = [
                {
                  team_membership_id = allquiet_team_membership.root_team_riemann.id
                },
              ]
            }
          ]
        }
      ]
    }
  ]
}

################################################################################
# Example 7: Team with weekly schedules
################################################################################

resource "allquiet_team" "my_team_with_weekly_schedules" {
  display_name = "My team with weekly schedules"
  time_zone_id = "America/Los_Angeles"
}


resource "allquiet_team_membership" "my_team_with_weekly_schedules_riemann" {
  team_id = allquiet_team.my_team_with_weekly_schedules.id
  user_id = allquiet_user.riemann.id
  role    = "Member"
}

resource "allquiet_team_escalations" "my_team_escalations_with_weekly_schedules" {
  team_id = allquiet_team.my_team_with_weekly_schedules.id
  escalation_tiers = [
    {
      schedules = [
        {
          schedule_settings = {
            weekly_schedules = [
              {
                selected_days = ["mon", "tue", "wed", "thu", "fri"]
                from          = "06:00"
                until         = "18:00"
              },
              {
                selected_days = ["sat", "sun"]
                from          = "10:00"
                until         = "14:00"
              }
            ]
          }
          rotations = [
            {
              members = [
                {
                  team_membership_id = allquiet_team_membership.my_team_with_weekly_schedules_riemann.id
                },
              ]
            }
          ]
        }
      ]
    }
  ]
}


################################################################################
# Example 8: Team with all tiers repeating
################################################################################

resource "allquiet_team" "my_team_with_all_tiers_repeating" {
  display_name = "My team with all tiers repeating"
  time_zone_id = "America/Los_Angeles"
}

resource "allquiet_team_membership" "my_team_with_all_tiers_repeating_riemann" {
  team_id = allquiet_team.my_team_with_all_tiers_repeating.id
  user_id = allquiet_user.riemann.id
  role    = "Member"
}

resource "allquiet_team_membership" "my_team_with_all_tiers_repeating_galois" {
  team_id = allquiet_team.my_team_with_all_tiers_repeating.id
  user_id = allquiet_user.galois.id
  role    = "Member"
}

resource "allquiet_team_escalations" "my_team_escalations_with_all_tiers_repeating" {
  team_id = allquiet_team.my_team_with_all_tiers_repeating.id
  tier_settings = {
    repeats               = 2
    repeats_after_minutes = 5
    repeats_stop_mode     = "resolved"
  }
  escalation_tiers = [
    {
      auto_escalation_enabled       = true
      auto_escalation_after_minutes = 5
      auto_escalation_stop_mode     = "acknowledged"
      repeats                       = 2
      repeats_after_minutes         = 1
      repeats_stop_mode             = "resolved"
      schedules = [
        {
          rotations = [
            {
              members = [
                {
                  team_membership_id = allquiet_team_membership.my_team_with_all_tiers_repeating_riemann.id
                },
              ]
            }
          ]
        }
      ]
    }
  ]
}










