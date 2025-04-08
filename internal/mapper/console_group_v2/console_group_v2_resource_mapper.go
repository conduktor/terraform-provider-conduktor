package console_group_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	groups "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_group_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *groups.ConsoleGroupV2Model) (console.GroupConsoleResource, error) {
	// ExternalGroups
	externalGroups, diag := schema.SetValueToStringArray(ctx, r.Spec.ExternalGroups)
	if diag.HasError() {
		return console.GroupConsoleResource{}, mapper.WrapDiagError(diag, "externalGroups", mapper.FromTerraform)
	}

	// Members
	members, diag := schema.SetValueToStringArray(ctx, r.Spec.Members)
	if diag.HasError() {
		return console.GroupConsoleResource{}, mapper.WrapDiagError(diag, "members", mapper.FromTerraform)
	}

	// Permissions
	permissions, err := schema.SetValueToPermissionArray(ctx, schema.GROUPS, r.Spec.Permissions)
	if err != nil {
		return console.GroupConsoleResource{}, err
	}

	return console.NewGroupConsoleResource(
		r.Name.ValueString(),
		console.GroupConsoleSpec{
			DisplayName:    r.Spec.DisplayName.ValueString(),
			Description:    r.Spec.Description.ValueString(),
			ExternalGroups: externalGroups,
			Members:        members,
			Permissions:    permissions,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *console.GroupConsoleResource) (groups.ConsoleGroupV2Model, error) {
	externalGroupsList, diag := schema.StringArrayToSetValue(r.Spec.ExternalGroups)
	if diag.HasError() {
		return groups.ConsoleGroupV2Model{}, mapper.WrapDiagError(diag, "external_groups", mapper.IntoTerraform)
	}

	membersList, diag := schema.StringArrayToSetValue(r.Spec.Members)
	if diag.HasError() {
		return groups.ConsoleGroupV2Model{}, mapper.WrapDiagError(diag, "members", mapper.IntoTerraform)
	}

	permissionsList, err := schema.PermissionArrayToSetValue(ctx, schema.GROUPS, r.Spec.Permissions)
	if err != nil {
		return groups.ConsoleGroupV2Model{}, err
	}

	specValue, diag := groups.NewSpecValue(
		map[string]attr.Type{
			"description":     basetypes.StringType{},
			"display_name":    basetypes.StringType{},
			"external_groups": externalGroupsList.Type(ctx),
			"members":         membersList.Type(ctx),
			"permissions":     permissionsList.Type(ctx),
		},
		map[string]attr.Value{
			"description":     schema.NewStringValue(r.Spec.Description),
			"display_name":    schema.NewStringValue(r.Spec.DisplayName),
			"external_groups": externalGroupsList,
			"members":         membersList,
			"permissions":     permissionsList,
		},
	)
	if diag.HasError() {
		return groups.ConsoleGroupV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return groups.ConsoleGroupV2Model{
		Name: types.StringValue(r.Metadata.Name),
		Spec: specValue,
	}, nil
}
