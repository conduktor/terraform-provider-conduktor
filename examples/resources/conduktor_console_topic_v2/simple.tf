
resource "conduktor_console_topic_v2" "simple" {
  name    = "simple"
  cluster = "kafka-cluster"
  labels = {
    domain  = "clickstream"
    appcode = "clk"
  }
  catalog_visibility      = "PUBLIC"
  description_is_editable = false
  description             = "# Event Stream from Click Application"
  spec = {
    partitions         = 3
    replication_factor = 1
    configs = {
      "min.insync.replicas" = "2",
      "cleanup.policy"      = "delete",
      "retention.ms"        = "60000"
    }
  }
}
