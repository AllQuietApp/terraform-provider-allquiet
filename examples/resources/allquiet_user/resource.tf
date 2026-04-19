resource "allquiet_user" "millie_brown" {
  display_name = "Millie Bobby Brown"
  email        = "acceptance-tests+millie@allquiet.app"
  phone_number = "+12035479055"
}

resource "allquiet_user" "taylor" {
  display_name = "Taylor Swift"
  email        = "acceptance-tests+taylor@allquiet.app"
}

# Incident notification settings have moved to the
# `allquiet_user_incident_notification_settings` resource. See its example for usage.
