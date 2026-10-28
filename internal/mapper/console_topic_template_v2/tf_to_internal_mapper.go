package console_topic_template_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_topic_template_v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TFToInternalModel(ctx context.Context, r *schema.ConsoleTopicTemplateV2Model) (console.TopicTemplateResource, error) {
	labels, diag := schemaUtils.MapValueToStringMap(ctx, r.Labels)
	if diag.HasError() {
		return console.TopicTemplateResource{}, mapper.WrapDiagError(diag, "labels", mapper.FromTerraform)
	}

	spec, err := specTFToInternalModel(ctx, &r.Spec)
	if err != nil {
		return console.TopicTemplateResource{}, err
	}

	return console.NewTopicTemplateResource(r.Name.ValueString(), labels, spec), nil
}

func specTFToInternalModel(ctx context.Context, r *schema.SpecValue) (console.TopicTemplateSpec, error) {
	defaultsValue := schema.NewDefaultsValueNull()
	if !r.Defaults.IsNull() && !r.Defaults.IsUnknown() {
		var d diag.Diagnostics
		defaultsValue, d = schema.NewDefaultsValue(r.Defaults.AttributeTypes(ctx), r.Defaults.Attributes())
		if d.HasError() {
			return console.TopicTemplateSpec{}, mapper.WrapDiagError(d, "defaults", mapper.FromTerraform)
		}
	}

	defaults, err := defaultsTFToInternalModel(ctx, &defaultsValue)
	if err != nil {
		return console.TopicTemplateSpec{}, err
	}

	return console.TopicTemplateSpec{
		DisplayName: r.DisplayName.ValueString(),
		Description: r.Description.ValueString(),
		Defaults:    defaults,
	}, nil
}

func defaultsTFToInternalModel(ctx context.Context, r *schema.DefaultsValue) (console.TopicTemplateDefaults, error) {
	defaultsMeta, err := defaultsMetaTFToInternalModel(ctx, r)
	if err != nil {
		return console.TopicTemplateDefaults{}, err
	}

	defaultsSpec, err := defaultsSpecTFToInternalModel(ctx, r)
	if err != nil {
		return console.TopicTemplateDefaults{}, err
	}

	return console.TopicTemplateDefaults{
		Metadata: defaultsMeta,
		Spec:     defaultsSpec,
	}, nil
}

func defaultsMetaTFToInternalModel(ctx context.Context, r *schema.DefaultsValue) (console.TopicTemplateDefaultsMetadata, error) {
	if r.IsNull() || r.IsUnknown() || !schemaUtils.AttrIsSet(r.Metadata) {
		return console.TopicTemplateDefaultsMetadata{}, nil
	}

	metaValue, diag := schema.NewMetadataValue(r.Metadata.AttributeTypes(ctx), r.Metadata.Attributes())
	if diag.HasError() {
		return console.TopicTemplateDefaultsMetadata{}, mapper.WrapDiagError(diag, "defaults.metadata", mapper.FromTerraform)
	}

	labels, diag := schemaUtils.MapValueToStringMap(ctx, metaValue.Labels)
	if diag.HasError() {
		return console.TopicTemplateDefaultsMetadata{}, mapper.WrapDiagError(diag, "defaults.metadata.labels", mapper.FromTerraform)
	}

	return console.TopicTemplateDefaultsMetadata{
		Name:   metaValue.Name.ValueString(),
		Labels: labels,
	}, nil
}

func defaultsSpecTFToInternalModel(ctx context.Context, r *schema.DefaultsValue) (console.TopicTemplateDefaultsSpec, error) {
	if r.IsNull() || r.IsUnknown() || !schemaUtils.AttrIsSet(r.Spec) {
		return console.TopicTemplateDefaultsSpec{}, nil
	}

	specValue, diag := schema.NewDefaultsSpecValue(r.Spec.AttributeTypes(ctx), r.Spec.Attributes())
	if diag.HasError() {
		return console.TopicTemplateDefaultsSpec{}, mapper.WrapDiagError(diag, "defaults.spec", mapper.FromTerraform)
	}

	configs, diag := schemaUtils.MapValueToStringMap(ctx, specValue.Configs)
	if diag.HasError() {
		return console.TopicTemplateDefaultsSpec{}, mapper.WrapDiagError(diag, "defaults.spec.configs", mapper.FromTerraform)
	}

	result := console.TopicTemplateDefaultsSpec{
		Configs: configs,
	}

	if !specValue.Partitions.IsNull() && !specValue.Partitions.IsUnknown() {
		v := specValue.Partitions.ValueInt64()
		result.Partitions = &v
	}

	if !specValue.ReplicationFactor.IsNull() && !specValue.ReplicationFactor.IsUnknown() {
		v := specValue.ReplicationFactor.ValueInt64()
		result.ReplicationFactor = &v
	}

	return result, nil
}
