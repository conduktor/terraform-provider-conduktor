package console_connector_v2

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

func TestConnectorV2ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonConnectorV2Resource := []byte(test.TestAccTestdata(t, "/console/connector_v2/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonConnectorV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "Connector", ctlResource.Kind)
	assert.Equal(t, "v2", ctlResource.Version)
	assert.Equal(t, "connector", ctlResource.Name)
	expectedLabels := map[string]any{
		"conduktor.io/application":          "test-app",
		"conduktor.io/application-instance": "test-app-instance",
		"kind":                              "connector",
		"data-criticality":                  "C0",
		"environment":                       "prod",
		"team":                              "analytics",
	}
	assert.Equal(t, map[string]any{"name": "connector", "cluster": "cluster", "connectCluster": "connect", "labels": expectedLabels, "description": "This is a connector", "autoRestart": map[string]any{"frequencySeconds": float64(500), "enabled": true}}, ctlResource.Metadata)
	config := map[string]any{
		"connector.class":                    "io.connect.jdbc.JdbcSourceConnector",
		"tasks.max":                          "1",
		"topic":                              "click.pageviews",
		"connection.url":                     "jdbc:mysql://127.0.0.1:3306/sample?verifyServerCertificate=false&useSSL=true&requireSSL=true",
		"consumer.override.sasl.jaas.config": "o.a.k.s.s.ScramLoginModule required username='<user>' password='<password>';",
	}
	assert.Equal(t, map[string]any{"config": config}, ctlResource.Spec)

	assert.Equal(t, jsonConnectorV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := console.NewConnectorConsoleResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "Connector", internal.Kind)
	assert.Equal(t, "v2", internal.ApiVersion)
	assert.Equal(t, "connector", internal.Metadata.Name)
	assert.Equal(t, "cluster", internal.Metadata.Cluster)
	assert.Equal(t, "connector", internal.Metadata.Labels["kind"])
	assert.Equal(t, "test-app", internal.Metadata.Labels["conduktor.io/application"])
	assert.Equal(t, "test-app-instance", internal.Metadata.Labels["conduktor.io/application-instance"])
	assert.Equal(t, "C0", internal.Metadata.Labels["data-criticality"])
	assert.Equal(t, "prod", internal.Metadata.Labels["environment"])
	assert.Equal(t, "analytics", internal.Metadata.Labels["team"])
	assert.Equal(t, "This is a connector", internal.Metadata.Description)
	assert.Equal(t, int64(500), internal.Metadata.AutoRestart.FrequencySeconds)
	assert.Equal(t, true, internal.Metadata.AutoRestart.Enabled)
	assert.Equal(t, "io.connect.jdbc.JdbcSourceConnector", internal.Spec.Config["connector.class"])
	assert.Equal(t, "1", internal.Spec.Config["tasks.max"])
	assert.Equal(t, "click.pageviews", internal.Spec.Config["topic"])
	assert.Equal(t, "jdbc:mysql://127.0.0.1:3306/sample?verifyServerCertificate=false&useSSL=true&requireSSL=true", internal.Spec.Config["connection.url"])
	assert.Equal(t, "o.a.k.s.s.ScramLoginModule required username='<user>' password='<password>';", internal.Spec.Config["consumer.override.sasl.jaas.config"])

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("connector"), tfModel.Name)
	assert.Equal(t, types.StringValue("cluster"), tfModel.Cluster)
	assert.Equal(t, false, tfModel.Labels.IsNull())
	assert.Equal(t, false, tfModel.Labels.IsUnknown())
	assert.Equal(t, types.StringValue("connector"), tfModel.Labels.Elements()["kind"])
	assert.Equal(t, types.StringValue("C0"), tfModel.Labels.Elements()["data-criticality"])
	assert.Equal(t, types.StringValue("prod"), tfModel.Labels.Elements()["environment"])
	assert.Equal(t, types.StringValue("analytics"), tfModel.Labels.Elements()["team"])
	assert.Equal(t, false, tfModel.ManagedLabels.IsNull())
	assert.Equal(t, false, tfModel.ManagedLabels.IsUnknown())
	assert.Equal(t, types.StringValue("test-app"), tfModel.ManagedLabels.Elements()["conduktor.io/application"])
	assert.Equal(t, types.StringValue("test-app-instance"), tfModel.ManagedLabels.Elements()["conduktor.io/application-instance"])
	assert.Equal(t, types.StringValue("This is a connector"), tfModel.Description)
	assert.Equal(t, types.Int64Value(500), tfModel.AutoRestart.FrequencySeconds)
	assert.Equal(t, types.BoolValue(true), tfModel.AutoRestart.Enabled)
	assert.Equal(t, false, tfModel.Spec.Config.IsNull())
	assert.Equal(t, false, tfModel.Spec.Config.IsUnknown())

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "Connector", internal2.Kind)
	assert.Equal(t, "v2", internal2.ApiVersion)
	assert.Equal(t, "connector", internal2.Metadata.Name)
	assert.Equal(t, "cluster", internal2.Metadata.Cluster)
	assert.Equal(t, "connector", internal2.Metadata.Labels["kind"])
	assert.Equal(t, "test-app", internal2.Metadata.Labels["conduktor.io/application"])
	assert.Equal(t, "test-app-instance", internal2.Metadata.Labels["conduktor.io/application-instance"])
	assert.Equal(t, "C0", internal2.Metadata.Labels["data-criticality"])
	assert.Equal(t, "prod", internal2.Metadata.Labels["environment"])
	assert.Equal(t, "analytics", internal2.Metadata.Labels["team"])
	assert.Equal(t, "This is a connector", internal2.Metadata.Description)
	assert.Equal(t, int64(500), internal2.Metadata.AutoRestart.FrequencySeconds)
	assert.Equal(t, true, internal2.Metadata.AutoRestart.Enabled)
	assert.Equal(t, "io.connect.jdbc.JdbcSourceConnector", internal2.Spec.Config["connector.class"])
	assert.Equal(t, "1", internal2.Spec.Config["tasks.max"])
	assert.Equal(t, "click.pageviews", internal2.Spec.Config["topic"])
	assert.Equal(t, "jdbc:mysql://127.0.0.1:3306/sample?verifyServerCertificate=false&useSSL=true&requireSSL=true", internal2.Spec.Config["connection.url"])
	assert.Equal(t, "o.a.k.s.s.ScramLoginModule required username='<user>' password='<password>';", internal2.Spec.Config["consumer.override.sasl.jaas.config"])

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
