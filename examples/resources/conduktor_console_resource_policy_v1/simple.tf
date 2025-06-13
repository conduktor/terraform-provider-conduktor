
resource "conduktor_console_resource_policy_v1" "simple" {
  name = "simple"
  labels = {
    "business-unit" = "delivery"
  }
  spec = {
    target_kind = "Topic"
    description = "A policy to check some basic rule for a topic"
    rules = [
      {
        condition     = "int(string(spec.configs[\"retention.ms\"])) >= 60000 && int(string(spec.configs[\"retention.ms\"])) <= 3600000"
        error_message = "retention should be between 1m and 1h"
      }
    ]
  }
}
