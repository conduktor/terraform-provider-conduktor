package customtypes

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ basetypes.StringValuable                   = (*SchemaNormalized)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*SchemaNormalized)(nil)
	_ xattr.ValidateableAttribute                = (*SchemaNormalized)(nil)
	_ function.ValidateableParameter             = (*SchemaNormalized)(nil)
)

// SchemaNormalized represents a schema string with format-aware semantic equality.
// It automatically detects the schema format (AVRO, PROTOBUF, JSON) from content
// and applies appropriate normalization to prevent drift caused by formatting differences.
type SchemaNormalized struct {
	basetypes.StringValue
}

// Type returns a SchemaNormalizedType.
func (v SchemaNormalized) Type(_ context.Context) attr.Type {
	return SchemaNormalizedType{}
}

// Equal returns true if the given value is equivalent.
func (v SchemaNormalized) Equal(o attr.Value) bool {
	other, ok := o.(SchemaNormalized)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals compares schemas for semantic equality based on detected format.
func (v SchemaNormalized) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(SchemaNormalized)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	// Detect format from schema content
	format := determineSchemaFormat(v.ValueString(), newValue.ValueString())

	result, err := schemaEqual(v.ValueString(), newValue.ValueString(), format)

	if err != nil {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected error occurred while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)

		return false, diags
	}

	return result, diags
}

// ValidateAttribute implements basic schema string validation.
func (v SchemaNormalized) ValidateAttribute(ctx context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	schema := v.ValueString()
	if strings.TrimSpace(schema) == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Schema String Value",
			"A schema string value cannot be empty.\n\n"+
				"Given Value: "+schema+"\n",
		)
	}
}

// ValidateParameter implements function parameter validation.
func (v SchemaNormalized) ValidateParameter(ctx context.Context, req function.ValidateParameterRequest, resp *function.ValidateParameterResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	schema := v.ValueString()
	if strings.TrimSpace(schema) == "" {
		resp.Error = function.NewArgumentFuncError(
			req.Position,
			"Invalid Schema String Value: "+
				"A schema string value cannot be empty.\n\n"+
				"Given Value: "+schema+"\n",
		)
	}
}

// schemaEqual compares two schema strings for semantic equality.
func schemaEqual(schema1, schema2, format string) (bool, error) {
	norm1, err := normalizeSchemaByFormat(schema1, format)
	if err != nil {
		return false, err
	}

	norm2, err := normalizeSchemaByFormat(schema2, format)
	if err != nil {
		return false, err
	}

	return norm1 == norm2, nil
}

// normalizeSchemaByFormat normalizes schema based on detected format.
func normalizeSchemaByFormat(schema, format string) (string, error) {
	switch strings.ToUpper(format) {
	case "AVRO":
		return normalizeAvroSchema(schema)
	case "PROTOBUF":
		return normalizeProtobufSchema(schema)
	case "JSON":
		// JSON schemas work fine as-is, just ensure consistent formatting.
		return normalizeJSONSchema(schema)
	default:
		// Unknown format, return as-is.
		return schema, nil
	}
}

// determineSchemaFormat determines schema format from schema content patterns.
func determineSchemaFormat(schema1, schema2 string) string {
	// Check both schemas and return the first detected format.
	for _, schema := range []string{schema1, schema2} {
		format := detectSchemaFormat(schema)
		if format != "UNKNOWN" {
			return format
		}
	}
	return "UNKNOWN"
}

// detectSchemaFormat analyzes a single schema string to determine its format.
func detectSchemaFormat(schemaStr string) string {
	trimmed := strings.TrimSpace(schemaStr)
	if trimmed == "" {
		return "UNKNOWN"
	}

	// 1. Check for Protobuf patterns first.
	if isProtobufSchema(trimmed) {
		return "PROTOBUF"
	}

	// 2. Check for Avro patterns.
	if isAvroSchema(trimmed) {
		return "AVRO"
	}

	// 3. If it's valid JSON but not Avro, it's JSON schema
	if isValidJSON(trimmed) {
		return "JSON"
	}

	return "UNKNOWN"
}

// isProtobufSchema checks for protobuf syntax patterns.
func isProtobufSchema(schema string) bool {
	// Look for protobuf syntax declaration.
	if strings.Contains(schema, `syntax =`) || strings.Contains(schema, `syntax=`) {
		return true
	}

	// Look for other protobuf-specific patterns (and ensure it's not JSON).
	if !strings.HasPrefix(schema, "{") && !strings.HasPrefix(schema, "[") {
		if strings.Contains(schema, "message ") ||
			strings.Contains(schema, "service ") ||
			strings.Contains(schema, "enum ") {
			return true
		}
	}

	return false
}

// isAvroSchema checks for Avro schema patterns.
func isAvroSchema(schema string) bool {
	// Avro primitive types as strings.
	avroPrimitives := []string{
		`"null"`, `"boolean"`, `"int"`, `"long"`, `"float"`, `"double"`,
		`"bytes"`, `"string"`, `"record"`, `"enum"`, `"array"`, `"map"`, `"union"`, `"fixed"`,
	}

	// Check if it's just a primitive type string.
	for _, primitive := range avroPrimitives {
		if schema == primitive {
			return true
		}
	}

	// If it's JSON, check for Avro-specific patterns.
	if isValidJSON(schema) {
		// Must have "type" field for Avro.
		if !strings.Contains(schema, `"type"`) {
			return false
		}

		// Check for Avro record pattern with fields.
		if strings.Contains(schema, `"fields"`) &&
			(strings.Contains(schema, `"type":"record"`) || strings.Contains(schema, `"type": "record"`)) {
			return true
		}

		// Check for Avro enum pattern.
		if strings.Contains(schema, `"symbols"`) &&
			(strings.Contains(schema, `"type":"enum"`) || strings.Contains(schema, `"type": "enum"`)) {
			return true
		}

		// Check for simple Avro type definition - but only if it doesn't have extra fields
		var parsed map[string]interface{}
		if json.Unmarshal([]byte(schema), &parsed) == nil {
			// If it has a "name" field, it's a named type (not a primitive)
			if _, hasName := parsed["name"]; hasName {
				return true // Named types are always AVRO
			}

			// For simple type definitions, check primitive patterns
			for _, primitive := range avroPrimitives {
				typePattern := `"type":` + primitive
				typePatternSpaced := `"type": ` + primitive
				if strings.Contains(schema, typePattern) || strings.Contains(schema, typePatternSpaced) {
					return true
				}
			}
		}
	}

	return false
}

// isValidJSON checks if a string is valid JSON.
func isValidJSON(str string) bool {
	var js interface{}
	return json.Unmarshal([]byte(str), &js) == nil
}

// normalizeAvroSchema normalizes Avro schema to canonical form.
func normalizeAvroSchema(schemaStr string) (string, error) {
	// If it's a primitive type string, return as-is.
	if !strings.HasPrefix(strings.TrimSpace(schemaStr), "{") && !strings.HasPrefix(strings.TrimSpace(schemaStr), "[") {
		return strings.TrimSpace(schemaStr), nil
	}

	var schema interface{}
	if err := json.Unmarshal([]byte(schemaStr), &schema); err != nil {
		return schemaStr, fmt.Errorf("failed to parse Avro schema JSON: %w", err)
	}

	normalized := normalizeAvroObject(schema)

	// Marshal with consistent formatting (no indentation, sorted keys).
	normalizedBytes, err := json.Marshal(normalized)
	if err != nil {
		return schemaStr, fmt.Errorf("failed to marshal normalized Avro schema: %w", err)
	}

	return string(normalizedBytes), nil
}

// normalizeAvroObject recursively normalizes Avro schema objects.
func normalizeAvroObject(obj interface{}) interface{} {
	switch v := obj.(type) {
	case map[string]interface{}:
		return normalizeAvroMap(v)
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = normalizeAvroObject(item)
		}
		return result
	default:
		return v
	}
}

// normalizeAvroMap normalizes map according to Schema Registry canonical form rules.
func normalizeAvroMap(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Schema Registry canonical form field ordering: name, type, fields, symbols, items, values, size.
	fieldOrder := []string{"name", "type", "fields", "symbols", "items", "values", "size"}

	// Add ordered fields first.
	for _, field := range fieldOrder {
		if val, exists := m[field]; exists {
			result[field] = normalizeAvroObject(val)
		}
	}

	// Add remaining fields in alphabetical order.
	var remainingFields []string
	for key := range m {
		found := false
		for _, orderedField := range fieldOrder {
			if key == orderedField {
				found = true
				break
			}
		}
		if !found {
			remainingFields = append(remainingFields, key)
		}
	}
	sort.Strings(remainingFields)

	for _, field := range remainingFields {
		result[field] = normalizeAvroObject(m[field])
	}

	return result
}

// normalizeProtobufSchema normalizes Protobuf schema.
func normalizeProtobufSchema(schemaStr string) (string, error) {
	// Normalize protobuf by: removing comments, normalizing whitespace, sorting imports.
	lines := strings.Split(schemaStr, "\n")
	var normalizedLines []string
	var imports []string
	inBlockComment := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines.
		if trimmed == "" {
			continue
		}

		// Handle block comments.
		if strings.Contains(line, "/*") {
			inBlockComment = true
		}
		if strings.Contains(line, "*/") {
			inBlockComment = false
			continue
		}
		if inBlockComment {
			continue
		}

		// Remove line comments.
		if strings.Contains(trimmed, "//") {
			parts := strings.SplitN(trimmed, "//", 2)
			trimmed = strings.TrimSpace(parts[0])
			if trimmed == "" {
				continue
			}
		}

		// Collect imports for sorting
		if strings.HasPrefix(trimmed, "import ") {
			imports = append(imports, trimmed)
			continue
		}

		// Normalize field definitions.
		trimmed = normalizeProtobufField(trimmed)

		normalizedLines = append(normalizedLines, trimmed)
	}

	// Sort imports and add them at the beginning (after syntax).
	sort.Strings(imports)

	// Rebuild the schema.
	var result []string
	syntaxAdded := false

	for _, line := range normalizedLines {
		if strings.HasPrefix(line, "syntax ") && !syntaxAdded {
			result = append(result, line)
			// Add sorted imports after syntax.
			result = append(result, imports...)
			syntaxAdded = true
		} else if !strings.HasPrefix(line, "syntax ") {
			result = append(result, line)
		}
	}

	// If no syntax was found, add imports at the beginning.
	if !syntaxAdded && len(imports) > 0 {
		result = append(imports, result...)
	}

	return strings.Join(result, "\n"), nil
}

// normalizeProtobufField normalizes protobuf field definitions.
func normalizeProtobufField(field string) string {
	// Normalize spacing around = and ;.
	field = regexp.MustCompile(`\s*=\s*`).ReplaceAllString(field, " = ")
	field = regexp.MustCompile(`\s*;\s*`).ReplaceAllString(field, ";")

	// Ensure single space between type and field name.
	field = regexp.MustCompile(`\s+`).ReplaceAllString(field, " ")

	return strings.TrimSpace(field)
}

// normalizeJSONSchema normalizes JSON schema (basic normalization).
func normalizeJSONSchema(schemaStr string) (string, error) {
	var schema interface{}
	if err := json.Unmarshal([]byte(schemaStr), &schema); err != nil {
		return schemaStr, fmt.Errorf("failed to parse JSON schema: %w", err)
	}

	// Marshal with consistent formatting (no indentation, sorted keys)
	normalizedBytes, err := json.Marshal(schema)
	if err != nil {
		return schemaStr, fmt.Errorf("failed to marshal normalized JSON schema: %w", err)
	}

	return string(normalizedBytes), nil
}

// NewSchemaNormalizedNull creates a SchemaNormalized with a null value.
func NewSchemaNormalizedNull() SchemaNormalized {
	return SchemaNormalized{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewSchemaNormalizedUnknown creates a SchemaNormalized with an unknown value.
func NewSchemaNormalizedUnknown() SchemaNormalized {
	return SchemaNormalized{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewSchemaNormalizedValue creates a SchemaNormalized with a known value.
func NewSchemaNormalizedValue(value string) SchemaNormalized {
	return SchemaNormalized{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewSchemaNormalizedPointerValue creates a SchemaNormalized with a null value if nil or a known value.
func NewSchemaNormalizedPointerValue(value *string) SchemaNormalized {
	return SchemaNormalized{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}

// Public wrapper functions for normalization (used by mappers)

// NormalizeAvroSchema normalizes an AVRO schema to canonical form.
func NormalizeAvroSchema(schema string) (string, error) {
	return normalizeAvroSchema(schema)
}

// NormalizeProtobufSchema normalizes a PROTOBUF schema.
func NormalizeProtobufSchema(schema string) (string, error) {
	return normalizeProtobufSchema(schema)
}

// NormalizeJSONSchema normalizes a JSON schema.
func NormalizeJSONSchema(schema string) (string, error) {
	return normalizeJSONSchema(schema)
}
