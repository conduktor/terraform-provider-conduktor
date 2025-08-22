package console_ksqldb_cluster_v2

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

func TestKsqlDBClusterV2ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonKsqlDBClusterV2Resource := []byte(test.TestAccTestdata(t, "console/ksqldb_cluster_v2/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonKsqlDBClusterV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "KsqlDBCluster", ctlResource.Kind)
	assert.Equal(t, "v2", ctlResource.Version)
	assert.Equal(t, "ksqldb-name", ctlResource.Name)
	assert.Equal(t, map[string]any{"name": "ksqldb-name", "cluster": "cluster-name"}, ctlResource.Metadata)
	assert.Equal(t, jsonKsqlDBClusterV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := console.NewKsqlDBClusterResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "KsqlDBCluster", internal.Kind)
	assert.Equal(t, "v2", internal.ApiVersion)
	assert.Equal(t, "ksqldb-name", internal.Metadata.Name)
	assert.Equal(t, "cluster-name", internal.Metadata.Cluster)
	assert.Equal(t, "KSQL 1", internal.Spec.DisplayName)
	assert.Equal(t, "http://localhost:8088", internal.Spec.Url)
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
	assert.Equal(t, types.StringValue("ksqldb-name"), tfModel.Name)
	assert.Equal(t, types.StringValue("KSQL 1"), tfModel.Spec.DisplayName)
	assert.Equal(t, false, tfModel.Spec.IsNull())
	assert.Equal(t, false, tfModel.Spec.IsUnknown())
	assert.Equal(t, types.StringValue("http://localhost:8088"), tfModel.Spec.Url)
	assert.Equal(t, types.BoolValue(false), tfModel.Spec.IgnoreUntrustedCertificate)
	assert.Equal(t, false, tfModel.Spec.Security.IsNull())
	assert.Equal(t, false, tfModel.Spec.Security.IsUnknown())

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "KsqlDBCluster", internal2.Kind)
	assert.Equal(t, "v2", internal2.ApiVersion)
	assert.Equal(t, "ksqldb-name", internal2.Metadata.Name)
	assert.Equal(t, "KSQL 1", internal2.Spec.DisplayName)
	assert.Equal(t, "http://localhost:8088", internal2.Spec.Url)
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
