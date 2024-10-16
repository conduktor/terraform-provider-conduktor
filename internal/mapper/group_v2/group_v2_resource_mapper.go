package group_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_group_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func TFToInternalModel(ctx context.Context, r *schema.GroupV2Model) (model.GroupConsoleResource, error) {

	externalGroups, diag := schemaUtils.ListValueToStringArray(ctx, r.Spec.ExternalGroups)
	if r.Spec.ExternalGroups.IsNull() {
		tflog.Debug(ctx, "ListValue externalGroups is null")
	}

	if r.Spec.ExternalGroups.IsUnknown() {
		tflog.Debug(ctx, "ListValue  externalGroups is unknown")
	}

	if diag.HasError() {
		return model.GroupConsoleResource{}, mapper.WrapDiagError(diag, "externalGroups", mapper.FromTerraform)
	}

	members, diag := schemaUtils.ListValueToStringArray(ctx, r.Spec.Members)
	if diag.HasError() {
		return model.GroupConsoleResource{}, mapper.WrapDiagError(diag, "members", mapper.FromTerraform)
	}

	membersFromExternalGroups, diag := schemaUtils.ListValueToStringArray(ctx, r.Spec.MembersFromExternalGroups)
	if diag.HasError() {
		return model.GroupConsoleResource{}, mapper.WrapDiagError(diag, "membersFromExternalGroups", mapper.FromTerraform)
	}

	permissions := make([]model.Permission, 0)
	if !r.Spec.Permissions.IsNull() && !r.Spec.Permissions.IsUnknown() {
		var tfPermissions []schema.PermissionsValue
		diag = r.Spec.Permissions.ElementsAs(ctx, &tfPermissions, false)
		if diag.HasError() {
			return model.GroupConsoleResource{}, mapper.WrapDiagError(diag, "permissions", mapper.FromTerraform)
		}
		for _, p := range tfPermissions {
			flags, diag := schemaUtils.ListValueToStringArray(ctx, p.Permissions)
			if diag.HasError() {
				return model.GroupConsoleResource{}, mapper.WrapDiagError(diag, "permissions.permissions", mapper.FromTerraform)
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

	externalGroupsList, diag := schemaUtils.StringArrayToListValue(r.Spec.ExternalGroups)
	if diag.HasError() {
		return schema.GroupV2Model{}, mapper.WrapDiagError(diag, "external_groups", mapper.IntoTerraform)
	}

	membersList, diag := schemaUtils.StringArrayToListValue(r.Spec.Members)
	if diag.HasError() {
		return schema.GroupV2Model{}, mapper.WrapDiagError(diag, "members", mapper.IntoTerraform)
	}

	membersFromExternalGroupsList, diag := schemaUtils.StringArrayToListValue(r.Spec.MembersFromExternalGroups)
	if diag.HasError() {
		return schema.GroupV2Model{}, mapper.WrapDiagError(diag, "members_from_external_groups", mapper.IntoTerraform)
	}

	var tfPermissions []attr.Value
	for _, p := range r.Spec.Permissions {
		flagsList, diag := schemaUtils.StringArrayToListValue(p.Permissions)
		if diag.HasError() {
			return schema.GroupV2Model{}, mapper.WrapDiagError(diag, "permissions.permissions", mapper.IntoTerraform)
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
			return schema.GroupV2Model{}, mapper.WrapDiagError(diag, "permissions", mapper.IntoTerraform)
		}

		tfPermissions = append(tfPermissions, permObj)
	}

	permissionsList, diag := types.ListValue(schema.PermissionsValue{}.Type(ctx), tfPermissions)
	if diag.HasError() {
		return schema.GroupV2Model{}, mapper.WrapDiagError(diag, "permissions", mapper.IntoTerraform)
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
