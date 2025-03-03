package console_application_instance_v1

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	appinstance "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_application_instance_v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *appinstance.ConsoleApplicationInstanceV1Model) (console.ApplicationInstanceConsoleResource, error) {
	topicPolicyRef, diag := schema.SetValueToStringArray(ctx, r.Spec.TopicPolicyRef)
	if diag.HasError() {
		return console.ApplicationInstanceConsoleResource{}, mapper.WrapDiagError(diag, "topic_policy_ref", mapper.FromTerraform)
	}
	resources, err := setValueToResourceArray(ctx, r.Spec.Resources)
	if err != nil {
		return console.ApplicationInstanceConsoleResource{}, err
	}

	return console.NewApplicationInstanceConsoleResource(
		r.Name.ValueString(),
		r.Application.ValueString(),
		console.ApplicationInstanceConsoleSpec{
			Cluster:                          r.Spec.Cluster.ValueString(),
			TopicPolicyRef:                   topicPolicyRef,
			Resources:                        resources,
			ApplicationManagedServiceAccount: r.Spec.ApplicationManagedServiceAccount.ValueBool(),
			ServiceAccount:                   r.Spec.ServiceAccount.ValueString(),
			DefaultCatalogVisibility:         r.Spec.DefaultCatalogVisibility.ValueString(),
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *console.ApplicationInstanceConsoleResource) (appinstance.ConsoleApplicationInstanceV1Model, error) {
	var diag diag.Diagnostics
	// Ideally StringArrayToSetValue() would return a SetNull if needed.
	// However as of now it would make tests to fail for other resources.
	topicPolicyRef := basetypes.NewSetNull(basetypes.StringType{})
	if r.Spec.TopicPolicyRef != nil {
		topicPolicyRef, diag = schema.StringArrayToSetValue(r.Spec.TopicPolicyRef)
		if diag.HasError() {
			return appinstance.ConsoleApplicationInstanceV1Model{}, mapper.WrapDiagError(diag, "topic_policy_ref", mapper.FromTerraform)
		}
	}

	resourcesSet, err := resourceArrayToSetValue(ctx, r.Spec.Resources)
	if err != nil {
		return appinstance.ConsoleApplicationInstanceV1Model{}, err
	}

	specValue, diag := appinstance.NewSpecValue(
		map[string]attr.Type{
			"cluster":                             basetypes.StringType{},
			"topic_policy_ref":                    topicPolicyRef.Type(ctx),
			"resources":                           resourcesSet.Type(ctx),
			"application_managed_service_account": basetypes.BoolType{},
			"service_account":                     basetypes.StringType{},
			"default_catalog_visibility":          basetypes.StringType{},
		},
		map[string]attr.Value{
			"cluster":                             schema.NewStringValue(r.Spec.Cluster),
			"topic_policy_ref":                    topicPolicyRef,
			"resources":                           resourcesSet,
			"application_managed_service_account": basetypes.NewBoolValue(r.Spec.ApplicationManagedServiceAccount),
			"service_account":                     schema.NewStringValue(r.Spec.ServiceAccount),
			"default_catalog_visibility":          schema.NewStringValue(r.Spec.DefaultCatalogVisibility),
		},
	)
	if diag.HasError() {
		return appinstance.ConsoleApplicationInstanceV1Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return appinstance.ConsoleApplicationInstanceV1Model{
		Name:        schema.NewStringValue(r.Metadata.Name),
		Application: schema.NewStringValue(r.Metadata.Application),
		Spec:        specValue,
	}, nil
}

// Parse a Resources Array into a Set.
func resourceArrayToSetValue(ctx context.Context, arr []console.ResourceWithOwnership) (basetypes.SetValue, error) {
	var tfResources []attr.Value
	var diag diag.Diagnostics

	for _, p := range arr {
		types := map[string]attr.Type{
			"type":            basetypes.StringType{},
			"name":            basetypes.StringType{},
			"pattern_type":    basetypes.StringType{},
			"connect_cluster": basetypes.StringType{},
			"ownership_mode":  basetypes.StringType{},
		}
		values := map[string]attr.Value{
			"type":            schema.NewStringValue(p.Type),
			"name":            schema.NewStringValue(p.Name),
			"pattern_type":    schema.NewStringValue(p.PatternType),
			"connect_cluster": schema.NewStringValue(p.ConnectCluster),
			"ownership_mode":  schema.NewStringValue(p.OwnershipMode),
		}

		permObj, diag := appinstance.NewResourcesValue(types, values)
		if diag.HasError() {
			return basetypes.SetValue{}, mapper.WrapDiagError(diag, "resources", mapper.FromTerraform)
		}
		tfResources = append(tfResources, permObj)

	}

	resourcesList, diag := types.SetValue(appinstance.ResourcesValue{}.Type(ctx), tfResources)
	if diag.HasError() {
		return basetypes.SetValue{}, mapper.WrapDiagError(diag, "resources", mapper.FromTerraform)
	}

	return resourcesList, nil
}

// Parse a Set into an array of Resources.
func setValueToResourceArray(ctx context.Context, set basetypes.SetValue) ([]console.ResourceWithOwnership, error) {
	resources := make([]console.ResourceWithOwnership, 0)
	var diag diag.Diagnostics

	if !set.IsNull() && !set.IsUnknown() {
		var tfResources []appinstance.ResourcesValue
		diag = set.ElementsAs(ctx, &tfResources, false)
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "resources", mapper.FromTerraform)
		}

		for _, p := range tfResources {
			resources = append(resources, console.ResourceWithOwnership{
				Type:           p.ResourcesType.ValueString(),
				Name:           p.Name.ValueString(),
				PatternType:    p.PatternType.ValueString(),
				ConnectCluster: p.ConnectCluster.ValueString(),
				OwnershipMode:  p.OwnershipMode.ValueString(),
			})
		}

	}
	return resources, nil
}
