
resource "conduktor_console_topic_policy_v1" "simple" {
  name = "simple"
  spec = {
    policies = {
      "spec.configs.retention.ms" = {
        range = {
          optional = true
          max      = 3600000
          min      = 60000
        }
      }
    }
  }
}

