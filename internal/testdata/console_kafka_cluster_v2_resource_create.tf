
resource "conduktor_console_kafka_cluster_v2" "test" {
  name = "test-cluster"
  labels = {
    "env" = "test"
  }
  spec {
    display_name      = "Test Cluster"
    bootstrap_servers = "localhost:9092"
    properties         = {
      "sasl.jaas.config" = "org.apache.kafka.common.security.plain.PlainLoginModule required username=admin password=admin-secret"
      "security.protocol" = "SASL_SSL"
      "sasl.mechanism" = "PLAIN"
    }
    icon                         = "kafka"
    color                        = "#FF0000"
    ignore_untrusted_certificate = true
    kafka_flavor = {
      type                     = "Confluent"
      key                      = "confluent-key"
      secret                   = "confluent-secret"
      confluent_cluster_id     = "confluent-cluster-id"
      confluent_environment_id = "confluent-environment-id"
    }
    schema_registry = {
      type                         = "ConfluentLike"
      url                          = "http://localhost:8081"
      ignore_untrusted_certificate = true
      security = {
        type  = "BearerToken"
        token = "auth-token"
      }
    }
  }
}
