resource "conduktor_kafka_cluster_v2" "aiven" {
  name = "aiven-cluster"
  labels = {
    "env" = "test"
  }
  spec {
    display_name      = "Aiven Cluster"
    bootstrap_servers = "cluster.aiven.io:9092"
    properties = {
      "sasl.jaas.config"  = "org.apache.kafka.common.security.plain.PlainLoginModule required username='admin' password='admin-secret';"
      "security.protocol" = "SASL_SSL"
      "sasl.mechanism"    = "PLAIN"
    }
    icon                         = "crab"
    ignore_untrusted_certificate = true
    kafka_flavor = {
      type         = "Aiven"
      api_token    = "a1b2c3d4e5f6g7h8i9j0"
      project      = "my-kafka-project"
      service_name = "my-kafka-service"
    }
    schema_registry = {
      type                         = "ConfluentLike"
      url                          = "https://sr.aiven.io:8081"
      ignore_untrusted_certificate = false
      security = {
        type     = "BasicAuth"
        username = "uuuuuuu"
        password = "ppppppp"
      }
    }
  }
}
