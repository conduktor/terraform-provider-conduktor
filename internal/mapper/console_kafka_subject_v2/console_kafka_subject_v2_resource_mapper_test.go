package console_kafka_subject_v2

import (
	"context"
	"testing"

	ctlresource "github.com/conduktor/ctl/resource"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	subject "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_kafka_subject_v2"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// When managedLabels are added we need to update the test

func TestKafkaSubjectV2ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonKafkaSubjectV2Resource := []byte(test.TestAccTestdata(t, "console/kafka_subject_v2/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonKafkaSubjectV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "Subject", ctlResource.Kind)
	assert.Equal(t, "v2", ctlResource.Version)
	assert.Equal(t, "api-json-example-subject.value", ctlResource.Name)
	assert.Equal(t, map[string]any{"name": "api-json-example-subject.value", "cluster": "kafka-cluster", "labels": map[string]any{"team": "test", "environment": "test"}}, ctlResource.Metadata)
	assert.Equal(t, jsonKafkaSubjectV2Resource, ctlResource.Json)

	internal, err := console.NewKafkaSubjectResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	// convert into internal model
	expectedRefs := []console.KafkaSubjectReferences{
		{
			Name:    "example-subject.value",
			Subject: "example-subject.value",
			Version: 1,
		},
	}
	schemaValue := "{\"$id\":\"https://mycompany.com/myrecord\",\"$schema\":\"https://json-schema.org/draft/2019-09/schema\",\"type\":\"object\",\"title\":\"MyRecord\",\"description\":\"Json schema for MyRecord\",\"properties\":{\"id\":{\"type\":\"string\"},\"name\":{\"type\":[\"string\",\"null\"]}},\"required\":[\"id\"],\"additionalProperties\":false}"

	assert.Equal(t, "Subject", internal.Kind)
	assert.Equal(t, "v2", internal.ApiVersion)
	assert.Equal(t, "api-json-example-subject.value", internal.Metadata.Name)
	assert.Equal(t, "kafka-cluster", internal.Metadata.Cluster)
	assert.Equal(t, "JSON", internal.Spec.Format)
	assert.Equal(t, "BACKWARD", internal.Spec.Compatibility)
	assert.Equal(t, schemaValue, internal.Spec.Schema)
	assert.Equal(t, expectedRefs, internal.Spec.References)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("api-json-example-subject.value"), tfModel.Name)
	assert.Equal(t, types.StringValue("kafka-cluster"), tfModel.Cluster)
	assert.Equal(t, types.StringValue("JSON"), tfModel.Spec.Format)
	assert.Equal(t, types.StringValue("BACKWARD"), tfModel.Spec.Compatibility)
	assert.Equal(t, types.StringValue(schemaValue), tfModel.Spec.Schema)
	var tfRefs []subject.ReferencesValue
	diag := tfModel.Spec.References.ElementsAs(ctx, &tfRefs, false)
	assert.False(t, diag.HasError())
	assert.Equal(t, 1, len(tfRefs))
	if len(tfRefs) > 0 {
		assert.Equal(t, "example-subject.value", tfRefs[0].Name.ValueString())
		assert.Equal(t, "example-subject.value", tfRefs[0].Subject.ValueString())
		assert.Equal(t, int64(1), tfRefs[0].Version.ValueInt64())
	}

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "Subject", internal2.Kind)
	assert.Equal(t, "v2", internal2.ApiVersion)
	assert.Equal(t, "api-json-example-subject.value", internal2.Metadata.Name)
	assert.Equal(t, "kafka-cluster", internal2.Metadata.Cluster)
	assert.Equal(t, "JSON", internal2.Spec.Format)
	assert.Equal(t, "BACKWARD", internal2.Spec.Compatibility)
	assert.Equal(t, schemaValue, internal2.Spec.Schema)
	assert.Equal(t, expectedRefs, internal2.Spec.References)

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
