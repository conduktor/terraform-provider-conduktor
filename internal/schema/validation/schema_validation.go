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

var ConfluentKafkaFlavor = "Confluent"
var AivenKafkaFlavor = "Aiven"
var GatewayKafkaFlavor = "Gateway"
var ValidKafkaFlavorTypes = []string{ConfluentKafkaFlavor, AivenKafkaFlavor, GatewayKafkaFlavor}

var ConfluentLikeSchemaRegistry = "ConfluentLike"
var GlueSchemaRegistry = "Glue"
var ValidSchemaRegistryTypes = []string{ConfluentLikeSchemaRegistry, GlueSchemaRegistry}

// ConfluentLike security.
var NoSecuritySchemaRegistrySecurity = "NoSecurity"
var BasicAuthSchemaRegistrySecurity = "BasicAuth"
var BearerTokenSchemaRegistrySecurity = "BearerToken"
var SSLAuthSchemaRegistrySecurity = "SSLAuth"

// Glue security.
var CredentialsSchemaRegistrySecurity = "Credentials"
var FromContextSchemaRegistrySecurity = "FromContext"
var FromRoleSchemaRegistrySecurity = "FromRole"
var IAMAnywhereSchemaRegistrySecurity = "IAMAnywhere"

var ValidSchemaRegistrySecurityTypes = []string{
	BasicAuthSchemaRegistrySecurity,
	BearerTokenSchemaRegistrySecurity,
	SSLAuthSchemaRegistrySecurity,
	NoSecuritySchemaRegistrySecurity,
	CredentialsSchemaRegistrySecurity,
	FromContextSchemaRegistrySecurity,
	FromRoleSchemaRegistrySecurity,
	IAMAnywhereSchemaRegistrySecurity,
}

// KafkaConnect security.
var BasicAuthKafkaConnectSecurity = "BasicAuth"
var BearerTokenKafkaConnectSecurity = "BearerToken"
var SSLAuthKafkaConnectSecurity = "SSLAuth"

var ValidKafkaConnectSecurityTypes = []string{
	BasicAuthKafkaConnectSecurity,
	BearerTokenKafkaConnectSecurity,
	SSLAuthKafkaConnectSecurity,
}

// Gateway Service Accounlts
var ValidServiceAccountTypes = []string{"LOCAL", "EXTERNAL"}
