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
	"datamaskingManage",
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
	"topicDataQualityManage",
	"chargebackManage",
	"sqlManage",
}

var ValidAppGroupPermissions = []string{
	"kafkaConnectorViewConfig",
	"kafkaConnectorStatus",
	"kafkaConnectorEditConfig",
	"kafkaConnectorDelete",
	"kafkaConnectorCreate",
	"kafkaConnectPauseResume",
	"kafkaConnectRestart",
	"consumerGroupView",
	"consumerGroupReset",
	"consumerGroupDelete",
	"consumerGroupCreate",
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
	"topicDataQualityManage",
}

// Provider modes.
var ValidProviderMode = []string{"console", "gateway"}

// Console Application Instance.
var ValidResourceTypes = []string{"TOPIC", "CONSUMER_GROUP", "SUBJECT", "CONNECTOR"}
var ValidPatternTypes = []string{"LITERAL", "PREFIXED"} // This seems to be a duplicate of ValidPermissionPatternTypes, keeping separate for now to avoid confusion.
var ValidOwnershipModes = []string{"ALL", "LIMITED"}
var ValidCatalogVisibilities = []string{"PRIVATE", "PUBLIC"}

var ValidPermissionTypes = []string{"CLUSTER", "CONSUMER_GROUP", "KAFKA_CONNECT", "KSQLDB", "PLATFORM", "SUBJECT", "TOPIC"}
var ValidAppGroupPermissionTypes = []string{"CONNECTOR", "CONSUMER_GROUP", "TOPIC", "SUBJECT"}
var ValidPermissionPatternTypes = []string{"LITERAL", "PREFIXED"}

var ValidResourcePolicyTargetKind = []string{"ApplicationGroup", "Connector", "Subject", "Topic"}

var ValidApplicationInstancePermission = []string{"NONE", "READ", "WRITE"}

var ValidPartnerZoneAuthenticationType = []string{"MTLS", "OAUTHBEARER", "PLAIN"}
var ValidPartnerZoneTopicPermission = []string{"READ", "WRITE"}

// Console Service Account.
var ValidServiceAccountAivenResourceType = []string{"SCHEMA", "TOPIC"}
var ValidServiceAccountAivenPermission = []string{"admin", "read", "readwrite", "schema_registry_read", "schema_registry_write", "write"}
var ValidServiceAccountKafkaType = []string{"CLUSTER", "CONSUMER_GROUP", "DELEGATION_TOKEN", "TOPIC", "TRANSACTIONAL_ID", "UNKNOWN", "USER"}
var ValidServiceAccountKafkaPatternType = []string{"LITERAL", "PREFIXED"}
var ValidServiceAccountKafkaOperations = []string{"All", "Alter", "AlterConfigs", "ClusterAction", "Create", "CreateTokens", "Delete", "Describe", "DescribeConfigs", "DescribeTokens", "IdempotentWrite", "Read", "Unknown", "Write"}
var ValidServiceAccountKafkaPermission = []string{"Deny", "Allow"}

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

// Gateway Service Accounts.
var ValidServiceAccountTypes = []string{"LOCAL", "EXTERNAL"}

// Gateway Virtual Clusters.
var ValidVirtualClusterACLModes = []string{"KAFKA_API", "REST_API"}
var ValidVirtualClusterTypes = []string{"Standard", "Partner"}
var ValidVirtualClusterACLsResourceTypes = []string{"ANY", "CLUSTER", "DELEGATION_TOKEN", "GROUP", "TOPIC", "TRANSACTIONAL_ID", "UNKNOWN"}
var ValidVirtualClusterACLsPatternTypes = []string{"ANY", "LITERAL", "MATCH", "PREFIXED", "UNKNOWN"}
var ValidVirtualClusterACLsOperations = []string{"ALL", "ALTER", "ALTER_CONFIGS", "CLUSTER_ACTION", "CREATE", "CREATE_TOKENS", "DELETE", "DESCRIBE", "DESCRIBE_CONFIGS", "DESCRIBE_TOKENS", "IDEMPOTENT_WRITE", "READ", "UNKNOWN", "WRITE"}
var ValidVirtualClusterACLsPermissionTypes = []string{"ALLOW", "DENY"}
