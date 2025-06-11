
resource "conduktor_console_resource_policy_v1" "test" {
  name = "resourcepolicy"
  labels = {
    "label1" = "value1"
    "label2" = "value2"
  }
  spec = {
    target_kind = "Topic"
    description = "This is an updated test resource policy"
    rules = [
      {
        condition     = "spec.replicationFactor == 3"
        error_message = "replication factor should be 3"
      }
    ]
  }
}
