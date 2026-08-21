package console_application_v1

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	apps "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_application_v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *apps.ConsoleApplicationV1Model) (console.ApplicationConsoleResource, error) {
	policyRef, diag := schema.SetValueToStringArray(ctx, r.Spec.PolicyRef)
	if diag.HasError() {
		return console.ApplicationConsoleResource{}, mapper.WrapDiagError(diag, "policy_ref", mapper.FromTerraform)
	}

	return console.NewApplicationConsoleResource(
		r.Name.ValueString(),
		console.ApplicationConsoleSpec{
			Title:       r.Spec.Title.ValueString(),
			Description: r.Spec.Description.ValueString(),
			Owner:       r.Spec.Owner.ValueString(),
			PolicyRef:   policyRef,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *console.ApplicationConsoleResource) (apps.ConsoleApplicationV1Model, error) {
	policyRef := basetypes.NewSetNull(basetypes.StringType{})
	if r.Spec.PolicyRef != nil {
		var diags diag.Diagnostics
		policyRef, diags = schema.StringArrayToSetValue(r.Spec.PolicyRef)
		if diags.HasError() {
			return apps.ConsoleApplicationV1Model{}, mapper.WrapDiagError(diags, "policy_ref", mapper.IntoTerraform)
		}
	}

	specValue, diags := apps.NewSpecValue(
		map[string]attr.Type{
			"title":       basetypes.StringType{},
			"description": basetypes.StringType{},
			"owner":       basetypes.StringType{},
			"policy_ref":  policyRef.Type(ctx),
		},
		map[string]attr.Value{
			"title":       schema.NewStringValue(r.Spec.Title),
			"description": schema.NewStringValue(r.Spec.Description),
			"owner":       schema.NewStringValue(r.Spec.Owner),
			"policy_ref":  policyRef,
		},
	)
	if diags.HasError() {
		return apps.ConsoleApplicationV1Model{}, mapper.WrapDiagError(diags, "spec", mapper.IntoTerraform)
	}

	return apps.ConsoleApplicationV1Model{
		Name: schema.NewStringValue(r.Metadata.Name),
		Spec: specValue,
	}, nil
}
