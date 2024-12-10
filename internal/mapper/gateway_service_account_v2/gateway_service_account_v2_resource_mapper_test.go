package gateway_service_account_v2

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

func TestGatewayServiceAccountV2ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonServiceAccountV2Resource := []byte(test.TestAccTestdata(t, "gateway_service_account_v2_api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonServiceAccountV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "GatewayServiceAccount", ctlResource.Kind)
	assert.Equal(t, "gateway/v2", ctlResource.Version)
	assert.Equal(t, "user1", ctlResource.Name)
	assert.Equal(t, map[string]interface{}{"name": "user1", "vCluster": "vcluster1"}, ctlResource.Metadata)
	assert.Equal(t, jsonServiceAccountV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := gateway.NewGatewayServiceAccountResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "GatewayServiceAccount", internal.Kind)
	assert.Equal(t, "gateway/v2", internal.ApiVersion)
	assert.Equal(t, "user1", internal.Metadata.Name)
	assert.Equal(t, "EXTERNAL", internal.Spec.Type)
	assert.Equal(t, []string{"externalName"}, internal.Spec.ExternalNames)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("user1"), tfModel.Name)
	assert.Equal(t, types.StringValue("EXTERNAL"), tfModel.Spec.SpecType)
	// do not test ExternalNames as it's a pain to parse SetValue

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "GatewayServiceAccount", internal2.Kind)
	assert.Equal(t, "gateway/v2", internal2.ApiVersion)
	assert.Equal(t, "user1", internal2.Metadata.Name)
	assert.Equal(t, "vcluster1", internal2.Metadata.VCluster)
	assert.Equal(t, "EXTERNAL", internal2.Spec.Type)
	assert.Equal(t, []string{"externalName"}, internal2.Spec.ExternalNames)

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
