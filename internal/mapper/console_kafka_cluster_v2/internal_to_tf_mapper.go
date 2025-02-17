package console_kafka_cluster_v2

import (
	"context"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_kafka_cluster_v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func InternalModelToTerraform(ctx context.Context, r *console.KafkaClusterResource) (schema.ConsoleKafkaClusterV2Model, error) {

	labels, diag := schemaUtils.StringMapToMapValue(ctx, r.Metadata.Labels)
	if diag.HasError() {
		return schema.ConsoleKafkaClusterV2Model{}, mapper.WrapDiagError(diag, "labels", mapper.IntoTerraform)
	}

	specValue, err := specInternalModelToTerraform(ctx, &r.Spec)
	if err != nil {
		return schema.ConsoleKafkaClusterV2Model{}, err
	}

	return schema.ConsoleKafkaClusterV2Model{
		Name:   types.StringValue(r.Metadata.Name),
		Labels: labels,
		Spec:   specValue,
	}, nil
}

func specInternalModelToTerraform(ctx context.Context, r *console.KafkaClusterSpec) (schema.SpecValue, error) {

	unknownSpecObjectValue, diag := schema.NewSpecValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schemaUtils.ValueMapFromTypes(ctx, typesMap)

	valuesMap["bootstrap_servers"] = schemaUtils.NewStringValue(r.BootstrapServers)
	valuesMap["display_name"] = schemaUtils.NewStringValue(r.DisplayName)
	valuesMap["color"] = schemaUtils.NewStringValue(r.Color)
	valuesMap["icon"] = schemaUtils.NewStringValue(r.Icon)
	valuesMap["ignore_untrusted_certificate"] = basetypes.NewBoolValue(r.IgnoreUntrustedCertificate)

	properties, diag := schemaUtils.StringMapToMapValue(ctx, r.Properties)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "properties", mapper.IntoTerraform)
	}
	valuesMap["properties"] = properties

	kafkaFlavor, err := kafkaFlavorInternalModelToTerraform(ctx, r.KafkaFlavor)
	if err != nil {
		return schema.SpecValue{}, err
	}
	kafkaFlavorValue, diag := kafkaFlavor.ToObjectValue(ctx)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "kafka_flavor", mapper.IntoTerraform)
	}
	valuesMap["kafka_flavor"] = kafkaFlavorValue

	schemaRegistry, err := schemaRegistryInternalModelToTerraform(ctx, r.SchemaRegistry)
	if err != nil {
		return schema.SpecValue{}, err
	}
	schemaRegistryValue, diag := schemaRegistry.ToObjectValue(ctx)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "schema_registry", mapper.IntoTerraform)
	}
	valuesMap["schema_registry"] = schemaRegistryValue

	value, diag := schema.NewSpecValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	return value, nil
}

func kafkaFlavorInternalModelToTerraform(ctx context.Context, r *console.KafkaFlavor) (schema.KafkaFlavorValue, error) {
	if r == nil || (r.Aiven == nil && r.Confluent == nil && r.Gateway == nil) {
		return schema.NewKafkaFlavorValueNull(), nil
	}

	unknownFlavorObjectValue, diag := schema.NewKafkaFlavorValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return schema.KafkaFlavorValue{}, mapper.WrapDiagError(diag, "kafka_flavor", mapper.IntoTerraform)
	}
	var typesMap = unknownFlavorObjectValue.AttributeTypes(ctx)
	var valuesMap = schemaUtils.ValueMapFromTypes(ctx, typesMap)

	if r.Aiven != nil {
		var aivenTypesMap = schema.NewAivenValueNull().AttributeTypes(ctx)
		var aivenValuesMap = schemaUtils.ValueMapFromTypes(ctx, aivenTypesMap)
		aivenValuesMap["api_token"] = schemaUtils.NewStringValue(r.Aiven.ApiToken)
		aivenValuesMap["project"] = schemaUtils.NewStringValue(r.Aiven.Project)
		aivenValuesMap["service_name"] = schemaUtils.NewStringValue(r.Aiven.ServiceName)
		valuesMap["aiven"], diag = types.ObjectValue(aivenTypesMap, aivenValuesMap)
		if diag.HasError() {
			return schema.KafkaFlavorValue{}, mapper.WrapDiagError(diag, "kafka_flavor", mapper.IntoTerraform)
		}
	}

	if r.Confluent != nil {
		var confluentTypesMap = schema.NewConfluentValueNull().AttributeTypes(ctx)
		var confluentValuesMap = schemaUtils.ValueMapFromTypes(ctx, confluentTypesMap)
		confluentValuesMap["key"] = schemaUtils.NewStringValue(r.Confluent.Key)
		confluentValuesMap["secret"] = schemaUtils.NewStringValue(r.Confluent.Secret)
		confluentValuesMap["confluent_environment_id"] = schemaUtils.NewStringValue(r.Confluent.ConfluentEnvironmentId)
		confluentValuesMap["confluent_cluster_id"] = schemaUtils.NewStringValue(r.Confluent.ConfluentClusterId)
		valuesMap["confluent"], diag = types.ObjectValue(confluentTypesMap, confluentValuesMap)
		if diag.HasError() {
			return schema.KafkaFlavorValue{}, mapper.WrapDiagError(diag, "kafka_flavor", mapper.IntoTerraform)
		}
	}

	if r.Gateway != nil {
		var gatewayTypesMap = schema.NewGatewayValueNull().AttributeTypes(ctx)
		var gatewayValuesMap = schemaUtils.ValueMapFromTypes(ctx, gatewayTypesMap)
		gatewayValuesMap["ignore_untrusted_certificate"] = basetypes.NewBoolValue(false) // default value
		gatewayValuesMap["virtual_cluster"] = basetypes.NewStringValue("passthrough")    // default value

		gatewayValuesMap["url"] = schemaUtils.NewStringValue(r.Gateway.Url)
		gatewayValuesMap["user"] = schemaUtils.NewStringValue(r.Gateway.User)
		gatewayValuesMap["password"] = schemaUtils.NewStringValue(r.Gateway.Password)
		gatewayValuesMap["virtual_cluster"] = schemaUtils.NewStringValue(r.Gateway.VirtualCluster)
		gatewayValuesMap["ignore_untrusted_certificate"] = basetypes.NewBoolValue(r.Gateway.IgnoreUntrustedCertificate)
		valuesMap["gateway"], diag = types.ObjectValue(gatewayTypesMap, gatewayValuesMap)
		if diag.HasError() {
			return schema.KafkaFlavorValue{}, mapper.WrapDiagError(diag, "kafka_flavor", mapper.IntoTerraform)
		}
	}

	value, diag := schema.NewKafkaFlavorValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.KafkaFlavorValue{}, mapper.WrapDiagError(diag, "kafka_flavor", mapper.IntoTerraform)
	}
	return value, nil
}

func schemaRegistryInternalModelToTerraform(ctx context.Context, r *model.SchemaRegistry) (schema.SchemaRegistryValue, error) {
	if r == nil || (r.Glue == nil && r.ConfluentLike == nil) {
		return schema.NewSchemaRegistryValueNull(), nil
	}

	var unknownSR = schema.NewSchemaRegistryValueUnknown()
	if r.ConfluentLike == nil && r.Glue == nil {
		return unknownSR, nil
	}
	unknownSRObjectValue, diag := unknownSR.ToObjectValue(ctx)
	if diag.HasError() {
		return schema.SchemaRegistryValue{}, mapper.WrapDiagError(diag, "schema_registry", mapper.IntoTerraform)
	}
	var typesMap = unknownSRObjectValue.AttributeTypes(ctx)
	var valuesMap = schemaUtils.ValueMapFromTypes(ctx, typesMap)

	if r.ConfluentLike != nil {
		var confluentTypesMap = schema.NewConfluentLikeValueNull().AttributeTypes(ctx)
		var confluentValuesMap = schemaUtils.ValueMapFromTypes(ctx, confluentTypesMap)
		confluentValuesMap["ignore_untrusted_certificate"] = basetypes.NewBoolValue(false) // default value

		security, err := confluentSecurityInternalModelToTerraform(ctx, &r.ConfluentLike.Security)
		if err != nil {
			return schema.SchemaRegistryValue{}, err
		}
		securityValue, diag2 := security.ToObjectValue(ctx)
		if diag2.HasError() {
			return schema.SchemaRegistryValue{}, mapper.WrapDiagError(diag2, "schema_registry.confluent_like.security", mapper.IntoTerraform)
		}

		confluentValuesMap["url"] = schemaUtils.NewStringValue(r.ConfluentLike.Url)
		confluentValuesMap["properties"] = schemaUtils.NewStringValue(r.ConfluentLike.Properties)
		confluentValuesMap["ignore_untrusted_certificate"] = basetypes.NewBoolValue(r.ConfluentLike.IgnoreUntrustedCertificate)
		confluentValuesMap["security"] = securityValue
		valuesMap["confluent_like"], diag = types.ObjectValue(confluentTypesMap, confluentValuesMap)
		if diag.HasError() {
			return schema.SchemaRegistryValue{}, mapper.WrapDiagError(diag, "schema_registry.confluent_like", mapper.IntoTerraform)
		}
	}

	if r.Glue != nil {
		var glueTypesMap = schema.NewGlueValueNull().AttributeTypes(ctx)
		var glueValuesMap = schemaUtils.ValueMapFromTypes(ctx, glueTypesMap)
		security, err := amazonSecurityInternalModelToTerraform(ctx, &r.Glue.Security)
		if err != nil {
			return schema.SchemaRegistryValue{}, err
		}
		securityValue, diag2 := security.ToObjectValue(ctx)
		if diag2.HasError() {
			return schema.SchemaRegistryValue{}, mapper.WrapDiagError(diag2, "schema_registry.glue.security", mapper.IntoTerraform)
		}

		glueValuesMap["registry_name"] = schemaUtils.NewStringValue(r.Glue.RegistryName)
		glueValuesMap["region"] = schemaUtils.NewStringValue(r.Glue.Region)
		glueValuesMap["security"] = securityValue
		valuesMap["glue"], diag = types.ObjectValue(glueTypesMap, glueValuesMap)
		if diag.HasError() {
			return schema.SchemaRegistryValue{}, mapper.WrapDiagError(diag, "schema_registry.glue", mapper.IntoTerraform)
		}
	}

	value, diag := schema.NewSchemaRegistryValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.SchemaRegistryValue{}, mapper.WrapDiagError(diag, "schema_registry", mapper.IntoTerraform)
	}
	return value, nil
}

func confluentSecurityInternalModelToTerraform(ctx context.Context, r *model.ConfluentLikeSchemaRegistrySecurity) (schema.ConfluentSecurityValue, error) {
	unknownSecurityObjectValue, diag := schema.NewConfluentSecurityValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return schema.ConfluentSecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.confluent_like.security", mapper.IntoTerraform)
	}
	var typesMap = unknownSecurityObjectValue.AttributeTypes(ctx)
	var valuesMap = schemaUtils.ValueMapFromTypes(ctx, typesMap)

	if r.BasicAuth != nil {
		var basicTypesMap = schema.NewBasicAuthValueNull().AttributeTypes(ctx)
		var basicValuesMap = schemaUtils.ValueMapFromTypes(ctx, basicTypesMap)
		basicValuesMap["username"] = schemaUtils.NewStringValue(r.BasicAuth.UserName)
		basicValuesMap["password"] = schemaUtils.NewStringValue(r.BasicAuth.Password)
		valuesMap["basic_auth"], diag = types.ObjectValue(basicTypesMap, basicValuesMap)
		if diag.HasError() {
			return schema.ConfluentSecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.confluent_like.security.basic_auth", mapper.IntoTerraform)
		}
	}

	if r.BearerToken != nil {
		var bearerTokenTypesMap = schema.NewBearerTokenValueNull().AttributeTypes(ctx)
		var bearerTokenValuesMap = schemaUtils.ValueMapFromTypes(ctx, bearerTokenTypesMap)
		bearerTokenValuesMap["token"] = schemaUtils.NewStringValue(r.BearerToken.Token)
		valuesMap["bearer_token"], diag = types.ObjectValue(bearerTokenTypesMap, bearerTokenValuesMap)
		if diag.HasError() {
			return schema.ConfluentSecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.confluent_like.security.bearer_token", mapper.IntoTerraform)
		}
	}
	if r.SSLAuth != nil {
		var sslAuthTypesMap = schema.NewSslAuthValueNull().AttributeTypes(ctx)
		var sslAuthValuesMap = schemaUtils.ValueMapFromTypes(ctx, sslAuthTypesMap)
		sslAuthValuesMap["certificate_chain"] = schemaUtils.NewStringValue(r.SSLAuth.CertificateChain)
		sslAuthValuesMap["key"] = schemaUtils.NewStringValue(r.SSLAuth.Key)
		valuesMap["ssl_auth"], diag = types.ObjectValue(sslAuthTypesMap, sslAuthValuesMap)
		if diag.HasError() {
			return schema.ConfluentSecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.confluent_like.security.ssl_auth", mapper.IntoTerraform)
		}
	}

	value, diag := schema.NewConfluentSecurityValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.ConfluentSecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.security", mapper.IntoTerraform)
	}
	return value, nil

}

func amazonSecurityInternalModelToTerraform(ctx context.Context, r *model.AmazonSecurity) (schema.GlueSecurityValue, error) {
	unknownSecurityObjectValue, diag := schema.NewGlueSecurityValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return schema.GlueSecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.security", mapper.IntoTerraform)
	}
	var typesMap = unknownSecurityObjectValue.AttributeTypes(ctx)
	var valuesMap = schemaUtils.ValueMapFromTypes(ctx, typesMap)

	if r.Credentials != nil {
		var credentialsTypesMap = schema.NewCredentialsValueNull().AttributeTypes(ctx)
		var credentialsValuesMap = schemaUtils.ValueMapFromTypes(ctx, credentialsTypesMap)
		credentialsValuesMap["access_key_id"] = schemaUtils.NewStringValue(r.Credentials.AccessKeyId)
		credentialsValuesMap["secret_key"] = schemaUtils.NewStringValue(r.Credentials.SecretKey)
		valuesMap["credentials"], diag = types.ObjectValue(credentialsTypesMap, credentialsValuesMap)
		if diag.HasError() {
			return schema.GlueSecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.glue.security.credentials", mapper.IntoTerraform)
		}
	}
	if r.FromContext != nil {
		var fromContextTypesMap = schema.NewFromContextValueNull().AttributeTypes(ctx)
		var fromContextValuesMap = schemaUtils.ValueMapFromTypes(ctx, fromContextTypesMap)
		fromContextValuesMap["profile"] = schemaUtils.NewStringValue(r.FromContext.Profile)
		valuesMap["from_context"], diag = types.ObjectValue(fromContextTypesMap, fromContextValuesMap)
		if diag.HasError() {
			return schema.GlueSecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.glue.security.from_context", mapper.IntoTerraform)
		}
	}
	if r.FromRole != nil {
		var fromRoleTypesMap = schema.NewFromRoleValueNull().AttributeTypes(ctx)
		var fromRoleValuesMap = schemaUtils.ValueMapFromTypes(ctx, fromRoleTypesMap)
		fromRoleValuesMap["role"] = schemaUtils.NewStringValue(r.FromRole.Role)
		valuesMap["from_role"], diag = types.ObjectValue(fromRoleTypesMap, fromRoleValuesMap)
		if diag.HasError() {
			return schema.GlueSecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.glue.security.from_role", mapper.IntoTerraform)
		}
	}
	if r.IAMAnywhere != nil {
		var iamAnywhereTypesMap = schema.NewIamAnywhereValueNull().AttributeTypes(ctx)
		var iamAnywhereValuesMap = schemaUtils.ValueMapFromTypes(ctx, iamAnywhereTypesMap)
		iamAnywhereValuesMap["trust_anchor_arn"] = schemaUtils.NewStringValue(r.IAMAnywhere.TrustAnchorArn)
		iamAnywhereValuesMap["profile_arn"] = schemaUtils.NewStringValue(r.IAMAnywhere.ProfileArn)
		iamAnywhereValuesMap["role_arn"] = schemaUtils.NewStringValue(r.IAMAnywhere.RoleArn)
		iamAnywhereValuesMap["certificate"] = schemaUtils.NewStringValue(r.IAMAnywhere.Certificate)
		iamAnywhereValuesMap["private_key"] = schemaUtils.NewStringValue(r.IAMAnywhere.PrivateKey)
		valuesMap["iam_anywhere"], diag = types.ObjectValue(iamAnywhereTypesMap, iamAnywhereValuesMap)
		if diag.HasError() {
			return schema.GlueSecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.glue.security.iam_anywhere", mapper.IntoTerraform)
		}
	}

	value, diag := schema.NewGlueSecurityValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.GlueSecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.security", mapper.IntoTerraform)
	}
	return value, nil
}
