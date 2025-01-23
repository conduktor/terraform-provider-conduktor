package customtypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	yaml "gopkg.in/yaml.v3"
)

var (
	_ basetypes.StringValuable                   = (*YAMLNormalized)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*YAMLNormalized)(nil)
	_ xattr.ValidateableAttribute                = (*YAMLNormalized)(nil)
	_ function.ValidateableParameter             = (*YAMLNormalized)(nil)
)

// YAMLNormalized represents a valid YAML string (RFC 9512). Semantic equality logic is defined for YAMLNormalized
// such that inconsequential differences between YAML strings are ignored (whitespace, property order, etc). If you
// need strict, byte-for-byte, string equality, consider using ExactType.
// Inspired by yamltypes.NormalizedType https://github.com/hashicorp/terraform-plugin-framework-yamltypes
type YAMLNormalized struct {
	basetypes.StringValue
}

// Type returns a YAMLNormalizedType.
func (v YAMLNormalized) Type(_ context.Context) attr.Type {
	return YAMLNormalizedType{}
}

// Equal returns true if the given value is equivalent.
func (v YAMLNormalized) Equal(o attr.Value) bool {
	other, ok := o.(YAMLNormalized)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals returns true if the given YAML string value is semantically equal to the current YAML string value. When compared,
// these YAML string values are "normalized" by marshalling them to empty Go structs. This prevents Terraform data consistency errors and
// resource drift due to inconsequential differences in the YAML strings (whitespace, property order, etc).
func (v YAMLNormalized) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(YAMLNormalized)
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

	result, err := yamlEqual(newValue.ValueString(), v.ValueString())

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

func yamlEqual(s1, s2 string) (bool, error) {
	s1, err := normalizeYAMLString(s1)
	if err != nil {
		return false, err
	}

	s2, err = normalizeYAMLString(s2)
	if err != nil {
		return false, err
	}

	return s1 == s2, nil
}

func normalizeYAMLString(yamlStr string) (string, error) {
	dec := yaml.NewDecoder(strings.NewReader(yamlStr))

	var temp interface{}
	if err := dec.Decode(&temp); err != nil {
		return "", err
	}

	yamlBytes, err := yaml.Marshal(&temp)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}

// ValidateAttribute implements attribute value validation. This type requires the value provided to be a String
// value that is valid YAML format (RFC 9512).
func (v YAMLNormalized) ValidateAttribute(ctx context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	var temp interface{}
	if err := yaml.Unmarshal([]byte(v.ValueString()), &temp); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid YAML String Value",
			"A string value was provided that is not valid YAML string format (RFC 9512).\n\n"+
				"Given Value: "+v.ValueString()+"\n",
		)

		return
	}
}

// ValidateParameter implements provider-defined function parameter value validation. This type requires the value
// provided to be a String value that is valid YAML format (RFC 9512).
func (v YAMLNormalized) ValidateParameter(ctx context.Context, req function.ValidateParameterRequest, resp *function.ValidateParameterResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	var temp interface{}
	if err := yaml.Unmarshal([]byte(v.ValueString()), &temp); err != nil {
		resp.Error = function.NewArgumentFuncError(
			req.Position,
			"Invalid YAML String Value: "+
				"A string value was provided that is not valid YAML string format (RFC 9512).\n\n"+
				"Given Value: "+v.ValueString()+"\n",
		)

		return
	}
}

// Unmarshal calls (encoding/yaml).Unmarshal with the YAMLNormalized StringValue and `target` input. A null or unknown value will produce an error diagnostic.
// See encoding/yaml docs for more on usage: https://pkg.go.dev/encoding/yaml#Unmarshal
func (v YAMLNormalized) Unmarshal(target any) diag.Diagnostics {
	var diags diag.Diagnostics

	if v.IsNull() {
		diags.Append(diag.NewErrorDiagnostic("YAMLNormalized YAML Unmarshal Error", "yaml string value is null"))
		return diags
	}

	if v.IsUnknown() {
		diags.Append(diag.NewErrorDiagnostic("YAMLNormalized YAML Unmarshal Error", "yaml string value is unknown"))
		return diags
	}

	err := yaml.Unmarshal([]byte(v.ValueString()), target)
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("YAMLNormalized YAML Unmarshal Error", err.Error()))
	}

	return diags
}

// NewNormalizedNull creates a YAMLNormalized with a null value. Determine whether the value is null via IsNull method.
func NewNormalizedNull() YAMLNormalized {
	return YAMLNormalized{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewNormalizedUnknown creates a YAMLNormalized with an unknown value. Determine whether the value is unknown via IsUnknown method.
func NewNormalizedUnknown() YAMLNormalized {
	return YAMLNormalized{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewNormalizedValue creates a YAMLNormalized with a known value. Access the value via ValueString method.
func NewNormalizedValue(value string) YAMLNormalized {
	return YAMLNormalized{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewNormalizedPointerValue creates a YAMLNormalized with a null value if nil or a known value. Access the value via ValueStringPointer method.
func NewNormalizedPointerValue(value *string) YAMLNormalized {
	return YAMLNormalized{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
