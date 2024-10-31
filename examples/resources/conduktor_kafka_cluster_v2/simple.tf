resource "conduktor_kafka_cluster_v2" "simple" {
  name = "simple-cluster"
  spec {
    display_name                 = "Simple kafka Cluster"
    icon                         = "kafka"
    color                        = "#000000"
    bootstrap_servers            = "localhost:9092"
    ignore_untrusted_certificate = true
  }
}
