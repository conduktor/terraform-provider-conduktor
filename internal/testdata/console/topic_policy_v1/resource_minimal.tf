
resource "conduktor_console_topic_policy_v1" "minimal" {
  name = "minimal"
  spec = {
    policies = {
      my-policy = {
        one_of = {
          values = [
            "value"
          ]
        }
      }
    }
  }
}

