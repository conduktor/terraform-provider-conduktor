resource "conduktor_gateway_alias_topic_v2" "complex_alias" {
  name     = "complex-alias"
  vcluster = "vcluster_sa"
  spec = {
    physical_name    = "my-topic"
    physical_cluster = "main"
  }
}
