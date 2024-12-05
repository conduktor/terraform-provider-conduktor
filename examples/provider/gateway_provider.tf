provider "conduktor" {
  mode = "gateway"
  # mandatory gateway URL
  base_url = "http://localhost:8888" # or env vars CDK_GATEWAY_URL

  # authentication with admin credentials
  admin_user     = "admin"          # or env var CDK_GATEWAY_USER
  admin_password = "admin-password" # or env var CDK_GATEWAY_PASSWORD

  # optional http client TLS configuration
  cert     = file("path/to/cert.pem") # or env var CDK_GATEWAY_CERT
  insecure = true                     # or env var CDK_GATEWAY_INSECURE

  # optional authentication via certificate
  key    = file("path/to/key.pem") # or env var CDK_GATEWAY_KEY
  cacert = file("path/to/ca.pem")  # or env var CDK_GATEWAY_CACERT
}
