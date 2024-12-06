provider "conduktor" {
  mode = "console"
  # mandatory console URL
  base_url = "http://localhost:8080" # or env vars CDK_CONSOLE_URL or CDK_BASE_URL

  # authentication either with api token or admin credentials
  api_token = "your-api-token" # or env var CDK_API_TOKEN or CDK_API_KEY
  #admin_user     = "admin@my-org.com" # or env var CDK_ADMIN_EMAIL
  #admin_password = "admin-password"   # or env var CDK_ADMIN_PASSWORD

  # optional http client TLS configuration
  cert     = file("path/to/cert.pem") # or env var CDK_CERT
  insecure = true                     # or env var CDK_INSECURE

  # optional authentication via certificate
  key    = file("path/to/key.pem") # or env var CDK_KEY
  cacert = file("path/to/ca.pem")  # or env var CDK_CA_CERT
}
