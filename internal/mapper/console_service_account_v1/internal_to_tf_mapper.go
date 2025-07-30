package console_service_account_v1

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	serviceacc "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_service_account_v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func InternalModelToTerraform(ctx context.Context, r *console.ServiceAccountResource) (serviceacc.ConsoleServiceAccountV1Model, error) {
	labels, diag := schema.StringMapToMapValue(ctx, r.Metadata.Labels)
	if diag.HasError() {
		return serviceacc.ConsoleServiceAccountV1Model{}, mapper.WrapDiagError(diag, "labels", mapper.IntoTerraform)
	}

	specValue, err := specInternalModelToTerraform(ctx, &r.Spec)
	if err != nil {
		return serviceacc.ConsoleServiceAccountV1Model{}, err
	}

	return serviceacc.ConsoleServiceAccountV1Model{
		Name:        schema.NewStringValue(r.Metadata.Name),
		Cluster:     schema.NewStringValue(r.Metadata.Cluster),
		Labels:      labels,
		AppInstance: schema.NewStringValue(r.Metadata.AppInstance),
		Spec:        specValue,
	}, nil
}

func specInternalModelToTerraform(ctx context.Context, r *console.ServiceAccountSpec) (serviceacc.SpecValue, error) {
	unknownSpecObjectValue, diag := serviceacc.NewSpecValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return serviceacc.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	authorization, err := authorizationInternalModelToTerraform(ctx, r.Authorization)
	if err != nil {
		return serviceacc.SpecValue{}, err
	}
	authValue, diag := authorization.ToObjectValue(ctx)
	if diag.HasError() {
		return serviceacc.SpecValue{}, mapper.WrapDiagError(diag, "authorization", mapper.IntoTerraform)
	}
	valuesMap["authorization"] = authValue

	value, diag := serviceacc.NewSpecValue(typesMap, valuesMap)
	if diag.HasError() {
		return serviceacc.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	return value, nil
}

func authorizationInternalModelToTerraform(ctx context.Context, r *console.ServiceAccountAuthorization) (serviceacc.AuthorizationValue, error) {
	if r == nil || (r.Aiven == nil && r.Kafka == nil) {
		return serviceacc.NewAuthorizationValueNull(), nil
	}

	unknownAuthObjectValue, diag := serviceacc.NewAuthorizationValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return serviceacc.AuthorizationValue{}, mapper.WrapDiagError(diag, "authorization", mapper.IntoTerraform)
	}
	var typesMap = unknownAuthObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	if r.Aiven != nil {
		aiven, err := aivenInternalModelToTerraform(ctx, r.Aiven)
		if err != nil {
			return serviceacc.AuthorizationValue{}, err
		}
		aivenValue, diag := aiven.ToObjectValue(ctx)
		if diag.HasError() {
			return serviceacc.AuthorizationValue{}, mapper.WrapDiagError(diag, "authorization.aiven", mapper.IntoTerraform)
		}
		valuesMap["aiven"] = aivenValue
	}

	if r.Kafka != nil {
		kafka, err := kafkaInternalModelToTerraform(ctx, r.Kafka)
		if err != nil {
			return serviceacc.AuthorizationValue{}, err
		}
		kafkaValue, diag := kafka.ToObjectValue(ctx)
		if diag.HasError() {
			return serviceacc.AuthorizationValue{}, mapper.WrapDiagError(diag, "authorization.kafka", mapper.IntoTerraform)
		}
		valuesMap["kafka"] = kafkaValue
	}

	value, diag := serviceacc.NewAuthorizationValue(typesMap, valuesMap)
	if diag.HasError() {
		return serviceacc.AuthorizationValue{}, mapper.WrapDiagError(diag, "authorization", mapper.IntoTerraform)
	}
	return value, nil
}

func aivenInternalModelToTerraform(ctx context.Context, r *console.ServiceAccountAuthAiven) (serviceacc.AivenValue, error) {
	if r == nil {
		return serviceacc.NewAivenValueNull(), nil
	}

	unknownAivenObjectValue, diag := serviceacc.NewAivenValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return serviceacc.AivenValue{}, mapper.WrapDiagError(diag, "aiven", mapper.IntoTerraform)
	}
	var typesMap = unknownAivenObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	acls, err := aivenAclsArrayToSetValue(ctx, r.ACLS)
	if err != nil {
		return serviceacc.AivenValue{}, err
	}
	valuesMap["acls"] = acls

	value, diag := serviceacc.NewAivenValue(typesMap, valuesMap)
	if diag.HasError() {
		return serviceacc.AivenValue{}, mapper.WrapDiagError(diag, "aiven", mapper.IntoTerraform)
	}
	return value, nil
}

func kafkaInternalModelToTerraform(ctx context.Context, r *console.ServiceAccountAuthKafka) (serviceacc.KafkaValue, error) {
	if r == nil {
		return serviceacc.NewKafkaValueNull(), nil
	}

	unknownKafkaObjectValue, diag := serviceacc.NewKafkaValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return serviceacc.KafkaValue{}, mapper.WrapDiagError(diag, "kafka", mapper.IntoTerraform)
	}
	var typesMap = unknownKafkaObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	acls, err := kafkaAclsArrayToSetValue(ctx, r.ACLS)
	if err != nil {
		return serviceacc.KafkaValue{}, err
	}
	valuesMap["acls"] = acls

	value, diag := serviceacc.NewKafkaValue(typesMap, valuesMap)
	if diag.HasError() {
		return serviceacc.KafkaValue{}, mapper.WrapDiagError(diag, "kafka", mapper.IntoTerraform)
	}
	return value, nil
}

func aivenAclsArrayToSetValue(ctx context.Context, arr []console.ServiceAccountAuthAivenACL) (basetypes.SetValue, error) {
	var aclsList basetypes.SetValue
	var tfACLs []attr.Value
	var diag diag.Diagnostics

	for _, p := range arr {
		types := map[string]attr.Type{
			"name":          basetypes.StringType{},
			"resource_type": basetypes.StringType{},
			"permission":    basetypes.StringType{},
		}
		values := map[string]attr.Value{
			"name":          schema.NewStringValue(p.Name),
			"resource_type": schema.NewStringValue(p.ResourceType),
			"permission":    schema.NewStringValue(p.Permission),
		}

		permObj, diag := serviceacc.NewAivenAclsValue(types, values)
		if diag.HasError() {
			return basetypes.SetValue{}, mapper.WrapDiagError(diag, "acls", mapper.FromTerraform)
		}
		tfACLs = append(tfACLs, permObj)

	}

	aclsList, diag = types.SetValue(serviceacc.AivenAclsValue{}.Type(ctx), tfACLs)
	if diag.HasError() {
		return basetypes.SetValue{}, mapper.WrapDiagError(diag, "acls", mapper.FromTerraform)
	}

	return aclsList, nil
}

func kafkaAclsArrayToSetValue(ctx context.Context, arr []console.ServiceAccountAuthKafkaACL) (basetypes.SetValue, error) {
	var aclsList basetypes.SetValue
	var tfACLs []attr.Value
	var diag diag.Diagnostics

	for _, p := range arr {
		operations, diag := schema.StringArrayToSetValue(p.Operations)
		if diag.HasError() {
			return basetypes.SetValue{}, mapper.WrapDiagError(diag, "acls.operations", mapper.FromTerraform)
		}

		types := map[string]attr.Type{
			"type":            basetypes.StringType{},
			"name":            basetypes.StringType{},
			"pattern_type":    basetypes.StringType{},
			"connect_cluster": basetypes.StringType{},
			"operations":      operations.Type(ctx),
			"host":            basetypes.StringType{},
			"permission":      basetypes.StringType{},
		}
		values := map[string]attr.Value{
			"type":            schema.NewStringValue(p.Type),
			"name":            schema.NewStringValue(p.Name),
			"pattern_type":    schema.NewStringValue(p.PatternType),
			"connect_cluster": schema.NewStringValue(p.ConnectCluster),
			"operations":      operations,
			"host":            schema.NewStringValue(p.Host),
			"permission":      schema.NewStringValue(p.Permission),
		}

		permObj, diag := serviceacc.NewKafkaAclsValue(types, values)
		if diag.HasError() {
			return basetypes.SetValue{}, mapper.WrapDiagError(diag, "acls", mapper.FromTerraform)
		}
		tfACLs = append(tfACLs, permObj)

	}

	aclsList, diag = types.SetValue(serviceacc.KafkaAclsValue{}.Type(ctx), tfACLs)
	if diag.HasError() {
		return basetypes.SetValue{}, mapper.WrapDiagError(diag, "kafka", mapper.FromTerraform)
	}

	return aclsList, nil
}
