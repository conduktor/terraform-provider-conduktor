package console_kafka_subject_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
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

func setValueToReferencesArray(ctx context.Context, set basetypes.SetValue) ([]console.KafkaSubjectReferences, error) {
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
	labels, diag := schema.StringMapToMapValue(ctx, r.Metadata.Labels)
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
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["schema"] = schema.NewStringValue(r.Schema)
	valuesMap["format"] = schema.NewStringValue(r.Format)
	valuesMap["version"] = types.Int64Value(int64(r.Version))
	valuesMap["compatibility"] = schema.NewStringValue(r.Compatibility)
	valuesMap["id"] = types.Int64Value(int64(r.Id))

	referencesValue, err := referencesToSetValue(ctx, r.References)
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

func referencesToSetValue(ctx context.Context, references []console.KafkaSubjectReferences) (basetypes.SetValue, error) {
	var referencesSet basetypes.SetValue
	var refs []attr.Value
	var diag diag.Diagnostics

	for _, p := range references {
		types := map[string]attr.Type{
			"name":    types.StringType,
			"subject": types.StringType,
			"version": types.Int64Type,
		}
		values := map[string]attr.Value{
			"name":    schema.NewStringValue(p.Name),
			"subject": schema.NewStringValue(p.Subject),
			"version": schema.NewInt64Value(int64(p.Version)),
		}
		value, diag := subject.NewReferencesValue(types, values)
		if diag.HasError() {
			return basetypes.SetValue{}, mapper.WrapDiagError(diag, "references", mapper.IntoTerraform)
		}
		refs = append(refs, value)
	}
	referencesSet, diag = types.SetValue(subject.ReferencesValue{}.Type(ctx), refs)
	if diag.HasError() {
		return basetypes.SetValue{}, mapper.WrapDiagError(diag, "references", mapper.IntoTerraform)
	}
	return referencesSet, nil
}
