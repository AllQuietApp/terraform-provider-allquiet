
resource "allquiet_user" "millie_brown" {
  display_name = "Millie Bobby Brown"
  email        = "acceptance-tests+millie@allquiet.app"
}

resource "allquiet_user" "taylor" {
  display_name = "Taylor Swift"
  email        = "acceptance-tests+taylor@allquiet.app"
}

resource "allquiet_user" "timothee" {
  display_name = "Timothée Chalamet"
  email        = "acceptance-tests+timothee@allquiet.app"
}



resource "allquiet_organization_membership" "my_organization_millie_brown" {
  user_id = allquiet_user.millie_brown.id
  role    = "Administrator"
}

resource "allquiet_organization_membership" "my_organization_taylor" {
  user_id = allquiet_user.taylor.id
  role    = "Owner"
}

resource "allquiet_organization_membership" "my_organization_timothee" {
  user_id = allquiet_user.timothee.id
  role    = "Member"
}