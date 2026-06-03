package validation

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// permissionFieldsByResourceType defines which optional fields are allowed for each resource_type.
// Based on the Console OpenAPI spec (ResourcePermissions1 schema, consistent since v1.24.0).
//
// Fields: name, pattern_type, cluster, kafka_connect, ksqldb
// All resource types always have: resource_type (required), permissions (required).
var permissionFieldsByResourceType = map[string]map[string]bool{
	"PLATFORM": {},
	"CLUSTER": {
		"name": true,
	},
	"TOPIC": {
		"name":         true,
		"pattern_type": true,
		"cluster":      true,
	},
	"SUBJECT": {
		"name":         true,
		"pattern_type": true,
		"cluster":      true,
	},
	"CONSUMER_GROUP": {
		"name":         true,
		"pattern_type": true,
		"cluster":      true,
	},
	"KAFKA_CONNECT": {
		"name":          true,
		"pattern_type":  true,
		"cluster":       true,
		"kafka_connect": true,
	},
	"KSQLDB": {
		"cluster": true,
		"ksqldb":  true,
	},
}

// allOptionalPermissionFields lists all optional fields in the permission object.
var allOptionalPermissionFields = []string{"name", "pattern_type", "cluster", "kafka_connect", "ksqldb"}

var _ validator.Set = permissionResourceTypeValidator{}

type permissionResourceTypeValidator struct{}

func (v permissionResourceTypeValidator) Description(_ context.Context) string {
	return "Validates that permission fields are compatible with the specified resource_type"
}

func (v permissionResourceTypeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v permissionResourceTypeValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()

	for i, elem := range elements {
		objVal, ok := elem.(types.Object)
		if !ok {
			continue
		}
		if objVal.IsNull() || objVal.IsUnknown() {
			continue
		}

		attrs := objVal.Attributes()

		// Extract resource_type
		resourceTypeAttr, exists := attrs["resource_type"]
		if !exists {
			continue
		}
		resourceTypeStr, ok := resourceTypeAttr.(types.String)
		if !ok || resourceTypeStr.IsNull() || resourceTypeStr.IsUnknown() {
			continue
		}
		resourceType := resourceTypeStr.ValueString()

		allowedFields, known := permissionFieldsByResourceType[resourceType]
		if !known {
			// Unknown resource_type; other validators handle this
			continue
		}

		// Check each optional field: if it's set but not allowed for this resource_type, report error
		for _, field := range allOptionalPermissionFields {
			if allowedFields[field] {
				continue
			}
			fieldAttr, exists := attrs[field]
			if !exists {
				continue
			}
			strVal, ok := fieldAttr.(types.String)
			if !ok {
				continue
			}
			if !strVal.IsNull() && !strVal.IsUnknown() {
				allowedFieldNames := allowedFieldNamesFor(resourceType)
				resp.Diagnostics.AddAttributeError(
					req.Path.AtSetValue(elem),
					fmt.Sprintf("Invalid field %q for resource_type %q", field, resourceType),
					fmt.Sprintf(
						"Permission at index %d has resource_type %q which does not support the %q field. "+
							"Allowed fields for %s are: %s. "+
							"Please remove the %q field or set it to null.",
						i, resourceType, field, resourceType, allowedFieldNames, field,
					),
				)
			}
		}
	}
}

// allowedFieldNamesFor returns a human-readable list of allowed fields for a resource type.
func allowedFieldNamesFor(resourceType string) string {
	fields := permissionFieldsByResourceType[resourceType]
	if len(fields) == 0 {
		return "resource_type, permissions"
	}
	names := make([]string, 0, len(fields)+2)
	names = append(names, "resource_type", "permissions")
	for _, f := range allOptionalPermissionFields {
		if fields[f] {
			names = append(names, f)
		}
	}
	return strings.Join(names, ", ")
}

// PermissionResourceType returns a set validator that validates permission fields
// are compatible with the specified resource_type.
//
// Each resource_type in Console permissions only supports specific optional fields:
//   - PLATFORM: (no optional fields)
//   - CLUSTER: name
//   - TOPIC: name, pattern_type, cluster
//   - SUBJECT: name, pattern_type, cluster
//   - CONSUMER_GROUP: name, pattern_type, cluster
//   - KAFKA_CONNECT: name, pattern_type, cluster, kafka_connect
//   - KSQLDB: cluster, ksqldb
func PermissionResourceType() validator.Set {
	return permissionResourceTypeValidator{}
}
