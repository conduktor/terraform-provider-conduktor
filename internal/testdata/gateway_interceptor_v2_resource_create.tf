
resource "conduktor_gateway_interceptor_v2" "test" {
  name = "yellow_cars_filter"
  scope {
    vcluster = "vcluster_sa"
  }

  spec {
    plugin_class = "io.conduktor.gateway.interceptor.VirtualSqlTopicPlugin"
    priority     = 1
    config = {
      virtual_topic = "yellow_cars"
      statement     = "SELECT '$.type' as type, '$.price' as price FROM cars WHERE '$.color' = 'yellow'"
    }
  }
}
