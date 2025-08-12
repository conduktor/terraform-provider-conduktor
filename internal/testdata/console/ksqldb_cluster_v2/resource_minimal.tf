
resource "conduktor_console_kafka_cluster_v2" "minimal" {
  name = "mini-cluster"
  spec = {
    display_name      = "Minimal Cluster"
    bootstrap_servers = "localhost:9092"
  }
}

resource "conduktor_console_ksqldb_cluster_v2" "minimal" {
  name    = "minimal-ksqldb"
  cluster = conduktor_console_kafka_cluster_v2.minimal.name
  spec = {
    display_name = "Minimal KSQLDB"
    url          = "http://localhost:8088"
  }
}
