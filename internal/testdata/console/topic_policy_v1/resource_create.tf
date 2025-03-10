
resource "conduktor_console_topic_policy_v1" "test" {
  name = "topicpolicy"
  spec = {
    policies = {
      my-policy = {
        one_of = {
          values = [
            "C0",
            "C1",
            "C2"
          ]
        }
      }
    }
  }
}
