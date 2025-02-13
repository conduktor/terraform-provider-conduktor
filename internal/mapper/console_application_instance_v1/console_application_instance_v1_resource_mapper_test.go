package console_application_instance_v1

import (
	"context"
	"testing"

	ctlresource "github.com/conduktor/ctl/resource"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	"github.com/conduktor/terraform-provider-conduktor/internal/schema"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestApplicationInstanceV1ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonApplicationInstanceV1Resource := []byte(test.TestAccTestdata(t, "console/application_instance_v1/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonApplicationInstanceV1Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ApplicationInstance", ctlResource.Kind)
	assert.Equal(t, "v1", ctlResource.Version)
	assert.Equal(t, "appinstance", ctlResource.Name)
	assert.Equal(t, map[string]interface{}{"name": "appinstance"}, ctlResource.Metadata)
	assert.Equal(t, jsonApplicationInstanceV1Resource, ctlResource.Json)

	// convert into internal model
	internal, err := console.NewApplicationInstanceConsoleResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ApplicationInstance", internal.Kind)
	assert.Equal(t, "v1", internal.ApiVersion)
	assert.Equal(t, "appinstance", internal.Metadata.Name)
	assert.Equal(t, "cluster", internal.Spec.Cluster)
	assert.Equal(t, []string{"ref1", "ref2"}, internal.Spec.TopicPolicyRef)
	assert.Equal(t, false, internal.Spec.ApplicationManagedServiceAccount)
	assert.Equal(t, "serviceaccount", internal.Spec.ServiceAccount)
	assert.Equal(t, "PRIVATE", internal.Spec.DefaultCatalogVisibility)
	expectedInternalResources := []model.ResourceWithOwnership{
		{
			Type:           "TOPIC",
			Name:           "resource-1",
			PatternType:    "LITERAL",
			ConnectCluster: "connectCluster",
			OwnershipMode:  "ALL",
		},
	}
	assert.Equal(t, expectedInternalResources, internal.Spec.Resources)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	topicPolicyRef, _ := schema.StringArrayToSetValue([]string{"ref1", "ref2"})
	assert.Equal(t, types.StringValue("sales"), tfModel.Name)
	assert.Equal(t, types.StringValue("cluster"), tfModel.Spec.Cluster)
	assert.Equal(t, topicPolicyRef, tfModel.Spec.TopicPolicyRef)
	assert.Equal(t, types.BoolValue(false), tfModel.Spec.ApplicationManagedServiceAccount)
	assert.Equal(t, types.StringValue("serviceaccount"), tfModel.Spec.ServiceAccount)
	assert.Equal(t, types.StringValue("PRIVATE"), tfModel.Spec.DefaultCatalogVisibility)
	// do not test Resources as it's a pain to parse ListValue

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ApplicationInstance", internal2.Kind)
	assert.Equal(t, "v1", internal2.ApiVersion)
	assert.Equal(t, "appinstance", internal.Metadata.Name)
	assert.Equal(t, "cluster", internal.Spec.Cluster)
	assert.Equal(t, []string{"ref1", "ref2"}, internal.Spec.TopicPolicyRef)
	assert.Equal(t, false, internal.Spec.ApplicationManagedServiceAccount)
	assert.Equal(t, "serviceaccount", internal.Spec.ServiceAccount)
	assert.Equal(t, "PRIVATE", internal.Spec.DefaultCatalogVisibility)
	assert.Equal(t, expectedInternalResources, internal2.Spec.Resources)
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
