
resource "conduktor_generic" "embedded" {
  kind    = "User"
  version = "v2"
  name    = "jim.halpert@dunder.mifflin.com"
  manifest = yamlencode({
    apiVersion = "v2"
    kind       = "User"
    metadata = {
      name = "jim.halpert@dunder.mifflin.com"
    }
    spec = {
      lastName  = "Canterbury"
      firstName = "Tim"
      permissions = [
        {
          resourceType = "PLATFORM"
          permissions = [
            "userView",
            "datamaskingView",
            "auditLogView"
          ]
        }
      ]
    }
  })
}
