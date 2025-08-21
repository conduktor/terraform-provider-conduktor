resource "conduktor_console_topic_v2" "backing_topic_a" {
  name    = "kafka-topic-a"
  cluster = "gw-cluster"
  spec = {
    partitions         = 1
    replication_factor = 1
  }
}

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
  depends_on = [
    conduktor_console_topic_v2.backing_topic_a
  ]
}
