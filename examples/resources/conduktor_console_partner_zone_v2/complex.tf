
resource "conduktor_console_partner_zone_v2" "complex" {
  name = "complex"
  spec = {
    cluster      = "gw-cluster"
    display_name = "Complex Partner Zone"
    description  = "An external partner to exchange data with."
    url          = "https://partner1.com"
    partner = {
      name  = "John Doe"
      role  = "Data analyst"
      email = "johndoe@partner.io"
      phone = "07827 837 177"
    }
    authentication_mode = {
      type            = "PLAIN"
      service_account = "external-partner"
    }
    topics = [
      {
        name          = "topic-a"
        backing_topic = "kafka-topic-a"
        permission    = "WRITE"
      },
      {
        name          = "topic-b"
        backing_topic = "kafka-topic-b"
        permission    = "READ"
      }
    ]
    traffic_control_policies = {
      max_produce_rate    = 1e+06
      max_consume_rate    = 1e+06
      limit_commit_offset = 30
    }
    headers = {
      add_on_produce = [
        {
          key                = "key"
          value              = "value"
          override_if_exists = false
        }
      ],
      remove_on_consume = [
        {
          key_regex = "my_org_prefix.*"
        }
      ]
    }
  }
}

