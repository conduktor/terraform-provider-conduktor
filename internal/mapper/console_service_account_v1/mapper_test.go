package console_service_account_v1

import (
	"context"
	"testing"

	ctlresource "github.com/conduktor/ctl/resource"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	"github.com/conduktor/terraform-provider-conduktor/internal/schema"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestServiceAccountV1ModelMapping(t *testing.T) {
	ctx := context.Background()

	jsonServiceAccountV1Resource := []byte(test.TestAccTestdata(t, "console/service_account_v1/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonServiceAccountV1Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ServiceAccount", ctlResource.Kind)
	assert.Equal(t, "v1", ctlResource.Version)
	assert.Equal(t, "sa-clicko-dev", ctlResource.Name)
	assert.Equal(t, map[string]any{"name": "sa-clicko-dev", "appInstance": "my-app-instance-dev", "cluster": "my-kafka-cluster", "labels": map[string]any{"key": "value"}}, ctlResource.Metadata)
	assert.Equal(t, jsonServiceAccountV1Resource, ctlResource.Json)

	// convert into internal model
	internal, err := console.NewServiceAccountResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ServiceAccount", internal.Kind)
	assert.Equal(t, "v1", internal.ApiVersion)
	assert.Equal(t, "sa-clicko-dev", internal.Metadata.Name)
	assert.Equal(t, "my-app-instance-dev", internal.Metadata.AppInstance)
	assert.Equal(t, "my-kafka-cluster", internal.Metadata.Cluster)
	assert.Equal(t, map[string]string{
		"key": "value",
	}, internal.Metadata.Labels)
	expectedACLs := []console.ServiceAccountAuthKafkaACL{
		{
			Type:           "TOPIC",
			Name:           "click.",
			PatternType:    "PREFIXED",
			ConnectCluster: "my-connect-cluster",
			Operations:     []string{"Write"},
			Host:           "*",
			Permission:     "Allow",
		},
	}
	assert.Equal(t, expectedACLs, internal.Spec.Authorization.Kafka.ACLS)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	labels, _ := schema.StringMapToMapValue(ctx, map[string]string{"key": "value"})
	assert.Equal(t, types.StringValue("sa-clicko-dev"), tfModel.Name)
	assert.Equal(t, types.StringValue("my-app-instance-dev"), tfModel.AppInstance)
	assert.Equal(t, types.StringValue("my-kafka-cluster"), tfModel.Cluster)
	assert.Equal(t, labels, tfModel.Labels)
	assert.Equal(t, false, tfModel.Spec.Authorization.IsNull())
	assert.Equal(t, false, tfModel.Spec.Authorization.IsUnknown())

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ServiceAccount", internal2.Kind)
	assert.Equal(t, "v1", internal2.ApiVersion)
	assert.Equal(t, "sa-clicko-dev", internal2.Metadata.Name)
	assert.Equal(t, "my-app-instance-dev", internal2.Metadata.AppInstance)
	assert.Equal(t, "my-kafka-cluster", internal2.Metadata.Cluster)
	assert.Equal(t, map[string]string{
		"key": "value",
	}, internal2.Metadata.Labels)
	assert.Equal(t, expectedACLs, internal2.Spec.Authorization.Kafka.ACLS)

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
