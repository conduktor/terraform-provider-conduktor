resource "conduktor_gateway_group_v2" "test" {
  name = "test-group"
  spec = {
    external_groups = ["ldap-admins", "oidc-devs"]
  }
}
