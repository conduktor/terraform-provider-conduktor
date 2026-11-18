resource "conduktor_gateway_group_v2" "test" {
  name = "test-group"
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
