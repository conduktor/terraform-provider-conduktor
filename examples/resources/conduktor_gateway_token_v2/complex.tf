resource "conduktor_gateway_token_v2" "complex" {
  vcluster         = "vcluster_sa"
  username         = "user10"
  lifetime_seconds = 3600
}
