resource "conduktor_console_connector_v2" "test" {
  name            = "connector-test"
  cluster         = "kafka-cluster"
  connect_cluster = "kafka-connect"
  labels = {
    "env" = "test"
    "sec" = "C1"
  }
  description = "description update"
  auto_restart = {
    enabled = false
  }
  spec = {
    config = {
      "connector.class" = "org.apache.kafka.connect.tools.MockSourceConnector"
      "tasks.max"       = "2"
      "topic"           = "click.pageviews.new"
      "file"            = "/etc/kafka/producer.properties"
    }
  }
}
