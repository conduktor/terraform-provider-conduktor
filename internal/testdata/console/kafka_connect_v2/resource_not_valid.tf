
resource "conduktor_console_kafka_cluster_v2" "minimal" {
  name = "mini-cluster"
  spec = {
    display_name      = "Minimal Cluster"
    bootstrap_servers = "localhost:9092"
  }
}

resource "conduktor_console_kafka_connect_v2" "not_valid" {
  name = "test-connect"
  cluster = conduktor_console_kafka_cluster_v2.minimal.name
  labels = {
    env = "test"
  }
  spec = {
    display_name      = "Test Connect"
    urls = "http://localhost:8083"
    headers         = {
      X-PROJECT-HEADER = "value"
      AnotherHeader = "test"
    }
    ignore_untrusted_certificate = true
    security = {
      bearer_token = {
        token = "auth-token"
      }
      basic_auth = {
        username = "user"
        password = "password"
      }
    }
  }
}
