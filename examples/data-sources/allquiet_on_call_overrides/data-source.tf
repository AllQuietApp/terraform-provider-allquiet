
resource "allquiet_user" "millie_brown" {
  display_name = "Millie Bobby Brown"
  email        = "acceptance-tests+millie@allquiet.app"
}

resource "allquiet_on_call_override" "millie_brown_override1" {
  user_id = allquiet_user.millie_brown.id
  type    = "online"
  start   = "2025-09-11T00:00:00Z"
  end     = "2025-09-11T00:20:00Z"
}

resource "allquiet_on_call_override" "millie_brown_override2" {
  user_id = allquiet_user.millie_brown.id
  type    = "offline"
  start   = "2025-10-01T00:00:00Z"
  end     = "2025-11-01T00:00:00Z"
}

data "allquiet_on_call_overrides" "example1" {
  user_id    = allquiet_user.millie_brown.id
  depends_on = [allquiet_on_call_override.millie_brown_override1, allquiet_on_call_override.millie_brown_override2, allquiet_user.millie_brown]
}
