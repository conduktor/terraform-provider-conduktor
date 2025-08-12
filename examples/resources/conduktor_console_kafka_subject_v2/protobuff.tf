resource "conduktor_console_kafka_subject_v2" "protobuff" {
  name    = "api-protobuff-example-subject.value"
  cluster = "kafka-cluster"
  labels = {
    "team"        = "test"
    "environment" = "test"
  }
  spec = {
    format        = "PROTOBUF"
    compatibility = "BACKWARD"
    schema        = "syntax = \"proto3\";\nmessage MyRecord {\n\tint32 id = 1;\n\tstring createdAt = 2;\n\tstring name = 3;\n}\n"
    id            = 2
    version       = 2
  }
}