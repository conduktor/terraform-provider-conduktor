
resource "conduktor_console_service_account_v1" "kafka_sa" {
  name    = "kafka-service-account"
  cluster = "kafka-cluster"
  labels = {
    domain  = "clickstream"
    appcode = "clk"
  }
  spec = {
    authorization = {
      kafka = {
        acls = [
          {
            name         = "*"
            pattern_type = "LITERAL"
            operations   = ["Describe"]
            type         = "TOPIC"
          },
          {
            name         = "click.event-stream.avro"
            pattern_type = "LITERAL"
            operations   = ["Write", "Read"]
            type         = "TOPIC"
          },
          {
            name         = "public_"
            pattern_type = "PREFIXED"
            operations   = ["Read"]
            type         = "TOPIC"
          },
          {
            name         = "click.event-stream."
            pattern_type = "PREFIXED"
            operations   = ["Read"]
            type         = "CONSUMER_GROUP"
          }
        ]
      }
    }
  }
}

