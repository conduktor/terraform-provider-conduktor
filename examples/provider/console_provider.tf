provider "conduktor" {
  mode = "console"
  # mandatory console URL
  base_url = "http://localhost:8080" # or env vars CDK_CONSOLE_BASE_URL or CDK_BASE_URL

  # authentication either with api token or admin credentials
  api_token = "your-api-token" # or env var CDK_API_TOKEN or CDK_API_KEY
  #admin_user     = "admin@my-org.com" # or env var CDK_CONSOLE_USER or CDK_ADMIN_EMAIL or CDK_ADMIN_USER
  #admin_password = "admin-password"   # or env var CDK_CONSOLE_PASSWORD or CDK_ADMIN_PASSWORD

  # optional http client TLS configuration
  cert     = file("path/to/cert.pem") # or env var CDK_CONSOLE_CERT or CDK_CERT
  insecure = true                     # or env var CDK_CONSOLE_INSECURE or CDK_INSECURE

  # optional authentication via certificate
  key    = file("path/to/key.pem") # or env var CDK_CONSOLE_KEY or CDK_KEY
  cacert = file("path/to/ca.pem")  # or env var CDK_CONSOLE_CA_CERT CDK_CA_CERT
}
