resource "conduktor_generic" "raw_yaml" {
  kind    = "User"
  version = "v2"
  name    = "alice@company.io"
  manifest = yamlencode(yamldecode(<<EOF
      apiVersion: v2
      kind: User
      metadata:
        name: "alice@company.io"
      spec:
        firstName: "Alice"
        lastName: "Smith"
        permissions:
          - resourceType: PLATFORM
            permissions: ["userView", "datamaskingView", "auditLogView"]
          - resourceType: TOPIC
            cluster: '*'
            name: test-topic
            patternType: LITERAL
            permissions: [ "topicViewConfig", "topicConsume", "topicProduce" ]
      EOF
  ))
}
