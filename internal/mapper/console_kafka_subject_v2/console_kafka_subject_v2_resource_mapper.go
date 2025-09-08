package console_kafka_subject_v2

import (
	"context"
	"github.com/conduktor/terraform-provider-conduktor/internal/customtypes"
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
	userLabels, diag := schema.MapValueToStringMap(ctx, r.Labels)
	if diag.HasError() {
		return console.KafkaSubjectResource{}, mapper.WrapDiagError(diag, "labels", mapper.FromTerraform)
	}
	managedLabels, diag := schema.MapValueToStringMap(ctx, r.ManagedLabels)
	if diag.HasError() {
		return console.KafkaSubjectResource{}, mapper.WrapDiagError(diag, "managed_labels", mapper.FromTerraform)
	}

	references, err := setValueToReferencesArray(ctx, r.Spec.References)
	if err != nil {
		return console.KafkaSubjectResource{}, err
	}

	var version *int
	if !r.Spec.Version.IsNull() && !r.Spec.Version.IsUnknown() {
		v := int(*r.Spec.Version.ValueInt64Pointer())
		version = &v
	}
	var id *int
	if !r.Spec.Id.IsNull() && !r.Spec.Id.IsUnknown() {
		v := int(*r.Spec.Id.ValueInt64Pointer())
		id = &v
	}

	return console.NewKafkaSubjectResource(
		r.Name.ValueString(),
		r.Cluster.ValueString(),
		mapper.MergeLabels(managedLabels, userLabels),
		console.KafkaSubjectSpec{
			Schema:        r.Spec.Schema.ValueString(),
			Format:        r.Spec.Format.ValueString(),
			Version:       version,
			Compatibility: r.Spec.Compatibility.ValueString(),
			Id:            id,
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
	var modelUserLabels, modelManagedLabels = mapper.SplitLabels(r.Metadata.Labels)
	labels, diag := schema.StringMapToMapValue(ctx, modelUserLabels)
	if diag.HasError() {
		return subject.ConsoleKafkaSubjectV2Model{}, mapper.WrapDiagError(diag, "labels", mapper.IntoTerraform)
	}

	managedLabels, diag := schema.StringMapToMapValue(ctx, modelManagedLabels)
	if diag.HasError() {
		return subject.ConsoleKafkaSubjectV2Model{}, mapper.WrapDiagError(diag, "managed_labels", mapper.IntoTerraform)
	}

	specValue, err := specInternalModelToTerraform(ctx, &r.Spec)
	if err != nil {
		return subject.ConsoleKafkaSubjectV2Model{}, err
	}

	return subject.ConsoleKafkaSubjectV2Model{
		Name:          types.StringValue(r.Metadata.Name),
		Cluster:       types.StringValue(r.Metadata.Cluster),
		Labels:        labels,
		ManagedLabels: managedLabels,
		Spec:          specValue,
	}, nil
}

func specInternalModelToTerraform(ctx context.Context, r *console.KafkaSubjectSpec) (subject.SpecValue, error) {
	unknownSpecObjectValue, diag := subject.NewSpecValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return subject.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["schema"] = customtypes.NewSchemaNormalizedValue(r.Schema)
	valuesMap["format"] = schema.NewStringValue(r.Format)
	valuesMap["compatibility"] = schema.NewStringValue(r.Compatibility)

	versionValue := types.Int64Null()
	if r.Version != nil {
		version := int64(*r.Version)
		versionValue = types.Int64Value(version)
	}
	valuesMap["version"] = versionValue

	idValue := types.Int64Null()
	if r.Id != nil {
		id := int64(*r.Id)
		idValue = types.Int64Value(id)
	}
	valuesMap["id"] = idValue

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
