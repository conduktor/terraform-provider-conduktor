resource "conduktor_console_connector_template_v2" "simple" {
  name = "simple-connector-template"
  spec = {
    display_name = "Simple Connector Template"
    description  = "A simple connector template for S3 sink connectors"
    defaults = {
      spec = {
        class = "io.confluent.connect.s3.S3SinkConnector"
        config = {
          "tasks.max" = "1"
          "s3.region" = "us-east-1"
        }
      }
    }
  }
}
