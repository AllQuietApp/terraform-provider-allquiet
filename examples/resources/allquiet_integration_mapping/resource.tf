resource "allquiet_team" "root" {
  display_name = "Root"
}

resource "allquiet_integration" "datadog" {
  display_name = "My Datadog Integration with custom mapping"
  team_id      = allquiet_team.root.id
  type         = "Datadog"
}

resource "allquiet_integration_mapping" "datadog_custom_mapping" {
  integration_id = allquiet_integration.datadog.id
  attributes_mapping = {
    attributes = [
      {
        name = "Severity",
        mappings = [
          { json_path = "$.jsonBody.title" },
          { map = "A->Critical,->Warning" }
        ]
      },
      {
        name = "Status",
        mappings = [
          { static = "Open" }
        ]
      },
      {
        name = "Title",
        mappings = [
          { xpath = "//json" },
          { json_path = "$.jsonBody.status" },
          { regex = "\\d+", replace = "$1" },
          { map = "->Open" },
        ]
      }
    ]
  }
}