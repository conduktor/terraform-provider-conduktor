
# This file contains a resource definition for creating an application instance targeting specifically Conduktor Console v1.34.
# Currently used to test the new `spec.policy_ref` field.
resource "conduktor_console_resource_policy_v1" "this" {
  name = "resource-policy"
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

resource "conduktor_console_application_instance_v1" "test" {
  name        = "appinstance"
  application = "myapp"
  spec = {
    cluster = "kafka-cluster"
    policy_ref = [
      conduktor_console_resource_policy_v1.this.name
    ]
    application_managed_service_account = false
    service_account                     = "my-service-account"
  }
}
