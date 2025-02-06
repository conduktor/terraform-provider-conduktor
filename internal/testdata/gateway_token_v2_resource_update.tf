
resource "conduktor_gateway_token_v2" "test" {
  vcluster         = "vcluster_sa"
  username         = "user10"
  lifetime_seconds = 3000
}
