
resource "conduktor_console_topic_v2" "test" {
  name    = "Kafka-1st-topic-test"
  cluster = "kafka-cluster"
  labels = {
    key1 = "value1"
  }
  catalog_visibility      = "PUBLIC"
  description_is_editable = true
  description             = "description"
  sql_storage = {
    retention_time_in_second = 86400
    enabled                  = true
  }
  spec = {
    partitions         = 10
    replication_factor = 1
    configs = {
      "cleanup.policy" = "delete"
    }
  }
}
