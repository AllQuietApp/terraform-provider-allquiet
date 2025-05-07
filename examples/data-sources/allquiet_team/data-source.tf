resource "allquiet_team" "test" {
  display_name = "Test Team"
}

data "allquiet_team" "test" {
  display_name = "Test Team"
  depends_on   = [allquiet_team.test]
}
