
resource "conduktor_gateway_service_account_v2" "test" {
  name = "user1"
  spec {
    type           = "EXTERNAL"
    external_names = ["externalName"]
  }
}
