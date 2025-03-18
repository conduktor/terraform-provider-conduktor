
resource "conduktor_console_kafka_cluster_v2" "cluster" {
  name = var.cluster_name
  spec = {
    display_name                 = "Simple kafka Cluster"
    icon                         = "kafka"
    color                        = "#000000"
    bootstrap_servers            = "redpanda:9092"
    ignore_untrusted_certificate = true
  }
}
