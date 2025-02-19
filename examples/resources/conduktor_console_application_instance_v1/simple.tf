
resource "conduktor_console_application_instance_v1" "simple" {
  name        = "simple"
  application = "myapp"
  spec = {
    cluster = "kafka-cluster"
    resources = [
      {
        type         = "TOPIC"
        name         = "topic"
        pattern_type = "PREFIXED"
      }
    ]
    application_managed_service_account = false
  }
}

