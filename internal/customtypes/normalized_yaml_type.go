package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable = (*YAMLNormalizedType)(nil)
)

// YAMLNormalizedType is an attribute type that represents a valid JSON string (RFC 7159). Semantic equality logic is defined for YAMLNormalizedType
// such that inconsequential differences between JSON strings are ignored (whitespace, property order, etc). If you need strict, byte-for-byte,
// string equality, consider using ExactType.
type YAMLNormalizedType struct {
	basetypes.StringType
}

// String returns a human readable string of the type name.
func (t YAMLNormalizedType) String() string {
	return "customtypes.YAMLNormalizedType"
}

// ValueType returns the Value type.
func (t YAMLNormalizedType) ValueType(ctx context.Context) attr.Value {
	return YAMLNormalized{}
}

// Equal returns true if the given type is equivalent.
func (t YAMLNormalizedType) Equal(o attr.Type) bool {
	other, ok := o.(YAMLNormalizedType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// ValueFromString returns a StringValuable type given a StringValue.
func (t YAMLNormalizedType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return YAMLNormalized{
		StringValue: in,
	}, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to convert the tftypes.Value into a more convenient Go type
// for the provider to consume the data with.
func (t YAMLNormalizedType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}
