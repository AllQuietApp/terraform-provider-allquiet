resource "allquiet_integration" "twilio_call_routing" {
  display_name = "On-call inbound line"
  team_id      = allquiet_team.example.id
  type         = "Twilio"

  integration_settings = {
    twilio = {
      access_settings = {
        account_sid  = var.twilio_account_sid
        auth_token   = var.twilio_auth_token
        phone_number = var.twilio_phone_number
      }
      call_flow_config = {
        welcome_message = {
          type = "Audio"
          audio_file = {
            source = var.welcome_audio_path
          }
        }
        route = {
          team_id = allquiet_team.on_call.id
          voicemail_message = {
            type = "Text"
            text = "Please leave a message after the tone."
            locale = "en-US"
          }
        }
      }
    }
  }
}

variable "twilio_account_sid" {
  type      = string
  sensitive = true
}

variable "twilio_auth_token" {
  type      = string
  sensitive = true
}

variable "twilio_phone_number" {
  type = string
}

variable "welcome_audio_path" {
  type        = string
  description = "Local path to an MP3 or WAV welcome message"
}
