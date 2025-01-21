
resource "allquiet_user" "millie_brown" {
  display_name = "Millie Bobby Brown"
  email        = "millie@acme.com"
}

resource "allquiet_user" "taylor" {
  display_name = "Taylor Swift"
  email        = "taylor@acme.com"
}

resource "allquiet_organization_membership" "my_organization_millie_brown" {
  user_id = allquiet_user.millie_brown.id
  role    = "Administrator"
}

resource "allquiet_organization_membership" "my_organization_taylor" {
  user_id = allquiet_user.taylor.id
  role    = "Owner"
}