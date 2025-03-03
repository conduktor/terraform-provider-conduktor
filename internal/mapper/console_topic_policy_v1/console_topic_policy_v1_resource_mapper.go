package console_topic_policy_v1

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	topicPolicy "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_topic_policy_v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *topicPolicy.ConsoleTopicPolicyV1Model) (console.TopicPolicyResource, error) {

	return console.NewTopicPolicyResource(
		r.Name.ValueString(),
		console.TopicPolicySpec{
			Policies: nil,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *console.TopicPolicyResource) (topicPolicy.ConsoleTopicPolicyV1Model, error) {

	specValue, diag := topicPolicy.NewSpecValue(
		//TODO
		map[string]attr.Type{
			"policies": basetypes.MapType{topicPolicy.PoliciesType{}},
		},
		map[string]attr.Value{
			"policies": basetypes.NewMapNull(topicPolicy.PoliciesType{}),
		},
	)
	if diag.HasError() {
		return topicPolicy.ConsoleTopicPolicyV1Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return topicPolicy.ConsoleTopicPolicyV1Model{
		Name: schema.NewStringValue(r.Metadata.Name),
		Spec: specValue,
	}, nil
}
