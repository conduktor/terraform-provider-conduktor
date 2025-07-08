
resource "conduktor_console_topic_v2" "managed_labels" {
  name    = "managed_labels"
  cluster = "kafka-cluster"
  labels = {
    "env"                  = "prod"
    "conduktor.io/managed" = "true"
  }
  spec = {
    partitions         = 3
    replication_factor = 1
  }
}
