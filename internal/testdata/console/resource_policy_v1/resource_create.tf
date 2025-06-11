
resource "conduktor_console_resource_policy_v1" "test" {
  name = "resourcepolicy"
  labels = {
    "label1" = "value1"
  }
  spec = {
    target_kind = "Topic"
    description = "This is a test resource policy"
    rules = [
      {
        condition     = "spec.replicationFactor == 3"
        error_message = "replication factor should be 3"
      }
    ]
  }
}
