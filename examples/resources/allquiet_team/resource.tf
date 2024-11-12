
resource "allquiet_team" "my_team" {
  display_name = "My Team"
  time_zone_id = "America/Los_Angeles"
  incident_engagement_report_settings = {
    day_of_week = "mon"
    time        = "09:00"
  }
}