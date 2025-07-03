
resource "conduktor_console_partner_zone_v2" "test" {
  name = "partner-zone"
  labels = {
    "label1" = "value1"
  }
  spec = {
    display_name = "Partner Zone"
    description  = "This is a test partner zone"
    cluster      = "gw-cluster"
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
