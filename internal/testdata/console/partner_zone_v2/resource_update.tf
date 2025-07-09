
resource "conduktor_console_partner_zone_v2" "test" {
  name = "partner-zone"
  labels = {
    "label1" = "new-value1"
  }
  spec = {
    display_name = "Updated Partner Zone"
    description  = "This is an updated test partner zone"
    cluster      = "gw-cluster"
    authentication_mode = {
      type            = "PLAIN"
      service_account = "service-account-234"
    }
    topics = [
      {
        name          = "topic-b"
        backing_topic = "kafka-topic-b"
        permission    = "READ"
      }
    ]
  }
}
