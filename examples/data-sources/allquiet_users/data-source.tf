## Setup
resource "allquiet_user" "test1" {
  email        = "acceptance-tests+ds+millie@allquiet.app"
  display_name = "Millie Bobby Brown"
}

resource "allquiet_user" "test2" {
  email        = "acceptance-tests+ds+miley@allquiet.app"
  display_name = "Miley Cyrus"
}

## Data sources
data "allquiet_users" "users_by_email" {
  email      = "acceptance-tests+ds"
  depends_on = [allquiet_user.test1, allquiet_user.test2]
}

data "allquiet_users" "all_users" {
  depends_on = [allquiet_user.test1, allquiet_user.test2]
}