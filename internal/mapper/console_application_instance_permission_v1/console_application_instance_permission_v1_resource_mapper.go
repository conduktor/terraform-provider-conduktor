package console_application_instance_permission_v1

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	permission "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_application_instance_permission_v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *permission.ConsoleApplicationInstancePermissionV1Model) (console.ApplicationInstancePermissionConsoleResource, error) {
	resource, err := resourceTFToInternalModel(ctx, r.Spec.Resource)
	if err != nil {
		return console.ApplicationInstancePermissionConsoleResource{}, err
	}

	return console.NewApplicationInstancePermissionConsoleResource(
		console.ApplicationInstancePermissionConsoleMetadata{
			Name:        r.Name.ValueString(),
			Application: r.Application.ValueString(),
			AppInstance: r.AppInstance.ValueString(),
		},
		console.ApplicationInstancePermissionConsoleSpec{
			Resource:                 resource,
			UserPermission:           r.Spec.UserPermission.ValueString(),
			ServiceAccountPermission: r.Spec.ServiceAccountPermission.ValueString(),
			GrantedTo:                r.Spec.GrantedTo.ValueString(),
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *console.ApplicationInstancePermissionConsoleResource) (permission.ConsoleApplicationInstancePermissionV1Model, error) {
	var diag diag.Diagnostics

	resource, err := resourceInternalModelToTerraform(r.Spec.Resource)
	if err != nil {
		return permission.ConsoleApplicationInstancePermissionV1Model{}, err
	}

	specValue, diag := permission.NewSpecValue(
		map[string]attr.Type{
			"resource":                   resource.Type(ctx),
			"user_permission":            basetypes.StringType{},
			"service_account_permission": basetypes.StringType{},
			"granted_to":                 basetypes.StringType{},
		},
		map[string]attr.Value{
			"resource":                   resource,
			"user_permission":            schema.NewStringValue(r.Spec.UserPermission),
			"service_account_permission": schema.NewStringValue(r.Spec.ServiceAccountPermission),
			"granted_to":                 schema.NewStringValue(r.Spec.GrantedTo),
		},
	)
	if diag.HasError() {
		return permission.ConsoleApplicationInstancePermissionV1Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return permission.ConsoleApplicationInstancePermissionV1Model{
		Name:        schema.NewStringValue(r.Metadata.Name),
		Application: schema.NewStringValue(r.Metadata.Application),
		AppInstance: schema.NewStringValue(r.Metadata.AppInstance),
		Spec:        specValue,
	}, nil
}

func resourceTFToInternalModel(ctx context.Context, r basetypes.ObjectValue) (console.AppInstancePermissionResource, error) {
	// Should never happen since the resource is required in the schema.
	if r.IsNull() {
		return console.AppInstancePermissionResource{}, nil
	}

	resourceValue, diag := permission.NewResourceValue(r.AttributeTypes(ctx), r.Attributes())
	if diag.HasError() {
		return console.AppInstancePermissionResource{}, mapper.WrapDiagError(diag, "resource", mapper.FromTerraform)
	}

	return console.AppInstancePermissionResource{
		Type:           resourceValue.ResourceType.ValueString(),
		Name:           resourceValue.Name.ValueString(),
		PatternType:    resourceValue.PatternType.ValueString(),
		ConnectCluster: resourceValue.ConnectCluster.ValueString(),
	}, nil
}

func resourceInternalModelToTerraform(r console.AppInstancePermissionResource) (basetypes.ObjectValue, error) {
	var diag diag.Diagnostics

	types := map[string]attr.Type{
		"type":            basetypes.StringType{},
		"name":            basetypes.StringType{},
		"pattern_type":    basetypes.StringType{},
		"connect_cluster": basetypes.StringType{},
	}
	values := map[string]attr.Value{
		"type":            schema.NewStringValue(r.Type),
		"name":            schema.NewStringValue(r.Name),
		"pattern_type":    schema.NewStringValue(r.PatternType),
		"connect_cluster": schema.NewStringValue(r.ConnectCluster),
	}

	resource, diag := basetypes.NewObjectValue(types, values)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "resource", mapper.FromTerraform)
	}

	return resource, nil

}
