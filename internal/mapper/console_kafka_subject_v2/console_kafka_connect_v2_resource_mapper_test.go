package console_kafka_subject_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	subject "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_kafka_subject_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *subject.ConsoleKafkaSubjectV2Model) (console.KafkaSubjectResource, error) {
	labels, diag := schema.MapValueToStringMap(ctx, r.Labels)
	if diag.HasError() {
		return console.KafkaSubjectResource{}, mapper.WrapDiagError(diag, "labels", mapper.FromTerraform)
	}

	references, err := setValueToReferencesArray(ctx, r.Spec.References)
	if err != nil {
		return console.KafkaSubjectResource{}, err
	}

	return console.NewKafkaSubjectResource(
		r.Name.ValueString(),
		r.Cluster.ValueString(),
		labels,
		console.KafkaSubjectSpec{
			Schema:        r.Spec.Schema.ValueString(),
			Format:        r.Spec.Format.ValueString(),
			Version:       int(r.Spec.Version.ValueInt64()),
			Compatibility: r.Spec.Compatibility.ValueString(),
			Id:            int(r.Spec.Id.ValueInt64()),
			References:    references,
		},
	), nil
}

func setValueToReferencesArray(ctx context.Context, set basetypes.ListValue) ([]console.KafkaSubjectReferences, error) {
	references := make([]console.KafkaSubjectReferences, 0)
	var diag diag.Diagnostics

	if !set.IsNull() && !set.IsUnknown() {
		var tfResources []subject.ReferencesValue
		diag = set.ElementsAs(ctx, &tfResources, false)
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "references", mapper.FromTerraform)
		}

		for _, p := range tfResources {
			references = append(references, console.KafkaSubjectReferences{
				Name:    p.Name.ValueString(),
				Subject: p.Subject.ValueString(),
				Version: int(p.Version.ValueInt64()),
			})
		}
	}
	return references, nil
}

func InternalModelToTerraform(ctx context.Context, r *console.KafkaSubjectResource) (subject.ConsoleKafkaSubjectV2Model, error) {
	labels, diag := schemaUtils.StringMapToMapValue(ctx, r.Metadata.Labels)
	if diag.HasError() {
		return subject.ConsoleKafkaSubjectV2Model{}, mapper.WrapDiagError(diag, "labels", mapper.IntoTerraform)
	}

	specValue, err := specInternalModelToTerraform(ctx, &r.Spec)
	if err != nil {
		return subject.ConsoleKafkaSubjectV2Model{}, err
	}

	return subject.ConsoleKafkaSubjectV2Model{
		Name:    types.StringValue(r.Metadata.Name),
		Cluster: types.StringValue(r.Metadata.Cluster),
		Labels:  labels,
		Spec:    specValue,
	}, nil
}

func specInternalModelToTerraform(ctx context.Context, r *console.KafkaSubjectSpec) (subject.SpecValue, error) {
	unknownSpecObjectValue, diag := subject.NewSpecValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return subject.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schemaUtils.ValueMapFromTypes(ctx, typesMap)

	valuesMap["schema"] = schemaUtils.NewStringValue(r.Schema)
	valuesMap["format"] = schemaUtils.NewStringValue(r.Format)
	valuesMap["version"] = basetypes.NewInt64Value(int64(r.Version))
	valuesMap["compatibility"] = schemaUtils.NewStringValue(r.Compatibility)
	valuesMap["id"] = basetypes.NewInt64Value(int64(r.Id))

	referencesValue, err := referencesToListValue(ctx, r.References)
	if err != nil {
		return subject.SpecValue{}, err
	}
	valuesMap["references"] = referencesValue

	value, diag := subject.NewSpecValue(typesMap, valuesMap)
	if diag.HasError() {
		return subject.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	return value, nil
}

func referencesToListValue(ctx context.Context, references []console.KafkaSubjectReferences) (basetypes.ListValue, error) {
	if len(references) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":    types.StringType,
				"subject": types.StringType,
				"version": types.Int64Type,
			},
		}), nil
	}

	referencesList := make([]subject.ReferencesValue, 0)
	for _, reference := range references {
		referencesList = append(referencesList, subject.ReferencesValue{
			Name:    types.StringValue(reference.Name),
			Subject: types.StringValue(reference.Subject),
			Version: types.Int64Value(int64(reference.Version)),
		})
	}

	referencesValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":    types.StringType,
			"subject": types.StringType,
			"version": types.Int64Type,
		},
	}, referencesList)

	if diags.HasError() {
		return basetypes.ListValue{}, mapper.WrapDiagError(diags, "references", mapper.IntoTerraform)
	}

	return referencesValue, nil
}
