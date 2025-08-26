resource "conduktor_gateway_virtual_cluster_v2" "complex" {
  name = "complex"
  spec = {
    acl_enabled = true
    acl_mode    = "REST_API"
    type        = "Standard"
    acls = [
      {
        resource_pattern = {
          resource_type = "TOPIC"
          name          = "customers"
          pattern_type  = "LITERAL"
        }
        principal       = "User:username1"
        host            = "*"
        operation       = "READ"
        permission_type = "ALLOW"
      },
      {
        resource_pattern = {
          resource_type = "TOPIC"
          name          = "customers"
          pattern_type  = "LITERAL"
        }
        principal       = "User:username1"
        host            = "*"
        operation       = "WRITE"
        permission_type = "ALLOW"
      }
    ]
  }
}

