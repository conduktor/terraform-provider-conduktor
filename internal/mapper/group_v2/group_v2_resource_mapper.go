package group_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_group_v2"
	helpers "github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func TFToInternalModel(ctx context.Context, r *schema.GroupV2Model) (model.GroupConsoleResource, error) {
	// ExternalGroups
	if r.Spec.ExternalGroups.IsNull() {
		tflog.Debug(ctx, "SetValue externalGroups is null")
	} else if r.Spec.ExternalGroups.IsUnknown() {
		tflog.Debug(ctx, "SetValue externalGroups is unknown")
	}
	externalGroups, diag := helpers.SetValueToStringArray(ctx, r.Spec.ExternalGroups)
	if diag.HasError() {
		return model.GroupConsoleResource{}, mapper.WrapDiagError(diag, "externalGroups", mapper.FromTerraform)
	}

	// Members
	members, diag := helpers.SetValueToStringArray(ctx, r.Spec.Members)
	if diag.HasError() {
		return model.GroupConsoleResource{}, mapper.WrapDiagError(diag, "members", mapper.FromTerraform)
	}

	membersFromExternalGroups, diag := helpers.SetValueToStringArray(ctx, r.Spec.MembersFromExternalGroups)
	if diag.HasError() {
		return model.GroupConsoleResource{}, mapper.WrapDiagError(diag, "membersFromExternalGroups", mapper.FromTerraform)
	}

	// Permissions
	if r.Spec.Permissions.IsNull() {
		tflog.Debug(ctx, "SetValue permissions is null")
	} else if r.Spec.Permissions.IsUnknown() {
		tflog.Debug(ctx, "SetValue permissions is unknown")
	}
	permissions, err := helpers.SetValueToPermissionArray(ctx, "groups", r.Spec.Permissions)
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

func InternalModelToTerraform(ctx context.Context, r *model.GroupConsoleResource) (schema.GroupV2Model, error) {
	externalGroupsList, diag := helpers.StringArrayToSetValue(r.Spec.ExternalGroups)
	if diag.HasError() {
		return schema.GroupV2Model{}, mapper.WrapDiagError(diag, "external_groups", mapper.IntoTerraform)
	}

	membersList, diag := helpers.StringArrayToSetValue(r.Spec.Members)
	if diag.HasError() {
		return schema.GroupV2Model{}, mapper.WrapDiagError(diag, "members", mapper.IntoTerraform)
	}

	membersFromExternalGroupsList, diag := helpers.StringArrayToSetValue(r.Spec.MembersFromExternalGroups)
	if diag.HasError() {
		return schema.GroupV2Model{}, mapper.WrapDiagError(diag, "members_from_external_groups", mapper.IntoTerraform)
	}

	permissionsList, err := helpers.PermissionArrayToSetValue(ctx, "groups", r.Spec.Permissions)
	if err != nil {
		return schema.GroupV2Model{}, err
	}

	specValue, diag := schema.NewSpecValue(
		map[string]attr.Type{
			"description":                  basetypes.StringType{},
			"display_name":                 basetypes.StringType{},
			"external_groups":              externalGroupsList.Type(ctx),
			"members":                      membersList.Type(ctx),
			"members_from_external_groups": membersFromExternalGroupsList.Type(ctx),
			"permissions":                  permissionsList.Type(ctx),
		},
		map[string]attr.Value{
			"description":                  schemaUtils.NewStringValue(r.Spec.Description),
			"display_name":                 schemaUtils.NewStringValue(r.Spec.DisplayName),
			"external_groups":              externalGroupsList,
			"members":                      membersList,
			"members_from_external_groups": membersFromExternalGroupsList,
			"permissions":                  permissionsList,
		},
	)
	if diag.HasError() {
		return schema.GroupV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return schema.GroupV2Model{
		Name: types.StringValue(r.Metadata.Name),
		Spec: specValue,
	}, nil
}
