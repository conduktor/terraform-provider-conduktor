resource "conduktor_console_connector_template_v2" "minimal" {
  name = "minimal-connector-template"
  spec = {
    display_name = "Minimal Connector Template"
    defaults = {
      spec = {
        class = "org.apache.kafka.connect.mirror.MirrorSourceConnector"
        config = {
          "tasks.max" = "1"
        }
      }
    }
  }
}
