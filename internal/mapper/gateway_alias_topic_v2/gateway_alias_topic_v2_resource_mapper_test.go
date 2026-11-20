package gateway_alias_topic_v2

import (
	"context"
	"testing"

	ctlresource "github.com/conduktor/ctl/resource"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestGatewayAliasTopicV2ModelMapping(t *testing.T) {
	ctx := context.Background()

	jsonAliasTopicV2Resource := []byte(test.TestAccTestdata(t, "gateway/alias_topic_v2/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonAliasTopicV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "AliasTopic", ctlResource.Kind)
	assert.Equal(t, "gateway/v2", ctlResource.Version)
	assert.Equal(t, "my-alias", ctlResource.Name)
	assert.Equal(t, map[string]any{"name": "my-alias", "vCluster": "vcluster_sa"}, ctlResource.Metadata)
	assert.Equal(t, jsonAliasTopicV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := gateway.NewGatewayAliasTopicResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "AliasTopic", internal.Kind)
	assert.Equal(t, "gateway/v2", internal.ApiVersion)
	assert.Equal(t, "my-alias", internal.Metadata.Name)
	assert.Equal(t, "vcluster_sa", internal.Metadata.VCluster)
	assert.Equal(t, "real-topic", internal.Spec.PhysicalName)
	assert.Equal(t, "main", internal.Spec.PhysicalCluster)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("my-alias"), tfModel.Name)
	assert.Equal(t, types.StringValue("vcluster_sa"), tfModel.Vcluster)
	assert.Equal(t, types.StringValue("real-topic"), tfModel.Spec.PhysicalName)
	assert.Equal(t, types.StringValue("main"), tfModel.Spec.PhysicalCluster)

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "AliasTopic", internal2.Kind)
	assert.Equal(t, "gateway/v2", internal2.ApiVersion)
	assert.Equal(t, "my-alias", internal2.Metadata.Name)
	assert.Equal(t, "vcluster_sa", internal2.Metadata.VCluster)
	assert.Equal(t, "real-topic", internal2.Spec.PhysicalName)
	assert.Equal(t, "main", internal2.Spec.PhysicalCluster)

	// convert back to ctl model
	ctlResource2, err := internal2.ToClientResource()
	if err != nil {
		t.Fatal(err)
		return
	}
	// compare without json field (json bytes may differ in key order)
	if !cmp.Equal(ctlResource, ctlResource2, cmpopts.IgnoreFields(ctlresource.Resource{}, "Json")) {
		t.Errorf("expected %+v, got %+v", ctlResource, ctlResource2)
	}
}
