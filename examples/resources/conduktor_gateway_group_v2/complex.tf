resource "conduktor_gateway_group_v2" "complex_group" {
  name = "complex-group"
  spec = {
    external_groups = ["ldap-admins"]
  }
}
