resource "conduktor_gateway_service_account_v2" "example" {
  name = "simple-service-account"
  spec {
    type = "LOCAL"
  }
}
