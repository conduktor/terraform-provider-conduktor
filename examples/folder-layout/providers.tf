terraform {
  required_version = ">= v1.10.4"
  required_providers {
    conduktor = {
      # source  = "terraform.local/conduktor/conduktor" # local provider
      source  = "conduktor/conduktor"
      version = ">= 0.4.0"
    }
  }
}

provider "conduktor" {
  alias          = "console"
  mode           = "console"
  base_url       = "http://localhost:8080" # or env vars CDK_CONSOLE_URL or CDK_BASE_URL
  admin_user     = "admin@conduktor.io"    # or env var CDK_ADMIN_EMAIL
  admin_password = "testP4ss!"             # or env var CDK_ADMIN_PASSWORD
  # api_token = "" # switch from admin user creds to api token once generated


  insecure = true # or env var CDK_INSECURE
}

provider "conduktor" {
  alias          = "gateway"
  mode           = "gateway"
  base_url       = "http://localhost:8888"
  admin_user     = "admin"
  admin_password = "conduktor"

  insecure = true # or env var CDK_INSECURE
}
