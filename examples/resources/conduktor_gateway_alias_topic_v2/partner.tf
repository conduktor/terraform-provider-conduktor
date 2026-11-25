resource "conduktor_gateway_virtual_cluster_v2" "partner" {
  name = "partner"
  spec = {
    type = "Partner"
  }
}

resource "conduktor_gateway_alias_topic_v2" "partner_alias" {
  name     = "trades"
  vcluster = conduktor_gateway_virtual_cluster_v2.partner.name
  spec = {
    physical_name    = "internal.trades.v1"
    physical_cluster = "main"
  }
}
