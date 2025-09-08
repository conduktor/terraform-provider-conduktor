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
      "$id"                  = "https://mycompany.com/myrecord"
      "$schema"              = "https://json-schema.org/draft/2019-09/schema"
      "additionalProperties" = false
      "description"          = "Json schema for MyRecord"
      "properties" = {
        "id" = {
          "type" = "string"
        }
        "name" = {
          "type" = ["string", "null"]
        }
        "ext_ref" = {
          "$ref" = "https://mycompany.com/example.json"
        }
      }
      "required" = ["id"]
      "title"    = "MyRecord"
      "type"     = "object"
    })
    references = [
      {
        name    = "https://mycompany.com/example.json"
        subject = "example-subject.value"
        version = 1
      }
    ]
  }
}