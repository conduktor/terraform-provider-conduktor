
resource "conduktor_console_application_group_v1" "test" {
  name        = "myappgroup"
  application = "myapp"
  spec = {
    display_name    = "My Application Group"
    description     = "test"
    external_groups = ["mygroup"]
    permissions = [
      {
        app_instance  = "my-app-instance"
        resource_type = "TOPIC"
        pattern_type  = "LITERAL"
        name          = "*"
        permissions   = ["topicViewConfig"]
      },
    ]
  }
}
