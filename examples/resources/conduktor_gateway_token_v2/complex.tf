resource "conduktor_gateway_token_v2" "complex" {
  vcluster         = "vcluster_sa"
  username         = "user10"
  lifetime_seconds = 3600
}


output "complex_token" {
  value     = conduktor_gateway_token_v2.complex.token
  sensitive = true
}

resource "local_file" "complex_token" {
  content  = conduktor_gateway_token_v2.complex.token
  filename = "${path.module}/complex_token.txt"
}