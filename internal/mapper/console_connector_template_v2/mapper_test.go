package console_connector_template_v2

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

func TestConnectorTemplateV2ModelMapping(t *testing.T) {
	ctx := context.Background()

	jsonResource := []byte(test.TestAccTestdata(t, "console/connector_template_v2/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonResource)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "ConnectorTemplate", ctlResource.Kind)
	assert.Equal(t, "v2", ctlResource.Version)
	assert.Equal(t, "my-connector-template", ctlResource.Name)

	// internal model
	internal, err := console.NewConnectorTemplateResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "ConnectorTemplate", internal.Kind)
	assert.Equal(t, "v2", internal.ApiVersion)
	assert.Equal(t, "my-connector-template", internal.Metadata.Name)
	assert.Equal(t, map[string]string{"env": "prod", "conduktor.io/managed": "true"}, internal.Metadata.Labels)
	assert.Equal(t, "My Connector Template", internal.Spec.DisplayName)
	assert.Equal(t, "A template for S3 sink connectors", internal.Spec.Description)
	assert.Equal(t, "connector-{{team}}", internal.Spec.Defaults.Metadata.Name)
	assert.Equal(t, map[string]string{"managed-by": "terraform"}, internal.Spec.Defaults.Metadata.Labels)
	assert.Equal(t, "io.confluent.connect.s3.S3SinkConnector", internal.Spec.Defaults.Spec.Class)
	assert.Equal(t, map[string]string{"tasks.max": "1", "s3.bucket.name": "my-bucket"}, internal.Spec.Defaults.Spec.Config)

	// internal → TF
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, types.StringValue("my-connector-template"), tfModel.Name)
	assert.False(t, tfModel.Labels.IsNull())
	assert.False(t, tfModel.ManagedLabels.IsNull())
	assert.False(t, tfModel.Spec.IsNull())
	assert.False(t, tfModel.Spec.IsUnknown())
	assert.Equal(t, types.StringValue("My Connector Template"), tfModel.Spec.DisplayName)

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
