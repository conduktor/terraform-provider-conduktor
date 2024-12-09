
resource "conduktor_console_group_v2" "test" {
  name = "sales"
  spec {
    display_name    = "Sales Department"
    description     = "Sales Department Group"
    external_groups = ["sales"]
    permissions = [
      {
        resource_type = "PLATFORM"
        permissions   = ["userView", "datamaskingView", "auditLogView", "clusterConnectionsManage"]
      }
    ]
  }
}
