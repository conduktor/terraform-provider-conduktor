
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
          max      = 3600000
          min      = 60000
        }
      },
      "spec.replicationFactor" = {
        none_of = {
          optional = true
          values = [
            "3",
          ]
        }
      },
      "metadata.name" = {
        match = {
          pattern = "^click.(?<event>[a-z0-9-]+).(avro|json)$"
        }
      }
    }
  }
}


