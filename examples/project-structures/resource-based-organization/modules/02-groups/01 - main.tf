
resource "conduktor_console_group_v2" "group" {
  name = var.group_name
  spec = {
    display_name = "Group"
    description  = "Group description"
    members      = var.users
    permissions = [
      {
        resource_type = "PLATFORM"
        permissions   = ["userView", "datamaskingView", "auditLogView"]
      },
      {
        resource_type = "TOPIC"
        name          = "website-analytics."
        cluster       = "*"
        pattern_type  = "PREFIXED"
        permissions   = ["topicViewConfig", "topicConsume", "topicProduce", "topicCreate"]
      }
    ]
  }
}
