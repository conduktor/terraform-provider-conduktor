package console_application_instance_permission_v1

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

func TestApplicationInstancePermissionV1ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonApplicationInstancePermissionV1Resource := []byte(test.TestAccTestdata(t, "console/application_instance_permission_v1/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonApplicationInstancePermissionV1Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ApplicationInstancePermission", ctlResource.Kind)
	assert.Equal(t, "v1", ctlResource.Version)
	assert.Equal(t, "appinstance-permission", ctlResource.Name)
	assert.Equal(t, map[string]any{"name": "appinstance-permission", "application": "app", "appInstance": "appinstance"}, ctlResource.Metadata)
	assert.Equal(t, jsonApplicationInstancePermissionV1Resource, ctlResource.Json)

	// convert into internal model
	internal, err := console.NewApplicationInstancePermissionConsoleResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ApplicationInstancePermission", internal.Kind)
	assert.Equal(t, "v1", internal.ApiVersion)
	assert.Equal(t, "appinstance-permission", internal.Metadata.Name)
	assert.Equal(t, "app", internal.Metadata.Application)
	assert.Equal(t, "appinstance", internal.Metadata.AppInstance)
	expectedInternalResource := console.AppInstancePermissionResource{
		Type:           "CONSUMER_GROUP",
		Name:           "resource",
		PatternType:    "PREFIXED",
		ConnectCluster: "connectCluster",
	}
	assert.Equal(t, expectedInternalResource, internal.Spec.Resource)
	assert.Equal(t, "READ", internal.Spec.UserPermission)
	assert.Equal(t, "READ", internal.Spec.ServiceAccountPermission)
	assert.Equal(t, "appinstance", internal.Spec.GrantedTo)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("appinstance-permission"), tfModel.Name)
	assert.Equal(t, types.StringValue("app"), tfModel.Application)
	assert.Equal(t, types.StringValue("appinstance"), tfModel.AppInstance)
	assert.Equal(t, false, tfModel.Spec.Resource.IsNull())
	assert.Equal(t, false, tfModel.Spec.Resource.IsUnknown())
	assert.Equal(t, types.StringValue("READ"), tfModel.Spec.UserPermission)
	assert.Equal(t, types.StringValue("READ"), tfModel.Spec.ServiceAccountPermission)
	assert.Equal(t, types.StringValue("appinstance"), tfModel.Spec.GrantedTo)

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "ApplicationInstancePermission", internal2.Kind)
	assert.Equal(t, "v1", internal2.ApiVersion)
	assert.Equal(t, "appinstance-permission", internal2.Metadata.Name)
	assert.Equal(t, "app", internal2.Metadata.Application)
	assert.Equal(t, "appinstance", internal2.Metadata.AppInstance)
	assert.Equal(t, expectedInternalResource, internal2.Spec.Resource)
	assert.Equal(t, "READ", internal2.Spec.UserPermission)
	assert.Equal(t, "READ", internal2.Spec.ServiceAccountPermission)
	assert.Equal(t, "appinstance", internal2.Spec.GrantedTo)

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
