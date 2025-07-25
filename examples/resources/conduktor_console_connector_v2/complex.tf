
resource "conduktor_console_connector_v2" "complex" {
  name            = "complex"
  cluster         = "kafka-cluster"
  connect_cluster = "kafka-connect"
  labels = {
    domain  = "clickstream"
    appcode = "clk"
  }
  description = "# Complex kafka connector"
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

