resource "conduktor_console_kafka_subject_v2" "avro" {
  name    = "avro.value"
  cluster = "kafka-cluster"
  labels = {
    "team"        = "test"
    "environment" = "test"
  }
  spec = {
    format        = "AVRO"
    compatibility = "FORWARD_TRANSITIVE"
    schema = jsonencode(
      {
        "type" : "record",
        "name" : "MyRecord",
        "namespace" : "com.mycompany",
        "fields" : [
          {
            "name" : "id",
            "type" : "long"
          }
        ]
      }
    )
  }
}