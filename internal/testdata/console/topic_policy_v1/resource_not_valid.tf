
resource "conduktor_console_topic_policy_v1" "test" {
  name = "notvalid"
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
        range = {
          max = 3600000
          min = 60000
        }
      }
    }
  }
}
