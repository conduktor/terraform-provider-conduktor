
resource "conduktor_console_topic_v2" "minimal" {
  name    = "minimal"
  cluster = "cluster"
  spec = {
    partitions         = 10
    replication_factor = 3
  }
}
