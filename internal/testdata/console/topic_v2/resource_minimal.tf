
resource "conduktor_console_topic_v2" "minimal" {
  name    = "minimal"
  cluster = "kafka-cluster"
  spec = {
    partitions         = 3
    replication_factor = 1
  }
}
