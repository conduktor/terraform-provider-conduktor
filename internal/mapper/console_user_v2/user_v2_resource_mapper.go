package console_user_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	users "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_user_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *users.ConsoleUserV2Model) (console.UserConsoleResource, error) {
	permissions, err := schema.SetValueToPermissionArray(ctx, schema.USERS, r.Spec.Permissions)
	if err != nil {
		return console.UserConsoleResource{}, err
	}

	return console.NewUserConsoleResource(
		r.Name.ValueString(),
		console.UserConsoleSpec{
			FirstName:   r.Spec.Firstname.ValueString(),
			LastName:    r.Spec.Lastname.ValueString(),
			Permissions: permissions,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *console.UserConsoleResource) (users.ConsoleUserV2Model, error) {
	permissionsList, err := schema.PermissionArrayToSetValue(ctx, schema.USERS, r.Spec.Permissions)
	if err != nil {
		return users.ConsoleUserV2Model{}, err
	}

	specValue, diag := users.NewSpecValue(
		map[string]attr.Type{
			"firstname":   basetypes.StringType{},
			"lastname":    basetypes.StringType{},
			"permissions": permissionsList.Type(ctx),
		},
		map[string]attr.Value{
			"firstname":   schema.NewStringValue(r.Spec.FirstName),
			"lastname":    schema.NewStringValue(r.Spec.LastName),
			"permissions": permissionsList,
		},
	)
	if diag.HasError() {
		return users.ConsoleUserV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return users.ConsoleUserV2Model{
		Name: types.StringValue(r.Metadata.Name),
		Spec: specValue,
	}, nil
}
