package console_user_v2

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

func TestUserV2ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonUserV2Resource := []byte(test.TestAccTestdata(t, "console_user_v2_api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonUserV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "User", ctlResource.Kind)
	assert.Equal(t, "v2", ctlResource.Version)
	assert.Equal(t, "michael.scott@dunder.mifflin.com", ctlResource.Name)
	assert.Equal(t, map[string]interface{}{"name": "michael.scott@dunder.mifflin.com"}, ctlResource.Metadata)
	assert.Equal(t, jsonUserV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := console.NewUserConsoleResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "User", internal.Kind)
	assert.Equal(t, "v2", internal.ApiVersion)
	assert.Equal(t, "michael.scott@dunder.mifflin.com", internal.Metadata.Name)
	assert.Equal(t, "Michael", internal.Spec.FirstName)
	assert.Equal(t, "Scott", internal.Spec.LastName)
	expectedInternalPermissions := []model.Permission{
		{
			ResourceType: "PLATFORM",
			Permissions:  []string{"userView", "clusterConnectionsManage"},
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
	assert.Equal(t, types.StringValue("michael.scott@dunder.mifflin.com"), tfModel.Name)
	assert.Equal(t, types.StringValue("Michael"), tfModel.Spec.Firstname)
	assert.Equal(t, types.StringValue("Scott"), tfModel.Spec.Lastname)
	// do not test permission as it's a pain to parse ListValue

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "User", internal2.Kind)
	assert.Equal(t, "v2", internal2.ApiVersion)
	assert.Equal(t, "michael.scott@dunder.mifflin.com", internal2.Metadata.Name)
	assert.Equal(t, "Michael", internal2.Spec.FirstName)
	assert.Equal(t, "Scott", internal2.Spec.LastName)
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
