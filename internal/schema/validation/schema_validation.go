package validation

var ValidPermissions = []string{
	"clusterViewBroker",
	"clusterEditSRCompatibility",
	"clusterEditBroker",
	"clusterViewACL",
	"clusterManageACL",
	"kafkaConnectorViewConfig",
	"kafkaConnectorStatus",
	"kafkaConnectorEditConfig",
	"kafkaConnectorDelete",
	"kafkaConnectorCreate",
	"kafkaConnectPauseResume",
	"kafkaConnectRestart",
	"ksqldbAccess",
	"consumerGroupView",
	"consumerGroupReset",
	"consumerGroupDelete",
	"consumerGroupCreate",
	"auditLogView",
	"taasView",
	"certificateManage",
	"userManage",
	"clusterConnectionsManage",
	"notificationChannelManage",
	"datamaskingView",
	"userView",
	"testingView",
	"datamaskingManage",
	"taasManage",
	"notificationChannelView",
	"subjectCreateUpdate",
	"subjectEditCompatibility",
	"subjectDelete",
	"subjectView",
	"topicViewConfig",
	"topicEmpty",
	"topicConsume",
	"topicProduce",
	"topicEditConfig",
	"topicCreate",
	"topicAddPartition",
	"topicDelete",
}
var ValidPermissionTypes = []string{"CLUSTER", "CONSUMER_GROUP", "KAFKA_CONNECT", "KSQLDB", "PLATFORM", "SUBJECT", "TOPIC"}
var ValidPermissionPatternTypes = []string{"LITERAL", "PREFIXED"}

var ValidKafkaFlavorTypes = []string{"Confluent", "Aiven", "Gateway"}
var ValidSchemaRegistryTypes = []string{"ConfluentLike", "Glue"}
var ValidSchemaRegistrySecurityTypes = []string{
	// ConfluentLike security
	"BasicAuth", "BearerToken", "SSLAuth", "NoSecurity",
	// Glue security
	"Credentials", "FromContext", "FromRole", "IAMAnywhere",
}
