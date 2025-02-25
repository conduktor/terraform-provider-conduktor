
resource "conduktor_console_application_instance_v1" "test" {
  name        = "appinstance"
  application = "myapp"
  spec = {
    cluster = "kafka-cluster"
    resources = [
      {
        type         = "TOPIC"
        name         = "mytopic"
        pattern_type = "LITERAL"
      },
      {
        type         = "TOPIC"
        name         = "mytopic2"
        pattern_type = "LITERAL"
      }
    ]
    application_managed_service_account = true
  }
}

