resource "allquiet_service" "payment_api" {
  display_name       = "Payment API"
  public_title       = "Payment API"
  public_description = "Payments"
  templates = [
    {
      display_name = "Refunds delayed"
      message      = "Refunds are currently delayed. All refunds will be processed but can currently take longer than usual to complete."
    },
    {
      display_name = "Payment gateway down"
      message      = "Our payment gateway is currently down. We are working to resolve the issue as soon as possible."
    }
  ]
}

resource "allquiet_team" "first_level_support" {
  display_name = "First Level Support"
}

resource "allquiet_team" "engineering" {
  display_name = "Engineering"
}

resource "allquiet_service" "test_with_team_connection_settings" {
  display_name       = "Test with team connection settings"
  public_title       = "Test with team connection settings"
  public_description = "Test with team connection settings"
  team_connection_settings = {
    team_connection_mode = "SelectedTeams"
    team_ids             = [allquiet_team.first_level_support.id, allquiet_team.engineering.id]
  }
}

resource "allquiet_service" "test_with_organization_connection_settings" {
  display_name       = "Test with organization connection settings"
  public_title       = "Test with organization connection settings"
  public_description = "Test with organization connection settings"
  team_connection_settings = {
    team_connection_mode = "OrganizationTeams"
  }
}