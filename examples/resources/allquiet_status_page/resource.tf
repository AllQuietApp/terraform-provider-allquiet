resource "allquiet_service" "payment_api" {
  display_name       = "Payment Provider"
  public_title       = "Payment Provider"
  public_description = "Payment Provider and integrations"
}

resource "allquiet_service" "chat_gpt" {
  display_name       = "Chat GPT"
  public_title       = "Chat GPT"
  public_description = "Chat GPT and integrations"
}

resource "allquiet_service" "shipping_api" {
  display_name       = "Shipping API"
  public_title       = "Shipping API"
  public_description = "Shipping APIs and integrations"
}


resource "allquiet_status_page" "public_status_page" {
  slug                                 = "public-status-page-test"
  display_name                         = "Public Status Page"
  public_title                         = "Public Status Page"
  public_description                   = "Here, we'll inform you about the status of our services"
  history_in_days                      = 30
  disable_public_subscription          = false
  public_company_url                   = "https://www.allquiet.app"
  public_company_name                  = "AllQuiet"
  public_support_url                   = "https://www.allquiet.app/support"
  public_support_email                 = "support@allquiet.app"
  public_severity_mapping_minor        = "Minor"
  public_severity_mapping_warning      = "Warning"
  public_severity_mapping_critical     = "Critical"
  banner_background_color              = "#000000"
  banner_background_color_dark_mode    = "#447788"
  banner_text_color                    = "#ffffff"
  banner_text_color_dark_mode          = "#ffffff"
  body_background_color                = "#f5f5f5"
  body_background_color_dark_mode      = "#111111"
  secondary_background_color           = "#e0e0e0"
  secondary_background_color_dark_mode = "#222222"
  primary_text_color                   = "#111111"
  primary_text_color_dark_mode         = "#f5f5f5"
  secondary_text_color                 = "#333333"
  secondary_text_color_dark_mode       = "#cccccc"
  button_background_color              = "#0066cc"
  button_background_color_dark_mode    = "#3399ff"
  button_text_color                    = "#ffffff"
  button_text_color_dark_mode          = "#ffffff"

  time_zone_id = "Europe/Amsterdam"
  service_groups = [
    {
      public_display_name = "External Services"
      public_description  = "External services and integrations"
      services = [
        allquiet_service.payment_api.id,
        allquiet_service.chat_gpt.id
      ]
    },
    {
      public_display_name = "Internal Services"
      public_description  = "Internal services and integrations"
      services = [
        allquiet_service.shipping_api.id
      ]
    }
  ]
}



resource "allquiet_status_page" "public_status_page_with_custom_host_settings" {
  custom_host_settings = {
    host = "status-page-test-resource.allquiet.com"
  }
  display_name                         = "Public Status Page with Custom Host Settings"
  public_title                         = "Public Status Page with Custom Host Settings"
  public_description                   = "Here, we'll inform you about the status of our services"
  history_in_days                      = 30
  disable_public_subscription          = false
  public_company_url                   = "https://www.allquiet.app"
  public_company_name                  = "AllQuiet"
  public_support_url                   = "https://www.allquiet.app/support"
  public_support_email                 = "support@allquiet.app"
  public_severity_mapping_minor        = "Minor"
  public_severity_mapping_warning      = "Warning"
  public_severity_mapping_critical     = "Critical"
  banner_background_color              = "#000000"
  banner_background_color_dark_mode    = "#447788"
  banner_text_color                    = "#ffffff"
  banner_text_color_dark_mode          = "#ffffff"
  body_background_color                = "#ffffff"
  body_background_color_dark_mode      = "#000000"
  secondary_background_color           = "#e0e0e0"
  secondary_background_color_dark_mode = "#222222"
  primary_text_color                   = "#111111"
  primary_text_color_dark_mode         = "#f5f5f5"
  secondary_text_color                 = "#333333"
  secondary_text_color_dark_mode       = "#cccccc"
  button_background_color              = "#0066cc"
  button_background_color_dark_mode    = "#3399ff"
  button_text_color                    = "#ffffff"
  button_text_color_dark_mode          = "#ffffff"
  decimal_places                       = 2

  time_zone_id = "Europe/Amsterdam"
  service_groups = [
    {
      public_display_name = "External Services"
      public_description  = "External services and integrations"
      services = [
        allquiet_service.payment_api.id,
        allquiet_service.chat_gpt.id
      ]
    },
    {
      public_display_name = "Internal Services"
      public_description  = "Internal services and integrations"
      services = [
        allquiet_service.shipping_api.id
      ]
    }
  ]
}

resource "allquiet_status_page" "private_status_page" {
  slug                                 = "private-status-page-test"
  display_name                         = "Private Status Page"
  public_title                         = "Private Status Page"
  public_description                   = "Internal status page for team members only"
  history_in_days                      = 30
  disable_public_subscription          = false
  disable_public_page                  = true
  disable_public_json                  = true
  private_ip_filter                    = "52.95.245.0/24,203.0.113.0/24"
  private_user_authentication_required = true
  enable_sms_subscription              = false

  time_zone_id = "Europe/Amsterdam"
  service_groups = [
    {
      public_display_name = "Internal Services"
      public_description  = "Internal services and integrations"
      services = [
        allquiet_service.shipping_api.id
      ]
    }
  ]
}


