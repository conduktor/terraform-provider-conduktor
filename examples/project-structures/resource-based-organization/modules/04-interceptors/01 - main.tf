

resource "conduktor_gateway_interceptor_v2" "partition-limit" {
  name = "enforce-partition-limit"
  scope = {
    vcluster = "passthrough"
  }
  spec = {
    plugin_class = "io.conduktor.gateway.interceptor.safeguard.CreateTopicPolicyPlugin"
    priority     = 1
    config = jsonencode({
      topic = "myprefix-.*"
      numPartition = {
        min    = 5
        max    = 5
        action = "INFO"
      }
    })
  }
}
