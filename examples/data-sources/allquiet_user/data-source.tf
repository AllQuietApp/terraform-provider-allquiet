
# Read a user by email
data "allquiet_user" "test_by_email" {
  email = "acceptance-tests+millie@allquiet.app"
}

# Read a user by display name
data "allquiet_user" "test_by_display_name" {
  display_name = "Millie Bobby Brown"
}