resource "conduktor_console_ksqldb_cluster_v2" "bearer" {
  name    = "bearer-ksqldb"
  cluster = "kafka-cluster"
  spec = {
    display_name = "Bearer KSQLDB cluster"
    url          = "http://localhost:8088"
    headers = {
      X-PROJECT-HEADER = "value"
      Cache-Control    = "no-cache"
    }
    ignore_untrusted_certificate = false
    security = {
      bearer_token = {
        token = "token"
      }
    }
  }
}
