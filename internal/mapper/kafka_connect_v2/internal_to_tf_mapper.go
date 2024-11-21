package kafka_connect_v2

import (
	"context"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_kafka_connect_v2"
	"github.com/conduktor/terraform-provider-conduktor/internal/schema/validation"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func InternalModelToTerraform(ctx context.Context, r *model.KafkaConnectResource) (schema.KafkaConnectV2Model, error) {

	labels, diag := schemaUtils.StringMapToMapValue(ctx, r.Metadata.Labels)
	if diag.HasError() {
		return schema.KafkaConnectV2Model{}, mapper.WrapDiagError(diag, "labels", mapper.IntoTerraform)
	}

	specValue, err := specInternalModelToTerraform(ctx, &r.Spec)
	if err != nil {
		return schema.KafkaConnectV2Model{}, err
	}

	return schema.KafkaConnectV2Model{
		Name:    types.StringValue(r.Metadata.Name),
		Cluster: types.StringValue(r.Metadata.Cluster),
		Labels:  labels,
		Spec:    specValue,
	}, nil
}

func specInternalModelToTerraform(ctx context.Context, r *model.KafkaConnectSpec) (schema.SpecValue, error) {

	unknownSpecObjectValue, diag := schema.NewSpecValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schemaUtils.ValueMapFromTypes(ctx, typesMap)

	valuesMap["urls"] = schemaUtils.NewStringValue(r.Urls)
	valuesMap["display_name"] = schemaUtils.NewStringValue(r.DisplayName)
	valuesMap["ignore_untrusted_certificate"] = basetypes.NewBoolValue(r.IgnoreUntrustedCertificate)

	properties, diag := schemaUtils.StringMapToMapValue(ctx, r.Headers)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "headers", mapper.IntoTerraform)
	}
	valuesMap["headers"] = properties

	security, err := securityInternalModelToTerraform(ctx, r.Security)
	if err != nil {
		return schema.SpecValue{}, err
	}
	securityValue, diag := security.ToObjectValue(ctx)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "security", mapper.IntoTerraform)
	}
	valuesMap["security"] = securityValue

	value, diag := schema.NewSpecValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	return value, nil
}

func securityInternalModelToTerraform(ctx context.Context, r *model.KafkaConnectSecurity) (schema.SecurityValue, error) {
	if r == nil {
		return schema.NewSecurityValueNull(), nil
	}

	unknownSecurityObjectValue, diag := schema.NewSecurityValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return schema.SecurityValue{}, mapper.WrapDiagError(diag, "security", mapper.IntoTerraform)
	}
	var typesMap = unknownSecurityObjectValue.AttributeTypes(ctx)
	var valuesMap = schemaUtils.ValueMapFromTypes(ctx, typesMap)

	if r.BasicAuth != nil {
		valuesMap["type"] = schemaUtils.NewStringValue(validation.BasicAuthSchemaRegistrySecurity)
		valuesMap["username"] = schemaUtils.NewStringValue(r.BasicAuth.Username)
		valuesMap["password"] = schemaUtils.NewStringValue(r.BasicAuth.Password)
	}
	if r.BearerToken != nil {
		valuesMap["type"] = schemaUtils.NewStringValue(validation.BearerTokenSchemaRegistrySecurity)
		valuesMap["token"] = schemaUtils.NewStringValue(r.BearerToken.Token)
	}
	if r.SSLAuth != nil {
		valuesMap["type"] = schemaUtils.NewStringValue(validation.SSLAuthSchemaRegistrySecurity)
		valuesMap["certificate_chain"] = schemaUtils.NewStringValue(r.SSLAuth.CertificateChain)
		valuesMap["key"] = schemaUtils.NewStringValue(r.SSLAuth.Key)
	}

	value, diag := schema.NewSecurityValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.SecurityValue{}, mapper.WrapDiagError(diag, "security", mapper.IntoTerraform)
	}
	return value, nil

}
