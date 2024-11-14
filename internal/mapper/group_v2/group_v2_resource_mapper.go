package group_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	groups "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_group_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *groups.GroupV2Model) (model.GroupConsoleResource, error) {
	// ExternalGroups
	externalGroups, diag := schema.SetValueToStringArray(ctx, r.Spec.ExternalGroups)
	if diag.HasError() {
		return model.GroupConsoleResource{}, mapper.WrapDiagError(diag, "externalGroups", mapper.FromTerraform)
	}

	// Members
	members, diag := schema.SetValueToStringArray(ctx, r.Spec.Members)
	if diag.HasError() {
		return model.GroupConsoleResource{}, mapper.WrapDiagError(diag, "members", mapper.FromTerraform)
	}

	membersFromExternalGroups, diag := schema.SetValueToStringArray(ctx, r.Spec.MembersFromExternalGroups)
	if diag.HasError() {
		return model.GroupConsoleResource{}, mapper.WrapDiagError(diag, "membersFromExternalGroups", mapper.FromTerraform)
	}

	// Permissions
	permissions, err := schema.SetValueToPermissionArray(ctx, schema.GROUPS, r.Spec.Permissions)
	if err != nil {
		return model.GroupConsoleResource{}, err
	}

	return model.NewGroupConsoleResource(
		r.Name.ValueString(),
		model.GroupConsoleSpec{
			DisplayName:               r.Spec.DisplayName.ValueString(),
			Description:               r.Spec.Description.ValueString(),
			ExternalGroups:            externalGroups,
			Members:                   members,
			MembersFromExternalGroups: membersFromExternalGroups,
			Permissions:               permissions,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *model.GroupConsoleResource) (groups.GroupV2Model, error) {
	externalGroupsList, diag := schema.StringArrayToSetValue(r.Spec.ExternalGroups)
	if diag.HasError() {
		return groups.GroupV2Model{}, mapper.WrapDiagError(diag, "external_groups", mapper.IntoTerraform)
	}

	membersList, diag := schema.StringArrayToSetValue(r.Spec.Members)
	if diag.HasError() {
		return groups.GroupV2Model{}, mapper.WrapDiagError(diag, "members", mapper.IntoTerraform)
	}

	membersFromExternalGroupsList, diag := schema.StringArrayToSetValue(r.Spec.MembersFromExternalGroups)
	if diag.HasError() {
		return groups.GroupV2Model{}, mapper.WrapDiagError(diag, "members_from_external_groups", mapper.IntoTerraform)
	}

	permissionsList, err := schema.PermissionArrayToSetValue(ctx, schema.GROUPS, r.Spec.Permissions)
	if err != nil {
		return groups.GroupV2Model{}, err
	}

	specValue, diag := groups.NewSpecValue(
		map[string]attr.Type{
			"description":                  basetypes.StringType{},
			"display_name":                 basetypes.StringType{},
			"external_groups":              externalGroupsList.Type(ctx),
			"members":                      membersList.Type(ctx),
			"members_from_external_groups": membersFromExternalGroupsList.Type(ctx),
			"permissions":                  permissionsList.Type(ctx),
		},
		map[string]attr.Value{
			"description":                  schema.NewStringValue(r.Spec.Description),
			"display_name":                 schema.NewStringValue(r.Spec.DisplayName),
			"external_groups":              externalGroupsList,
			"members":                      membersList,
			"members_from_external_groups": membersFromExternalGroupsList,
			"permissions":                  permissionsList,
		},
	)
	if diag.HasError() {
		return groups.GroupV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return groups.GroupV2Model{
		Name: types.StringValue(r.Metadata.Name),
		Spec: specValue,
	}, nil
}
