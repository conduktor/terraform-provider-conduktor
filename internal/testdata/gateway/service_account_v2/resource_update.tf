
resource "conduktor_gateway_service_account_v2" "test" {
  name     = "test-sa"
  vcluster = "vcluster_sa"
  spec = {
    type           = "EXTERNAL"
    external_names = ["newExternalName"]
  }
}

