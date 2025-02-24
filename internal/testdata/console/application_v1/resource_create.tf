
resource "conduktor_console_application_v1" "test" {
  name = "my-application"
  spec = {
    title = "My Application"
    owner = "admin"
  }
}

