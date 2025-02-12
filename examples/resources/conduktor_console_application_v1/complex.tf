resource "conduktor_console_application_v1" "example" {
  name = "complex-app"
  spec = {
    title       = "Complex Application"
    description = "Complex Application description"
    owner       = "admin"
  }
}

