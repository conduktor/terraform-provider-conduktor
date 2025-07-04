
resource "conduktor_console_partner_zone_v2" "minimal" {
  name = "minimal"
  spec = {
    cluster = "gw-cluster"
    authentication_mode = {
      type            = "PLAIN"
      service_account = "service-account-123"
    }
    topics = [
      {
        name          = "topic-a"
        backing_topic = "kafka-topic-a"
        permission    = "WRITE"
      }
    ]
  }
}
