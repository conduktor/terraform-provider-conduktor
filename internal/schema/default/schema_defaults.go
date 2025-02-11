package _default

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var DefaultScope = types.ObjectValueMust(
	map[string]attr.Type{
		"vcluster": types.StringType,
		"username": types.StringType,
		"group":    types.StringType,
	},
	map[string]attr.Value{
		"vcluster": types.StringValue("passthrough"),
		"username": types.StringNull(),
		"group":    types.StringNull(),
	},
)
