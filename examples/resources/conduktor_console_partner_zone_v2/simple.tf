
resource "conduktor_console_partner_zone_v2" "simple" {
  name = "simple"
  spec = {
    cluster      = "gw-cluster"
    display_name = "Simple Partner Zone"
    url          = "https://partner1.com"
    authentication_mode = {
      type            = "PLAIN"
      service_account = "simple-partner"
    }
    topics = [
      {
        name          = "topic"
        backing_topic = "kafka-topic"
        permission    = "WRITE"
      }
    ]
    traffic_control_policies = {
      max_produce_rate    = 1e+06
      max_consume_rate    = 1e+06
      limit_commit_offset = 30
    }
  }
}
