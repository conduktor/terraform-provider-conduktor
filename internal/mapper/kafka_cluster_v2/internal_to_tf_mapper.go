package kafka_cluster_v2

import (
	"context"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_kafka_cluster_v2"
	"github.com/conduktor/terraform-provider-conduktor/internal/schema/validation"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func InternalModelToTerraform(ctx context.Context, r *model.KafkaClusterResource) (schema.KafkaClusterV2Model, error) {

	labels, diag := schemaUtils.StringMapToMapValue(ctx, r.Metadata.Labels)
	if diag.HasError() {
		return schema.KafkaClusterV2Model{}, mapper.WrapDiagError(diag, "labels", mapper.IntoTerraform)
	}

	specValue, err := specInternalModelToTerraform(ctx, &r.Spec)
	if err != nil {
		return schema.KafkaClusterV2Model{}, err
	}

	return schema.KafkaClusterV2Model{
		Name:   types.StringValue(r.Metadata.Name),
		Labels: labels,
		Spec:   specValue,
	}, nil
}

func specInternalModelToTerraform(ctx context.Context, r *model.KafkaClusterSpec) (schema.SpecValue, error) {

	properties, diag := schemaUtils.StringMapToMapValue(ctx, r.Properties)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "properties", mapper.IntoTerraform)
	}

	kafkaFlavor, err := kafkaFlavorInternalModelToTerraform(ctx, r.KafkaFlavor)
	if err != nil {
		return schema.SpecValue{}, err
	}
	kafkaFlavorValue, diag := kafkaFlavor.ToObjectValue(ctx)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "kafka_flavor", mapper.IntoTerraform)
	}

	schemaRegistry, err := schemaRegistryInternalModelToTerraform(ctx, r.SchemaRegistry)
	if err != nil {
		return schema.SpecValue{}, err
	}
	schemaRegistryValue, diag := schemaRegistry.ToObjectValue(ctx)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "schema_registry", mapper.IntoTerraform)
	}

	return schema.SpecValue{
		BootstrapServers:           schemaUtils.NewStringValue(r.BootstrapServers),
		DisplayName:                schemaUtils.NewStringValue(r.DisplayName),
		Color:                      schemaUtils.NewStringValue(r.Color),
		Icon:                       schemaUtils.NewStringValue(r.Icon),
		IgnoreUntrustedCertificate: basetypes.NewBoolValue(r.IgnoreUntrustedCertificate),
		Properties:                 properties,
		KafkaFlavor:                kafkaFlavorValue,
		SchemaRegistry:             schemaRegistryValue,
	}, nil
}

func kafkaFlavorInternalModelToTerraform(ctx context.Context, r *model.KafkaFlavor) (schema.KafkaFlavorValue, error) {
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
		valuesMap["type"] = schemaUtils.NewStringValue(validation.AivenKafkaFlavor)
		valuesMap["api_token"] = schemaUtils.NewStringValue(r.Aiven.ApiToken)
		valuesMap["project"] = schemaUtils.NewStringValue(r.Aiven.Project)
		valuesMap["service_name"] = schemaUtils.NewStringValue(r.Aiven.ServiceName)
	}
	if r.Confluent != nil {
		valuesMap["type"] = schemaUtils.NewStringValue(validation.ConfluentKafkaFlavor)
		valuesMap["key"] = schemaUtils.NewStringValue(r.Confluent.Key)
		valuesMap["secret"] = schemaUtils.NewStringValue(r.Confluent.Secret)
		valuesMap["confluent_environment_id"] = schemaUtils.NewStringValue(r.Confluent.ConfluentEnvironmentId)
		valuesMap["confluent_cluster_id"] = schemaUtils.NewStringValue(r.Confluent.ConfluentClusterId)
	}
	if r.Gateway != nil {
		valuesMap["type"] = schemaUtils.NewStringValue(validation.GatewayKafkaFlavor)
		valuesMap["url"] = schemaUtils.NewStringValue(r.Gateway.Url)
		valuesMap["user"] = schemaUtils.NewStringValue(r.Gateway.User)
		valuesMap["password"] = schemaUtils.NewStringValue(r.Gateway.Password)
		valuesMap["virtual_cluster"] = schemaUtils.NewStringValue(r.Gateway.VirtualCluster)
		valuesMap["ignore_untrusted_certificate"] = basetypes.NewBoolValue(r.Gateway.IgnoreUntrustedCertificate)
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
		security, err := confluentSecurityInternalModelToTerraform(ctx, &r.ConfluentLike.Security)
		if err != nil {
			return schema.SchemaRegistryValue{}, err
		}
		securityValue, diag2 := security.ToObjectValue(ctx)
		if diag2.HasError() {
			return schema.SchemaRegistryValue{}, mapper.WrapDiagError(diag2, "schema_registry", mapper.IntoTerraform)
		}

		valuesMap["type"] = schemaUtils.NewStringValue(validation.ConfluentLikeSchemaRegistry)
		valuesMap["url"] = schemaUtils.NewStringValue(r.ConfluentLike.Url)
		valuesMap["properties"] = schemaUtils.NewStringValue(r.ConfluentLike.Properties)
		valuesMap["ignore_untrusted_certificate"] = basetypes.NewBoolValue(r.ConfluentLike.IgnoreUntrustedCertificate)
		valuesMap["security"] = securityValue
	}
	if r.Glue != nil {
		security, err := ammazonSecurityInternalModelToTerraform(ctx, &r.Glue.Security)
		if err != nil {
			return schema.SchemaRegistryValue{}, err
		}
		securityValue, diag2 := security.ToObjectValue(ctx)
		if diag2.HasError() {
			return schema.SchemaRegistryValue{}, mapper.WrapDiagError(diag2, "schema_registry", mapper.IntoTerraform)
		}

		valuesMap["type"] = schemaUtils.NewStringValue(validation.GlueSchemaRegistry)
		valuesMap["registry_name"] = schemaUtils.NewStringValue(r.Glue.RegistryName)
		valuesMap["region"] = schemaUtils.NewStringValue(r.Glue.Region)
		valuesMap["security"] = securityValue
	}
	typesMap["security"] = valuesMap["security"].Type(ctx)

	value, diag := schema.NewSchemaRegistryValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.SchemaRegistryValue{}, mapper.WrapDiagError(diag, "schema_registry", mapper.IntoTerraform)
	}
	return value, nil
}

func confluentSecurityInternalModelToTerraform(ctx context.Context, r *model.ConfluentLikeSchemaRegistrySecurity) (schema.SecurityValue, error) {
	unknownSecurityObjectValue, diag := schema.NewSecurityValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return schema.SecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.security", mapper.IntoTerraform)
	}
	var typesMap = unknownSecurityObjectValue.AttributeTypes(ctx)
	var valuesMap = schemaUtils.ValueMapFromTypes(ctx, typesMap)

	if r.NoSecurity != nil {
		valuesMap["type"] = schemaUtils.NewStringValue(validation.NoSecuritySchemaRegistrySecurity)
	}
	if r.BasicAuth != nil {
		valuesMap["type"] = schemaUtils.NewStringValue(validation.BasicAuthSchemaRegistrySecurity)
		valuesMap["username"] = schemaUtils.NewStringValue(r.BasicAuth.UserName)
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
		return schema.SecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.security", mapper.IntoTerraform)
	}
	return value, nil

}

func ammazonSecurityInternalModelToTerraform(ctx context.Context, r *model.AmazonSecurity) (schema.SecurityValue, error) {
	unknownSecurityObjectValue, diag := schema.NewSecurityValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return schema.SecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.security", mapper.IntoTerraform)
	}
	var typesMap = unknownSecurityObjectValue.AttributeTypes(ctx)
	var valuesMap = schemaUtils.ValueMapFromTypes(ctx, typesMap)

	if r.Credentials != nil {
		valuesMap["type"] = schemaUtils.NewStringValue(validation.CredentialsSchemaRegistrySecurity)
		valuesMap["access_key_id"] = schemaUtils.NewStringValue(r.Credentials.AccessKeyId)
		valuesMap["secret_key"] = schemaUtils.NewStringValue(r.Credentials.SecretKey)
	}
	if r.FromContext != nil {
		valuesMap["type"] = schemaUtils.NewStringValue(validation.FromContextSchemaRegistrySecurity)
		valuesMap["profile"] = schemaUtils.NewStringValue(r.FromContext.Profile)
	}
	if r.FromRole != nil {
		valuesMap["type"] = schemaUtils.NewStringValue(validation.FromRoleSchemaRegistrySecurity)
		valuesMap["role"] = schemaUtils.NewStringValue(r.FromRole.Role)
	}
	if r.IAMAnywhere != nil {
		valuesMap["type"] = schemaUtils.NewStringValue(validation.IAMAnywhereSchemaRegistrySecurity)
		valuesMap["trust_anchor_arn"] = schemaUtils.NewStringValue(r.IAMAnywhere.TrustAnchorArn)
		valuesMap["profile_arn"] = schemaUtils.NewStringValue(r.IAMAnywhere.ProfileArn)
		valuesMap["role_arn"] = schemaUtils.NewStringValue(r.IAMAnywhere.RoleArn)
		valuesMap["certificate"] = schemaUtils.NewStringValue(r.IAMAnywhere.Certificate)
		valuesMap["private_key"] = schemaUtils.NewStringValue(r.IAMAnywhere.PrivateKey)
	}

	value, diag := schema.NewSecurityValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.SecurityValue{}, mapper.WrapDiagError(diag, "schema_registry.security", mapper.IntoTerraform)
	}
	return value, nil
}
