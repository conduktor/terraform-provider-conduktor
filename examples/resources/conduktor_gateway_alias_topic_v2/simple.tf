resource "conduktor_gateway_alias_topic_v2" "simple_alias" {
  name     = "simple-alias"
  vcluster = "passthrough"
  spec = {
    physical_name = "my-topic"
  }
}
