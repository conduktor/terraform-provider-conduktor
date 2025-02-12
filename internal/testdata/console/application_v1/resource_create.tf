
resource "conduktor_console_application_v1" "app" {
  name = "myapp"
  spec = {
    title       = "Application"
    description = "My Application description"
    owner       = "Application owner"
  }
}

