
resource "allquiet_user" "millie_brown" {
  display_name = "Millie Bobby Brown"
  email        = "acceptance-tests+millie@allquiet.app"
}

resource "allquiet_user" "taylor_swift" {
  display_name = "Taylor Swift"
  email        = "acceptance-tests+taylor@allquiet.app"
}

resource "allquiet_on_call_override" "millie_brown_override1" {
  user_id = allquiet_user.millie_brown.id
  type    = "online"
  start   = "2025-09-11T00:00:00Z"
  end     = "2025-09-11T00:20:00Z"
}

resource "allquiet_on_call_override" "millie_brown_override2" {
  user_id              = allquiet_user.millie_brown.id
  type                 = "offline"
  start                = "2025-10-01T00:00:00Z"
  end                  = "2025-11-01T00:00:00Z"
  replacement_user_ids = [allquiet_user.taylor_swift.id]
}

resource "allquiet_team" "example_team" {
  display_name = "Example Team"
}

resource "allquiet_on_call_override" "team_scoped_override" {
  user_id              = allquiet_user.millie_brown.id
  team_id              = allquiet_team.example_team.id
  type                 = "offline"
  start                = "2025-12-01T00:00:00Z"
  end                  = "2025-12-02T00:00:00Z"
  replacement_user_ids = [allquiet_user.taylor_swift.id]
}
