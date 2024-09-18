
resource "conduktor_generic" "test" {
  kind     = "KafkaCluster"
  version  = "v2"
  name     = "cluserA"
  manifest = file("clusterA.yaml")
}
