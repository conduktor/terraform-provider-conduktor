
resource "conduktor_console_topic_v2" "production_topic" {
  name    = "production-topic"
  cluster = "kafka-cluster"
  spec = {
    partitions         = 10
    replication_factor = 1
  }

  lifecycle {
    prevent_destroy = true
  }
}
