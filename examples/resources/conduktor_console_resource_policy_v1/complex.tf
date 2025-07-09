
resource "conduktor_console_resource_policy_v1" "complex" {
  name = "complex"
  labels = {
    "business-unit" = "delivery"
  }
  spec = {
    target_kind = "Topic"
    description = "A policy to check some basic rule for a topic"
    rules = [
      {
        condition     = "metadata.name.matches(\"^click\\\\.[a-z0-9-]+\\\\.(avro|json)$\")" # Note: \\\\ to escape in Terraform string to end up as \\ in api call
        error_message = "topic name should match ^click\\.(?<event>[a-z0-9-]+)\\.(avro|json)$"
      },
      {
        condition     = "metadata.labels[\"data-criticality\"] in [\"C0\", \"C1\", \"C2\"]"
        error_message = "data-criticality should be one of C0, C1, C2"
      }
    ]
  }
}
