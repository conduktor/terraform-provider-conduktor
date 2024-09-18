package schema

import (
	"context"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Provider string configuration extracted from the schema and environment variables.
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

// Provider bool configuration extracted from the schema and environment variables.
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

// Convert a string to a basetypes.StringValue.
func NewStringValue(s string) basetypes.StringValue {
	if s == "" {
		return basetypes.NewStringNull()
	}
	return basetypes.NewStringValue(s)
}

// Convert a ListValue to a string array.
func ListValueToStringArray(ctx context.Context, list basetypes.ListValue) ([]string, diag.Diagnostics) {
	if list.IsNull() {
		return nil, diag.Diagnostics{}
	}

	var result []string
	diag := list.ElementsAs(ctx, &result, false)
	return result, diag
}

// Convert a string array to a ListValue.
func StringArrayToListValue(array []string) (basetypes.ListValue, diag.Diagnostics) {

	var values []attr.Value
	for _, f := range array {
		values = append(values, types.StringValue(f))
	}

	return types.ListValue(types.StringType, values)
}
