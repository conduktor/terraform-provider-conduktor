resource "conduktor_console_kafka_cluster_v2" "test" {
  name = "test-cluster"
  labels = {
    "env" = "test"
  }
  spec = {
    display_name      = "Test Cluster"
    bootstrap_servers = "localhost:9092"
    properties = {
      "sasl.jaas.config"  = "org.apache.kafka.common.security.plain.PlainLoginModule required username=admin password=admin-secret"
      "security.protocol" = "SASL_SSL"
      "sasl.mechanism"    = "PLAIN"
    }
    icon                         = "kafka"
    color                        = "#FF0000"
    ignore_untrusted_certificate = true
    kafka_flavor = {
      confluent = {
        key                      = "confluent-key"
        secret                   = "confluent-secret"
        confluent_cluster_id     = "confluent-cluster-id"
        confluent_environment_id = "confluent-environment-id"
      }
    }
    schema_registry = {
      glue = {
        region        = "eu-west-1"
        registry_name = "default"
        security = {
          credentials = {
            access_key_id = "accessKey"
            secret_key    = "secretKey"
          }
        }
      }
      confluent_like = {
        url                          = "http://localhost:8081"
        ignore_untrusted_certificate = true
        security = {
          bearer_token = {
            token = "auth-token"
          }
        }
      }
    }
  }
}
