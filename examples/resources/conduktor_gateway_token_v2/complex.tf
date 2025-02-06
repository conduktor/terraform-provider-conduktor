resource "conduktor_gateway_token_v2" "complex" {
  vcluster         = "vcluster_1"
  username         = "user10"
  lifetime_seconds = 3600
}

output "complex_token" {
  value     = conduktor_gateway_token_v2.complex.token
  sensitive = true
}

