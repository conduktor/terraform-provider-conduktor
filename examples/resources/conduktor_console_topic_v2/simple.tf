
resource "conduktor_console_topic_v2" "simple" {
  name    = "simple"
  cluster = "kafka-cluster"
  labels = {
    domain = "clickstream"
  }
  description = "# Simple kafka topic"
  spec = {
    partitions         = 3
    replication_factor = 1
    configs = {
      "cleanup.policy" = "delete"
    }
  }
}
