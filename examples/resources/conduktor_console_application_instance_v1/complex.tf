
resource "conduktor_console_application_instance_v1" "complex" {
  name        = "complex"
  application = "myapp"
  spec = {
    cluster         = "kafka-cluster"
    service_account = "my-service-account"
    topic_policy_ref = [
      "topic-policy"
    ]
    default_catalog_visibility = "PUBLIC"
    resources = [
      {
        type         = "TOPIC"
        name         = "click."
        pattern_type = "PREFIXED"
      },
      {
        type         = "CONSUMER_GROUP"
        name         = "click."
        pattern_type = "PREFIXED"
      },
      {
        type         = "SUBJECT"
        name         = "click."
        pattern_type = "PREFIXED"
      },
      {
        type            = "CONNECTOR"
        connect_cluster = "kafka-connect"
        name            = "click."
        pattern_type    = "PREFIXED"
      },
      {
        type           = "TOPIC"
        name           = "legacy-click."
        pattern_type   = "PREFIXED"
        ownership_mode = "LIMITED"
      }
    ]
    application_managed_service_account = false
  }
}


