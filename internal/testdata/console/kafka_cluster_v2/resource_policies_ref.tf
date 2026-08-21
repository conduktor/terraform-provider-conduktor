resource "conduktor_console_resource_policy_v1" "cluster_test_policy" {
  name = "cluster-test-policy"
  spec = {
    target_kind = "Topic"
    description = "Test policy for kafka cluster policies_ref"
    rules = [
      {
        condition     = "spec.replicationFactor == 3"
        error_message = "replication factor should be 3"
      }
    ]
  }
}

resource "conduktor_console_kafka_cluster_v2" "test_policies_ref" {
  name = "kafka-cluster-with-policies"
  spec = {
    display_name      = "Kafka Cluster With Policies"
    bootstrap_servers = "localhost:9092"
    policies_ref      = [conduktor_console_resource_policy_v1.cluster_test_policy.name]
  }
}
