resource "conduktor_console_kafka_subject_v2" "json_file" {
  name    = "json-file-reference.value"
  cluster = "kafka-cluster"
  labels = {
    "team"        = "test"
    "environment" = "test"
  }
  spec = {
    format        = "JSON"
    compatibility = "BACKWARD"
    schema        = file("${path.module}/schema.json")
    id            = 2
    version       = 1
  }
}