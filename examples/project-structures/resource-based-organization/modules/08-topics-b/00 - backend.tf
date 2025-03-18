
terraform {
  required_version = ">= v1.10.4"
  required_providers {
    conduktor = {
      source  = "conduktor/conduktor" # local provider
      version = ">= 0.4.0"
    }
  }
}

provider "conduktor" {
  mode      = "console"
  base_url  = "http://localhost:8080" # or env vars CDK_CONSOLE_URL or CDK_BASE_URL
  api_token = ""                      # Use an appInstance api token

  insecure = true # or env var CDK_INSECURE
}