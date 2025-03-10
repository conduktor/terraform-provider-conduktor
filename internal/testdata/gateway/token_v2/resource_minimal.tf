
resource "conduktor_gateway_token_v2" "minimal" {
  username         = "user_passthrough"
  lifetime_seconds = 3600
}
