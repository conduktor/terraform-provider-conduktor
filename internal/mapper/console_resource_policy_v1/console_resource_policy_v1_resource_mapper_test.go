package console_resource_policy_v1

import (
	"context"
	"testing"

	ctlresource "github.com/conduktor/ctl/resource"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	"github.com/conduktor/terraform-provider-conduktor/internal/schema"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestResourcePolicyV1ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonResourcePolicyV1Resource := []byte(test.TestAccTestdata(t, "console/resource_policy_v1/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonResourcePolicyV1Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ResourcePolicy", ctlResource.Kind)
	assert.Equal(t, "v1", ctlResource.Version)
	assert.Equal(t, "resourcepolicy", ctlResource.Name)
	assert.Equal(t, map[string]any{"name": "resourcepolicy", "labels": map[string]any{"key": "value"}}, ctlResource.Metadata)
	assert.Equal(t, map[string]any{"rules": []any{map[string]any{"condition": "condition-1", "errorMessage": "error-1"}, map[string]any{"condition": "condition-2", "errorMessage": "error-2"}}, "targetKind": "Connector", "description": "This is a resource policy"}, ctlResource.Spec)
	assert.Equal(t, jsonResourcePolicyV1Resource, ctlResource.Json)

	// convert into internal model
	internal, err := console.NewResourcePolicyConsoleResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ResourcePolicy", internal.Kind)
	assert.Equal(t, "v1", internal.ApiVersion)
	assert.Equal(t, "resourcepolicy", internal.Metadata.Name)
	assert.Equal(t, map[string]string{
		"key": "value",
	}, internal.Metadata.Labels)
	assert.Equal(t, "Connector", internal.Spec.TargetKind)
	assert.Equal(t, "This is a resource policy", internal.Spec.Description)
	expectedRules := []console.ResourcePolicyConsoleRule{
		{
			Condition:    "condition-1",
			ErrorMessage: "error-1",
		},
		{
			Condition:    "condition-2",
			ErrorMessage: "error-2",
		},
	}
	assert.Equal(t, expectedRules, internal.Spec.Rules)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	labels, _ := schema.StringMapToMapValue(ctx, map[string]string{"key": "value"})
	assert.Equal(t, types.StringValue("resourcepolicy"), tfModel.Name)
	assert.Equal(t, labels, tfModel.Labels)
	assert.Equal(t, types.StringValue("Connector"), tfModel.Spec.TargetKind)
	assert.Equal(t, types.StringValue("This is a resource policy"), tfModel.Spec.Description)
	assert.Equal(t, false, tfModel.Spec.Rules.IsNull())
	assert.Equal(t, false, tfModel.Spec.Rules.IsUnknown())

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ResourcePolicy", internal2.Kind)
	assert.Equal(t, "v1", internal2.ApiVersion)
	assert.Equal(t, "resourcepolicy", internal2.Metadata.Name)
	assert.Equal(t, map[string]string{
		"key": "value",
	}, internal2.Metadata.Labels)
	assert.Equal(t, "Connector", internal2.Spec.TargetKind)
	assert.Equal(t, "This is a resource policy", internal2.Spec.Description)
	assert.Equal(t, expectedRules, internal2.Spec.Rules)

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
