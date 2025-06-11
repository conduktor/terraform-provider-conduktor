
resource "conduktor_console_resource_policy_v1" "minimal" {
  name = "minimal"
  spec = {
    target_kind = "Connector"
    rules = [
      {
        condition     = "spec.replicationFactor == 3"
        error_message = "replication factor should be 3"
      }
    ]
  }
}
