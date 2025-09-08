resource "conduktor_console_kafka_subject_v2" "unknown_schema" {
  name    = "bad-schema"
  cluster = "kafka-cluster"
  spec = {
    format = "JSON"
    schema = "this is not a valid schema"
  }
}