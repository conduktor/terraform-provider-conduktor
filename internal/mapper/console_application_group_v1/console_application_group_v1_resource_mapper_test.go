package console_application_group_v1

import (
	"context"
	"testing"

	ctlresource "github.com/conduktor/ctl/resource"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestApplicationGroupV1ModelMapping(t *testing.T) {
	ctx := context.Background()

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
	expectedInternalResources := []console.ApplicationGroupPermission{
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

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	externalGroups := []string{}
	tfModel.Spec.ExternalGroups.ElementsAs(ctx, &externalGroups, false)
	members := []string{}
	tfModel.Spec.Members.ElementsAs(ctx, &members, false)

	assert.Equal(t, types.StringValue("test-application-group"), tfModel.Name)
	assert.Equal(t, types.StringValue("test-application"), tfModel.Application)
	assert.Equal(t, types.StringValue("Test Application Group"), tfModel.Spec.DisplayName)
	assert.Equal(t, types.StringValue("A great test application group"), tfModel.Spec.Description)
	assert.Equal(t, []string{"COMPANY-SUPPORT"}, externalGroups)
	assert.Equal(t, []string{"tatum@conduktor.io"}, members)

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ApplicationGroup", internal2.Kind)
	assert.Equal(t, "v1", internal2.ApiVersion)
	assert.Equal(t, "test-application-group", internal2.Metadata.Name)
	assert.Equal(t, "test-application", internal2.Metadata.Application)
	assert.Equal(t, "Test Application Group", internal2.Spec.DisplayName)
	assert.Equal(t, "A great test application group", internal2.Spec.Description)
	assert.Equal(t, []string{"COMPANY-SUPPORT"}, internal2.Spec.ExternalGroups)
	assert.Equal(t, []string{"tatum@conduktor.io"}, internal2.Spec.Members)
	// assert.Equal(t, expectedInternalResources, internal2.Spec.Permissions)
	// assert.Equal(t, internal, internal2)

	// // convert back to ctl model
	// ctlResource2, err := internal2.ToClientResource()
	// if err != nil {
	// 	t.Fatal(err)
	// 	return
	// }
	// // compare without json
	// if !cmp.Equal(ctlResource, ctlResource2, cmpopts.IgnoreFields(ctlresource.Resource{}, "Json")) {
	// 	t.Errorf("expected %+v, got %+v", ctlResource, ctlResource2)
	// }
}
