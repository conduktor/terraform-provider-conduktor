package console_topic_policy_v1

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	topicPolicy "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_topic_policy_v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func InternalModelToTerraform(ctx context.Context, r *console.TopicPolicyResource) (topicPolicy.ConsoleTopicPolicyV1Model, error) {
	policiesMap, err := policiesMapToPoliciesValue(ctx, r.Spec.Policies)
	if err != nil {
		return topicPolicy.ConsoleTopicPolicyV1Model{}, err
	}

	specValue, diag := topicPolicy.NewSpecValue(
		map[string]attr.Type{
			"policies": policiesMap.Type(ctx),
		},
		map[string]attr.Value{
			"policies": policiesMap,
		},
	)
	if diag.HasError() {
		return topicPolicy.ConsoleTopicPolicyV1Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return topicPolicy.ConsoleTopicPolicyV1Model{
		Name: schema.NewStringValue(r.Metadata.Name),
		Spec: specValue,
	}, nil
}

// Parse a Resources Array into a Set.
func policiesMapToPoliciesValue(ctx context.Context, m map[string]console.Constraint) (basetypes.MapValue, error) {
	var tfPolicies map[string]attr.Value
	var diag diag.Diagnostics

	for k, v := range m {

		permObj, err := constraintInternalModelToTerraform(ctx, &v)
		if err != nil {
			return basetypes.MapValue{}, err
		}

		tfPolicies[k] = permObj
	}

	policiesMap, diag := types.MapValue(topicPolicy.PoliciesValue{}.Type(ctx), tfPolicies)
	if diag.HasError() {
		return basetypes.MapValue{}, mapper.WrapDiagError(diag, "policies", mapper.FromTerraform)
	}

	return policiesMap, nil
}

func constraintInternalModelToTerraform(ctx context.Context, r *console.Constraint) (topicPolicy.PoliciesValue, error) {
	if r == nil || (r.Match == nil && r.NoneOf == nil && r.OneOf == nil && r.Range == nil) {
		return topicPolicy.NewPoliciesValueNull(), nil
	}

	unknownConstraintObjectValue, diag := topicPolicy.NewPoliciesValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return topicPolicy.PoliciesValue{}, mapper.WrapDiagError(diag, "policies", mapper.IntoTerraform)
	}
	var typesMap = unknownConstraintObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	if r.Match != nil {
		var matchTypesMap = topicPolicy.NewMatchValueNull().AttributeTypes(ctx)
		var matchValuesMap = schema.ValueMapFromTypes(ctx, matchTypesMap)
		matchValuesMap["optional"] = basetypes.NewBoolValue(r.Match.Optional)
		matchValuesMap["pattern"] = schema.NewStringValue(r.Match.Pattern)

		valuesMap["match"], diag = types.ObjectValue(matchTypesMap, matchValuesMap)
		if diag.HasError() {
			return topicPolicy.PoliciesValue{}, mapper.WrapDiagError(diag, "policies", mapper.IntoTerraform)
		}
	}

	if r.NoneOf != nil {
		var noneofTypesMap = topicPolicy.NewNoneOfValueNull().AttributeTypes(ctx)
		var noneofValuesMap = schema.ValueMapFromTypes(ctx, noneofTypesMap)
		noneofValuesMap["optional"] = basetypes.NewBoolValue(r.NoneOf.Optional)
		valuesList, diag := schema.StringArrayToSetValue(r.NoneOf.Values)
		if diag.HasError() {
			return topicPolicy.PoliciesValue{}, mapper.WrapDiagError(diag, "policies", mapper.IntoTerraform)
		}
		noneofValuesMap["values"] = valuesList

		valuesMap["none_of"], diag = types.ObjectValue(noneofTypesMap, noneofValuesMap)
		if diag.HasError() {
			return topicPolicy.PoliciesValue{}, mapper.WrapDiagError(diag, "policies", mapper.IntoTerraform)
		}
	}

	if r.OneOf != nil {
		var oneofTypesMap = topicPolicy.NewOneOfValueNull().AttributeTypes(ctx)
		var oneofValuesMap = schema.ValueMapFromTypes(ctx, oneofTypesMap)
		oneofValuesMap["optional"] = basetypes.NewBoolValue(r.OneOf.Optional)
		valuesList, diag := schema.StringArrayToSetValue(r.OneOf.Values)
		if diag.HasError() {
			return topicPolicy.PoliciesValue{}, mapper.WrapDiagError(diag, "policies", mapper.IntoTerraform)
		}
		oneofValuesMap["values"] = valuesList

		valuesMap["one_of"], diag = types.ObjectValue(oneofTypesMap, oneofValuesMap)
		if diag.HasError() {
			return topicPolicy.PoliciesValue{}, mapper.WrapDiagError(diag, "policies", mapper.IntoTerraform)
		}
	}

	if r.Range != nil {
		var rangeTypesMap = topicPolicy.NewRangeValueNull().AttributeTypes(ctx)
		var rangeValuesMap = schema.ValueMapFromTypes(ctx, rangeTypesMap)
		rangeValuesMap["optional"] = basetypes.NewBoolValue(r.OneOf.Optional)
		rangeValuesMap["min"] = schema.NewInt64Value(r.Range.Min)
		rangeValuesMap["max"] = schema.NewInt64Value(r.Range.Max)

		valuesMap["range"], diag = types.ObjectValue(rangeTypesMap, rangeValuesMap)
		if diag.HasError() {
			return topicPolicy.PoliciesValue{}, mapper.WrapDiagError(diag, "policies", mapper.IntoTerraform)
		}
	}

	value, diag := topicPolicy.NewPoliciesValue(typesMap, valuesMap)
	if diag.HasError() {
		return topicPolicy.PoliciesValue{}, mapper.WrapDiagError(diag, "policies", mapper.IntoTerraform)
	}
	return value, nil
}
