package console_topic_v2

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

func TestTopicV2ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonTopicV2Resource := []byte(test.TestAccTestdata(t, "/console/topic_v2/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonTopicV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "Topic", ctlResource.Kind)
	assert.Equal(t, "v2", ctlResource.Version)
	assert.Equal(t, "topic", ctlResource.Name)
	expectedLabels := map[string]interface{}{
		"conduktor.io/application":          "test-app",
		"conduktor.io/application-instance": "test-app-instance",
		"kind":                              "topic",
		"data-criticality":                  "C0",
		"environment":                       "prod",
		"team":                              "analytics",
	}
	assert.Equal(t, map[string]interface{}{"name": "topic", "cluster": "cluster", "labels": expectedLabels, "catalogVisibility": "PRIVATE", "descriptionIsEditable": true, "description": "This is a topic", "sqlStorage": map[string]interface{}{"retentionTimeInSecond": float64(86400), "enabled": true}}, ctlResource.Metadata)
	assert.Equal(t, map[string]interface{}{"configs": map[string]interface{}{"cleanup.policy": "delete", "min.insync.replicas": "2", "retention.ms": "60000"}, "partitions": float64(1), "replicationFactor": float64(1)}, ctlResource.Spec)

	assert.Equal(t, jsonTopicV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := console.NewTopicConsoleResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "Topic", internal.Kind)
	assert.Equal(t, "v2", internal.ApiVersion)
	assert.Equal(t, "topic", internal.Metadata.Name)
	assert.Equal(t, "cluster", internal.Metadata.Cluster)
	assert.Equal(t, "topic", internal.Metadata.Labels["kind"])
	assert.Equal(t, "test-app", internal.Metadata.Labels["conduktor.io/application"])
	assert.Equal(t, "test-app-instance", internal.Metadata.Labels["conduktor.io/application-instance"])
	assert.Equal(t, "C0", internal.Metadata.Labels["data-criticality"])
	assert.Equal(t, "prod", internal.Metadata.Labels["environment"])
	assert.Equal(t, "analytics", internal.Metadata.Labels["team"])
	assert.Equal(t, "PRIVATE", internal.Metadata.CatalogVisibility)
	assert.Equal(t, true, internal.Metadata.DescriptionIsEditable)
	assert.Equal(t, "This is a topic", internal.Metadata.Description)
	assert.Equal(t, int64(86400), internal.Metadata.SqlStorage.RetentionTimeInSecond)
	assert.Equal(t, true, internal.Metadata.SqlStorage.Enabled)
	assert.Equal(t, int64(1), internal.Spec.Partitions)
	assert.Equal(t, int64(1), internal.Spec.ReplicationFactor)
	assert.Equal(t, "2", internal.Spec.Configs["min.insync.replicas"])
	assert.Equal(t, "delete", internal.Spec.Configs["cleanup.policy"])
	assert.Equal(t, "60000", internal.Spec.Configs["retention.ms"])

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("topic"), tfModel.Name)
	assert.Equal(t, types.StringValue("cluster"), tfModel.Cluster)
	assert.Equal(t, false, tfModel.Labels.IsNull())
	assert.Equal(t, false, tfModel.Labels.IsUnknown())
	assert.Equal(t, types.StringValue("topic"), tfModel.Labels.Elements()["kind"])
	assert.Equal(t, types.StringValue("C0"), tfModel.Labels.Elements()["data-criticality"])
	assert.Equal(t, types.StringValue("prod"), tfModel.Labels.Elements()["environment"])
	assert.Equal(t, types.StringValue("analytics"), tfModel.Labels.Elements()["team"])
	assert.Equal(t, false, tfModel.ManagedLabels.IsNull())
	assert.Equal(t, false, tfModel.ManagedLabels.IsUnknown())
	assert.Equal(t, types.StringValue("test-app"), tfModel.ManagedLabels.Elements()["conduktor.io/application"])
	assert.Equal(t, types.StringValue("test-app-instance"), tfModel.ManagedLabels.Elements()["conduktor.io/application-instance"])
	assert.Equal(t, types.StringValue("PRIVATE"), tfModel.CatalogVisibility)
	assert.Equal(t, types.BoolValue(true), tfModel.DescriptionIsEditable)
	assert.Equal(t, types.StringValue("This is a topic"), tfModel.Description)
	assert.Equal(t, types.Int64Value(86400), tfModel.SqlStorage.RetentionTimeInSecond)
	assert.Equal(t, types.BoolValue(true), tfModel.SqlStorage.Enabled)
	assert.Equal(t, types.Int64Value(1), tfModel.Spec.Partitions)
	assert.Equal(t, types.Int64Value(1), tfModel.Spec.ReplicationFactor)
	assert.Equal(t, false, tfModel.Spec.Configs.IsNull())
	assert.Equal(t, false, tfModel.Spec.Configs.IsUnknown())

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "Topic", internal2.Kind)
	assert.Equal(t, "v2", internal2.ApiVersion)
	assert.Equal(t, "topic", internal2.Metadata.Name)
	assert.Equal(t, "cluster", internal2.Metadata.Cluster)
	assert.Equal(t, "topic", internal2.Metadata.Labels["kind"])
	assert.Equal(t, "topic", internal2.Metadata.Labels["kind"])
	assert.Equal(t, "test-app", internal2.Metadata.Labels["conduktor.io/application"])
	assert.Equal(t, "test-app-instance", internal2.Metadata.Labels["conduktor.io/application-instance"])
	assert.Equal(t, "C0", internal2.Metadata.Labels["data-criticality"])
	assert.Equal(t, "prod", internal2.Metadata.Labels["environment"])
	assert.Equal(t, "analytics", internal2.Metadata.Labels["team"])
	assert.Equal(t, "PRIVATE", internal2.Metadata.CatalogVisibility)
	assert.Equal(t, true, internal2.Metadata.DescriptionIsEditable)
	assert.Equal(t, "This is a topic", internal2.Metadata.Description)
	assert.Equal(t, int64(86400), internal2.Metadata.SqlStorage.RetentionTimeInSecond)
	assert.Equal(t, true, internal2.Metadata.SqlStorage.Enabled)
	assert.Equal(t, int64(1), internal2.Spec.Partitions)
	assert.Equal(t, int64(1), internal2.Spec.ReplicationFactor)
	assert.Equal(t, "2", internal2.Spec.Configs["min.insync.replicas"])
	assert.Equal(t, "delete", internal2.Spec.Configs["cleanup.policy"])
	assert.Equal(t, "60000", internal2.Spec.Configs["retention.ms"])
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
