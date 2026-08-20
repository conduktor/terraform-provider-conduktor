
resource "conduktor_console_resource_policy_v1" "test_policy" {
  name = "app-test-policy"
  spec = {
    target_kind = "Topic"
    description = "Test policy for application policy_ref"
    rules = [
      {
        condition     = "spec.replicationFactor == 3"
        error_message = "replication factor should be 3"
      }
    ]
  }
}

resource "conduktor_console_application_v1" "test_policy_ref" {
  name = "my-application-with-policy"
  spec = {
    title      = "My Application With Policy"
    owner      = "admin"
    policy_ref = [conduktor_console_resource_policy_v1.test_policy.name]
  }
}
