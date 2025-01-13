package gateway_token_v2

import (
	"context"

	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	token "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_token_v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TFToInternalModel(ctx context.Context, r *token.GatewayTokenV2Model) (gateway.GatewayTokenResource, error) {
	vCluster := r.Vcluster.ValueString()
	username := r.Username.ValueString()
	lifetimeSeconds := r.LifetimeSeconds.ValueInt64()

	return gateway.NewGatewayTokenResource(vCluster, username, lifetimeSeconds), nil
}

func InternalModelToTerraform(ctx context.Context, r *gateway.GatewayTokenResource) (token.GatewayTokenV2Model, error) {
	// Configuring default value for vcluster
	if r.VCluster == "" {
		r.VCluster = "passthrough"
	}

	return token.GatewayTokenV2Model{
		Vcluster:        types.StringValue(r.VCluster),
		Username:        types.StringValue(r.Username),
		LifetimeSeconds: types.Int64Value(r.LifetimeSeconds),
	}, nil
}
