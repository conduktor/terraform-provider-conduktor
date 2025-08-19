resource "conduktor_console_kafka_subject_v2" "this" {
  name    = "minimal_subject"
  cluster = "kafka-cluster"
  spec = {
    format = "JSON"
    schema = file("${path.module}/schema.json")
  }
}

resource "conduktor_console_kafka_subject_v2" "complex" {
  name    = "api-json-example-subject.value"
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
    references = [
      {
        name    = "example-reference"
        subject = conduktor_console_kafka_subject_v2.this.name
        version = 1
      }
    ]
  }
}