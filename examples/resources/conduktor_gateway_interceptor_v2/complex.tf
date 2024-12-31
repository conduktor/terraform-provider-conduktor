resource "conduktor_gateway_interceptor_v2" "external_sa" {
  name     = "complex-service-account"
  vcluster = "vcluster_sa"
  spec {
    type           = "EXTERNAL"
    external_names = ["externalName"]
  }
}
