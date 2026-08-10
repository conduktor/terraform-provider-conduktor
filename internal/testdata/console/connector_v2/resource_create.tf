resource "conduktor_console_connector_v2" "test" {
  name            = "connector-test"
  cluster         = "kafka-cluster"
  connect_cluster = "kafka-connect"
  labels = {
    key1 = "value1"
  }
  description = "description"
  auto_restart = {
    enabled           = true
    frequency_seconds = 800
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
