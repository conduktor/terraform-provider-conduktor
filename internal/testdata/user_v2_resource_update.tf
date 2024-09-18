
resource "conduktor_user_v2" "test" {
  name = "pam.beesly@dunder.mifflin.com"
  spec {
    firstname = "Pam"
    lastname  = "Halpert"
    permissions = [
      {
        resource_type = "PLATFORM"
        permissions   = ["userView", "datamaskingView", "auditLogView", "clusterConnectionsManage"]
      },
      {
        resource_type = "TOPIC"
        permissions   = ["topicViewConfig", "topicConsume", "topicProduce"]
        name          = "test-topic"
        pattern_type  = "LITERAL"
        cluster       = "*"
      }
    ]
  }
}
