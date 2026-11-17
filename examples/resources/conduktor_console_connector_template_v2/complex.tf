resource "conduktor_console_connector_template_v2" "complex" {
  name = "complex-connector-template"
  labels = {
    env    = "production"
    domain = "data-platform"
  }
  spec = {
    display_name = "Complex Connector Template"
    description  = "A complex connector template with metadata defaults"
    defaults = {
      metadata = {
        name = "{{team}}-{{connector_type}}-connector"
        labels = {
          managed-by = "terraform"
          team       = "{{team}}"
        }
      }
      spec = {
        class = "io.confluent.connect.s3.S3SinkConnector"
        config = {
          "tasks.max"            = "1"
          "s3.region"            = "us-east-1"
          "s3.bucket.name"       = "{{s3_bucket}}"
          "flush.size"           = "1000"
          "rotate.interval.ms"   = "3600000"
          "storage.class"        = "io.confluent.connect.s3.storage.S3Storage"
          "format.class"         = "io.confluent.connect.s3.format.avro.AvroFormat"
          "schema.compatibility" = "FULL"
        }
      }
    }
  }
}
