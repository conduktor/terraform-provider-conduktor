
resource "conduktor_console_application_group_v1" "minimal" {
  name        = "minimal"
  application = "myapp"
  spec = {
    display_name = "Minimal Application Group"
  }
}
