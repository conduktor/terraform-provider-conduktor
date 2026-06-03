package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeWithPlannedPermissions_KsqldbFieldPreserved(t *testing.T) {
	// Customer scenario: KSQLDB resource with ksqldb="*" set, but API strips ksqldb field
	planned := []Permission{
		{
			ResourceType: "KSQLDB",
			Name:         "*",
			Cluster:      "my-cluster",
			PatternType:  "LITERAL",
			KsqlDB:       "*",
			Permissions:  []string{"ksqldbAccess"},
		},
	}
	response := []Permission{
		{
			ResourceType: "KSQLDB",
			Name:         "*",
			Cluster:      "my-cluster",
			Permissions:  []string{"ksqldbAccess"},
			// API stripped: KsqlDB, PatternType
		},
	}

	merged := MergeWithPlannedPermissions(planned, response)

	assert.Equal(t, 1, len(merged))
	assert.Equal(t, "*", merged[0].KsqlDB)
	assert.Equal(t, "LITERAL", merged[0].PatternType)
	assert.Equal(t, "*", merged[0].Name)
	assert.Equal(t, "my-cluster", merged[0].Cluster)
}

func TestMergeWithPlannedPermissions_KafkaConnectFieldPreserved(t *testing.T) {
	planned := []Permission{
		{
			ResourceType: "KAFKA_CONNECT",
			Name:         "*",
			Cluster:      "my-cluster",
			PatternType:  "LITERAL",
			KafkaConnect: "*",
			Permissions:  []string{"kafkaConnectorStatus", "kafkaConnectorViewConfig"},
		},
	}
	response := []Permission{
		{
			ResourceType: "KAFKA_CONNECT",
			Name:         "*",
			PatternType:  "LITERAL",
			KafkaConnect: "*",
			Cluster:      "my-cluster",
			Permissions:  []string{"kafkaConnectorStatus", "kafkaConnectorViewConfig"},
		},
	}

	merged := MergeWithPlannedPermissions(planned, response)

	assert.Equal(t, 1, len(merged))
	assert.Equal(t, "*", merged[0].KafkaConnect)
	assert.Equal(t, "", merged[0].KsqlDB)
}

func TestMergeWithPlannedPermissions_ClusterFieldsPreserved(t *testing.T) {
	planned := []Permission{
		{
			ResourceType: "CLUSTER",
			Name:         "my-cluster",
			Cluster:      "my-cluster",
			PatternType:  "LITERAL",
			Permissions:  []string{"clusterViewACL", "clusterViewBroker"},
		},
	}
	response := []Permission{
		{
			ResourceType: "CLUSTER",
			Name:         "my-cluster",
			Permissions:  []string{"clusterViewACL", "clusterViewBroker"},
			// API stripped: Cluster, PatternType
		},
	}

	merged := MergeWithPlannedPermissions(planned, response)

	assert.Equal(t, 1, len(merged))
	assert.Equal(t, "my-cluster", merged[0].Cluster)
	assert.Equal(t, "LITERAL", merged[0].PatternType)
	assert.Equal(t, "my-cluster", merged[0].Name)
}

func TestMergeWithPlannedPermissions_MultiplePermissions(t *testing.T) {
	planned := []Permission{
		{
			ResourceType: "KSQLDB",
			Name:         "*",
			Cluster:      "my-cluster",
			PatternType:  "LITERAL",
			KsqlDB:       "*",
			Permissions:  []string{"ksqldbAccess"},
		},
		{
			ResourceType: "CLUSTER",
			Name:         "my-cluster",
			Cluster:      "my-cluster",
			PatternType:  "LITERAL",
			Permissions:  []string{"clusterViewACL", "clusterViewBroker"},
		},
		{
			ResourceType: "TOPIC",
			Name:         "sales-*",
			Cluster:      "scranton",
			PatternType:  "PREFIXED",
			Permissions:  []string{"topicViewConfig", "topicConsume"},
		},
	}
	response := []Permission{
		{
			ResourceType: "KSQLDB",
			Name:         "*",
			Cluster:      "my-cluster",
			Permissions:  []string{"ksqldbAccess"},
		},
		{
			ResourceType: "CLUSTER",
			Name:         "my-cluster",
			Permissions:  []string{"clusterViewACL", "clusterViewBroker"},
		},
		{
			ResourceType: "TOPIC",
			Name:         "sales-*",
			Cluster:      "scranton",
			PatternType:  "PREFIXED",
			Permissions:  []string{"topicViewConfig", "topicConsume"},
		},
	}

	merged := MergeWithPlannedPermissions(planned, response)

	assert.Equal(t, 3, len(merged))

	assert.Equal(t, "KSQLDB", merged[0].ResourceType)
	assert.Equal(t, "*", merged[0].KsqlDB)
	assert.Equal(t, "LITERAL", merged[0].PatternType)

	assert.Equal(t, "CLUSTER", merged[1].ResourceType)
	assert.Equal(t, "my-cluster", merged[1].Cluster)
	assert.Equal(t, "LITERAL", merged[1].PatternType)

	assert.Equal(t, "TOPIC", merged[2].ResourceType)
	assert.Equal(t, "sales-*", merged[2].Name)
	assert.Equal(t, "scranton", merged[2].Cluster)
	assert.Equal(t, "PREFIXED", merged[2].PatternType)
}

func TestMergeWithPlannedPermissions_NoPlannedPermissions(t *testing.T) {
	response := []Permission{
		{
			ResourceType: "TOPIC",
			Name:         "test",
			Cluster:      "cluster1",
			PatternType:  "LITERAL",
			Permissions:  []string{"topicViewConfig"},
		},
	}

	merged := MergeWithPlannedPermissions(nil, response)

	assert.Equal(t, response, merged)
}

func TestMergeWithPlannedPermissions_EmptyResponse(t *testing.T) {
	planned := []Permission{
		{
			ResourceType: "TOPIC",
			Name:         "test",
			Cluster:      "cluster1",
			Permissions:  []string{"topicViewConfig"},
		},
	}

	merged := MergeWithPlannedPermissions(planned, nil)

	assert.Equal(t, 0, len(merged))
}

func TestMergeWithPlannedPermissions_ResponseFieldsNotOverwritten(t *testing.T) {
	planned := []Permission{
		{
			ResourceType: "KAFKA_CONNECT",
			Name:         "planned-name",
			Cluster:      "planned-cluster",
			KafkaConnect: "planned-connect",
			PatternType:  "PREFIXED",
			Permissions:  []string{"kafkaConnectorDelete"},
		},
	}
	response := []Permission{
		{
			ResourceType: "KAFKA_CONNECT",
			Name:         "actual-name",
			Cluster:      "actual-cluster",
			KafkaConnect: "actual-connect",
			PatternType:  "LITERAL",
			Permissions:  []string{"kafkaConnectorDelete"},
		},
	}

	merged := MergeWithPlannedPermissions(planned, response)

	assert.Equal(t, 1, len(merged))
	assert.Equal(t, "actual-name", merged[0].Name)
	assert.Equal(t, "actual-cluster", merged[0].Cluster)
	assert.Equal(t, "actual-connect", merged[0].KafkaConnect)
	assert.Equal(t, "LITERAL", merged[0].PatternType)
}

func TestMergeWithPlannedPermissions_DuplicateResourceTypes(t *testing.T) {
	planned := []Permission{
		{
			ResourceType: "TOPIC",
			Name:         "topic-a",
			Cluster:      "cluster1",
			PatternType:  "LITERAL",
			Permissions:  []string{"topicViewConfig"},
		},
		{
			ResourceType: "TOPIC",
			Name:         "topic-b",
			Cluster:      "cluster1",
			PatternType:  "PREFIXED",
			Permissions:  []string{"topicConsume"},
		},
	}
	response := []Permission{
		{
			ResourceType: "TOPIC",
			Name:         "topic-a",
			Cluster:      "cluster1",
			PatternType:  "LITERAL",
			Permissions:  []string{"topicViewConfig"},
		},
		{
			ResourceType: "TOPIC",
			Name:         "topic-b",
			Cluster:      "cluster1",
			PatternType:  "PREFIXED",
			Permissions:  []string{"topicConsume"},
		},
	}

	merged := MergeWithPlannedPermissions(planned, response)

	assert.Equal(t, 2, len(merged))
	assert.Equal(t, "topic-a", merged[0].Name)
	assert.Equal(t, "topic-b", merged[1].Name)
}

func TestStringSlicesEqual(t *testing.T) {
	assert.True(t, stringSlicesEqual([]string{"a", "b"}, []string{"a", "b"}))
	assert.True(t, stringSlicesEqual([]string{"b", "a"}, []string{"a", "b"}))
	assert.True(t, stringSlicesEqual(nil, nil))
	assert.True(t, stringSlicesEqual([]string{}, []string{}))
	assert.False(t, stringSlicesEqual([]string{"a"}, []string{"b"}))
	assert.False(t, stringSlicesEqual([]string{"a"}, []string{"a", "b"}))
	assert.False(t, stringSlicesEqual(nil, []string{"a"}))
}
