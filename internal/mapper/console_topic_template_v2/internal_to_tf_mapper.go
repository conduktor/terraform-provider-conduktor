package console_topic_template_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_topic_template_v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func InternalModelToTerraform(ctx context.Context, r *console.TopicTemplateResource) (schema.ConsoleTopicTemplateV2Model, error) {
	labels, diag := schemaUtils.StringMapToMapValue(ctx, r.Metadata.Labels)
	if diag.HasError() {
		return schema.ConsoleTopicTemplateV2Model{}, mapper.WrapDiagError(diag, "labels", mapper.IntoTerraform)
	}

	specValue, err := specInternalModelToTerraform(ctx, &r.Spec)
	if err != nil {
		return schema.ConsoleTopicTemplateV2Model{}, err
	}

	return schema.ConsoleTopicTemplateV2Model{
		Name:   types.StringValue(r.Metadata.Name),
		Labels: labels,
		Spec:   specValue,
	}, nil
}

func specInternalModelToTerraform(ctx context.Context, r *console.TopicTemplateSpec) (schema.SpecValue, error) {
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

func defaultsInternalModelToTerraform(ctx context.Context, r *console.TopicTemplateDefaults) (schema.DefaultsValue, error) {
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

func defaultsMetaInternalModelToTerraform(ctx context.Context, r *console.TopicTemplateDefaultsMetadata) (schema.MetadataValue, error) {
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

func defaultsSpecInternalModelToTerraform(ctx context.Context, r *console.TopicTemplateDefaultsSpec) (schema.DefaultsSpecValue, error) {
	unknownSpecObj, diag := schema.NewDefaultsSpecValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return schema.DefaultsSpecValue{}, mapper.WrapDiagError(diag, "defaults.spec", mapper.IntoTerraform)
	}
	typesMap := unknownSpecObj.AttributeTypes(ctx)
	valuesMap := schemaUtils.ValueMapFromTypes(ctx, typesMap)

	if r.Partitions != nil {
		valuesMap["partitions"] = basetypes.NewInt64Value(*r.Partitions)
	} else {
		valuesMap["partitions"] = basetypes.NewInt64Null()
	}

	if r.ReplicationFactor != nil {
		valuesMap["replication_factor"] = basetypes.NewInt64Value(*r.ReplicationFactor)
	} else {
		valuesMap["replication_factor"] = basetypes.NewInt64Null()
	}

	configs, diag := schemaUtils.StringMapToMapValue(ctx, r.Configs)
	if diag.HasError() {
		return schema.DefaultsSpecValue{}, mapper.WrapDiagError(diag, "defaults.spec.configs", mapper.IntoTerraform)
	}
	valuesMap["configs"] = configs

	value, diag := schema.NewDefaultsSpecValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.DefaultsSpecValue{}, mapper.WrapDiagError(diag, "defaults.spec", mapper.IntoTerraform)
	}
	return value, nil
}
