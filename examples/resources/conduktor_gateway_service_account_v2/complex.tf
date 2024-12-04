resource "conduktor_gateway_service_account_v2" "external_sa" {
  name     = "complex-service-account"
  vcluster = "vcluster_sa"
  spec {
    type           = "EXTERNAL"
    external_names = ["externalName"]
  }
}
