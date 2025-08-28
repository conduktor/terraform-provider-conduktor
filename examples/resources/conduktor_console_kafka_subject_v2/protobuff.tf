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
    schema        = <<EOF
syntax = "proto3";
message MyRecord {
  int32 id = 1;
  string createdAt = 2;
  string name = 3;
}
EOF
  }
}