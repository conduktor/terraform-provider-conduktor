# website analytics dev
resource "conduktor_console_application_instance_v1" "website-analytics-dev" {
  name        = "website-analytics-dev"
  application = "website-analytics"
  spec = {
    cluster = "my-cluster"
    resources = [
      {
        type         = "TOPIC"
        name         = "website-analytics."
        pattern_type = "PREFIXED"
      },
      {
        type         = "SUBJECT"
        name         = "website-analytics."
        pattern_type = "PREFIXED"
      },
      {
        type         = "CONSUMER_GROUP"
        name         = "website-analytics."
        pattern_type = "PREFIXED"
      }
    ]
    topic_policy_ref = [
      "generic-dev-topic-policy"
    ]
    application_managed_service_account = false
  }
}
