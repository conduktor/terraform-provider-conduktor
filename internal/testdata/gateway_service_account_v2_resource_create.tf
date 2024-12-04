
resource "conduktor_gateway_service_account_v2" "test" {
  name     = "user1"
  vcluster = "vcluster_sa"
  spec {
    type           = "EXTERNAL"
    external_names = ["externalName"]
  }
}
