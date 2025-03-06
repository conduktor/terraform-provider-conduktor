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
	assert.Equal(t, map[string]interface{}{"name": "topicPolicy"}, ctlResource.Metadata)
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

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("topicPolicy"), tfModel.Name)
	// TODO
	assert.Equal(t, true, tfModel.Spec.Policies.IsNull())
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
