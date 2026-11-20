resource "conduktor_gateway_alias_topic_v2" "minimal" {
  name     = "minimal-alias"
  vcluster = "passthrough"
  spec = {
    physical_name = "real-topic"
  }
}
