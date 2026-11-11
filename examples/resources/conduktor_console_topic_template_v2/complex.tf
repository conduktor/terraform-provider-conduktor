resource "conduktor_console_topic_template_v2" "complex" {
  name = "complex-template"
  labels = {
    team = "platform"
    env  = "prod"
  }
  spec = {
    display_name = "Complex Production Template"
    description  = "High-throughput template for production services with Handlebars naming"
    defaults = {
      metadata = {
        name = "svc-{{department}}-{{env}}"
        labels = {
          managed-by = "terraform"
        }
      }
      spec = {
        partitions         = 12
        replication_factor = 3
        configs = {
          "retention.ms"        = "604800000"
          "cleanup.policy"      = "delete"
          "min.insync.replicas" = "2"
        }
      }
    }
  }
}
