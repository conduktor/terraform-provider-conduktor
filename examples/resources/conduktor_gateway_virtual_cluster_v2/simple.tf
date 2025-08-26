resource "conduktor_gateway_virtual_cluster_v2" "simple" {
  name = "simple"
  spec = {
    acl_enabled = false
    type        = "Standard"
    acl_mode    = "KAFKA_API"
    super_users = [
      "user1"
    ]
  }
}

