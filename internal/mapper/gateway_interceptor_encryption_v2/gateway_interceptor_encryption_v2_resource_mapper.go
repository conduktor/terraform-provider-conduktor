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

	KMSconfig, err := ConfigValueValueToKMSConfig(ctx, &configValue)
	if err != nil {
		return &gateway.GatewayInterceptorEncryptionConfig{}, err
	}

	additionalConfigs, diag := schema.MapValueToStringMap(ctx, schemaRegistryValue.AdditionalConfigs)
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionConfig{}, mapper.WrapDiagError(diag, "config", mapper.FromTerraform)
	}

	return &gateway.GatewayInterceptorEncryptionConfig{
		ExternalStorage:       configValue.ExternalStorage.ValueBool(),
		Topic:                 configValue.Topic.ValueString(),
		EnableAuditLogOnError: configValue.EnableAuditLogOnError.ValueBool(),
		SchemaDataMode:        configValue.SchemaDataMode.ValueString(),
		KmsConfig:             KMSconfig,
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

	KMSConfigValue, diag := gwinterceptor.NewKmsConfigValue(configValue.KmsConfig.AttributeTypes(ctx), configValue.KmsConfig.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, mapper.WrapDiagError(diag, "kms_config", mapper.FromTerraform)
	}

	AzureConfigValue, diag := gwinterceptor.NewAzureValue(KMSConfigValue.Azure.AttributeTypes(ctx), KMSConfigValue.Azure.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, mapper.WrapDiagError(diag, "azure", mapper.FromTerraform)
	}

	RetryPolicyValue, diag := gwinterceptor.NewRetryPolicyValue(AzureConfigValue.RetryPolicy.AttributeTypes(ctx), AzureConfigValue.RetryPolicy.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, mapper.WrapDiagError(diag, "retry_policy", mapper.FromTerraform)
	}

	TokenCredentialValue, diag := gwinterceptor.NewTokenCredentialValue(AzureConfigValue.TokenCredential.AttributeTypes(ctx), AzureConfigValue.TokenCredential.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, mapper.WrapDiagError(diag, "token_credential", mapper.FromTerraform)
	}

	UsernamePasswordCredentialValue, diag := gwinterceptor.NewUsernamePasswordCredentialValue(AzureConfigValue.UsernamePasswordCredential.AttributeTypes(ctx), AzureConfigValue.UsernamePasswordCredential.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, mapper.WrapDiagError(diag, "username_password_credential", mapper.FromTerraform)
	}

	AWSConfigValue, diag := gwinterceptor.NewAwsValue(KMSConfigValue.Aws.AttributeTypes(ctx), KMSConfigValue.Aws.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, mapper.WrapDiagError(diag, "aws", mapper.FromTerraform)
	}
	BasicCredentialsConfigValue, diag := gwinterceptor.NewBasicCredentialsValue(AWSConfigValue.BasicCredentials.AttributeTypes(ctx), AWSConfigValue.BasicCredentials.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, mapper.WrapDiagError(diag, "basic_credentials", mapper.FromTerraform)
	}
	BasicSessionConfigValue, diag := gwinterceptor.NewSessionCredentialsValue(AWSConfigValue.SessionCredentials.AttributeTypes(ctx), AWSConfigValue.SessionCredentials.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, mapper.WrapDiagError(diag, "aws", mapper.FromTerraform)
	}

	GCPConfigValue, diag := gwinterceptor.NewGcpValue(KMSConfigValue.Gcp.AttributeTypes(ctx), KMSConfigValue.Gcp.Attributes())
	if diag.HasError() {
		return &gateway.GatewayInterceptorEncryptionKMSConfig{}, mapper.WrapDiagError(diag, "gcp", mapper.FromTerraform)
	}

	return &gateway.GatewayInterceptorEncryptionKMSConfig{
		KeyTtlMs: KMSConfigValue.KeyTtlMs.ValueInt64(),
		Azure: &gateway.GatewayInterceptorEncryptionAzureKMSConfig{
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
		AWS: &gateway.GatewayInterceptorEncryptionAWSKMSConfig{
			BasicCredentials: &gateway.GatewayInterceptorEncryptionBasicAWSCredentialsConfig{
				AccessKey: BasicCredentialsConfigValue.AccessKey.ValueString(),
				SecretKey: BasicCredentialsConfigValue.SecretKey.ValueString(),
			},
			SessionCredentials: &gateway.GatewayInterceptorEncryptionBasicSessionCredentialsConfig{
				AccessKey:    BasicSessionConfigValue.AccessKey.ValueString(),
				SecretKey:    BasicSessionConfigValue.SecretKey.ValueString(),
				SessionToken: BasicSessionConfigValue.SessionToken.ValueString(),
			},
		},
		GCP: &gateway.GatewayInterceptorEncryptionGCPKMSConfig{
			ServiceAccountCredentialsFilePath: GCPConfigValue.ServiceAccountCredentialsFilePath.ValueString(),
		},
	}, nil
}

//
//
//
//
// Internal to TF
//
//
//
//

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

	valuesMap["topic"] = schema.NewStringValue(r.Topic)
	valuesMap["external_storage"] = schema.NewBoolValue(r.ExternalStorage)
	valuesMap["schema_data_mode"] = schema.NewStringValue(r.SchemaDataMode)
	valuesMap["enable_audit_log_on_error"] = schema.NewBoolValue(r.EnableAuditLogOnError)

	kmsConfig, err := kmsConfigInternalModelToTerraform(ctx, r.KmsConfig)
	if err != nil {
		return basetypes.ObjectValue{}, err
	}
	valuesMap["kms_config"] = kmsConfig

	schemaRegistryConfig, err := schemaRegistryConfigInternalModelToTerraform(ctx, r.SchemaRegistryConfig)
	if err != nil {
		return basetypes.ObjectValue{}, err
	}
	valuesMap["schema_registry_config"] = schemaRegistryConfig

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

	azure, err := azureInternalModelToTerraform(ctx, r.Azure)
	if err != nil {
		return basetypes.ObjectValue{}, err
	}
	valuesMap["azure"] = azure
	aws, err := awsInternalModelToTerraform(ctx, r.AWS)
	if err != nil {
		return basetypes.ObjectValue{}, err
	}
	valuesMap["aws"] = aws
	gcp, err := gcpInternalModelToTerraform(ctx, r.GCP)
	if err != nil {
		return basetypes.ObjectValue{}, err
	}
	valuesMap["gcp"] = gcp

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

func azureInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorEncryptionAzureKMSConfig) (basetypes.ObjectValue, error) {
	unknownSpecObjectValue, diag := gwinterceptor.NewAzureValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "azure", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	retryPolicy, err := retryPolicyConfigInternalModelToTerraform(ctx, r.RetryPolicy)
	if err != nil {
		return basetypes.ObjectValue{}, err
	}
	valuesMap["retry_policy_config"] = retryPolicy
	tokenCredentials, err := azureTokenCredentialInternalModelToTerraform(ctx, r.TokenCredential)
	if err != nil {
		return basetypes.ObjectValue{}, err
	}
	valuesMap["azure_token_credential"] = tokenCredentials
	usernamePassword, err := azureUsernamePasswordCredentialInternalModelToTerraform(ctx, r.UsernamePasswordCredential)
	if err != nil {
		return basetypes.ObjectValue{}, err
	}
	valuesMap["azure_username_password_credential"] = usernamePassword

	config, diag := gwinterceptor.NewAzureValue(typesMap, valuesMap)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "azure", mapper.IntoTerraform)
	}
	value, diag := config.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "azure", mapper.IntoTerraform)
	}
	return value, nil
}
func retryPolicyConfigInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorEncryptionRetryPolicyConfig) (basetypes.ObjectValue, error) {
	unknownSpecObjectValue, diag := gwinterceptor.NewRetryPolicyValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "retry_policy_config", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["max_retries"] = schema.NewInt64Value(r.MaxRetries)
	valuesMap["delay_ms"] = schema.NewInt64Value(r.DelayMs)
	valuesMap["max_delay_ms"] = schema.NewInt64Value(r.MaxDelayMs)

	config, diag := gwinterceptor.NewKmsConfigValue(typesMap, valuesMap)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "retry_policy_config", mapper.IntoTerraform)
	}
	value, diag := config.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "retry_policy_config", mapper.IntoTerraform)
	}
	return value, nil
}
func azureTokenCredentialInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorEncryptionAzureTokenCredential) (basetypes.ObjectValue, error) {
	unknownSpecObjectValue, diag := gwinterceptor.NewKmsConfigValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "azure_token_credential", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["client_id"] = schema.NewStringValue(r.ClientId)
	valuesMap["tenant_id"] = schema.NewStringValue(r.TenantId)
	valuesMap["client_secret"] = schema.NewStringValue(r.ClientSecret)

	config, diag := gwinterceptor.NewKmsConfigValue(typesMap, valuesMap)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "azure_token_credential", mapper.IntoTerraform)
	}
	value, diag := config.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "azure_token_credential", mapper.IntoTerraform)
	}
	return value, nil
}
func azureUsernamePasswordCredentialInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorEncryptionAzureUsernamePasswordCredential) (basetypes.ObjectValue, error) {
	unknownSpecObjectValue, diag := gwinterceptor.NewKmsConfigValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "azure_username_password_credential", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["client_id"] = schema.NewStringValue(r.ClientId)
	valuesMap["tenant_id"] = schema.NewStringValue(r.TenantId)
	valuesMap["username"] = schema.NewStringValue(r.Username)
	valuesMap["password"] = schema.NewStringValue(r.Password)

	config, diag := gwinterceptor.NewKmsConfigValue(typesMap, valuesMap)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "azure_username_password_credential", mapper.IntoTerraform)
	}
	value, diag := config.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "azure_username_password_credential", mapper.IntoTerraform)
	}
	return value, nil
}

func awsInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorEncryptionAWSKMSConfig) (basetypes.ObjectValue, error) {
	unknownSpecObjectValue, diag := gwinterceptor.NewAwsValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "aws", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	basic, err := basicAWSCredentialsConfigInternalModelToTerraform(ctx, r.BasicCredentials)
	if err != nil {
		return basetypes.ObjectValue{}, err
	}
	valuesMap["basic_credentials"] = basic
	basicSession, err := basicSessionCredentialsConfigInternalModelToTerraform(ctx, r.SessionCredentials)
	if err != nil {
		return basetypes.ObjectValue{}, err
	}
	valuesMap["session_credentials"] = basicSession

	config, diag := gwinterceptor.NewKmsConfigValue(typesMap, valuesMap)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "aws", mapper.IntoTerraform)
	}
	value, diag := config.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "kms_config", mapper.IntoTerraform)
	}
	return value, nil
}
func basicAWSCredentialsConfigInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorEncryptionBasicAWSCredentialsConfig) (basetypes.ObjectValue, error) {
	unknownSpecObjectValue, diag := gwinterceptor.NewBasicCredentialsValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "basic_credentials", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["access_key"] = schema.NewStringValue(r.AccessKey)
	valuesMap["secret_key"] = schema.NewStringValue(r.SecretKey)

	config, diag := gwinterceptor.NewKmsConfigValue(typesMap, valuesMap)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "basic_credentials", mapper.IntoTerraform)
	}
	value, diag := config.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "basic_credentials", mapper.IntoTerraform)
	}
	return value, nil
}
func basicSessionCredentialsConfigInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorEncryptionBasicSessionCredentialsConfig) (basetypes.ObjectValue, error) {
	unknownSpecObjectValue, diag := gwinterceptor.NewSessionCredentialsValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "session_credentials", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["access_key"] = schema.NewStringValue(r.AccessKey)
	valuesMap["secret_key"] = schema.NewStringValue(r.SecretKey)
	valuesMap["session_token"] = schema.NewStringValue(r.SessionToken)

	config, diag := gwinterceptor.NewKmsConfigValue(typesMap, valuesMap)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "session_credentials", mapper.IntoTerraform)
	}
	value, diag := config.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "session_credentials", mapper.IntoTerraform)
	}
	return value, nil
}

func gcpInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorEncryptionGCPKMSConfig) (basetypes.ObjectValue, error) {
	unknownSpecObjectValue, diag := gwinterceptor.NewGcpValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "gcp", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["service_account_credentials_file_path"] = schema.NewStringValue(r.ServiceAccountCredentialsFilePath)

	config, diag := gwinterceptor.NewKmsConfigValue(typesMap, valuesMap)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "gcp", mapper.IntoTerraform)
	}
	value, diag := config.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "gcp", mapper.IntoTerraform)
	}
	return value, nil
}
