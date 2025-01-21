resource "allquiet_user" "millie_brown" {
  display_name = "Millie Bobby Brown"
  email        = "acceptance-tests+millie@allquiet.app"
  phone_number = "+12035479055"
}

resource "allquiet_user" "taylor" {
  display_name = "Taylor Swift"
  email        = "acceptance-tests+taylor@allquiet.app"
  incident_notification_settings = {
    should_send_sms  = true
    delay_in_min_sms = 5
    severities_sms   = ["Critical"]

    should_call_voice  = false
    delay_in_min_voice = 0
    severities_voice   = []

    should_send_push  = true
    delay_in_min_push = 0
    severities_push   = ["Critical", "Warning"]

    should_send_email  = true
    delay_in_min_email = 0
    severities_email   = ["Critical", "Warning", "Minor"]
  }
}

resource "allquiet_user" "billie_eilish" {
  display_name = "Billie Eilish"
  email        = "acceptance-tests+billie@allquiet.app"
  incident_notification_settings = {
    should_send_sms      = true
    delay_in_min_sms     = 5
    severities_sms       = ["Critical"]
    disabled_intents_sms = ["Resolved"]

    should_call_voice      = false
    delay_in_min_voice     = 0
    severities_voice       = []
    disabled_intents_voice = ["Resolved"]

    should_send_push  = true
    delay_in_min_push = 0
    severities_push   = ["Critical", "Warning"]

    should_send_email  = true
    delay_in_min_email = 0
    severities_email   = ["Critical", "Warning", "Minor"]
  }
}