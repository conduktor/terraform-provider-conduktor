package console_application_group_v1

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	appinstance "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_application_group_v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *appinstance.ConsoleApplicationGroupV1Model) (console.ApplicationGroupConsoleResource, error) {
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
	permissions, err := schema.SetValueToApplicationGroupPermissionArray(ctx, r.Spec.Permissions)
	if err != nil {
		return console.ApplicationGroupConsoleResource{}, err
	}

	return console.NewApplicationGroupConsoleResource(
		r.Name.ValueString(),
		r.Application.ValueString(),
		console.ApplicationGroupSpec{
			DisplayName:           r.Spec.DisplayName.ValueString(),
			Description:           r.Spec.Description.ValueString(),
			ExternalGroups:        externalGroups,
			Members:               members,
			ExternalGroupMemebers: membersFromExternalGroups,
			Permissions:           permissions,
		},
	), nil

}

func InternalModelToTerraform(ctx context.Context, r *console.ApplicationGroupConsoleResource) (appinstance.ConsoleApplicationGroupV1Model, error) {
	externalGroupsList, diag := schema.StringArrayToSetValue(r.Spec.ExternalGroups)
	if diag.HasError() {
		return appinstance.ConsoleApplicationGroupV1Model{}, mapper.WrapDiagError(diag, "external_groups", mapper.IntoTerraform)
	}

	membersList, diag := schema.StringArrayToSetValue(r.Spec.Members)
	if diag.HasError() {
		return appinstance.ConsoleApplicationGroupV1Model{}, mapper.WrapDiagError(diag, "members", mapper.IntoTerraform)
	}

	membersFromExternalGroupsList, diag := schema.StringArrayToSetValue(r.Spec.ExternalGroupMemebers)
	if diag.HasError() {
		return appinstance.ConsoleApplicationGroupV1Model{}, mapper.WrapDiagError(diag, "members_from_external_groups", mapper.IntoTerraform)
	}

	permissionsList, err := schema.ApplicationGroupPermissionArrayToSetValue(ctx, r.Spec.Permissions)
	if err != nil {
		return appinstance.ConsoleApplicationGroupV1Model{}, err
	}

	specValue, diag := appinstance.NewSpecValue(
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
		return appinstance.ConsoleApplicationGroupV1Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return appinstance.ConsoleApplicationGroupV1Model{
		Name:        schema.NewStringValue(r.Metadata.Name),
		Application: schema.NewStringValue(r.Metadata.Application),
		Spec:        specValue,
	}, nil
}
