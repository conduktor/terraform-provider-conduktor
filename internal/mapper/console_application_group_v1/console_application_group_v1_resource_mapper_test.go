package console_application_group_v1

import (
	"testing"

	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestApplicationGroupV1ModelMapping(t *testing.T) {
	// ctx := context.Background()

	jsonApplicationGroupV1Resource := []byte(test.TestAccTestdata(t, "console/application_group_v1/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonApplicationGroupV1Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ApplicationGroup", ctlResource.Kind)
	assert.Equal(t, "v1", ctlResource.Version)
	assert.Equal(t, "test-application-group", ctlResource.Name)
	assert.Equal(t, map[string]interface{}{"name": "test-application-group", "application": "test-application"}, ctlResource.Metadata)
	assert.Equal(t, jsonApplicationGroupV1Resource, ctlResource.Json)

	// convert into internal model
	internal, err := console.NewApplicationGroupConsoleResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ApplicationGroup", internal.Kind)
	assert.Equal(t, "v1", internal.ApiVersion)
	assert.Equal(t, "test-application-group", internal.Metadata.Name)
	assert.Equal(t, "test-application", internal.Metadata.Application)
	assert.Equal(t, "Test Application Group", internal.Spec.DisplayName)
	assert.Equal(t, "A great test application group", internal.Spec.Description)
	assert.Equal(t, []string{"COMPANY-SUPPORT"}, internal.Spec.ExternalGroups)
	assert.Equal(t, []string{"tatum@conduktor.io"}, internal.Spec.Members)
	expectedInternalResources := []model.ApplicationGroupPermission{
		{
			AppInstance:    "test-application-dev",
			PatternType:    "LITERAL",
			ConnectCluster: "kafka-connect",
			Name:           "*",
			Permissions:    []string{"kafkaConnectPauseResume", "kafkaConnectRestart", "kafkaConnectorStatus", "kafkaConnectorViewConfig"},
			ResourceType:   "CONNECTOR",
		},
		{
			AppInstance:  "test-application-dev",
			PatternType:  "LITERAL",
			Name:         "*",
			Permissions:  []string{"consumerGroupCreate", "consumerGroupDelete", "consumerGroupReset", "consumerGroupView"},
			ResourceType: "CONSUMER_GROUP",
		},
		{
			AppInstance:  "test-application-dev",
			PatternType:  "LITERAL",
			Name:         "*",
			Permissions:  []string{"topicConsume", "topicViewConfig"},
			ResourceType: "TOPIC",
		},
	}
	assert.Equal(t, expectedInternalResources, internal.Spec.Permissions)

}
