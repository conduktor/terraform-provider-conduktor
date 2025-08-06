resource "conduktor_console_kafka_subject_v2" "avro" {
  name    = "api-avro-example-subject.value"
  cluster = "kafka-cluster"
  labels = {
    "team"        = "test"
    "environment" = "test"
  }
  spec = {
    format        = "AVRO"
    compatibility = "FORWARD_TRANSITIVE"
    schema        = <<EOF
{
  "type": "record",
  "name": "MyRecord",
  "namespace": "com.mycompany",
  "fields": [
    {
      "name": "id",
      "type": "long"
    }
  ]
}
EOF
    id            = 1
    version       = 1
  }
}