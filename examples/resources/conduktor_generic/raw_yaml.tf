resource "conduktor_generic" "raw_yaml" {
  kind    = "User"
  version = "v2"
  name    = "bob@company.io"
  manifest = yamlencode(yamldecode(<<EOF
      apiVersion: v2
      kind: User
      metadata:
        name: "bob@company.io"
      spec:
        firstName: "Bob"
        lastName: Smith
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
