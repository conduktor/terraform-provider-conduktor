
resource "conduktor_console_connector_v2" "managed_labels_ro" {
  name            = "connector-test"
  cluster         = "kafka-cluster"
  connect_cluster = "kafka-connect"
  labels = {
    "env" = "prod"
  }
  managed_labels = {
    "conduktor.io/managed" = "true"
  }
  spec = {
    config = {
      "connector.class" = "org.apache.kafka.connect.tools.MockSourceConnector"
      "tasks.max"       = "1"
      "topic"           = "click.pageviews"
      "file"            = "/etc/kafka/consumer.properties"
    }
  }
}
