
resource "conduktor_console_user_v2" "coworkers1" {
  name = "michael.scott@dunder.mifflin.com"
  spec {
    firstname   = "Michael"
    lastname    = "Scott"
    permissions = []
  }
}

resource "conduktor_console_group_v2" "test" {
  name = "sales"
  spec {
    display_name    = "New Sales Department"
    description     = "New Sales Department Group"
    external_groups = ["sales", "scranton"]
    members         = ["michael.scott@dunder.mifflin.com"]
    permissions = [
      {
        resource_type = "PLATFORM"
        permissions   = ["userView", "datamaskingView", "auditLogView", "clusterConnectionsManage"]
      },
      {
        resource_type = "TOPIC"
        name          = "test-topic"
        cluster       = "*"
        pattern_type  = "LITERAL"
        permissions   = ["topicViewConfig", "topicConsume", "topicProduce"]
      }
    ]
  }
  depends_on = [conduktor_console_user_v2.coworkers1]
}
