provider "conduktor" {
  alias    = "console"
  mode     = "console"
  base_url = "http://localhost:8080"

  api_token = "your-api-token"
  #admin_user     = "admin@my-org.com"
  #admin_password = "admin-password"

  insecure = true
}

provider "conduktor" {
  alias    = "gateway"
  mode     = "gateway"
  base_url = "http://localhost:8888"

  admin_user     = "admin"
  admin_password = "admin-password"

  insecure = true
}

# And how to use them with example resources
resource "conduktor_console_user_v2" "user" {
  provider = conduktor.console
  # ...
}

resource "conduktor_gateway_service_account_v2" "gateway_sa" {
  provider = conduktor.gateway
  # ...
}
