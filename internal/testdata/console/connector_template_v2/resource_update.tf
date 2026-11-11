resource "conduktor_console_connector_template_v2" "test" {
  name = "my-connector-template"
  labels = {
    env = "prod"
  }
  spec = {
    display_name = "My Connector Template Updated"
    description  = "Updated description"
    defaults = {
      metadata = {
        name = "connector-{{team}}"
        labels = {
          managed-by = "terraform"
        }
      }
      spec = {
        class = "io.confluent.connect.s3.S3SinkConnector"
        config = {
          "tasks.max"      = "2"
          "s3.bucket.name" = "my-bucket-updated"
        }
      }
    }
  }
}
