locals {
  file_content = file("${path.module}/schema.json")
}


resource "conduktor_console_kafka_subject_v2" "json_file" {
  name    = "api-json-example-subject.value"
  cluster = "kafka-cluster"
  labels = {
    "team"        = "test"
    "environment" = "test"
  }
  spec = {
    format        = "JSON"
    compatibility = "BACKWARD"
    schema        = local.file_content
    id            = 2
    version       = 2
  }
}