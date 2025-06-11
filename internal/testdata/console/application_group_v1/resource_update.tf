
resource "conduktor_console_application_group_v1" "test" {
  name        = "myappgroup"
  application = "myapp"
  spec = {
    display_name    = "My Updated Application Group"
    description     = "update test"
    external_groups = ["mygroup"]
    permissions = [
      {
        app_instance  = "my-app-instance"
        resource_type = "TOPIC"
        pattern_type  = "LITERAL"
        name          = "*"
        permissions   = ["topicViewConfig", "topicConsume"]
      },
      {
        app_instance  = "my-app-instance"
        resource_type = "SUBJECT"
        pattern_type  = "LITERAL"
        name          = "*"
        permissions   = ["subjectView"]
      }
    ]
  }
}
