
resource "conduktor_console_service_account_v1" "test" {
  name    = "test-service-account"
  cluster = "kafka-cluster"
  spec = {
    authorization = {
      kafka = {
        acls = [
          {
            name         = "test-topic"
            pattern_type = "LITERAL"
            operations   = ["Write"]
            type         = "TOPIC"
            host         = "*"
            permission   = "Deny",
          },
          {
            name         = "test-topic-2"
            pattern_type = "PREFIXED"
            operations   = ["Write"]
            type         = "TOPIC"
          },
        ]
      }
    }
  }
}
