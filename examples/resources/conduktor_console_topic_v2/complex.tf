
resource "conduktor_console_topic_v2" "complex" {
  name    = "complex"
  cluster = "kafka-cluster"
  labels = {
    domain  = "clickstream"
    appcode = "clk"
  }
  catalog_visibility      = "PRIVATE"
  description_is_editable = false
  description             = "# Complex kafka topic"
  sql_storage = {
    retention_time_in_second = 60000
    enabled                  = true
  }
  spec = {
    partitions         = 3
    replication_factor = 1
    configs = {
      "cleanup.policy" = "delete",
      "retention.ms"   = "60000"
    }
  }
}

