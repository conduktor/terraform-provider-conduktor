package gateway_token_v2

import (
	"context"
	"testing"

	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestGatewayTokenV2ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonTokenV2Resource := []byte(test.TestAccTestdata(t, "gateway_token_v2_api.json"))

	token := gateway.GatewayTokenResource{}
	err := token.FromRawJson(jsonTokenV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "passthrough", token.VCluster)
	assert.Equal(t, "user1", token.Username)
	assert.Equal(t, int64(3600), token.LifetimeSeconds)

	// convert into internal model
	internal := gateway.NewGatewayTokenResource(token.VCluster, token.Username, token.LifetimeSeconds)

	assert.Equal(t, "passthrough", internal.VCluster)
	assert.Equal(t, "user1", internal.Username)
	assert.Equal(t, int64(3600), internal.LifetimeSeconds)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("passthrough"), tfModel.Vcluster)
	assert.Equal(t, types.StringValue("user1"), tfModel.Username)
	assert.Equal(t, types.Int64Value(3600), tfModel.LifetimeSeconds)

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "passthrough", internal2.VCluster)
	assert.Equal(t, "user1", internal2.Username)
	assert.Equal(t, int64(3600), internal2.LifetimeSeconds)

	// compare with original
	if !cmp.Equal(token, internal2) {
		t.Errorf("expected %+v, got %+v", token, internal2)
	}
}
