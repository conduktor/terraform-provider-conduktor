package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// AlwaysUseStateForSet returns a plan modifier that always uses the state value for a set attribute.
// This effectively ignores any changes to the attribute during planning.
func AlwaysUseStateForSet() planmodifier.Set {
	return &alwaysUseStateForSetModifier{}
}

// alwaysUseStateForSetModifier implements the plan modifier.
type alwaysUseStateForSetModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m *alwaysUseStateForSetModifier) Description(ctx context.Context) string {
	return "Always uses the state value for this attribute, ignoring any changes."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m *alwaysUseStateForSetModifier) MarkdownDescription(ctx context.Context) string {
	return "Always uses the state value for this attribute, ignoring any changes."
}

// PlanModifySet implements the plan modification logic.
func (m *alwaysUseStateForSetModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// If there's no state value, don't modify the plan
	if req.StateValue.IsNull() {
		return
	}

	// Always use the state value, ignoring any changes
	resp.PlanValue = req.StateValue
}
