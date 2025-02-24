package console_application_v1

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	apps "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_application_v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *apps.ConsoleApplicationV1Model) (console.ApplicationConsoleResource, error) {
	return console.NewApplicationConsoleResource(
		r.Name.ValueString(),
		console.ApplicationConsoleSpec{
			Title:       r.Spec.Title.ValueString(),
			Description: r.Spec.Description.ValueString(),
			Owner:       r.Spec.Owner.ValueString(),
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *console.ApplicationConsoleResource) (apps.ConsoleApplicationV1Model, error) {
	specValue, diag := apps.NewSpecValue(
		map[string]attr.Type{
			"title":       basetypes.StringType{},
			"description": basetypes.StringType{},
			"owner":       basetypes.StringType{},
		},
		map[string]attr.Value{
			"title":       schema.NewStringValue(r.Spec.Title),
			"description": schema.NewStringValue(r.Spec.Description),
			"owner":       schema.NewStringValue(r.Spec.Owner),
		},
	)
	if diag.HasError() {
		return apps.ConsoleApplicationV1Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return apps.ConsoleApplicationV1Model{
		Name: schema.NewStringValue(r.Metadata.Name),
		Spec: specValue,
	}, nil
}
