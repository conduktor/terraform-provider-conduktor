resource "conduktor_console_resource_policy_v1" "connect_test_policy" {
  name = "connect-test-policy"
  spec = {
    target_kind = "Connector"
    description = "Test policy for kafka connect policies_ref"
    rules = [
      {
        condition     = "spec.config[\"connector.class\"] != \"\""
        error_message = "connector class must be set"
      }
    ]
  }
}

resource "conduktor_console_kafka_cluster_v2" "connect_cluster" {
  name = "connect-test-cluster"
  spec = {
    display_name      = "Connect Test Cluster"
    bootstrap_servers = "localhost:9092"
  }
}

resource "conduktor_console_kafka_connect_v2" "test_policies_ref" {
  name    = "kafka-connect-with-policies"
  cluster = conduktor_console_kafka_cluster_v2.connect_cluster.name
  spec = {
    display_name = "Kafka Connect With Policies"
    urls         = "http://localhost:8083"
    policies_ref = [conduktor_console_resource_policy_v1.connect_test_policy.name]
  }
}
