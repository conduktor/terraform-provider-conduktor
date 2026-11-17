package console_connector_template_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_connector_template_v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func InternalModelToTerraform(ctx context.Context, r *console.ConnectorTemplateResource) (schema.ConsoleConnectorTemplateV2Model, error) {
	var modelUserLabels, modelManagedLabels = mapper.SplitLabels(r.Metadata.Labels)

	labels, diag := schemaUtils.StringMapToMapValue(ctx, modelUserLabels)
	if diag.HasError() {
		return schema.ConsoleConnectorTemplateV2Model{}, mapper.WrapDiagError(diag, "labels", mapper.IntoTerraform)
	}

	managedLabels, diag := schemaUtils.StringMapToMapValue(ctx, modelManagedLabels)
	if diag.HasError() {
		return schema.ConsoleConnectorTemplateV2Model{}, mapper.WrapDiagError(diag, "managed_labels", mapper.IntoTerraform)
	}

	specValue, err := specInternalModelToTerraform(ctx, &r.Spec)
	if err != nil {
		return schema.ConsoleConnectorTemplateV2Model{}, err
	}

	return schema.ConsoleConnectorTemplateV2Model{
		Name:          types.StringValue(r.Metadata.Name),
		Labels:        labels,
		ManagedLabels: managedLabels,
		Spec:          specValue,
	}, nil
}

func specInternalModelToTerraform(ctx context.Context, r *console.ConnectorTemplateSpec) (schema.SpecValue, error) {
	unknownSpecObj, diag := schema.NewSpecValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	typesMap := unknownSpecObj.AttributeTypes(ctx)
	valuesMap := schemaUtils.ValueMapFromTypes(ctx, typesMap)

	valuesMap["display_name"] = schemaUtils.NewStringValue(r.DisplayName)
	valuesMap["description"] = schemaUtils.NewStringValue(r.Description)

	defaultsValue, err := defaultsInternalModelToTerraform(ctx, &r.Defaults)
	if err != nil {
		return schema.SpecValue{}, err
	}
	defaultsObj, diag := defaultsValue.ToObjectValue(ctx)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "spec.defaults", mapper.IntoTerraform)
	}
	valuesMap["defaults"] = defaultsObj

	value, diag := schema.NewSpecValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	return value, nil
}

func defaultsInternalModelToTerraform(ctx context.Context, r *console.ConnectorTemplateDefaults) (schema.DefaultsValue, error) {
	unknownDefaultsObj, diag := schema.NewDefaultsValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return schema.DefaultsValue{}, mapper.WrapDiagError(diag, "defaults", mapper.IntoTerraform)
	}
	typesMap := unknownDefaultsObj.AttributeTypes(ctx)
	valuesMap := schemaUtils.ValueMapFromTypes(ctx, typesMap)

	metaValue, err := defaultsMetaInternalModelToTerraform(ctx, &r.Metadata)
	if err != nil {
		return schema.DefaultsValue{}, err
	}
	metaObj, diag := metaValue.ToObjectValue(ctx)
	if diag.HasError() {
		return schema.DefaultsValue{}, mapper.WrapDiagError(diag, "defaults.metadata", mapper.IntoTerraform)
	}
	valuesMap["metadata"] = metaObj

	defaultsSpecValue, err := defaultsSpecInternalModelToTerraform(ctx, &r.Spec)
	if err != nil {
		return schema.DefaultsValue{}, err
	}
	defaultsSpecObj, diag := defaultsSpecValue.ToObjectValue(ctx)
	if diag.HasError() {
		return schema.DefaultsValue{}, mapper.WrapDiagError(diag, "defaults.spec", mapper.IntoTerraform)
	}
	valuesMap["spec"] = defaultsSpecObj

	value, diag := schema.NewDefaultsValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.DefaultsValue{}, mapper.WrapDiagError(diag, "defaults", mapper.IntoTerraform)
	}
	return value, nil
}

func defaultsMetaInternalModelToTerraform(ctx context.Context, r *console.ConnectorTemplateDefaultsMetadata) (schema.MetadataValue, error) {
	if r.Name == "" && r.Labels == nil {
		return schema.NewMetadataValueNull(), nil
	}

	unknownMetaObj, diag := schema.NewMetadataValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return schema.MetadataValue{}, mapper.WrapDiagError(diag, "defaults.metadata", mapper.IntoTerraform)
	}
	typesMap := unknownMetaObj.AttributeTypes(ctx)
	valuesMap := schemaUtils.ValueMapFromTypes(ctx, typesMap)

	valuesMap["name"] = schemaUtils.NewStringValue(r.Name)

	labels, diag := schemaUtils.StringMapToMapValue(ctx, r.Labels)
	if diag.HasError() {
		return schema.MetadataValue{}, mapper.WrapDiagError(diag, "defaults.metadata.labels", mapper.IntoTerraform)
	}
	valuesMap["labels"] = labels

	value, diag := schema.NewMetadataValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.MetadataValue{}, mapper.WrapDiagError(diag, "defaults.metadata", mapper.IntoTerraform)
	}
	return value, nil
}

func defaultsSpecInternalModelToTerraform(ctx context.Context, r *console.ConnectorTemplateDefaultsSpec) (schema.ConnectorDefaultsSpecValue, error) {
	unknownSpecObj, diag := schema.NewConnectorDefaultsSpecValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return schema.ConnectorDefaultsSpecValue{}, mapper.WrapDiagError(diag, "defaults.spec", mapper.IntoTerraform)
	}
	typesMap := unknownSpecObj.AttributeTypes(ctx)
	valuesMap := schemaUtils.ValueMapFromTypes(ctx, typesMap)

	valuesMap["class"] = schemaUtils.NewStringValue(r.Class)

	config, diag := schemaUtils.StringMapToMapValue(ctx, r.Config)
	if diag.HasError() {
		return schema.ConnectorDefaultsSpecValue{}, mapper.WrapDiagError(diag, "defaults.spec.config", mapper.IntoTerraform)
	}
	valuesMap["config"] = config

	value, diag := schema.NewConnectorDefaultsSpecValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.ConnectorDefaultsSpecValue{}, mapper.WrapDiagError(diag, "defaults.spec", mapper.IntoTerraform)
	}
	return value, nil
}
