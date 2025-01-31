resource "allquiet_service" "payment_api" {
  display_name       = "Payment API"
  public_title       = "Payment API"
  public_description = "Payment APIs and integrations"
}

resource "allquiet_status_page" "public_status_page" {
  slug                              = "public-status-page-test"
  display_name                      = "Public Status Page"
  public_title                      = "Public Status Page"
  public_description                = "Here, we'll inform you about the status of our services"
  history_in_days                   = 30
  disable_public_subscription       = false
  public_company_url                = "https://www.allquiet.app"
  public_company_name               = "AllQuiet"
  public_support_url                = "https://www.allquiet.app/support"
  public_support_email              = "support@allquiet.app"
  public_severity_mapping_minor     = "Minor"
  public_severity_mapping_warning   = "Warning"
  public_severity_mapping_critical  = "Critical"
  banner_background_color           = "#000000"
  banner_background_color_dark_mode = "#447788"
  banner_text_color                 = "#ffffff"
  banner_text_color_dark_mode       = "#ffffff"

  time_zone_id = "Europe/Amsterdam"
  services = [
    allquiet_service.payment_api.id
  ]
}
