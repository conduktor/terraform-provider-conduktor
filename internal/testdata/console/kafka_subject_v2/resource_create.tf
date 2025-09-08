resource "conduktor_console_kafka_subject_v2" "test" {
  name    = "api-json-example-subject.value"
  cluster = "kafka-cluster"
  labels = {
    "team"        = "test"
    "environment" = "test"
  }
  spec = {
    format        = "JSON"
    compatibility = "BACKWARD"
    schema = jsonencode({
      "$id"         = "https://mycompany.com/myrecord"
      "$schema"     = "https://json-schema.org/draft/2019-09/schema"
      "type"        = "object"
      "title"       = "MyRecord"
      "description" = "Json schema for MyRecord"
      "properties" = {
        "id" = {
          "type" = "string"
        }
        "name" = {
          "type" = ["string", "null"]
        }
      }
      "required"             = ["id"]
      "additionalProperties" = false
    })
  }
}