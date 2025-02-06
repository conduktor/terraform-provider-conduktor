resource "conduktor_gateway_token_v2" "simple" {
  username         = "user_passthrough"
  lifetime_seconds = 3600
}

output "simple_token" {
  value     = conduktor_gateway_token_v2.simple.token
  sensitive = true
}

resource "local_file" "simple_token" {
  content  = conduktor_gateway_token_v2.simple.token
  filename = "${path.module}/simple_token.txt"
}