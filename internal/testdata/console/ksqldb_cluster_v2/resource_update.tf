
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
    display_name = "Test KSQLDB updated"
    url          = "https://localhost:8088"
    headers = {
      X-PROJECT-HEADER = "value"
      AnotherHeader    = "test"
      Cache-Control    = "no-store"
    }
    ignore_untrusted_certificate = false
    security = {
      basic_auth = {
        username = "user"
        password = "password"
      }
    }
  }
}
