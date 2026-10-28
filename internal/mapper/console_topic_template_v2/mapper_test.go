package console_topic_template_v2

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

func TestTopicTemplateV2ModelMapping(t *testing.T) {
	ctx := context.Background()

	jsonResource := []byte(test.TestAccTestdata(t, "console/topic_template_v2/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonResource)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "TopicTemplate", ctlResource.Kind)
	assert.Equal(t, "v2", ctlResource.Version)
	assert.Equal(t, "my-template", ctlResource.Name)

	// internal model
	internal, err := console.NewTopicTemplateResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "TopicTemplate", internal.Kind)
	assert.Equal(t, "v2", internal.ApiVersion)
	assert.Equal(t, "my-template", internal.Metadata.Name)
	assert.Equal(t, map[string]string{"env": "prod"}, internal.Metadata.Labels)
	assert.Equal(t, "My Topic Template", internal.Spec.DisplayName)
	assert.Equal(t, "A template for production topics", internal.Spec.Description)
	assert.Equal(t, "topic-{{department}}", internal.Spec.Defaults.Metadata.Name)
	assert.Equal(t, map[string]string{"throughput": "high"}, internal.Spec.Defaults.Metadata.Labels)
	assert.Equal(t, int64(6), *internal.Spec.Defaults.Spec.Partitions)
	assert.Equal(t, int64(3), *internal.Spec.Defaults.Spec.ReplicationFactor)
	assert.Equal(t, map[string]string{"retention.ms": "604800000", "cleanup.policy": "delete"}, internal.Spec.Defaults.Spec.Configs)

	// internal → TF
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, types.StringValue("my-template"), tfModel.Name)
	assert.False(t, tfModel.Spec.IsNull())
	assert.False(t, tfModel.Spec.IsUnknown())
	assert.Equal(t, types.StringValue("My Topic Template"), tfModel.Spec.DisplayName)

	// TF → internal (round-trip)
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, internal, internal2)

	// internal → ctl resource (compare without Json field)
	ctlResource2, err := internal2.ToClientResource()
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(ctlResource, ctlResource2, cmpopts.IgnoreFields(ctlresource.Resource{}, "Json")) {
		t.Errorf("expected %+v, got %+v", ctlResource, ctlResource2)
	}
}
