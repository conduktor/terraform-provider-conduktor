resource "conduktor_console_kafka_subject_v2" "json_full" {
  name    = "api-json-example-subject.value"
  cluster = "kafka-cluster"
  labels = {
    "team"        = "test"
    "environment" = "test"
  }
  spec = {
    format        = "JSON"
    compatibility = "BACKWARD"
    schema = <<EOF
{
  "$id": "https://mycompany.com/myrecord",
  "$schema": "https://json-schema.org/draft/2019-09/schema",
  "type": "object",
  "title": "MyRecord",
  "description": "Json schema for MyRecord",
  "properties": {
    "id": {
      "type": "string"
    },
    "name": {
      "type": ["string", "null"]
    }
  },
  "required": ["id"],
  "additionalProperties": false
}
EOF
    id            = 2
    version       = 2
    references = [
      {
        name    = "example-subject.value"
        subject = "example-subject.value"
        version = 1
      }
    ]
  }
}