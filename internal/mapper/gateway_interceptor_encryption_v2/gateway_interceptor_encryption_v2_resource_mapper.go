package gateway_interceptor_encryption_v2

import (
	"context"
	"fmt"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	gwinterceptor "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_interceptor_encryption_v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *gwinterceptor.GatewayInterceptorEncryptionV2Model) (gateway.GatewayInterceptorEncryptionResource, error) {
	scope := gateway.GatewayInterceptorEncryptionScope{
		Group:    r.Scope.Group.ValueString(),
		VCluster: r.Scope.Vcluster.ValueString(),
		Username: r.Scope.Username.ValueString(),
	}

	config, err := ObjectValueToInterceptorEncryptionConfig(ctx, &r.Spec.Config)
	if err != nil {
		return gateway.GatewayInterceptorEncryptionResource{}, err
	}

	return gateway.NewGatewayInterceptorEncryptionResource(
		gateway.GatewayInterceptorEncryptionMetadata{
			Name:  r.Name.ValueString(),
			Scope: scope,
		},
		gateway.GatewayInterceptorEncryptionSpec{
			Comment:     r.Spec.Comment.ValueString(),
			PluginClass: r.Spec.PluginClass.ValueString(),
			Priority:    r.Spec.Priority.ValueInt64(),
			Config:      config,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorEncryptionResource) (gwinterceptor.GatewayInterceptorEncryptionV2Model, error) {
	specValue, err := specInternalModelToTerraform(ctx, r.Spec)
	if err != nil {
		return gwinterceptor.GatewayInterceptorEncryptionV2Model{}, err
	}

	return gwinterceptor.GatewayInterceptorEncryptionV2Model{
		Name: types.StringValue(r.Metadata.Name),
		Scope: gwinterceptor.ScopeValue{
			Group:    schema.NewStringValue(r.Metadata.Scope.Group),
			Vcluster: schema.NewStringValue(r.Metadata.Scope.VCluster),
			Username: schema.NewStringValue(r.Metadata.Scope.Username),
		},
		Spec: specValue,
	}, nil
}

func specInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorEncryptionSpec) (gwinterceptor.SpecValue, error) {
	unknownSpecObjectValue, diag := gwinterceptor.NewSpecValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return gwinterceptor.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["comment"] = schema.NewStringValue(r.Comment)
	valuesMap["plugin_class"] = schema.NewStringValue(r.PluginClass)
	valuesMap["priority"] = schema.NewInt64Value(r.Priority)

	config, err := configInternalModelToTerraform(ctx, r.Config)
	if err != nil {
		return gwinterceptor.SpecValue{}, err
	}
	valuesMap["config"] = config

	value, diag := gwinterceptor.NewSpecValue(typesMap, valuesMap)
	if diag.HasError() {
		return gwinterceptor.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	return value, nil
}

func configInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorEncryptionConfig) (basetypes.ObjectValue, error) {
	unknownSpecObjectValue, diag := gwinterceptor.NewConfigValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["external_storage"] = schema.NewBoolValue(r.ExternalStorage)
	valuesMap["topic"] = schema.NewStringValue(r.Topic)
	valuesMap["schema_data_mode"] = schema.NewStringValue(r.SchemaDataMode)
	valuesMap["enable_audit_log_on_error"] = schema.NewBoolValue(r.EnableAuditLogOnError)

	schemaRegistryConfig, err := schemaRegistryConfigInternalModelToTerraform(ctx, r.SchemaRegistryConfig)
	if err != nil {
		return basetypes.ObjectValue{}, err
	}
	valuesMap["schema_registry_config"] = schemaRegistryConfig

	// kmsConfig, err := kmsConfigInternalModelToTerraform(ctx, r.KmsConfig)
	// if err != nil {
	// 	return basetypes.ObjectValue{}, err
	// }
	// valuesMap["kms_config"] = kmsConfig

	config, diag := gwinterceptor.NewConfigValue(typesMap, valuesMap)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}
	configValue, diag := config.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}
	return configValue, nil
}

func schemaRegistryConfigInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorEncryptionSchemaRegistryConfig) (basetypes.ObjectValue, error) {
	unknownSpecObjectValue, diag := gwinterceptor.NewSchemaRegistryConfigValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "schema_registry_config", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["host"] = schema.NewStringValue(r.Host)
	valuesMap["cache_size"] = schema.NewInt64Value(r.CacheSize)

	additionalConfigs, diag := schema.StringMapToMapValue(ctx, r.AdditionalConfigs)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "additional_configs", mapper.IntoTerraform)
	}
	valuesMap["additional_configs"] = additionalConfigs

	schemaRegistry, diag := gwinterceptor.NewSchemaRegistryConfigValue(typesMap, valuesMap)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "schema_registry_config", mapper.IntoTerraform)
	}
	schemaRegistryValue, diag := schemaRegistry.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "schema_registry_config", mapper.IntoTerraform)
	}
	return schemaRegistryValue, nil
}

func kmsConfigInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorEncryptionKMSConfig) (basetypes.ObjectValue, error) {
	unknownSpecObjectValue, diag := gwinterceptor.NewKmsConfigValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["key_ttl_ms"] = schema.NewInt64Value(r.KeyTtlMs)
	// valuesMap["topic"] = schema.NewStringValue(r.Topic)
	// valuesMap["schema_data_mode"] = schema.NewStringValue(r.SchemaDataMode)
	// valuesMap["enable_audit_log_on_error"] = schema.NewBoolValue(r.EnableAuditLogOnError)

	// schemaRegistryConfig, err := schemaRegistryConfigInternalModelToTerraform(ctx, *r.Config)
	// if err != nil {
	// 	return gwinterceptor.KmsConfigValue{}, err
	// }
	// valuesMap["schema_registry_config"] = schemaRegistryConfig
	kms, diag := gwinterceptor.NewKmsConfigValue(typesMap, valuesMap)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "kms_config", mapper.IntoTerraform)
	}
	kmsValue, diag := kms.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "kms_config", mapper.IntoTerraform)
	}
	return kmsValue, nil
}

func ObjectValueToInterceptorEncryptionConfig(ctx context.Context, r *basetypes.ObjectValue) (*gateway.GatewayInterceptorEncryptionConfig, error) {
	if r.IsNull() {
		return &gateway.GatewayInterceptorEncryptionConfig{}, nil
	}

	configValue, diag := gwinterceptor.NewConfigValue(r.AttributeTypes(ctx), r.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionConfig{}, mapper.WrapDiagError(diag, "config", mapper.FromTerraform)
	}

	schemaRegistryValue, diag := gwinterceptor.NewSchemaRegistryConfigValue(configValue.SchemaRegistryConfig.AttributeTypes(ctx), configValue.SchemaRegistryConfig.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionConfig{}, mapper.WrapDiagError(diag, "schema_registry_config", mapper.FromTerraform)
	}

	// KMSconfig, err := ConfigValueValueToKMSConfig(ctx, &configValue)
	// if err != nil {
	// 	return &gateway.GatewayInterceptorEncryptionConfig{}, err
	// }

	additionalConfigs, diag := schema.MapValueToStringMap(ctx, schemaRegistryValue.AdditionalConfigs)
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionConfig{}, mapper.WrapDiagError(diag, "config", mapper.FromTerraform)
	}

	return &gateway.GatewayInterceptorEncryptionConfig{
		ExternalStorage:       configValue.ExternalStorage.ValueBool(),
		Topic:                 configValue.Topic.ValueString(),
		EnableAuditLogOnError: configValue.EnableAuditLogOnError.ValueBool(),
		SchemaDataMode:        configValue.SchemaDataMode.ValueString(),
		// KmsConfig:             KMSconfig,
		SchemaRegistryConfig: &gateway.GatewayInterceptorEncryptionSchemaRegistryConfig{
			Host:              schemaRegistryValue.Host.ValueString(),
			CacheSize:         schemaRegistryValue.CacheSize.ValueInt64(),
			AdditionalConfigs: additionalConfigs,
		},
	}, nil
}

func ConfigValueValueToKMSConfig(ctx context.Context, configValue *gwinterceptor.ConfigValue) (*gateway.GatewayInterceptorEncryptionKMSConfig, error) {
	if configValue.IsNull() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, fmt.Errorf("Empty config")
	}

	KMSconfigValue, diag := gwinterceptor.NewKmsConfigValue(configValue.KmsConfig.AttributeTypes(ctx), configValue.KmsConfig.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, mapper.WrapDiagError(diag, "kms_config", mapper.FromTerraform)
	}

	AzureconfigValue, diag := gwinterceptor.NewAzureValue(KMSconfigValue.Azure.AttributeTypes(ctx), KMSconfigValue.Azure.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, mapper.WrapDiagError(diag, "azure", mapper.FromTerraform)
	}

	RetryPolicyValue, diag := gwinterceptor.NewRetryPolicyValue(AzureconfigValue.RetryPolicy.AttributeTypes(ctx), AzureconfigValue.RetryPolicy.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, mapper.WrapDiagError(diag, "retry_policy", mapper.FromTerraform)
	}

	TokenCredentialValue, diag := gwinterceptor.NewTokenCredentialValue(AzureconfigValue.TokenCredential.AttributeTypes(ctx), AzureconfigValue.TokenCredential.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, mapper.WrapDiagError(diag, "token_credential", mapper.FromTerraform)
	}

	UsernamePasswordCredentialValue, diag := gwinterceptor.NewUsernamePasswordCredentialValue(AzureconfigValue.UsernamePasswordCredential.AttributeTypes(ctx), AzureconfigValue.UsernamePasswordCredential.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, mapper.WrapDiagError(diag, "username_password_credential", mapper.FromTerraform)
	}

	return &gateway.GatewayInterceptorEncryptionKMSConfig{
		KeyTtlMs: KMSconfigValue.KeyTtlMs.ValueInt64(),
		Azure: gateway.GatewayInterceptorEncryptionAzureKMSConfig{
			RetryPolicy: &gateway.GatewayInterceptorEncryptionRetryPolicyConfig{
				MaxRetries: RetryPolicyValue.MaxRetries.ValueInt64(),
				DelayMs:    RetryPolicyValue.DelayMs.ValueInt64(),
				MaxDelayMs: RetryPolicyValue.MaxDelayMs.ValueInt64(),
			},
			TokenCredential: &gateway.GatewayInterceptorEncryptionAzureTokenCredential{
				ClientId:     TokenCredentialValue.ClientId.ValueString(),
				TenantId:     TokenCredentialValue.TenantId.ValueString(),
				ClientSecret: TokenCredentialValue.ClientSecret.ValueString(),
			},
			UsernamePasswordCredential: &gateway.GatewayInterceptorEncryptionAzureUsernamePasswordCredential{
				ClientId: UsernamePasswordCredentialValue.ClientId.ValueString(),
				TenantId: UsernamePasswordCredentialValue.TenantId.ValueString(),
				Username: UsernamePasswordCredentialValue.Username.ValueString(),
				Password: UsernamePasswordCredentialValue.Password.ValueString(),
			},
		},
		AWS: gateway.GatewayInterceptorEncryptionAWSKMSConfig{},
		GCP: gateway.GatewayInterceptorEncryptionGCPKMSConfig{},
	}, nil
}
