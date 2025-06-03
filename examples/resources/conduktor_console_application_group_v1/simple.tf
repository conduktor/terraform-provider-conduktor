
resource "conduktor_console_application_group_v1" "simple" {
  name        = "simple"
  application = "myapp"
  spec = {
    display_name = "Simple Application Group"
    description  = "Simple Description"
    permissions = [
      {
        app_instance  = "my-app-instance"
        resource_type = "TOPIC"
        pattern_type  = "LITERAL"
        name          = "*"
        permissions   = ["topicViewConfig", "topicConsume"]
      },
    ]
  }
}
