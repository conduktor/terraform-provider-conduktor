resource "conduktor_generic" "alert" {
  kind    = "Alert"
  version = "v3"
  name    = "test-generic-alert"
  manifest = yamlencode({
    apiVersion = "v3"
    kind       = "Alert"
    metadata = {
      name = "test-generic-alert"
      user = "admin@conduktor.io"
    }
    spec = {
      type      = "TopicAlert"
      cluster   = "kafka-cluster"
      topicName = "my-topic"
      metric    = "MessageCount"
      operator  = "GreaterThan"
      threshold = 500
      destination = {
        type   = "Webhook"
        method = "POST"
        url    = "https://example.com/hook"
        body   = "{}"
      }
    }
  })
}
