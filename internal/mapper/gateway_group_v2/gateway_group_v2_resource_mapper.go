package gateway_group_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_group_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *schema.GatewayGroupV2Model) (gateway.GatewayGroupResource, error) {
	members, err := listValueToMembers(ctx, r.Spec.Members)
	if err != nil {
		return gateway.GatewayGroupResource{}, err
	}

	externalGroups, diag := schemaUtils.SetValueToStringArray(ctx, r.Spec.ExternalGroups)
	if diag.HasError() {
		return gateway.GatewayGroupResource{}, mapper.WrapDiagError(diag, "external_groups", mapper.FromTerraform)
	}

	return gateway.NewGatewayGroupResource(
		gateway.GatewayGroupMetadata{
			Name: r.Name.ValueString(),
		},
		gateway.GatewayGroupSpec{
			Members:        members,
			ExternalGroups: externalGroups,
		},
	), nil
}

func listValueToMembers(ctx context.Context, list basetypes.ListValue) ([]gateway.GatewayServiceAccountId, error) {
	if list.IsNull() || list.IsUnknown() {
		return []gateway.GatewayServiceAccountId{}, nil
	}
	var membersValue []schema.MembersValue
	diags := list.ElementsAs(ctx, &membersValue, false)
	if diags.HasError() {
		return nil, mapper.WrapDiagError(diags, "members", mapper.FromTerraform)
	}
	result := make([]gateway.GatewayServiceAccountId, 0, len(membersValue))
	for _, m := range membersValue {
		result = append(result, gateway.GatewayServiceAccountId{
			Name:     m.Name.ValueString(),
			VCluster: m.Vcluster.ValueString(),
		})
	}
	return result, nil
}

func InternalModelToTerraform(ctx context.Context, r *gateway.GatewayGroupResource) (schema.GatewayGroupV2Model, error) {
	specValue, err := specInternalModelToTerraform(ctx, &r.Spec)
	if err != nil {
		return schema.GatewayGroupV2Model{}, err
	}
	return schema.GatewayGroupV2Model{
		Name: schemaUtils.NewStringValue(r.Metadata.Name),
		Spec: specValue,
	}, nil
}

func specInternalModelToTerraform(ctx context.Context, r *gateway.GatewayGroupSpec) (schema.SpecValue, error) {
	members, err := membersToListValue(ctx, r.Members)
	if err != nil {
		return schema.SpecValue{}, err
	}

	externalGroups, diag := schemaUtils.StringArrayToSetValue(r.ExternalGroups)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "external_groups", mapper.IntoTerraform)
	}

	specValue, diags := schema.NewSpecValue(
		map[string]attr.Type{
			"external_groups": externalGroups.Type(ctx),
			"members":         members.Type(ctx),
		},
		map[string]attr.Value{
			"external_groups": externalGroups,
			"members":         members,
		},
	)
	if diags.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diags, "spec", mapper.IntoTerraform)
	}
	return specValue, nil
}

func membersToListValue(ctx context.Context, members []gateway.GatewayServiceAccountId) (basetypes.ListValue, error) {
	if len(members) == 0 {
		memberType := schema.MembersValue{}.Type(ctx)
		return types.ListValueMust(memberType, []attr.Value{}), nil
	}

	tfMembers := make([]attr.Value, 0, len(members))
	for _, m := range members {
		attrTypes := map[string]attr.Type{
			"name":     basetypes.StringType{},
			"vcluster": basetypes.StringType{},
		}
		attrValues := map[string]attr.Value{
			"name":     schemaUtils.NewStringValue(m.Name),
			"vcluster": schemaUtils.NewStringValue(m.VCluster),
		}
		memberVal, diags := schema.NewMembersValue(attrTypes, attrValues)
		if diags.HasError() {
			return basetypes.ListValue{}, mapper.WrapDiagError(diags, "members", mapper.IntoTerraform)
		}
		tfMembers = append(tfMembers, memberVal)
	}

	listVal, diags := types.ListValue(schema.MembersValue{}.Type(ctx), tfMembers)
	if diags.HasError() {
		return basetypes.ListValue{}, mapper.WrapDiagError(diags, "members", mapper.IntoTerraform)
	}
	return listVal, nil
}
