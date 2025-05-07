resource "allquiet_user" "test" {
  display_name = "Millie Bobby Brown"
  email        = "acceptance-tests+millie@allquiet.app"
  phone_number = "+12035479055"
}

data "allquiet_user" "test_by_email" {
  email      = "acceptance-tests+millie@allquiet.app"
  depends_on = [allquiet_user.test]
}

data "allquiet_user" "test_by_display_name" {
  display_name = "Millie Bobby Brown"
  depends_on   = [allquiet_user.test]
}