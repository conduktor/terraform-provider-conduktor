package schema

import (
	"context"
	"os"
	"strconv"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	groups "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_group_v2"
	users "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_user_v2"
	gwinterceptor "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_interceptor_encryption_v2"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Enum used to perform different actions based on resource type.
type Resource int

const (
	GROUPS Resource = iota // 0
	USERS                  // 1
)

// GetStringConfig Provider string configuration extracted from the schema and environment variables.
// Priority order: configValue > envs.
func GetStringConfig(configValue basetypes.StringValue, envs []string) string {
	if !configValue.IsNull() {
		return configValue.ValueString()
	}
	for _, env := range envs {
		if value := os.Getenv(env); value != "" {
			return value
		}
	}
	return ""
}

// GetBooleanConfig Provider bool configuration extracted from the schema and environment variables.
// Priority order: configValue > envs > fallback.
func GetBooleanConfig(configValue basetypes.BoolValue, envs []string, fallback bool) bool {
	if !configValue.IsNull() {
		return configValue.ValueBool()
	}
	for _, env := range envs {
		if value := os.Getenv(env); value != "" {
			if b, err := strconv.ParseBool(value); err == nil {
				return b
			}
		}
	}

	return fallback
}

// NewBoolValue Convert a string to a basetypes.BoolValue, with default to false.
func NewBoolValue(b bool) basetypes.BoolValue {
	if &b == nil {
		return basetypes.NewBoolValue(false)
	}
	return basetypes.NewBoolValue(b)
}

// NewStringValue Convert a string to a basetypes.StringValue.
func NewStringValue(s string) basetypes.StringValue {
	if s == "" {
		return basetypes.NewStringNull()
	}
	return basetypes.NewStringValue(s)
}

// NewInt64Value Convert a int64 to a basetypes.Int64Value.
func NewInt64Value(i int64) basetypes.Int64Value {
	if i == 0 {
		return basetypes.NewInt64Null()
	}
	return basetypes.NewInt64Value(i)
}

// ListValueToStringArray Convert a ListValue to a string array.
func ListValueToStringArray(ctx context.Context, list basetypes.ListValue) ([]string, diag.Diagnostics) {
	var result []string
	diagnostic := list.ElementsAs(ctx, &result, true)
	return result, diagnostic
}

// StringArrayToListValue Convert a string array to a ListValue.
func StringArrayToListValue(array []string) (basetypes.ListValue, diag.Diagnostics) {
	var values []attr.Value
	for _, f := range array {
		values = append(values, types.StringValue(f))
	}
	return types.ListValue(types.StringType, values)
}

// Helper to parse array of strings to SetValue.
func SetValueToStringArray(ctx context.Context, set basetypes.SetValue) ([]string, diag.Diagnostics) {
	var flags []string
	diag := set.ElementsAs(ctx, &flags, false)
	if diag.HasError() {
		return nil, diag
	}
	return flags, nil
}

// Helper to parse array of strings to SetValue.
func StringArrayToSetValue(arr []string) (basetypes.SetValue, diag.Diagnostics) {
	var flags []attr.Value
	for _, f := range arr {
		flags = append(flags, types.StringValue(f))
	}

	flagsList, diag := types.SetValue(types.StringType, flags)
	if diag.HasError() {
		return basetypes.SetValue{}, diag
	}

	return flagsList, nil
}

// Parse a Permissions Array into a Set based on resource type.
func PermissionArrayToSetValue(ctx context.Context, resource Resource, arr []model.Permission) (basetypes.SetValue, error) {
	var permissionsList basetypes.SetValue
	var tfPermissions []attr.Value
	var diag diag.Diagnostics

	for _, p := range arr {
		flagsList, diag := StringArrayToSetValue(p.Permissions)
		if diag.HasError() {
			return basetypes.SetValue{}, mapper.WrapDiagError(diag, "permissions.permissions", mapper.FromTerraform)
		}

		types := map[string]attr.Type{
			"name":          basetypes.StringType{},
			"resource_type": basetypes.StringType{},
			"permissions":   flagsList.Type(ctx),
			"pattern_type":  basetypes.StringType{},
			"cluster":       basetypes.StringType{},
			"kafka_connect": basetypes.StringType{},
		}
		values := map[string]attr.Value{
			"name":          NewStringValue(p.Name),
			"resource_type": NewStringValue(p.ResourceType),
			"permissions":   flagsList,
			"pattern_type":  NewStringValue(p.PatternType),
			"cluster":       NewStringValue(p.Cluster),
			"kafka_connect": NewStringValue(p.KafkaConnect),
		}

		if resource == GROUPS {
			permObj, diag := groups.NewPermissionsValue(types, values)
			if diag.HasError() {
				return basetypes.SetValue{}, mapper.WrapDiagError(diag, "permissions", mapper.FromTerraform)
			}
			tfPermissions = append(tfPermissions, permObj)
		} else if resource == USERS {
			permObj, diag := users.NewPermissionsValue(types, values)
			if diag.HasError() {
				return basetypes.SetValue{}, mapper.WrapDiagError(diag, "permissions", mapper.FromTerraform)
			}
			tfPermissions = append(tfPermissions, permObj)
		}

	}

	if resource == GROUPS {
		permissionsList, diag = types.SetValue(groups.PermissionsValue{}.Type(ctx), tfPermissions)
	} else if resource == USERS {
		permissionsList, diag = types.SetValue(users.PermissionsValue{}.Type(ctx), tfPermissions)
	}

	if diag.HasError() {
		return basetypes.SetValue{}, mapper.WrapDiagError(diag, "permissions", mapper.FromTerraform)
	}

	return permissionsList, nil
}

// Parse a Set into an array of Permissions based on resource type.
func SetValueToPermissionArray(ctx context.Context, resource Resource, set basetypes.SetValue) ([]model.Permission, error) {
	permissions := make([]model.Permission, 0)
	var diag diag.Diagnostics

	// Ideally the switch within groups and users would have less replication.
	// This might be worth a re-work in the future.
	// NOTE: an idea would be to use ObjectValue instead of user/group PermissionsValue.
	if !set.IsNull() && !set.IsUnknown() {
		// Case for groups
		if resource == GROUPS {
			var tfPermissions []groups.PermissionsValue
			diag = set.ElementsAs(ctx, &tfPermissions, false)
			if diag.HasError() {
				return nil, mapper.WrapDiagError(diag, "permissions", mapper.FromTerraform)
			}
			for _, p := range tfPermissions {
				flags, diag := SetValueToStringArray(ctx, p.Permissions)
				if diag.HasError() {
					return nil, mapper.WrapDiagError(diag, "permissions.permissions", mapper.FromTerraform)
				}

				permissions = append(permissions, model.Permission{
					Name:         p.Name.ValueString(),
					ResourceType: p.ResourceType.ValueString(),
					Permissions:  flags,
					PatternType:  p.PatternType.ValueString(),
					Cluster:      p.Cluster.ValueString(),
					KafkaConnect: p.KafkaConnect.ValueString(),
				})
			}

			// Case for users
		} else if resource == USERS {
			var tfPermissions []users.PermissionsValue
			diag = set.ElementsAs(ctx, &tfPermissions, false)
			if diag.HasError() {
				return nil, mapper.WrapDiagError(diag, "permissions", mapper.FromTerraform)
			}

			for _, p := range tfPermissions {
				flags, diag := SetValueToStringArray(ctx, p.Permissions)
				if diag.HasError() {
					return nil, mapper.WrapDiagError(diag, "permissions.permissions", mapper.FromTerraform)
				}

				permissions = append(permissions, model.Permission{
					Name:         p.Name.ValueString(),
					ResourceType: p.ResourceType.ValueString(),
					Permissions:  flags,
					PatternType:  p.PatternType.ValueString(),
					Cluster:      p.Cluster.ValueString(),
					KafkaConnect: p.KafkaConnect.ValueString(),
				})
			}
		}
	}
	return permissions, nil
}

// MapValueToStringMap Convert a MapValue to a map[string]string.
func MapValueToStringMap(ctx context.Context, mapValue basetypes.MapValue) (map[string]string, diag.Diagnostics) {
	result := make(map[string]string)
	diagnostic := mapValue.ElementsAs(ctx, &result, true)
	return result, diagnostic
}

// StringMapToMapValue Convert a MapValue to a map[string]string.
func StringMapToMapValue(_ context.Context, in map[string]string) (basetypes.MapValue, diag.Diagnostics) {
	if in == nil {
		return types.MapNull(types.StringType), nil
	}
	var values = make(map[string]attr.Value)
	for k, v := range in {
		values[k] = NewStringValue(v)
	}
	return types.MapValue(types.StringType, values)
}

func ValueMapFromTypes(ctx context.Context, types map[string]attr.Type) map[string]attr.Value {
	result := make(map[string]attr.Value)
	for k, v := range types {
		result[k] = v.ValueType(ctx)
	}
	return result
}

// MapValueToStringMap Convert a MapValue to a map[string]string.
func StringToNormalizedJson(ctx context.Context, input string) (jsontypes.Normalized, diag.Diagnostics) {
	return jsontypes.NewNormalizedValue(input), nil
}

func JsonToNormalizedString(ctx context.Context, input basetypes.StringValue) (string, diag.Diagnostics) {
	// jsontypes.NewNormalizedUnknown()
	return jsontypes.NewExactValue(input.String()).String(), nil
	// return "", nil
}

func InterceptorSchemaRegistryConfigToObjectValue(ctx context.Context, config gateway.GatewayInterceptorEncryptionSchemaRegistryConfig) (basetypes.ObjectValue, error) {

	// unknownSpecObjectValue, diag := gwinterceptor.NewSpecValueUnknown().ToObjectValue(ctx)
	// if diag.HasError() {
	// 	return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "schema_registry_config", mapper.IntoTerraform)
	// }
	// var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	// var valuesMap = ValueMapFromTypes(ctx, typesMap)
	//
	// valuesMap["host"] = NewStringValue(config.Host)
	// valuesMap["cache_size"] = NewInt64Value(config.CacheSize)
	//
	// additionalConfigs, diag := StringMapToMapValue(ctx, *config.AdditionalConfigs)
	// if diag.HasError() {
	// 	return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "additional_configs", mapper.IntoTerraform)
	// }
	// // valuesMap["additional_configs"] = additionalConfigs
	//
	// configValue, diag := gwinterceptor.NewSchemaRegistryConfigValue(
	// 	map[string]attr.Type{
	// 		"host":               basetypes.StringType{},
	// 		"cache_size":         basetypes.Int64Type{},
	// 		"additional_configs": additionalConfigs.Type(ctx),
	// 	},
	// 	map[string]attr.Value{
	// 		"host":               NewStringValue(config.Host),
	// 		"cache_size":         NewInt64Value(config.CacheSize),
	// 		"additional_configs": additionalConfigs,
	// 	},
	// )
	// if diag.HasError() {
	// 	return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	// }

	configValue := gwinterceptor.SchemaRegistryConfigValue{
		Host:              NewStringValue(config.Host),
		CacheSize:         NewInt64Value(config.CacheSize),
		AdditionalConfigs: basetypes.NewMapNull(basetypes.StringType{}),
	}

	objectValue, diag := configValue.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}

	return objectValue, nil
}

func InterceptorAzureConfigToObjectValue(ctx context.Context, config gateway.GatewayInterceptorEncryptionAzureKMSConfig) (basetypes.ObjectValue, error) {

	// configValue, diag := gwinterceptor.NewSchemaRegistryConfigValue(
	// 	map[string]attr.Type{
	// 		"host":       basetypes.StringType{},
	// 		"cache_size": basetypes.Int64Type{},
	// 	},
	// 	map[string]attr.Value{
	// 		"host":       NewStringValue(config.Host),
	// 		"cache_size": NewInt64Value(config.CacheSize),
	// 	},
	// )
	// if diag.HasError() {
	// 	return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	// }

	configValue := gwinterceptor.AzureValue{
		RetryPolicy:                basetypes.ObjectValue{},
		TokenCredential:            basetypes.ObjectValue{},
		UsernamePasswordCredential: basetypes.ObjectValue{},
	}

	objectValue, diag := configValue.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}

	return objectValue, nil
}

func InterceptorKMSConfigToObjectValue(ctx context.Context, config gateway.GatewayInterceptorEncryptionKMSConfig) (basetypes.ObjectValue, error) {

	configValue, diag := gwinterceptor.NewSchemaRegistryConfigValue(
		map[string]attr.Type{
			"key_ttl_ms": basetypes.Int64Type{},
			// "cache_size": basetypes.Int64Type{},
		},
		map[string]attr.Value{
			"key_ttl_ms": NewInt64Value(config.KeyTtlMs),
			// "cache_size": NewInt64Value(config.CacheSize),
			// KeyTtlMs: NewInt64Value(config.KeyTtlMs),
			// Azure:    basetypes.ObjectValue{},
			// Aws:      basetypes.ObjectValue{},
			// Gcp:      basetypes.ObjectValue{},
		},
	)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}

	// configValue := gwinterceptor.KmsConfigValue{
	// 	KeyTtlMs: NewInt64Value(config.KeyTtlMs),
	// 	Azure:    basetypes.ObjectValue{},
	// 	Aws:      basetypes.ObjectValue{},
	// 	Gcp:      basetypes.ObjectValue{},
	// }

	objectValue, diag := configValue.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}

	return objectValue, nil
}

// Helper function used to parse Interceptor Config block into an Object Value.
func InterceptorConfigToObjectValue(ctx context.Context, config gateway.GatewayInterceptorEncryptionConfig) (basetypes.ObjectValue, error) {
	schemaRegistryConfig, err := InterceptorSchemaRegistryConfigToObjectValue(ctx, *config.SchemaRegistryConfig)
	if err != nil {
		return basetypes.ObjectValue{}, err
	}
	var KMSConfig basetypes.ObjectValue

	if config.KmsConfig != nil {
		KMSConfig, err = InterceptorKMSConfigToObjectValue(ctx, *config.KmsConfig)
		if err != nil {
			return basetypes.ObjectValue{}, err
		}
	}

	configValue, diag := gwinterceptor.NewConfigValue(
		map[string]attr.Type{
			"schema_data_mode":          basetypes.StringType{},
			"external_storage":          basetypes.BoolType{},
			"topic":                     basetypes.StringType{},
			"schema_registry_config":    schemaRegistryConfig.Type(ctx),
			"enable_audit_log_on_error": basetypes.BoolType{},
			"kms":                       KMSConfig.Type(ctx),
		},
		map[string]attr.Value{
			"external_storage":          NewBoolValue(config.ExternalStorage),
			"topic":                     NewStringValue(config.Topic),
			"schema_data_mode":          NewStringValue(config.SchemaDataMode),
			"schema_registry_config":    schemaRegistryConfig,
			"kms":                       KMSConfig,
			"enable_audit_log_on_error": NewBoolValue(config.EnableAuditLogOnError),
		},
	)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}

	// configValue := gwinterceptor.ConfigValue{
	// 	ExternalStorage:       NewBoolValue(config.ExternalStorage),
	// 	Topic:                 NewStringValue(config.Topic),
	// 	SchemaDataMode:        NewStringValue(config.SchemaDataMode),
	// 	SchemaRegistryConfig:  schemaRegistryConfig,
	// 	KmsConfig:             KMSConfig,
	// 	EnableAuditLogOnError: NewBoolValue(config.EnableAuditLogOnError),
	// }

	objectValue, diag := configValue.ToObjectValue(ctx)
	if diag.HasError() {
		return basetypes.ObjectValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}

	return objectValue, nil
}
