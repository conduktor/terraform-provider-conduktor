resource "conduktor_console_kafka_subject_v2" "protobuf" {
  name    = "protobuf.value"
  cluster = "kafka-cluster"
  labels = {
    "team"        = "test"
    "environment" = "test"
  }
  spec = {
    format        = "PROTOBUF"
    compatibility = "BACKWARD"
    schema = trimspace(<<-EOF
      syntax = "proto3";

      message MyRecord {
        int32 id = 1;
        string createdAt = 2;
        string name = 3;
      }
    EOF
    )
  }
}