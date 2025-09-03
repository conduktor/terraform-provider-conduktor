package customtypes

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/emicklei/proto"
	"github.com/hamba/avro/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/qri-io/jsonschema"
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

	format := detectSchemaFormat(v.ValueString())
	if format == "UNKNOWN" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Unknown Schema Format",
			"The schema format could not be determined from the content. "+
				"Ensure the schema is valid and in a supported format (AVRO, PROTOBUF, JSON).\n\n"+
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

	format := detectSchemaFormat(v.ValueString())
	if format == "UNKNOWN" {
		resp.Error = function.NewArgumentFuncError(
			req.Position,
			"Unknown Schema Format: "+
				"The schema format could not be determined from the content. "+
				"Ensure the schema is valid and in a supported format (AVRO, PROTOBUF, JSON).\n\n"+
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

	// 1. Check for Protobuf first
	if isProtobufSchema(trimmed) {
		return "PROTOBUF"
	}

	// 2. Check for AVRO schema
	if isAvroSchema(trimmed) {
		return "AVRO"
	}

	// 3. Check for JSON Schema
	if isJSONSchema(trimmed) {
		return "JSON"
	}

	return "UNKNOWN"
}

// isProtobufSchema checks if the string is a valid Protobuf schema.
func isProtobufSchema(schema string) bool {
	reader := strings.NewReader(schema)
	parser := proto.NewParser(reader)

	// Try to parse as protobuf
	_, err := parser.Parse()
	return err == nil
}

// isAvroSchema checks if the string is a valid AVRO schema using hamba/avro library.
func isAvroSchema(schema string) bool {
	// Try to parse as AVRO schema
	_, err := avro.Parse(schema)
	return err == nil
}

// isJSONSchema checks if the string is a valid JSON Schema.
func isJSONSchema(schema string) bool {
	// First check if it's valid JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(schema), &jsonData); err != nil {
		return false
	}

	// Try to parse as JSON Schema
	rs := &jsonschema.Schema{}
	if err := json.Unmarshal([]byte(schema), rs); err != nil {
		return false
	}

	// Additional validation - check for common JSON Schema properties
	var schemaMap map[string]interface{}
	if err := json.Unmarshal([]byte(schema), &schemaMap); err != nil {
		return false
	}

	// Look for JSON Schema indicators
	hasSchemaProperty := false
	hasJsonSchemaProps := false

	for key := range schemaMap {
		switch key {
		case "$schema", "$id", "$ref":
			hasSchemaProperty = true
		case "type", "properties", "items", "additionalProperties", "required",
			"minimum", "maximum", "pattern", "enum", "anyOf", "oneOf", "allOf":
			hasJsonSchemaProps = true
		}
	}

	// It's a JSON Schema if it has explicit schema markers or typical JSON Schema properties
	return hasSchemaProperty || hasJsonSchemaProps
}

// normalizeAvroSchema normalizes AVRO schema using hamba/avro library.
func normalizeAvroSchema(schemaStr string) (string, error) {
	// Use hamba/avro to parse and normalize the schema
	schema, err := avro.Parse(schemaStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse AVRO schema: %w", err)
	}

	// Get the canonical schema string representation
	return schema.String(), nil
}

// normalizeProtobufSchema normalizes Protobuf schema.
func normalizeProtobufSchema(schemaStr string) (string, error) {
	reader := strings.NewReader(schemaStr)
	parser := proto.NewParser(reader)

	definition, err := parser.Parse()
	if err != nil {
		return "", fmt.Errorf("failed to parse Protobuf schema: %w", err)
	}

	var result strings.Builder

	// Handle syntax declaration
	proto.Walk(definition, func(v proto.Visitee) {
		switch element := v.(type) {
		case *proto.Syntax:
			result.WriteString(fmt.Sprintf("syntax = \"%s\";\n", element.Value))
		}
	})

	// Collect and sort imports
	var imports []string
	proto.Walk(definition, func(v proto.Visitee) {
		if imp, ok := v.(*proto.Import); ok {
			imports = append(imports, fmt.Sprintf("import \"%s\";", imp.Filename))
		}
	})
	sort.Strings(imports)
	for _, imp := range imports {
		result.WriteString(imp + "\n")
	}

	// Handle messages, services, enums with simplified approach
	proto.Walk(definition, func(v proto.Visitee) {
		switch element := v.(type) {
		case *proto.Message:
			result.WriteString(normalizeProtoMessage(element))
		case *proto.Service:
			result.WriteString(normalizeProtoService(element))
		case *proto.Enum:
			result.WriteString(normalizeProtoEnum(element))
		}
	})

	return strings.TrimSpace(result.String()), nil
}

// normalizeProtoMessage normalizes a protobuf message definition.
func normalizeProtoMessage(msg *proto.Message) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("message %s {\n", msg.Name))

	for _, element := range msg.Elements {
		switch field := element.(type) {
		case *proto.NormalField:
			result.WriteString(fmt.Sprintf("%s %s = %d;\n", field.Type, field.Name, field.Sequence))
		case *proto.MapField:
			result.WriteString(fmt.Sprintf("map<%s, %s> %s = %d;\n", field.KeyType, field.Type, field.Name, field.Sequence))
		case *proto.OneOfField:
			result.WriteString(fmt.Sprintf("oneof %s {\n", field.Name))
		}
		result.WriteString("}\n")
	}

	result.WriteString("}\n")
	return result.String()
}

// normalizeProtoService normalizes a protobuf service definition.
func normalizeProtoService(svc *proto.Service) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("service %s {\n", svc.Name))

	for _, element := range svc.Elements {
		if rpc, ok := element.(*proto.RPC); ok {
			result.WriteString(fmt.Sprintf("rpc %s(%s) returns (%s);\n",
				rpc.Name, rpc.RequestType, rpc.ReturnsType))
		}
	}

	result.WriteString("}\n")
	return result.String()
}

// normalizeProtoEnum normalizes a protobuf enum definition.
func normalizeProtoEnum(enum *proto.Enum) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("enum %s {\n", enum.Name))

	for _, element := range enum.Elements {
		if field, ok := element.(*proto.EnumField); ok {
			result.WriteString(fmt.Sprintf("%s = %d;\n", field.Name, field.Integer))
		}
	}

	result.WriteString("}\n")
	return result.String()
}

// normalizeJSONSchema normalizes JSON schema to canonical form.
func normalizeJSONSchema(schemaStr string) (string, error) {
	// Parse the schema
	rs := &jsonschema.Schema{}
	if err := json.Unmarshal([]byte(schemaStr), rs); err != nil {
		return "", fmt.Errorf("failed to parse JSON schema: %w", err)
	}

	// Marshal back to get consistent formatting
	normalized, err := json.Marshal(rs)
	if err != nil {
		return "", fmt.Errorf("failed to normalize JSON schema: %w", err)
	}

	return string(normalized), nil
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

// Public wrapper functions for external use.
func NormalizeAvroSchema(schema string) (string, error) {
	return normalizeAvroSchema(schema)
}

func NormalizeProtobufSchema(schema string) (string, error) {
	return normalizeProtobufSchema(schema)
}

func NormalizeJSONSchema(schema string) (string, error) {
	return normalizeJSONSchema(schema)
}
