package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// alwaysUseStateForStringModifier implements the plan modifier.
type alwaysUseStateForStringModifier struct{}

// AlwaysUseStateForString returns a plan modifier that always uses the state value for a string attribute.
// Use this for write-only fields that are sent to the API on creation but never returned on read.
func AlwaysUseStateForString() planmodifier.String {
	return &alwaysUseStateForStringModifier{}
}

// Description returns a human-readable description of the plan modifier.
func (m *alwaysUseStateForStringModifier) Description(ctx context.Context) string {
	return "Always uses the state value for this attribute, ignoring any changes from the API response."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m *alwaysUseStateForStringModifier) MarkdownDescription(ctx context.Context) string {
	return "Always uses the state value for this attribute, ignoring any changes from the API response."
}

// PlanModifyString implements the plan modification logic.
func (m *alwaysUseStateForStringModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If there's no state value, don't modify the plan (first apply)
	if req.StateValue.IsNull() {
		return
	}

	// Always use the state value, so API not returning this field doesn't cause a diff
	resp.PlanValue = req.StateValue
}
