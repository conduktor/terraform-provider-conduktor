
resource "conduktor_console_topic_v2" "test" {
  name    = "topic-test"
  cluster = "kafka-cluster"
  labels = {
    key1 = "value1"
    key2 = "value2"
  }
  catalog_visibility      = "PRIVATE"
  description_is_editable = false
  description             = "new description"
  sql_storage = {
    retention_time_in_second = 86400
    enabled                  = true
  }
  spec = {
    partitions         = 10
    replication_factor = 1
    configs = {
      "cleanup.policy" = "delete"
      "retention.ms"   = "60000"
    }
  }
}

