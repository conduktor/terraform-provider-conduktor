resource "conduktor_generic" "embedded" {
  kind    = "User"
  version = "v2"
  name    = "martin@company.io"
  manifest = yamlencode({
    apiVersion = "v2"
    kind       = "User"
    metadata = {
      name = "martin@company.io"
    }
    spec = {
      firstName = "Martin"
      lastName  = "Smith"
      permissions = [
        {
          permissions = [
            "userView",
            "datamaskingView",
            "auditLogView"
          ]
          resourceType = "PLATFORM"
        },
        {
          permissions = [
            "topicViewConfig",
            "topicConsume",
            "topicProduce"
          ]
          resourceType = "TOPIC"
          name         = "test-topic"
          cluster      = "*"
          patternType  = "LITERAL"
        }
      ]
    }
  })
}
