package console_kafka_connect_v2

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

func TestKafkaConnectV2ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonKafkaConnectV2Resource := []byte(test.TestAccTestdata(t, "console/kafka_connect_v2/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonKafkaConnectV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "KafkaConnectCluster", ctlResource.Kind)
	assert.Equal(t, "v2", ctlResource.Version)
	assert.Equal(t, "connect-name", ctlResource.Name)
	assert.Equal(t, map[string]interface{}{"name": "connect-name", "cluster": "cluster-name", "labels": map[string]interface{}{"key1": "value1"}}, ctlResource.Metadata)
	assert.Equal(t, jsonKafkaConnectV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := console.NewKafkaConnectResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "KafkaConnectCluster", internal.Kind)
	assert.Equal(t, "v2", internal.ApiVersion)
	assert.Equal(t, "connect-name", internal.Metadata.Name)
	assert.Equal(t, "cluster-name", internal.Metadata.Cluster)
	assert.Equal(t, "Connect 1", internal.Spec.DisplayName)
	assert.Equal(t, "http://localhost:8083", internal.Spec.Urls)
	assert.Equal(t, false, internal.Spec.IgnoreUntrustedCertificate)
	assert.Equal(t, map[string]string{
		"X-PROJECT-HEADER": "value",
		"AnotherHeader":    "test",
	}, internal.Spec.Headers)
	assert.Equal(t, "some_user", internal.Spec.Security.BasicAuth.Username)
	assert.Equal(t, "some_password", internal.Spec.Security.BasicAuth.Password)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("connect-name"), tfModel.Name)
	assert.Equal(t, types.StringValue("Connect 1"), tfModel.Spec.DisplayName)
	assert.Equal(t, false, tfModel.Spec.IsNull())
	assert.Equal(t, false, tfModel.Spec.IsUnknown())
	assert.Equal(t, types.StringValue("http://localhost:8083"), tfModel.Spec.Urls)
	assert.Equal(t, types.BoolValue(false), tfModel.Spec.IgnoreUntrustedCertificate)
	assert.Equal(t, false, tfModel.Spec.Security.IsNull())
	assert.Equal(t, false, tfModel.Spec.Security.IsUnknown())

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "KafkaConnectCluster", internal2.Kind)
	assert.Equal(t, "v2", internal2.ApiVersion)
	assert.Equal(t, "connect-name", internal2.Metadata.Name)
	assert.Equal(t, "Connect 1", internal2.Spec.DisplayName)
	assert.Equal(t, "http://localhost:8083", internal2.Spec.Urls)
	assert.Equal(t, false, internal2.Spec.IgnoreUntrustedCertificate)
	assert.Equal(t, map[string]string{
		"X-PROJECT-HEADER": "value",
		"AnotherHeader":    "test",
	}, internal2.Spec.Headers)
	assert.Equal(t, "some_user", internal2.Spec.Security.BasicAuth.Username)
	assert.Equal(t, "some_password", internal2.Spec.Security.BasicAuth.Password)
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
