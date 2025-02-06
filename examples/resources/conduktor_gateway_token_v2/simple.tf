resource "conduktor_gateway_token_v2" "simple" {
  username         = "user_passthrough"
  lifetime_seconds = 3600
}
