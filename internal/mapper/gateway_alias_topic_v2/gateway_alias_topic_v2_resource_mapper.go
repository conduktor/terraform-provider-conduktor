package gateway_alias_topic_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	gwaliasTopics "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_alias_topic_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *gwaliasTopics.GatewayAliasTopicV2Model) (gateway.GatewayAliasTopicResource, error) {
	return gateway.NewGatewayAliasTopicResource(
		gateway.GatewayAliasTopicMetadata{
			Name:     r.Name.ValueString(),
			VCluster: r.Vcluster.ValueString(),
		},
		gateway.GatewayAliasTopicSpec{
			PhysicalName:    r.Spec.PhysicalName.ValueString(),
			PhysicalCluster: r.Spec.PhysicalCluster.ValueString(),
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *gateway.GatewayAliasTopicResource) (gwaliasTopics.GatewayAliasTopicV2Model, error) {
	specValue, diag := gwaliasTopics.NewSpecValue(
		map[string]attr.Type{
			"physical_name":    basetypes.StringType{},
			"physical_cluster": basetypes.StringType{},
		},
		map[string]attr.Value{
			"physical_name":    schema.NewStringValue(r.Spec.PhysicalName),
			"physical_cluster": schema.NewStringValue(r.Spec.PhysicalCluster),
		},
	)
	if diag.HasError() {
		return gwaliasTopics.GatewayAliasTopicV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	vcluster := r.Metadata.VCluster
	if vcluster == "" {
		vcluster = "passthrough"
	}

	return gwaliasTopics.GatewayAliasTopicV2Model{
		Name:     types.StringValue(r.Metadata.Name),
		Vcluster: types.StringValue(vcluster),
		Spec:     specValue,
	}, nil
}
