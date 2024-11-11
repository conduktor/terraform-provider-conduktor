package user_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_user_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *schema.UserV2Model) (model.UserConsoleResource, error) {

	permissions := make([]model.Permission, 0)

	if !r.Spec.Permissions.IsNull() && !r.Spec.Permissions.IsUnknown() {
		var tfPermissions []schema.PermissionsValue
		diag := r.Spec.Permissions.ElementsAs(ctx, &tfPermissions, false)
		if diag.HasError() {
			return model.UserConsoleResource{}, mapper.WrapDiagError(diag, "permissions", mapper.FromTerraform)
		}
		for _, p := range tfPermissions {
			var flags []string
			diag := p.Permissions.ElementsAs(ctx, &flags, false)
			if diag.HasError() {
				return model.UserConsoleResource{}, mapper.WrapDiagError(diag, "permissions.permissions", mapper.FromTerraform)
			}

			permissions = append(permissions, model.Permission{
				Name:         p.Name.ValueString(),
				ResourceType: p.ResourceType.ValueString(),
				Permissions:  flags,
				PatternType:  p.PatternType.ValueString(),
				Cluster:      p.Cluster.ValueString(),
				KafkaConnect: p.KafkaConnect.ValueString(),
			})
		}
	}

	return model.NewUserConsoleResource(
		r.Name.ValueString(),
		model.UserConsoleSpec{
			FirstName:   r.Spec.Firstname.ValueString(),
			LastName:    r.Spec.Lastname.ValueString(),
			Permissions: permissions,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *model.UserConsoleResource) (schema.UserV2Model, error) {
	var tfPermissions []attr.Value
	for _, p := range r.Spec.Permissions {
		var flags []attr.Value
		for _, f := range p.Permissions {
			flags = append(flags, types.StringValue(f))
		}

		flagsList, diag := types.SetValue(types.StringType, flags)
		if diag.HasError() {
			return schema.UserV2Model{}, mapper.WrapDiagError(diag, "permissions.permissions", mapper.IntoTerraform)
		}
		permObj, diag := schema.NewPermissionsValue(
			map[string]attr.Type{
				"name":          basetypes.StringType{},
				"resource_type": basetypes.StringType{},
				"permissions":   flagsList.Type(ctx),
				"pattern_type":  basetypes.StringType{},
				"cluster":       basetypes.StringType{},
				"kafka_connect": basetypes.StringType{},
			},
			map[string]attr.Value{
				"name":          schemaUtils.NewStringValue(p.Name),
				"resource_type": schemaUtils.NewStringValue(p.ResourceType),
				"permissions":   flagsList,
				"pattern_type":  schemaUtils.NewStringValue(p.PatternType),
				"cluster":       schemaUtils.NewStringValue(p.Cluster),
				"kafka_connect": schemaUtils.NewStringValue(p.KafkaConnect),
			},
		)
		if diag.HasError() {
			return schema.UserV2Model{}, mapper.WrapDiagError(diag, "permissions (value)", mapper.IntoTerraform)
		}

		tfPermissions = append(tfPermissions, permObj)
	}

	permissionsList, diag := types.SetValue(schema.PermissionsValue{}.Type(ctx), tfPermissions)
	if diag.HasError() {
		return schema.UserV2Model{}, mapper.WrapDiagError(diag, "permissions (SetValue)", mapper.IntoTerraform)
	}

	specValue, diag := schema.NewSpecValue(
		map[string]attr.Type{
			"firstname":   basetypes.StringType{},
			"lastname":    basetypes.StringType{},
			"permissions": permissionsList.Type(ctx),
		},
		map[string]attr.Value{
			"firstname":   schemaUtils.NewStringValue(r.Spec.FirstName),
			"lastname":    schemaUtils.NewStringValue(r.Spec.LastName),
			"permissions": permissionsList,
		},
	)
	if diag.HasError() {
		return schema.UserV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return schema.UserV2Model{
		Name: types.StringValue(r.Metadata.Name),
		Spec: specValue,
	}, nil
}
