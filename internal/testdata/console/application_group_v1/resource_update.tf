resource "conduktor_console_user_v2" "coworkers1" {
  name = "uriel.septim@tamriel.com"
  spec = {
    firstname   = "Uriel"
    lastname    = "Septim"
    permissions = []
  }
}

resource "conduktor_console_application_group_v1" "test" {
  name        = "myappgroup"
  application = "myapp"
  spec = {
    display_name    = "My Application Group"
    description     = "test"
    external_groups = ["mygroup"]
    members         = ["uriel.septim@tamriel.com"]
    permissions = [
      {
        app_instance  = "my-app-instance"
        resource_type = "TOPIC"
        pattern_type  = "LITERAL"
        name          = "*"
        permissions   = ["topicViewConfig"]
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
  depends_on = [conduktor_console_user_v2.coworkers1]
}
