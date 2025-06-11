package console_application_group_v1

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	appgroup "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_application_group_v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *appgroup.ConsoleApplicationGroupV1Model) (console.ApplicationGroupConsoleResource, error) {
	// ExternalGroups
	externalGroups, diag := schema.SetValueToStringArray(ctx, r.Spec.ExternalGroups)
	if diag.HasError() {
		return console.ApplicationGroupConsoleResource{}, mapper.WrapDiagError(diag, "externalGroups", mapper.FromTerraform)
	}

	// Members
	members, diag := schema.SetValueToStringArray(ctx, r.Spec.Members)
	if diag.HasError() {
		return console.ApplicationGroupConsoleResource{}, mapper.WrapDiagError(diag, "members", mapper.FromTerraform)
	}

	// MembersFromExternalGroups
	membersFromExternalGroups, diag := schema.SetValueToStringArray(ctx, r.Spec.MembersFromExternalGroups)
	if diag.HasError() {
		return console.ApplicationGroupConsoleResource{}, mapper.WrapDiagError(diag, "membersFromExternalGroups", mapper.FromTerraform)
	}

	// Permissions
	permissions, err := SetValueToApplicationGroupPermissionArray(ctx, r.Spec.Permissions)
	if err != nil {
		return console.ApplicationGroupConsoleResource{}, err
	}

	return console.NewApplicationGroupConsoleResource(
		r.Name.ValueString(),
		r.Application.ValueString(),
		console.ApplicationGroupSpec{
			DisplayName:               r.Spec.DisplayName.ValueString(),
			Description:               r.Spec.Description.ValueString(),
			ExternalGroups:            externalGroups,
			Members:                   members,
			MembersFromExternalGroups: membersFromExternalGroups,
			Permissions:               permissions,
		},
	), nil

}

func InternalModelToTerraform(ctx context.Context, r *console.ApplicationGroupConsoleResource) (appgroup.ConsoleApplicationGroupV1Model, error) {
	externalGroupsList, diag := schema.StringArrayToSetValue(r.Spec.ExternalGroups)
	if diag.HasError() {
		return appgroup.ConsoleApplicationGroupV1Model{}, mapper.WrapDiagError(diag, "external_groups", mapper.IntoTerraform)
	}

	membersList, diag := schema.StringArrayToSetValue(r.Spec.Members)
	if diag.HasError() {
		return appgroup.ConsoleApplicationGroupV1Model{}, mapper.WrapDiagError(diag, "members", mapper.IntoTerraform)
	}

	membersFromExternalGroupsList, diag := schema.StringArrayToSetValue(r.Spec.MembersFromExternalGroups)
	if diag.HasError() {
		return appgroup.ConsoleApplicationGroupV1Model{}, mapper.WrapDiagError(diag, "members_from_external_groups", mapper.IntoTerraform)
	}

	permissionsList, err := ApplicationGroupPermissionArrayToSetValue(ctx, r.Spec.Permissions)
	if err != nil {
		return appgroup.ConsoleApplicationGroupV1Model{}, err
	}

	specValue, diag := appgroup.NewSpecValue(
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
		return appgroup.ConsoleApplicationGroupV1Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return appgroup.ConsoleApplicationGroupV1Model{
		Name:        schema.NewStringValue(r.Metadata.Name),
		Application: schema.NewStringValue(r.Metadata.Application),
		Spec:        specValue,
	}, nil
}

// Parse a ApplicationGroupPermissions Array into a Set.
func ApplicationGroupPermissionArrayToSetValue(ctx context.Context, arr []console.ApplicationGroupPermission) (basetypes.SetValue, error) {
	var permissionsList basetypes.SetValue
	var tfPermissions []attr.Value
	var diag diag.Diagnostics

	for _, p := range arr {
		flagsList, diag := schema.StringArrayToSetValue(p.Permissions)
		if diag.HasError() {
			return basetypes.SetValue{}, mapper.WrapDiagError(diag, "permissions.permissions", mapper.FromTerraform)
		}

		types := map[string]attr.Type{
			"name":            basetypes.StringType{},
			"resource_type":   basetypes.StringType{},
			"permissions":     flagsList.Type(ctx),
			"pattern_type":    basetypes.StringType{},
			"connect_cluster": basetypes.StringType{},
			"app_instance":    basetypes.StringType{},
		}
		values := map[string]attr.Value{
			"name":            schema.NewStringValue(p.Name),
			"resource_type":   schema.NewStringValue(p.ResourceType),
			"permissions":     flagsList,
			"pattern_type":    schema.NewStringValue(p.PatternType),
			"connect_cluster": schema.NewStringValue(p.ConnectCluster),
			"app_instance":    schema.NewStringValue(p.AppInstance),
		}

		appPermObj, diag := appgroup.NewPermissionsValue(types, values)
		if diag.HasError() {
			return basetypes.SetValue{}, mapper.WrapDiagError(diag, "permissions", mapper.FromTerraform)
		}
		tfPermissions = append(tfPermissions, appPermObj)
	}

	permissionsList, diag = types.SetValue(appgroup.PermissionsValue{}.Type(ctx), tfPermissions)
	if diag.HasError() {
		return basetypes.SetValue{}, mapper.WrapDiagError(diag, "permissions", mapper.FromTerraform)
	}

	return permissionsList, nil
}

// Parse a Set into an array of ApplicationGroupPermissions.
func SetValueToApplicationGroupPermissionArray(ctx context.Context, set basetypes.SetValue) ([]console.ApplicationGroupPermission, error) {
	permission := make([]console.ApplicationGroupPermission, 0)
	var diag diag.Diagnostics

	if !set.IsNull() && !set.IsUnknown() {
		var tfPermissions []appgroup.PermissionsValue
		diag = set.ElementsAs(ctx, &tfPermissions, false)
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "permissions", mapper.FromTerraform)
		}
		for _, p := range tfPermissions {
			flags, diag := schema.SetValueToStringArray(ctx, p.Permissions)
			if diag.HasError() {
				return nil, mapper.WrapDiagError(diag, "permissions.permissions", mapper.FromTerraform)
			}

			permission = append(permission, console.ApplicationGroupPermission{
				AppInstance:    p.AppInstance.ValueString(),
				ResourceType:   p.ResourceType.ValueString(),
				Name:           p.Name.ValueString(),
				Permissions:    flags,
				PatternType:    p.PatternType.ValueString(),
				ConnectCluster: p.ConnectCluster.ValueString(),
			})
		}
	}
	return permission, nil
}
