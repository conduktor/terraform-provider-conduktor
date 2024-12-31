resource "conduktor_gateway_interceptor_v2" "local_sa" {
  name = "simple-service-account"
  spec {
    type = "LOCAL"
  }
}
