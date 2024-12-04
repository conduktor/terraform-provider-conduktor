resource "conduktor_gateway_service_account_v2" "local_sa" {
  name = "simple-service-account"
  spec {
    type = "LOCAL"
  }
}
