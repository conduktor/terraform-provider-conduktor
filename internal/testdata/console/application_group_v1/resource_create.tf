resource "conduktor_console_application_group_v1" "test" {
  name = "test-application-group"
  application = "test-application"
  spec = {
    display_name = "Test Application Group"
    description = "A great test application group"
    members = "tatum@conduktor.io"
    external_group = "COMPANY-SUPPORT"
    permissions = [
        {
            app_instance = "test-application-dev"
            pattern_type = "LITERAL"
            connect_cluster = "kafka-connect"
            name = "*"
            permissions = [
                "kafkaConnectPauseResume",
                "kafkaConnectRestart",
                "kafkaConnectorStatus",
                "kafkaConnectorViewConfig"
            ]
            resource_type = "CONNECTOR"
        },
        {
            app_instance = "test-application-dev"
            pattern_type = "LITERAL"
            connect_cluster = "kafka-connect"
            name = "*"
            permissions = [
                "consumerGroupCreate",
                "consumerGroupDelete",
                "consumerGroupReset",
                "consumerGroupView"
            ]
            resource_type = "CONSUMER_GROUP"
        },
        {
            app_instance = "test-application-dev"
            pattern_type = "LITERAL"
            name = "*"
            permissions = [
                "topicConsume",
                "topicViewConfig"
            ]
            resource_type = "TOPIC"
        }
    ]
  }
}