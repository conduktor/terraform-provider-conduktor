
resource "conduktor_console_user_v2" "user1" {
  name = var.user1
  spec = {
    firstname = "Bob"
    lastname  = "Smith"
  }
}

resource "conduktor_console_user_v2" "user2" {
  name = var.user2
  spec = {
    firstname = "Tim"
    lastname  = "Smith"
  }
}
