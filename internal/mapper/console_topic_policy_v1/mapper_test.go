package console_topic_policy_v1

import (
	"context"
	"testing"

	ctlresource "github.com/conduktor/ctl/resource"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestTopicPolicyV1ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonTopicPolicyV1Resource := []byte(test.TestAccTestdata(t, "console/topic_policy_v1/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonTopicPolicyV1Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "TopicPolicy", ctlResource.Kind)
	assert.Equal(t, "v1", ctlResource.Version)
	assert.Equal(t, "topicPolicy", ctlResource.Name)
	assert.Equal(t, map[string]any{"name": "topicPolicy"}, ctlResource.Metadata)
	assert.Equal(t, jsonTopicPolicyV1Resource, ctlResource.Json)

	// convert into internal model
	internal, err := console.NewTopicPolicyResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "TopicPolicy", internal.Kind)
	assert.Equal(t, "v1", internal.ApiVersion)
	assert.Equal(t, "topicPolicy", internal.Metadata.Name)
	assert.Equal(t, "OneOf", internal.Spec.Policies["metadata.labels.data-criticality"].OneOf.Constraint)
	assert.Equal(t, []string{"C0", "C1", "C2"}, internal.Spec.Policies["metadata.labels.data-criticality"].OneOf.Values)
	assert.Equal(t, "Range", internal.Spec.Policies["spec.configs.retention.ms"].Range.Constraint)
	assert.Equal(t, int64(60000), internal.Spec.Policies["spec.configs.retention.ms"].Range.Min)
	assert.Equal(t, int64(3600000), internal.Spec.Policies["spec.configs.retention.ms"].Range.Max)
	assert.Equal(t, "Match", internal.Spec.Policies["metadata.name"].Match.Constraint)
	assert.Equal(t, "^click\\.(?<event>[a-z0-9-]+)\\.(avro|json)$", internal.Spec.Policies["metadata.name"].Match.Pattern)
	assert.Equal(t, "AllowedKeys", internal.Spec.Policies["spec.name"].AllowedKeys.Constraint)
	assert.Equal(t, []string{"k1", "k2"}, internal.Spec.Policies["spec.name"].AllowedKeys.Keys)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("topicPolicy"), tfModel.Name)
	assert.Equal(t, false, tfModel.Spec.IsNull())
	assert.Equal(t, false, tfModel.Spec.IsUnknown())
	assert.Equal(t, false, tfModel.Spec.Policies.IsNull())
	assert.Equal(t, false, tfModel.Spec.Policies.IsUnknown())

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "TopicPolicy", internal2.Kind)
	assert.Equal(t, "v1", internal2.ApiVersion)
	assert.Equal(t, "topicPolicy", internal2.Metadata.Name)
	assert.Equal(t, internal, internal2)
	assert.Equal(t, "OneOf", internal2.Spec.Policies["metadata.labels.data-criticality"].OneOf.Constraint)
	assert.Equal(t, []string{"C0", "C1", "C2"}, internal2.Spec.Policies["metadata.labels.data-criticality"].OneOf.Values)
	assert.Equal(t, "Range", internal2.Spec.Policies["spec.configs.retention.ms"].Range.Constraint)
	assert.Equal(t, int64(60000), internal2.Spec.Policies["spec.configs.retention.ms"].Range.Min)
	assert.Equal(t, int64(3600000), internal2.Spec.Policies["spec.configs.retention.ms"].Range.Max)
	assert.Equal(t, "Match", internal2.Spec.Policies["metadata.name"].Match.Constraint)
	assert.Equal(t, "^click\\.(?<event>[a-z0-9-]+)\\.(avro|json)$", internal2.Spec.Policies["metadata.name"].Match.Pattern)
	assert.Equal(t, "AllowedKeys", internal2.Spec.Policies["spec.name"].AllowedKeys.Constraint)
	assert.Equal(t, []string{"k1", "k2"}, internal2.Spec.Policies["spec.name"].AllowedKeys.Keys)

	// convert back to ctl model
	ctlResource2, err := internal2.ToClientResource()
	if err != nil {
		t.Fatal(err)
		return
	}
	// compare without json
	if !cmp.Equal(ctlResource, ctlResource2, cmpopts.IgnoreFields(ctlresource.Resource{}, "Json")) {
		t.Errorf("expected %+v, got %+v", ctlResource, ctlResource2)
	}
}
