resource "conduktor_gateway_alias_topic_v2" "test" {
  name     = "test-alias"
  vcluster = "vcluster_sa"
  spec = {
    physical_name = "my-topic"
  }
}
