package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.StringTypable = (*SchemaNormalizedType)(nil)

// SchemaNormalizedType is an attribute type that represents a schema string with semantic equality support.
// This type detects schema format from content and applies appropriate normalization.
type SchemaNormalizedType struct {
	basetypes.StringType
}

// String returns a human readable string of the type name.
func (t SchemaNormalizedType) String() string {
	return "SchemaNormalizedType"
}

// ValueType returns the Value type.
func (t SchemaNormalizedType) ValueType(_ context.Context) attr.Value {
	return SchemaNormalized{}
}

// Equal returns true if the given type is equivalent.
func (t SchemaNormalizedType) Equal(o attr.Type) bool {
	other, ok := o.(SchemaNormalizedType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// ValueFromString returns a StringValuable type given a StringValue.
func (t SchemaNormalizedType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return SchemaNormalized{
		StringValue: in,
	}, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.
func (t SchemaNormalizedType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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
