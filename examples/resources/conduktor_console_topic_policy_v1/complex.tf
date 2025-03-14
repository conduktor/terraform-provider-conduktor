
resource "conduktor_console_topic_policy_v1" "complex" {
  name = "complex"
  spec = {
    policies = {
      "metadata.labels.data-criticality" = {
        one_of = {
          values = [
            "C0",
            "C1",
            "C2"
          ]
        }
      },
      "spec.configs.retention.ms" = {
        range = {
          optional = false
          max      = 604800000
          min      = 3600000
        }
      },
      "spec.replicationFactor" = {
        none_of = {
          optional = true
          values = [
            "1",
            "2"
          ]
        }
      },
      "metadata.name" = {
        match = {
          pattern = "^website-analytics.(?<event>[a-z0-9-]+).(avro|json)$"
        }
      },
      "spec.configs" = {
        allowed_keys = {
          keys = [
            "retention.ms",
            "cleanup.policy"
          ]
        }
      }
    }
  }
}