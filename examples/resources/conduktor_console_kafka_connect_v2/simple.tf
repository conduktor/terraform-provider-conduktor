resource "conduktor_console_kafka_cluster_v2" "minimal" {
  name = "mini-cluster"
  spec = {
    display_name      = "Minimal Cluster"
    bootstrap_servers = "localhost:9092"
  }
}

resource "conduktor_console_kafka_connect_v2" "simple" {
  name    = "simple-connect"
  cluster = conduktor_console_kafka_cluster_v2.minimal.name
  spec = {
    display_name = "Simple Connect Server"
    urls         = "http://localhost:8083"
  }
}
