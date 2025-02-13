package console_kafka_connect_v2

import (
	"context"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_kafka_connect_v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func InternalModelToTerraform(ctx context.Context, r *console.KafkaConnectResource) (schema.ConsoleKafkaConnectV2Model, error) {

	labels, diag := schemaUtils.StringMapToMapValue(ctx, r.Metadata.Labels)
	if diag.HasError() {
		return schema.ConsoleKafkaConnectV2Model{}, mapper.WrapDiagError(diag, "labels", mapper.IntoTerraform)
	}

	specValue, err := specInternalModelToTerraform(ctx, &r.Spec)
	if err != nil {
		return schema.ConsoleKafkaConnectV2Model{}, err
	}

	return schema.ConsoleKafkaConnectV2Model{
		Name:    types.StringValue(r.Metadata.Name),
		Cluster: types.StringValue(r.Metadata.Cluster),
		Labels:  labels,
		Spec:    specValue,
	}, nil
}

func specInternalModelToTerraform(ctx context.Context, r *console.KafkaConnectSpec) (schema.SpecValue, error) {

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

func securityInternalModelToTerraform(ctx context.Context, r *console.KafkaConnectSecurity) (schema.SecurityValue, error) {
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
		var basicTypesMap = schema.NewBasicAuthValueNull().AttributeTypes(ctx)
		var basicValuesMap = schemaUtils.ValueMapFromTypes(ctx, basicTypesMap)
		basicValuesMap["username"] = schemaUtils.NewStringValue(r.BasicAuth.Username)
		basicValuesMap["password"] = schemaUtils.NewStringValue(r.BasicAuth.Password)
		valuesMap["basic_auth"], diag = types.ObjectValue(basicTypesMap, basicValuesMap)
		if diag.HasError() {
			return schema.SecurityValue{}, mapper.WrapDiagError(diag, "security.basic_auth", mapper.IntoTerraform)
		}
	}

	if r.BearerToken != nil {
		var bearerTokenTypesMap = schema.NewBearerTokenValueNull().AttributeTypes(ctx)
		var bearerTokenValuesMap = schemaUtils.ValueMapFromTypes(ctx, bearerTokenTypesMap)
		bearerTokenValuesMap["token"] = schemaUtils.NewStringValue(r.BearerToken.Token)
		valuesMap["bearer_token"], diag = types.ObjectValue(bearerTokenTypesMap, bearerTokenValuesMap)
		if diag.HasError() {
			return schema.SecurityValue{}, mapper.WrapDiagError(diag, "security.bearer_token", mapper.IntoTerraform)
		}
	}

	if r.SSLAuth != nil {
		var sslAuthTypesMap = schema.NewSslAuthValueNull().AttributeTypes(ctx)
		var sslAuthValuesMap = schemaUtils.ValueMapFromTypes(ctx, sslAuthTypesMap)
		sslAuthValuesMap["certificate_chain"] = schemaUtils.NewStringValue(r.SSLAuth.CertificateChain)
		sslAuthValuesMap["key"] = schemaUtils.NewStringValue(r.SSLAuth.Key)
		valuesMap["ssl_auth"], diag = types.ObjectValue(sslAuthTypesMap, sslAuthValuesMap)
		if diag.HasError() {
			return schema.SecurityValue{}, mapper.WrapDiagError(diag, "security.ssl_auth", mapper.IntoTerraform)
		}
	}

	value, diag := schema.NewSecurityValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.SecurityValue{}, mapper.WrapDiagError(diag, "security", mapper.IntoTerraform)
	}
	return value, nil

}
