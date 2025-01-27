
resource "conduktor_gateway_interceptor_v2" "header-removal" {
  name = "remove-headers"
  spec {
    plugin_class = "io.conduktor.gateway.interceptor.safeguard.MessageHeaderRemovalPlugin"
    priority     = 100
    config = jsonencode(jsondecode(<<EOF
{
  "topic": "topic-.*",
  "headerKeyRegex": "headerKey.*"
}
EOF
    ))
  }
}
