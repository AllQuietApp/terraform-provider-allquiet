
resource "allquiet_team" "my_team_with_weekend_rotation" {
  display_name = "My team with weekly weekend rotation"
  time_zone_id = "America/Los_Angeles"
  incident_engagement_report_settings = {
    day_of_week = "mon"
    time        = "09:00"
  }
  members = [
    {
      email = "cantor@acme.com"
      role  = "Administrator"
    },
    {
      email = "riemann@acme.com"
      role  = "Member"
    },
    {
      email = "galois@acme.com"
      role  = "Member"
    },
    {
      email = "gauss@acme.com"
      role  = "Member"
    },
    {
      email = "kolmogorov@acme.com"
      role  = "Member"
    }
  ]
  tiers = [
    {
      auto_escalation_after_minutes = 5
      schedules = [
        {
          schedule_settings = {
            start : "08:00"
            end : "20:00"
          }
          rotations = [
            {
              members = [
                {
                  email = "riemann@acme.com"
                },
                {
                  email = "galois@acme.com"
                },
                {
                  email = "gauss@acme.com"
                },
                {
                  email = "kolmogorov@acme.com"
                }
              ]
            }
          ]
        },
        {
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
                  email = "riemann@acme.com"
                },
                {
                  email = "galois@acme.com"
                }
              ]
            },
            {
              members = [
                {
                  email = "gauss@acme.com"
                },
                {
                  email = "kolmogorov@acme.com"
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}


resource "allquiet_team" "my_team_with_day_and_night_rotation" {
  display_name = "My team with day and night rotation"
  time_zone_id = "America/Los_Angeles"
  incident_engagement_report_settings = {
    day_of_week = "mon"
    time        = "09:00"
  }
  members = [
    {
      email = "cantor@acme.com"
      role  = "Administrator"
    },
    {
      email = "riemann@acme.com"
      role  = "Member"
    },
    {
      email = "galois@acme.com"
      role  = "Member"
    },
    {
      email = "gauss@acme.com"
      role  = "Member"
    },
    {
      email = "kolmogorov@acme.com"
      role  = "Member"
    }
  ]
  tiers = [
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
                  email = "riemann@acme.com"
                },
                {
                  email = "galois@acme.com"
                },
                {
                  email = "gauss@acme.com"
                },
                {
                  email = "kolmogorov@acme.com"
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
            repeats               = "weekly"
            starts_on_day_of_week = "sat"
          }
          rotations = [
            {
              members = [
                {
                  email = "riemann@acme.com"
                },
                {
                  email = "galois@acme.com"
                }
              ]
            },
            {
              members = [
                {
                  email = "gauss@acme.com"
                },
                {
                  email = "kolmogorov@acme.com"
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}

resource "allquiet_team" "my_team_with_hourly_rotation" {
  display_name = "My team with hourly rotation"
  time_zone_id = "America/Los_Angeles"
  members = [
    {
      email = "cantor@acme.com"
      role  = "Administrator"
    },
    {
      email = "riemann@acme.com"
      role  = "Member"
    },
    {
      email = "galois@acme.com"
      role  = "Member"
    }
  ]
  tiers = [
    {
      auto_escalation_after_minutes = 5
      schedules = [
        {
          rotation_settings = {
            repeats             = "custom"
            custom_repeat_unit  = "hours"
            custom_repeat_value = 6
            starts_on_time      = "00:00"
          }
          rotations = [
            {
              members = [
                {
                  email = "riemann@acme.com"
                },
                {
                  email = "galois@acme.com"
                },
              ]
            }
          ]
        },
      ]
    }
  ]
}
