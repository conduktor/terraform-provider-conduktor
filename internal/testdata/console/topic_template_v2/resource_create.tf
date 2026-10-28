resource "conduktor_console_topic_template_v2" "test" {
  name = "my-template"
  labels = {
    env = "prod"
  }
  spec = {
    display_name = "My Topic Template"
    description  = "A template for production topics"
    defaults = {
      metadata = {
        name = "topic-{{department}}"
        labels = {
          throughput = "high"
        }
      }
      spec = {
        partitions         = 6
        replication_factor = 3
        configs = {
          "retention.ms"   = "604800000"
          "cleanup.policy" = "delete"
        }
      }
    }
  }
}
