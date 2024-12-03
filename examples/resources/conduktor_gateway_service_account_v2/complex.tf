resource "conduktor_gateway_service_account_v2" "example" {
  name = "complex-service-account"
  # vcluster = "vcluster1"
  spec {
    type           = "EXTERNAL"
    external_names = ["externalName"]
  }
}
