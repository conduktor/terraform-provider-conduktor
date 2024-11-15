
resource "conduktor_kafka_cluster_v2" "test" {
  name = "test-cluster"
  labels = {
    "env" = "test"
    "sec" = "C1"
  }
  spec {
    display_name      = "Test Cluster"
    bootstrap_servers = "cluster.aiven.io:9092"
    properties         = {
      "sasl.jaas.config" = "org.apache.kafka.common.security.plain.PlainLoginModule required username=admin-update password=admin-secret-update"
      "security.protocol" = "SASL_SSL"
      "sasl.mechanism" = "PLAIN"
    }
    icon                         = "kafka"
    color                        = "#FF0000"
    ignore_untrusted_certificate = true
    kafka_flavor = {
      type                     = "Aiven"
      api_token = "aiven-api-token"
      project = "aiven-project"
      service_name = "aiven-service-name"
    }
    schema_registry = {
      type                         = "ConfluentLike"
      url                          = "http://localhost:8081"
      ignore_untrusted_certificate = false
      security = {
        type  = "BasicAuth"
        username = "user"
        password = "password"
      }
    }
  }
}