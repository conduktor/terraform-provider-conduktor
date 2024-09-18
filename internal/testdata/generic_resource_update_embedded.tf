
resource "conduktor_generic" "embedded" {
  kind     = "User"
  version  = "v2"
  name     = "jim.halpert@dunder.mifflin.com"
  manifest = <<EOF
apiVersion: v2
kind: User
metadata:
  name: jim.halpert@dunder.mifflin.com
spec:
  firstName: Tim
  lastName: Canterbury
  permissions:
  - permissions:
    - userView
    - datamaskingView
    - auditLogView
    resourceType: PLATFORM
  EOF
}
