package console_group_v2

import (
	"context"
	"testing"

	ctlresource "github.com/conduktor/ctl/resource"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestGroupV2ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonGroupV2Resource := []byte(test.TestAccTestdata(t, "console/group_v2/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonGroupV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "Group", ctlResource.Kind)
	assert.Equal(t, "v2", ctlResource.Version)
	assert.Equal(t, "sales", ctlResource.Name)
	assert.Equal(t, map[string]any{"name": "sales"}, ctlResource.Metadata)
	assert.Equal(t, jsonGroupV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := console.NewGroupConsoleResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "Group", internal.Kind)
	assert.Equal(t, "v2", internal.ApiVersion)
	assert.Equal(t, "sales", internal.Metadata.Name)
	assert.Equal(t, "Sales Department", internal.Spec.DisplayName)
	assert.Equal(t, "Sales Department Group", internal.Spec.Description)
	assert.Equal(t, []string{"sales"}, internal.Spec.ExternalGroups)
	assert.Equal(t, []string{".*"}, internal.Spec.ExternalGroupRegex)
	assert.Equal(t, []string{"jim.halpert@dunder.mifflin.com", "dwight.schrute@dunder.mifflin.com"}, internal.Spec.Members)
	expectedInternalPermissions := []model.Permission{
		{
			ResourceType: "PLATFORM",
			Permissions:  []string{"groupView", "clusterConnectionsManage"},
		},
		{
			ResourceType: "CLUSTER",
			Name:         "scranton",
			Permissions:  []string{"clusterViewBroker", "clusterEditBroker"},
		},
		{
			ResourceType: "TOPIC",
			Name:         "sales-*",
			PatternType:  "PREFIXED",
			Cluster:      "scranton",
			Permissions:  []string{"topicViewConfig", "topicConsume", "topicProduce"},
		},
		{
			ResourceType: "SUBJECT",
			Name:         "sales-*",
			PatternType:  "PREFIXED",
			Cluster:      "scranton",
			Permissions:  []string{"subjectView", "subjectEditCompatibility"},
		},
		{
			ResourceType: "CONSUMER_GROUP",
			Name:         "sales-*",
			PatternType:  "PREFIXED",
			Cluster:      "scranton",
			Permissions:  []string{"consumerGroupView"},
		},
		{
			ResourceType: "KAFKA_CONNECT",
			Name:         "sales-*",
			PatternType:  "PREFIXED",
			KafkaConnect: "scranton",
			Cluster:      "scranton",
			Permissions:  []string{"subjectView", "kafkaConnectorDelete"},
		},
		{
			ResourceType: "KSQLDB",
			Name:         "sales-ksqldb",
			Cluster:      "scranton",
			Permissions:  []string{"ksqldbAccess"},
		},
	}
	assert.Equal(t, expectedInternalPermissions, internal.Spec.Permissions)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("sales"), tfModel.Name)
	assert.Equal(t, types.StringValue("Sales Department"), tfModel.Spec.DisplayName)
	assert.Equal(t, types.StringValue("Sales Department Group"), tfModel.Spec.Description)
	// do not test permission as it's a pain to parse ListValue

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "Group", internal2.Kind)
	assert.Equal(t, "v2", internal2.ApiVersion)
	assert.Equal(t, "sales", internal2.Metadata.Name)
	assert.Equal(t, "Sales Department", internal2.Spec.DisplayName)
	assert.Equal(t, "Sales Department Group", internal2.Spec.Description)
	assert.Equal(t, []string{"sales"}, internal2.Spec.ExternalGroups)
	assert.Equal(t, []string{".*"}, internal2.Spec.ExternalGroupRegex)
	assert.Equal(t, []string{"jim.halpert@dunder.mifflin.com", "dwight.schrute@dunder.mifflin.com"}, internal2.Spec.Members)
	assert.Equal(t, expectedInternalPermissions, internal2.Spec.Permissions)
	assert.Equal(t, internal, internal2)

	// convert back to ctl model
	ctlResource2, err := internal2.ToClientResource()
	if err != nil {
		t.Fatal(err)
		return
	}
	// compare without json
	if !cmp.Equal(ctlResource, ctlResource2, cmpopts.IgnoreFields(ctlresource.Resource{}, "Json")) {
		t.Errorf("expected %+v, got %+v", ctlResource, ctlResource2)
	}
}
