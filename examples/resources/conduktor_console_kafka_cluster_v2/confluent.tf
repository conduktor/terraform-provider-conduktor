resource "conduktor_console_kafka_cluster_v2" "confluent" {
  name = "confluent-cluster"
  labels = {
    "env" = "staging"
  }
  spec = {
    display_name      = "Confluent Cluster"
    bootstrap_servers = "aaa-aaaa.us-west4.gcp.confluent.cloud:9092"
    properties = {
      "sasl.jaas.config"  = "org.apache.kafka.common.security.plain.PlainLoginModule required username='admin' password='admin-secret';"
      "security.protocol" = "SASL_PLAINTEXT"
      "sasl.mechanism"    = "PLAIN"
    }
    icon                         = "kafka"
    ignore_untrusted_certificate = false
    kafka_flavor = {
      confluent = {
        key                      = "yourApiKey123456"
        secret                   = "yourApiSecret123456"
        confluent_environment_id = "env-12345"
        confluent_cluster_id     = "lkc-67890"
      }
    }
    schema_registry = {
      confluent_like = {
        url                          = "https://bbb-bbbb.us-west4.gcp.confluent.cloud:8081"
        ignore_untrusted_certificate = false
        security = {
          ssl_auth = {
            key               = <<EOT
-----BEGIN PRIVATE KEY-----
MIIOXzCCDUegAwIBAgIRAPRytMVYJNUgCbhnA+eYumgwDQYJKoZIhvcNAQELBQAw
...
IFyCs+xkcgvHFtBjjel4pnIET0agtbGJbGDEQBNxX+i4MDA=
-----END PRIVATE KEY-----
EOT
            certificate_chain = <<EOT
-----BEGIN CERTIFICATE-----
MIIOXzCCDUegAwIBAgIRAPRytMVYJNUgCbhnA+eYumgwDQYJKoZIhvcNAQELBQAw
...
IFyCs+xkcgvHFtBjjel4pnIET0agtbGJbGDEQBNxX+i4MDA=
-----END CERTIFICATE-----
EOT
          }
        }
      }
    }
  }
}
