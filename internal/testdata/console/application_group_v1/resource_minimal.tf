resource "conduktor_console_application_group_v1" "minimal" {
  name        = "test-application-group-minimal"
  application = "test-application"
  spec = {
    display_name = "Test Application Group Minimal"
  }
}
