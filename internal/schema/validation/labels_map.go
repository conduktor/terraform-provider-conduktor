package validation

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strings"
)

// Prefix of labels managed by Conduktor Console.
const ManagedLabelsPrefix = "conduktor.io/"

var _ validator.Map = labelsValidator{}

// labelsValidator validates that each map key are valid labels keys.
type labelsValidator struct{}

// Description describes the validation in plain text formatting.
func (v labelsValidator) Description(ctx context.Context) string {
	return "All key in the map must be valid labels keys"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v labelsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateMap performs the validation.
func (v labelsValidator) ValidateMap(_ context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()

	for key := range elements {
		if strings.HasPrefix(key, ManagedLabelsPrefix) {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Managed Label Key",
				"Keys starting with '"+ManagedLabelsPrefix+"' are reserved for Conduktor managed labels.",
			)
			continue
		}
	}
}

// Labels returns a map validator which ensures that all keys in the map
// are valid label keys and do not start with the reserved prefix "conduktor.io/".
func Labels() validator.Map {
	return labelsValidator{}
}
