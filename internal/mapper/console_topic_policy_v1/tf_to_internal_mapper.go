package console_topic_policy_v1

import (
	"context"
	"fmt"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	topicPolicy "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_topic_policy_v1"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *topicPolicy.ConsoleTopicPolicyV1Model) (console.TopicPolicyResource, error) {
	policies, err := mapValueToPoliciesMap(ctx, r.Spec.Policies)
	if err != nil {
		return console.TopicPolicyResource{}, err
	}

	return console.NewTopicPolicyResource(
		r.Name.ValueString(),
		console.TopicPolicySpec{
			Policies: policies,
		},
	), nil
}

// Parse a Map nested into an map of Constraints.
func mapValueToPoliciesMap(ctx context.Context, m basetypes.MapValue) (map[string]console.Constraint, error) {
	policies := make(map[string]console.Constraint, 0)

	if !m.IsNull() && !m.IsUnknown() {
		var tfPolicies map[string]topicPolicy.PoliciesValue
		diag := m.ElementsAs(ctx, &tfPolicies, false)
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "policies", mapper.FromTerraform)
		}

		for k, v := range tfPolicies {
			constraint, err := policiesValueToConstraint(ctx, k, v)
			if err != nil {
				return nil, err
			}

			policies[k] = constraint
		}
	}

	return policies, nil
}

func policiesValueToConstraint(ctx context.Context, key string, policy topicPolicy.PoliciesValue) (console.Constraint, error) {
	if schema.AttrIsSet(policy.Match) {
		matchValue, diag := topicPolicy.NewMatchValue(policy.Match.AttributeTypes(ctx), policy.Match.Attributes())
		if diag.HasError() {
			return console.Constraint{}, mapper.WrapDiagError(diag, "policies."+key, mapper.FromTerraform)
		}
		Match := &console.Match{
			Optional: matchValue.Optional.ValueBool(),
			Pattern:  matchValue.Pattern.ValueString(),
		}
		return console.Constraint{Match: Match}, nil
	}

	if schema.AttrIsSet(policy.NoneOf) {
		noneOfValue, diag := topicPolicy.NewNoneOfValue(policy.NoneOf.AttributeTypes(ctx), policy.NoneOf.Attributes())
		if diag.HasError() {
			return console.Constraint{}, mapper.WrapDiagError(diag, "policies."+key, mapper.FromTerraform)
		}
		values, diag := schema.SetValueToStringArray(ctx, noneOfValue.Values)
		if diag.HasError() {
			return console.Constraint{}, mapper.WrapDiagError(diag, "policies."+key+".values", mapper.FromTerraform)
		}
		NoneOf := &console.NoneOf{
			Optional: noneOfValue.Optional.ValueBool(),
			Values:   values,
		}
		return console.Constraint{NoneOf: NoneOf}, nil
	}

	if schema.AttrIsSet(policy.OneOf) {
		oneOfValue, diag := topicPolicy.NewNoneOfValue(policy.OneOf.AttributeTypes(ctx), policy.OneOf.Attributes())
		if diag.HasError() {
			return console.Constraint{}, mapper.WrapDiagError(diag, "policies."+key, mapper.FromTerraform)
		}
		values, diag := schema.SetValueToStringArray(ctx, oneOfValue.Values)
		if diag.HasError() {
			return console.Constraint{}, mapper.WrapDiagError(diag, "policies."+key+".values", mapper.FromTerraform)
		}
		OneOf := &console.OneOf{
			Optional: oneOfValue.Optional.ValueBool(),
			Values:   values,
		}
		return console.Constraint{OneOf: OneOf}, nil
	}

	if schema.AttrIsSet(policy.Range) {
		rangeValue, diag := topicPolicy.NewRangeValue(policy.Range.AttributeTypes(ctx), policy.Range.Attributes())
		if diag.HasError() {
			return console.Constraint{}, mapper.WrapDiagError(diag, "policies."+key, mapper.FromTerraform)
		}
		Range := &console.Range{
			Optional: rangeValue.Optional.ValueBool(),
			Min:      rangeValue.Min.ValueInt64(),
			Max:      rangeValue.Max.ValueInt64(),
		}
		return console.Constraint{Range: Range}, nil
	}

	// Should never get here.
	return console.Constraint{}, fmt.Errorf("Unknown constraint type")
}
