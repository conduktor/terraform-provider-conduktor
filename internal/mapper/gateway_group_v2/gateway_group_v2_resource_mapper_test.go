package gateway_group_v2

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

func TestGatewayGroupV2ModelMapping(t *testing.T) {
	ctx := context.Background()

	jsonGroupResource := []byte(test.TestAccTestdata(t, "gateway/group_v2/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonGroupResource)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "GatewayGroup", ctlResource.Kind)
	assert.Equal(t, "gateway/v2", ctlResource.Version)
	assert.Equal(t, "test-group", ctlResource.Name)

	// ctlresource → internal model
	internal, err := gateway.NewGatewayGroupResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "GatewayGroup", internal.Kind)
	assert.Equal(t, "gateway/v2", internal.ApiVersion)
	assert.Equal(t, "test-group", internal.Metadata.Name)
	assert.Len(t, internal.Spec.Members, 1)
	assert.Equal(t, "sa1", internal.Spec.Members[0].Name)
	assert.Equal(t, "vc1", internal.Spec.Members[0].VCluster)
	assert.Equal(t, []string{"ldap-admins"}, internal.Spec.ExternalGroups)

	// internal → TF model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, types.StringValue("test-group"), tfModel.Name)

	// TF model → internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "GatewayGroup", internal2.Kind)
	assert.Equal(t, "gateway/v2", internal2.ApiVersion)
	assert.Equal(t, "test-group", internal2.Metadata.Name)
	assert.Len(t, internal2.Spec.Members, 1)
	assert.Equal(t, "sa1", internal2.Spec.Members[0].Name)
	assert.Equal(t, "vc1", internal2.Spec.Members[0].VCluster)
	assert.Equal(t, []string{"ldap-admins"}, internal2.Spec.ExternalGroups)

	// internal → ctl resource (round-trip)
	ctlResource2, err := internal2.ToClientResource()
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(ctlResource, ctlResource2, cmpopts.IgnoreFields(ctlresource.Resource{}, "Json")) {
		t.Errorf("expected %+v, got %+v", ctlResource, ctlResource2)
	}
}
