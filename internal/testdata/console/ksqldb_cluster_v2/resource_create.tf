
resource "conduktor_console_kafka_cluster_v2" "minimal" {
  name = "mini-cluster"
  spec = {
    display_name      = "Minimal Cluster"
    bootstrap_servers = "localhost:9092"
  }
}

resource "conduktor_console_ksqldb_cluster_v2" "test" {
  name    = "test-ksqldb"
  cluster = conduktor_console_kafka_cluster_v2.minimal.name
  spec = {
    display_name = "Test KSQLDB"
    url          = "http://localhost:8088"
    headers = {
      X-PROJECT-HEADER = "value"
      AnotherHeader    = "test"
    }
    ignore_untrusted_certificate = true
    security = {
      bearer_token = {
        token = "auth-token"
      }
    }
  }
}
