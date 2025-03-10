
resource "conduktor_gateway_service_account_v2" "minimal" {
  name = "minimal"
  spec = {
    type = "LOCAL"
  }
}
