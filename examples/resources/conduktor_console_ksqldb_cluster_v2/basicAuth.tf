resource "conduktor_console_ksqldb_cluster_v2" "basic" {
  name    = "basic-ksqldb"
  cluster = "kafka-cluster"
  spec = {
    display_name = "Basic KSQLDB cluster"
    url          = "http://localhost:8088"
    headers = {
      X-PROJECT-HEADER = "value"
      Cache-Control    = "no-cache"
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
