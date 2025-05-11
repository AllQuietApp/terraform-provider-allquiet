
# Read a user by email
data "allquiet_user" "test_by_email" {
  email = "acceptance-tests+millie@allquiet.app"
}