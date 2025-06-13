package console_resource_policy_v1

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	resourcePolicy "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_resource_policy_v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *resourcePolicy.ConsoleResourcePolicyV1Model) (console.ResourcePolicyConsoleResource, error) {
	labels, diag := schema.MapValueToStringMap(ctx, r.Labels)
	if diag.HasError() {
		return console.ResourcePolicyConsoleResource{}, mapper.WrapDiagError(diag, "labels", mapper.FromTerraform)
	}

	rules, err := setValueToRulesArray(ctx, r.Spec.Rules)
	if err != nil {
		return console.ResourcePolicyConsoleResource{}, err
	}

	return console.NewResourcePolicyConsoleResource(
		console.ResourcePolicyConsoleMetadata{
			Name:   r.Name.ValueString(),
			Labels: labels,
		},
		console.ResourcePolicyConsoleSpec{
			TargetKind:  r.Spec.TargetKind.ValueString(),
			Description: r.Spec.Description.ValueString(),
			Rules:       rules,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *console.ResourcePolicyConsoleResource) (resourcePolicy.ConsoleResourcePolicyV1Model, error) {
	var diag diag.Diagnostics

	rulesSet, err := rulesArrayToSetValue(ctx, r.Spec.Rules)
	if err != nil {
		return resourcePolicy.ConsoleResourcePolicyV1Model{}, err
	}

	labels, diag := schema.StringMapToMapValue(ctx, r.Metadata.Labels)
	if diag.HasError() {
		return resourcePolicy.ConsoleResourcePolicyV1Model{}, mapper.WrapDiagError(diag, "labels", mapper.IntoTerraform)
	}

	specValue, diag := resourcePolicy.NewSpecValue(
		map[string]attr.Type{
			"target_kind": basetypes.StringType{},
			"description": basetypes.StringType{},
			"rules":       rulesSet.Type(ctx),
		},
		map[string]attr.Value{
			"target_kind": schema.NewStringValue(r.Spec.TargetKind),
			"description": schema.NewStringValue(r.Spec.Description),
			"rules":       rulesSet,
		},
	)
	if diag.HasError() {
		return resourcePolicy.ConsoleResourcePolicyV1Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return resourcePolicy.ConsoleResourcePolicyV1Model{
		Name:   schema.NewStringValue(r.Metadata.Name),
		Labels: labels,
		Spec:   specValue,
	}, nil
}

// Parse a Rules Array into a Set.
func rulesArrayToSetValue(ctx context.Context, arr []console.ResourcePolicyConsoleRule) (basetypes.SetValue, error) {
	var tfResources []attr.Value
	var diag diag.Diagnostics

	for _, p := range arr {
		types := map[string]attr.Type{
			"condition":     basetypes.StringType{},
			"error_message": basetypes.StringType{},
		}
		values := map[string]attr.Value{
			"condition":     schema.NewStringValue(p.Condition),
			"error_message": schema.NewStringValue(p.ErrorMessage),
		}

		permObj, diag := resourcePolicy.NewRulesValue(types, values)
		if diag.HasError() {
			return basetypes.SetValue{}, mapper.WrapDiagError(diag, "rules", mapper.FromTerraform)
		}
		tfResources = append(tfResources, permObj)

	}

	resourcesList, diag := types.SetValue(resourcePolicy.RulesValue{}.Type(ctx), tfResources)
	if diag.HasError() {
		return basetypes.SetValue{}, mapper.WrapDiagError(diag, "rules", mapper.FromTerraform)
	}

	return resourcesList, nil
}

// Parse a Set into an array of Rules.
func setValueToRulesArray(ctx context.Context, set basetypes.SetValue) ([]console.ResourcePolicyConsoleRule, error) {
	rules := make([]console.ResourcePolicyConsoleRule, 0)
	var diag diag.Diagnostics

	if !set.IsNull() && !set.IsUnknown() {
		var tfResources []resourcePolicy.RulesValue
		diag = set.ElementsAs(ctx, &tfResources, false)
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "rules", mapper.FromTerraform)
		}

		for _, p := range tfResources {
			rules = append(rules, console.ResourcePolicyConsoleRule{
				Condition:    p.Condition.ValueString(),
				ErrorMessage: p.ErrorMessage.ValueString(),
			})
		}

	}
	return rules, nil
}
