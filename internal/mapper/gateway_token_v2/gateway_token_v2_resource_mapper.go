package gateway_token_v2

import (
	"context"

	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	token "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_token_v2"
)

func TFToInternalModel(ctx context.Context, r *token.GatewayTokenV2Model) (gateway.GatewayTokenResource, error) {
	vCluster := r.Vcluster.ValueString()
	username := r.Username.ValueString()
	lifetimeSeconds := r.LifetimeSeconds.ValueInt64()

	return gateway.NewGatewayTokenResource(vCluster, username, lifetimeSeconds), nil
}

func InternalModelToTerraform(ctx context.Context, r *gateway.GatewayTokenResource) (token.GatewayTokenV2Model, error) {
	return token.GatewayTokenV2Model{
		Vcluster:        schema.NewStringValue(r.VCluster),
		Username:        schema.NewStringValue(r.Username),
		LifetimeSeconds: schema.NewInt64Value(r.LifetimeSeconds),
		Token:           schema.NewStringValue(r.Token),
	}, nil
}
