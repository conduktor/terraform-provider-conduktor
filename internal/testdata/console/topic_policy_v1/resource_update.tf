
resource "conduktor_console_topic_policy_v1" "test" {
  name = "topicpolicy"
  spec = {
    policies = {
      my-policy = {
        range = {
          max = 3600000
          min = 60000
        }
      }
    }
  }
}

