
resource "conduktor_console_application_v1" "test" {
  name = "my-application"
  spec = {
    title       = "My Application"
    description = "My Application description"
    owner       = "admin"
  }
}

