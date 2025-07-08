
resource "conduktor_console_topic_v2" "managed_labels_ro" {
  name    = "managed_labels_ro"
  cluster = "kafka-cluster"
  labels = {
    "env" = "prod"
  }
  managed_labels = {
    "conduktor.io/managed" = "true"
  }
  spec = {
    partitions         = 3
    replication_factor = 1
  }
}
