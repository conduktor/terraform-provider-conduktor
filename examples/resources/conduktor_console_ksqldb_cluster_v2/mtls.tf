resource "conduktor_console_ksqldb_cluster_v2" "mtls" {
  name    = "mtls-ksqldb"
  cluster = "kafka-cluster"
  spec = {
    display_name = "mTLS KSQLDB cluster"
    url          = "https://localhost:8088"
    headers = {
      X-PROJECT-HEADER = "value"
      Cache-Control    = "no-cache"
    }
    ignore_untrusted_certificate = false
    security = {
      ssl_auth = {
        key               = <<EOT
-----BEGIN PRIVATE KEY-----
MIIOXzCCDUegAwIBAgIRAPRytMVYJNUgCbhnA+eYumgwDQYJKoZIhvcNAQELBQAw
...
IFyCs+xkcgvHFtBjjel4pnIET0agtbGJbGDEQBNxX+i4MDA=
-----END PRIVATE KEY-----
EOT
        certificate_chain = <<EOT
-----BEGIN CERTIFICATE-----
MIIOXzCCDUegAwIBAgIRAPRytMVYJNUgCbhnA+eYumgwDQYJKoZIhvcNAQELBQAw
...
IFyCs+xkcgvHFtBjjel4pnIET0agtbGJbGDEQBNxX+i4MDA=
-----END CERTIFICATE-----
EOT
      }
    }
  }
}
