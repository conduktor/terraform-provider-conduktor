resource "conduktor_console_topic_template_v2" "test" {
  name = "my-template"
  labels = {
    env = "prod"
  }
  spec = {
    display_name = "My Topic Template Updated"
    description  = "Updated description"
    defaults = {
      metadata = {
        name = "topic-{{department}}"
        labels = {
          throughput = "high"
        }
      }
      spec = {
        partitions         = 12
        replication_factor = 3
        configs = {
          "retention.ms"   = "1209600000"
          "cleanup.policy" = "delete"
        }
      }
    }
  }
}
