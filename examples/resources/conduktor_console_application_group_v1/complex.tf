
resource "conduktor_console_application_group_v1" "complex" {
  name        = "complex"
  application = "myapp"
  spec = {
    display_name = "Complex Application Group"
    description  = "Complex Description"
    permissions = [
      {
        app_instance  = "my-app-instance"
        resource_type = "TOPIC"
        pattern_type  = "LITERAL"
        name          = "*"
        permissions   = ["topicViewConfig", "topicConsume"]
      },
      {
        app_instance  = "my-app-instance"
        resource_type = "CONSUMER_GROUP"
        pattern_type  = "LITERAL"
        name          = "*"
        permissions   = ["consumerGroupCreate", "consumerGroupReset", "consumerGroupDelete", "consumerGroupView"]
      },
      {
        app_instance    = "my-app-instance"
        resource_type   = "CONNECTOR"
        pattern_type    = "LITERAL"
        name            = "*"
        connect_cluster = "kafka-connect"
        permissions     = ["kafkaConnectorViewConfig", "kafkaConnectorStatus", "kafkaConnectRestart"]
      },
    ]
    members = [
      "user1@company.org",
      "user2@company.org"
    ]
    external_groups = ["GP-COMPANY-CLICKSTREAM-SUPPORT"]
  }
}

