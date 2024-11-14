package test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	groups "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_group_v2"
	users "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_user_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Helper to read testdata files into string.
func TestAccTestdata(t *testing.T, path string) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not get current file")
	}
	example, err := os.ReadFile(filepath.Join(filepath.Dir(currentFile), "..", "testdata", path))
	if err != nil {
		t.Fatal(err)
	}
	return string(example)
}

// Helper to read examples files into string.
// path is defined relative to examples directory.
func TestAccExample(t *testing.T, path ...string) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not get current file")
	}
	pathFragments := append([]string{filepath.Dir(currentFile), "..", "..", "examples"}, path...)
	example, err := os.ReadFile(filepath.Join(pathFragments...))
	if err != nil {
		t.Fatal(err)
	}
	return string(example)
}

// Check if a string contains all expected values.
func TestCheckResourceAttrContainsStringsFunc(expected ...string) func(value string) error {
	return func(value string) error {
		for _, e := range expected {
			if !strings.Contains(value, e) {
				return fmt.Errorf("expected manifest to contain %q", e)
			}
		}
		return nil
	}
}

// Check if license is setup in env to enable some tests behind license.
func CheckEnterpriseEnabled(t *testing.T) {
	if !(os.Getenv("CDK_LICENSE") != "") {
		t.Skip("Skipping TestAccGroupV2Resource tests in free mode as it requires a license set on CDK_LICENSE env var")
	}
}

// Provider configuration pre-checks.
func TestAccPreCheck(t *testing.T) {
	// check that the environment variables are set
	if os.Getenv("CDK_BASE_URL") == "" {
		t.Fatal("CDK_BASE_URL must be set for acceptance tests")
	}
	if os.Getenv("CDK_ADMIN_EMAIL") == "" {
		t.Fatal("CDK_ADMIN_EMAIL must be set for acceptance tests")
	}
	if os.Getenv("CDK_ADMIN_PASSWORD") == "" {
		t.Fatal("CDK_ADMIN_PASSWORD must be set for acceptance tests")
	}
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
func PermissionArrayToSetValue(ctx context.Context, resource string, arr []model.Permission) (basetypes.SetValue, error) {
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
			"name":          schema.NewStringValue(p.Name),
			"resource_type": schema.NewStringValue(p.ResourceType),
			"permissions":   flagsList,
			"pattern_type":  schema.NewStringValue(p.PatternType),
			"cluster":       schema.NewStringValue(p.Cluster),
			"kafka_connect": schema.NewStringValue(p.KafkaConnect),
		}

		if resource == "groups" {
			permObj, diag := groups.NewPermissionsValue(types, values)
			if diag.HasError() {
				return basetypes.SetValue{}, mapper.WrapDiagError(diag, "permissions", mapper.FromTerraform)
			}
			tfPermissions = append(tfPermissions, permObj)
		} else if resource == "users" {
			permObj, diag := users.NewPermissionsValue(types, values)
			if diag.HasError() {
				return basetypes.SetValue{}, mapper.WrapDiagError(diag, "permissions", mapper.FromTerraform)
			}
			tfPermissions = append(tfPermissions, permObj)
		}

	}

	if resource == "groups" {
		permissionsList, diag = types.SetValue(groups.PermissionsValue{}.Type(ctx), tfPermissions)
	} else if resource == "users" {
		permissionsList, diag = types.SetValue(users.PermissionsValue{}.Type(ctx), tfPermissions)
	}

	if diag.HasError() {
		return basetypes.SetValue{}, mapper.WrapDiagError(diag, "permissions", mapper.FromTerraform)
	}

	return permissionsList, nil
}

// Parse a Set into an array of Permissions based on resource type.
func SetValueToPermissionArray(ctx context.Context, resource string, set basetypes.SetValue) ([]model.Permission, error) {
	permissions := make([]model.Permission, 0)
	var diag diag.Diagnostics

	// Ideally the switch within groups and users would have less replication.
	// This might be worth a re-work in the future.
	// NOTE: an idea would be to use ObjectValue instead of user/group PermissionsValue.
	if !set.IsNull() && !set.IsUnknown() {
		// Case for groups
		if resource == "groups" {
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
		} else if resource == "users" {
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
