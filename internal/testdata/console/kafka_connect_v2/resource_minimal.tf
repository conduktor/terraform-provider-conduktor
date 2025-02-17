
resource "conduktor_console_kafka_cluster_v2" "minimal" {
  name = "mini-cluster"
  spec = {
    display_name      = "Minimal Cluster"
    bootstrap_servers = "localhost:9092"
  }
}

resource "conduktor_console_kafka_connect_v2" "minimal" {
  name    = "minimal-connect"
  cluster = conduktor_console_kafka_cluster_v2.minimal.name
  spec = {
    display_name = "Minimal Connect"
    urls         = "http://localhost:8083"
  }
}
