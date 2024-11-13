package user_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_user_v2"
	helpers "github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func TFToInternalModel(ctx context.Context, r *schema.UserV2Model) (model.UserConsoleResource, error) {
	if r.Spec.Permissions.IsNull() {
		tflog.Debug(ctx, "ListValue externalGroups is null")
	} else if r.Spec.Permissions.IsUnknown() {
		tflog.Debug(ctx, "ListValue externalGroups is unknown")
	}
	permissions, err := helpers.SetValueToPermissionArray(ctx, "users", r.Spec.Permissions)
	if err != nil {
		return model.UserConsoleResource{}, err
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
	permissionsList, err := helpers.PermissionArrayToSetValue(ctx, "users", r.Spec.Permissions)
	if err != nil {
		return schema.UserV2Model{}, err
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
