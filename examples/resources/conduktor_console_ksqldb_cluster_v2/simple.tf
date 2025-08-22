resource "conduktor_console_ksqldb_cluster_v2" "simple" {
  name    = "simple-ksqldb"
  cluster = "kafka-cluster"
  spec = {
    display_name = "Simple KSQLDB cluster"
    url          = "http://localhost:8088"
  }
}
