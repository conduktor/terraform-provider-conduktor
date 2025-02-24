package console_application_v1

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

func TestApplicationV1ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonApplicationV1Resource := []byte(test.TestAccTestdata(t, "/console/application_v1/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonApplicationV1Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "Application", ctlResource.Kind)
	assert.Equal(t, "v1", ctlResource.Version)
	assert.Equal(t, "application", ctlResource.Name)
	assert.Equal(t, map[string]interface{}{"title": "application title", "description": "application description", "owner": "application owner"}, ctlResource.Spec)
	assert.Equal(t, jsonApplicationV1Resource, ctlResource.Json)

	// convert into internal model
	internal, err := console.NewApplicationConsoleResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "Application", internal.Kind)
	assert.Equal(t, "v1", internal.ApiVersion)
	assert.Equal(t, "application", internal.Metadata.Name)
	assert.Equal(t, "application title", internal.Spec.Title)
	assert.Equal(t, "application description", internal.Spec.Description)
	assert.Equal(t, "application owner", internal.Spec.Owner)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("application"), tfModel.Name)
	assert.Equal(t, types.StringValue("application title"), tfModel.Spec.Title)
	assert.Equal(t, types.StringValue("application description"), tfModel.Spec.Description)
	assert.Equal(t, types.StringValue("application owner"), tfModel.Spec.Owner)

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "Application", internal2.Kind)
	assert.Equal(t, "v1", internal2.ApiVersion)
	assert.Equal(t, "application", internal2.Metadata.Name)
	assert.Equal(t, "application title", internal2.Spec.Title)
	assert.Equal(t, "application description", internal2.Spec.Description)
	assert.Equal(t, "application owner", internal2.Spec.Owner)
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
