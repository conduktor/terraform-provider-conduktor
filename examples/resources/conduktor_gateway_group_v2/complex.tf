resource "conduktor_gateway_group_v2" "complex_group" {
  name = "complex-group"
  spec = {
    members = [
      {
        name     = "admin"
        vcluster = "vcluster_sa"
      }
    ]
    external_groups = ["ldap-admins"]
  }
}
