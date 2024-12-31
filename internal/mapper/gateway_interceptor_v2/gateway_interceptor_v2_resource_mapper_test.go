package gateway_interceptor_v2

import (
	"context"
	"testing"

	ctlresource "github.com/conduktor/ctl/resource"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestGatewayInterceptorV2ModelMapping(t *testing.T) {
	ctx := context.Background()

	jsonInterceptorV2Resource := []byte(test.TestAccTestdata(t, "gateway_interceptor_v2_api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonInterceptorV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "GatewayInterceptor", ctlResource.Kind)
	assert.Equal(t, "gateway/v2", ctlResource.Version)
	assert.Equal(t, "interceptor", ctlResource.Name)
	assert.Equal(t, map[string]interface{}{"name": "interceptor", "scope": map[string]interface{}{"vCluster": "vcluster", "group": "group", "username": "username"}}, ctlResource.Metadata)
	assert.Equal(t, jsonInterceptorV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := gateway.NewGatewayInterceptorResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "GatewayInterceptor", internal.Kind)
	assert.Equal(t, "gateway/v2", internal.ApiVersion)
	assert.Equal(t, "interceptor", internal.Metadata.Name)
	assert.Equal(t, "vcluster", internal.Metadata.Scope.VCluster)
	assert.Equal(t, "group", internal.Metadata.Scope.Group)
	assert.Equal(t, "username", internal.Metadata.Scope.Username)
	assert.Equal(t, int64(1), internal.Spec.Priority)
	assert.Equal(t, "io.conduktor.gateway.interceptor.VirtualSqlTopicPlugin", internal.Spec.PluginClass)
	// assert.Equal(t, "    {\n      \"virtualTopic\": \"yellow_cars\",\n      \"statement\": \"SELECT '$.type' as type, '$.price' as price FROM cars WHERE '$.color' = 'yellow'\"\n    }\n", internal.Spec.Config)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("interceptor"), tfModel.Name)
	assert.Equal(t, types.StringValue("vcluster"), tfModel.Scope.Vcluster)
	assert.Equal(t, types.StringValue("group"), tfModel.Scope.Group)
	assert.Equal(t, types.StringValue("username"), tfModel.Scope.Username)
	assert.Equal(t, types.Int64Value(1), tfModel.Spec.Priority)
	assert.Equal(t, types.StringValue("io.conduktor.gateway.interceptor.VirtualSqlTopicPlugin"), tfModel.Spec.PluginClass)
	// assert.Equal(t, types.StringValue("    {\n      \"virtualTopic\": \"yellow_cars\",\n      \"statement\": \"SELECT '$.type' as type, '$.price' as price FROM cars WHERE '$.color' = 'yellow'\"\n    }\n"), tfModel.Spec.Config)

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "GatewayInterceptor", internal2.Kind)
	assert.Equal(t, "gateway/v2", internal2.ApiVersion)
	assert.Equal(t, "interceptor", internal2.Metadata.Name)
	assert.Equal(t, "vcluster", internal2.Metadata.Scope.VCluster)
	assert.Equal(t, "group", internal2.Metadata.Scope.Group)
	assert.Equal(t, "username", internal2.Metadata.Scope.Username)
	assert.Equal(t, int64(1), internal2.Spec.Priority)
	// assert.Equal(t, "    {\n      \"virtualTopic\": \"yellow_cars\",\n      \"statement\": \"SELECT '$.type' as type, '$.price' as price FROM cars WHERE '$.color' = 'yellow'\"\n    }\n", internal2.Spec.Config)

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
