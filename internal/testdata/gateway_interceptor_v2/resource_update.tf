
resource "conduktor_gateway_interceptor_v2" "topic-policy-test" {
  name = "enforce-partition-limit-test"
  scope = {
    vcluster = "passthrough"
    username = "my.user2"
  }
  spec = {
    plugin_class = "io.conduktor.gateway.interceptor.safeguard.CreateTopicPolicyPlugin"
    priority     = 100
    config = jsonencode({
      topic = "updatemyprefix-.*"
      numPartition = {
        min    = 5
        max    = 10
        action = "BLOCK"
      }
      retentionMs = {
        min = 10
        max = 100
      }
    })
  }
}
