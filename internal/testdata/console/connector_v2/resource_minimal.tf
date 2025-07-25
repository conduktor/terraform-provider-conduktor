resource "conduktor_console_connector_v2" "minimal" {
  name            = "minimal"
  cluster         = "kafka-cluster"
  connect_cluster = "kafka-connect"
  spec = {
    config = {
      "connector.class" = "org.apache.kafka.connect.tools.MockSourceConnector"
      "tasks.max"       = "1"
      "topic"           = "click.pageviews"
      "file"            = "/etc/kafka/consumer.properties"
    }
  }
}
