resource "conduktor_console_user_v2" "user1" {
  name = "user1@company.com"
  spec {
    firstname   = "User"
    lastname    = "1"
    permissions = []
  }
}

resource "conduktor_console_group_v2" "example" {
  name = "complex-group"
  spec {
    display_name    = "Complex group"
    description     = "Complex group description"
    external_groups = ["sso-group1"]
    members         = [conduktor_console_user_v2.user1.name]
    permissions = [
      {
        resource_type = "PLATFORM"
        permissions   = ["userView", "datamaskingView", "auditLogView"]
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
}
