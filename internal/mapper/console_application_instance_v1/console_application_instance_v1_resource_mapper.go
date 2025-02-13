package console_application_instance_v1

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	appinstance "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_application_instance_v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *appinstance.ConsoleApplicationInstanceV1Model) (console.ApplicationInstanceConsoleResource, error) {

	// TODO Resources
	// permissions, err := schema.SetValueToPermissionArray(ctx, schema.GROUPS, r.Spec.Permissions)
	// if err != nil {
	// 	return console.ApplicationInstanceConsoleResource{}, err
	// }

	return console.NewApplicationInstanceConsoleResource(
		r.Name.ValueString(),
		console.ApplicationInstanceConsoleSpec{
			Cluster:                          r.Spec.Cluster.ValueString(),
			TopicPolicyRef:                   nil, // TODO
			Resources:                        nil, // TODO
			ApplicationManagedServiceAccount: r.Spec.ApplicationManagedServiceAccount.ValueBool(),
			ServiceAccount:                   r.Spec.ServiceAccount.ValueString(),
			DefaultCatalogVisibility:         r.Spec.DefaultCatalogVisibility.ValueString(),
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *console.ApplicationInstanceConsoleResource) (appinstance.ConsoleApplicationInstanceV1Model, error) {
	// permissionsList, err := schema.PermissionArrayToSetValue(ctx, schema.GROUPS, r.Spec.Permissions)
	// if err != nil {
	// 	return appinstance.ConsoleGroupV2Model{}, err
	// }

	specValue, diag := appinstance.NewSpecValue(
		map[string]attr.Type{
			"cluster":                             basetypes.StringType{},
			"topic_policy_ref":                    nil, // TODO
			"resources":                           nil, // TODO permissionsList.Type(ctx),
			"application_managed_service_account": basetypes.BoolType{},
			"service_account":                     basetypes.StringType{},
			"default_catalog_visibility":          basetypes.StringType{},
		},
		map[string]attr.Value{
			"cluster":                             schema.NewStringValue(r.Spec.Cluster),
			"topic_policy_ref":                    nil, // TODO
			"resources":                           nil, // TODO permissionsList.Type(ctx),
			"application_managed_service_account": basetypes.NewBoolValue(r.Spec.ApplicationManagedServiceAccount),
			"service_account":                     schema.NewStringValue(r.Spec.ServiceAccount),
			"default_catalog_visibility":          schema.NewStringValue(r.Spec.DefaultCatalogVisibility),
		},
	)
	if diag.HasError() {
		return appinstance.ConsoleApplicationInstanceV1Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return appinstance.ConsoleApplicationInstanceV1Model{
		Name: types.StringValue(r.Metadata.Name),
		Spec: specValue,
	}, nil
}
