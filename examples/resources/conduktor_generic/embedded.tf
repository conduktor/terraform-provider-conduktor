resource "conduktor_generic" "example" {
  kind    = "User"
  version = "v2"
  name    = "bob@company.io"
  manifest = yamlencode({
    apiVersion = "v2"
    kind       = "User"
    metadata = {
      name = "bob@company.io"
    }
    spec = {
      firstName = "Bob"
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
