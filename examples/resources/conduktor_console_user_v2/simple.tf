resource "conduktor_console_user_v2" "example" {
  name = "bob@company.io"
  spec = {
    firstname = "Bob"
    lastname  = "Smith"
  }
}
