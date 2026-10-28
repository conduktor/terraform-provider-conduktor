resource "conduktor_console_kafka_cluster_v2" "test_cmc" {
  name = "kafka-cluster-with-cmc"
  spec = {
    display_name      = "Cluster with Metrics Collectors"
    bootstrap_servers = "localhost:9092"
    client_metrics_collectors = [
      {
        url         = "https://prometheus.example.com/api/v1/write"
        no_security = {}
      }
    ]
  }
}
