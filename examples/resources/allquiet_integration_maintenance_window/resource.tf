resource "allquiet_team" "root" {
  display_name = "Root"
}

resource "allquiet_integration" "some_integration" {
  display_name = "My Datadog Integration"
  team_id      = allquiet_team.root.id
  type         = "Datadog"
}

resource "allquiet_integration_maintenance_window" "test1" {
  integration_id = allquiet_integration.some_integration.id
  start          = "2025-01-01T00:00:00Z"
  end            = "2025-01-02T00:00:00Z"
  description    = "My Maintenance Window (January)"
  type           = "maintenance"
}

resource "allquiet_integration_maintenance_window" "test2" {
  integration_id = allquiet_integration.some_integration.id
  start          = "2025-02-01T00:00:00Z"
  end            = "2025-02-02T00:00:00Z"
  description    = "My Muted Window (February)"
  type           = "muted"
}

resource "allquiet_integration_maintenance_window" "test3" {
  integration_id = allquiet_integration.some_integration.id
  end            = "2026-01-01T00:00:00Z"
  description    = "My Maintenance Window (Until First of 2026)"
  type           = "maintenance"
}

resource "allquiet_integration_maintenance_window" "test4" {
  integration_id = allquiet_integration.some_integration.id
  start          = "2025-01-01T00:00:00Z"
  description    = "My Muted Window (Until End of Time)"
  type           = "muted"
}