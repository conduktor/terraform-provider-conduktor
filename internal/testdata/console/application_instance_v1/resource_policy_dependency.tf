# Satisfies the `spec.policy_ref` of the complex application instance example.
# Kept out of the example HCL so it doesn't show up in the generated docs.
# Uses ResourcePolicy, the replacement for the TopicPolicy kind removed from Console 1.47.0.
resource "conduktor_console_resource_policy_v1" "topic_policy_ref" {
  name = "topic-policy"
  spec = {
    target_kind = "Topic"
    rules = [
      {
        condition     = "spec.replicationFactor >= 1"
        error_message = "replication factor should be at least 1"
      }
    ]
  }
}
