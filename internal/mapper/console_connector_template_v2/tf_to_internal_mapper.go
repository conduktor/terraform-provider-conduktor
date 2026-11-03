package console_connector_template_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_connector_template_v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TFToInternalModel(ctx context.Context, r *schema.ConsoleConnectorTemplateV2Model) (console.ConnectorTemplateResource, error) {
	labels, diag := schemaUtils.MapValueToStringMap(ctx, r.Labels)
	if diag.HasError() {
		return console.ConnectorTemplateResource{}, mapper.WrapDiagError(diag, "labels", mapper.FromTerraform)
	}

	spec, err := specTFToInternalModel(ctx, &r.Spec)
	if err != nil {
		return console.ConnectorTemplateResource{}, err
	}

	return console.NewConnectorTemplateResource(r.Name.ValueString(), labels, spec), nil
}

func specTFToInternalModel(ctx context.Context, r *schema.SpecValue) (console.ConnectorTemplateSpec, error) {
	defaultsValue := schema.NewDefaultsValueNull()
	if !r.Defaults.IsNull() && !r.Defaults.IsUnknown() {
		var d diag.Diagnostics
		defaultsValue, d = schema.NewDefaultsValue(r.Defaults.AttributeTypes(ctx), r.Defaults.Attributes())
		if d.HasError() {
			return console.ConnectorTemplateSpec{}, mapper.WrapDiagError(d, "defaults", mapper.FromTerraform)
		}
	}

	defaults, err := defaultsTFToInternalModel(ctx, &defaultsValue)
	if err != nil {
		return console.ConnectorTemplateSpec{}, err
	}

	return console.ConnectorTemplateSpec{
		DisplayName: r.DisplayName.ValueString(),
		Description: r.Description.ValueString(),
		Defaults:    defaults,
	}, nil
}

func defaultsTFToInternalModel(ctx context.Context, r *schema.DefaultsValue) (console.ConnectorTemplateDefaults, error) {
	defaultsMeta, err := defaultsMetaTFToInternalModel(ctx, r)
	if err != nil {
		return console.ConnectorTemplateDefaults{}, err
	}

	defaultsSpec, err := defaultsSpecTFToInternalModel(ctx, r)
	if err != nil {
		return console.ConnectorTemplateDefaults{}, err
	}

	return console.ConnectorTemplateDefaults{
		Metadata: defaultsMeta,
		Spec:     defaultsSpec,
	}, nil
}

func defaultsMetaTFToInternalModel(ctx context.Context, r *schema.DefaultsValue) (console.ConnectorTemplateDefaultsMetadata, error) {
	if r.IsNull() || r.IsUnknown() || !schemaUtils.AttrIsSet(r.Metadata) {
		return console.ConnectorTemplateDefaultsMetadata{}, nil
	}

	metaValue, diag := schema.NewMetadataValue(r.Metadata.AttributeTypes(ctx), r.Metadata.Attributes())
	if diag.HasError() {
		return console.ConnectorTemplateDefaultsMetadata{}, mapper.WrapDiagError(diag, "defaults.metadata", mapper.FromTerraform)
	}

	labels, diag := schemaUtils.MapValueToStringMap(ctx, metaValue.Labels)
	if diag.HasError() {
		return console.ConnectorTemplateDefaultsMetadata{}, mapper.WrapDiagError(diag, "defaults.metadata.labels", mapper.FromTerraform)
	}

	return console.ConnectorTemplateDefaultsMetadata{
		Name:   metaValue.Name.ValueString(),
		Labels: labels,
	}, nil
}

func defaultsSpecTFToInternalModel(ctx context.Context, r *schema.DefaultsValue) (console.ConnectorTemplateDefaultsSpec, error) {
	// defaults.spec is required — it will always be set
	specValue, diag := schema.NewConnectorDefaultsSpecValue(r.Spec.AttributeTypes(ctx), r.Spec.Attributes())
	if diag.HasError() {
		return console.ConnectorTemplateDefaultsSpec{}, mapper.WrapDiagError(diag, "defaults.spec", mapper.FromTerraform)
	}

	config, diag := schemaUtils.MapValueToStringMap(ctx, specValue.Config)
	if diag.HasError() {
		return console.ConnectorTemplateDefaultsSpec{}, mapper.WrapDiagError(diag, "defaults.spec.config", mapper.FromTerraform)
	}

	return console.ConnectorTemplateDefaultsSpec{
		Class:  specValue.Class.ValueString(),
		Config: config,
	}, nil
}
