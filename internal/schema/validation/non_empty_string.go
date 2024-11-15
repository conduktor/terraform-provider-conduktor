package validation

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.String = nonEmpty{}

// stringNonEmpty validates that a string Attribute's length is at least a certain value.
type nonEmpty struct {
}

// Description describes the validation in plain text formatting.
func (validator nonEmpty) Description(_ context.Context) string {
	return "string should not be empty"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator nonEmpty) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (v nonEmpty) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if len(value) == 0 {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			"\"\"",
		))
		return
	}
}

// NonEmptyString returns a validator which ensures that any configured
// attribute value a not empty string. Null (unconfigured) and unknown (known after apply)
// values are skipped.
func NonEmptyString() validator.String {
	return nonEmpty{}
}
