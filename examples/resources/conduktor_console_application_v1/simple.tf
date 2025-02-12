resource "conduktor_console_application_v1" "example" {
  name = "simple-app"
  spec = {
    title = "Simple Application"
    owner = "admin"
  }
}

