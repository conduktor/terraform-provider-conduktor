
resource "conduktor_console_application_instance_v1" "minimal" {
  name        = "minimal"
  application = "myapp"
  spec = {
    cluster = "kafka-cluster"
  }
}

