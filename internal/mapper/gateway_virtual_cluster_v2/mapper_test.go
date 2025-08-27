package gateway_virtual_cluster_v2

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

func TestFullVirtualClusterV2KafkaModelMapping(t *testing.T) {
	ctx := context.Background()

	jsonVirtualClusterV2Resource := []byte(test.TestAccTestdata(t, "gateway/virtual_cluster_v2/full_api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonVirtualClusterV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "VirtualCluster", ctlResource.Kind)
	assert.Equal(t, "gateway/v2", ctlResource.Version)
	assert.Equal(t, "vcluster-full", ctlResource.Name)
	assert.Equal(t, jsonVirtualClusterV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := gateway.NewVirtualClusterResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}

	expectedACLs := []gateway.VirtualClusterACL{
		{
			ResourcePattern: gateway.VirtualClusterACLResourcePattern{
				ResourceType: "TOPIC",
				Name:         "topic1",
				PatternType:  "ANY",
			},
			Principal:      "User:username1",
			Host:           "*",
			Operation:      "READ",
			PermissionType: "ALLOW",
		},
		{
			ResourcePattern: gateway.VirtualClusterACLResourcePattern{
				ResourceType: "TOPIC",
				Name:         "topic2",
				PatternType:  "LITERAL",
			},
			Principal:      "User:username2",
			Host:           "*",
			Operation:      "WRITE",
			PermissionType: "DENY",
		},
	}
	expectedClientProperties := map[string]map[string]string{
		"PLAIN": {
			"security.protocol": "SASL_PLAINTEXT",
			"sasl.mechanism":    "PLAIN",
			"sasl.jaas.config":  "org.apache.kafka.common.security.plain.PlainLoginModule required username='{{username}}' password='{{password}}';",
		},
	}

	assert.Equal(t, "VirtualCluster", internal.Kind)
	assert.Equal(t, "gateway/v2", internal.ApiVersion)
	assert.Equal(t, "vcluster-full", internal.Metadata.Name)
	assert.Equal(t, true, internal.Spec.AclEnabled)
	assert.Equal(t, []string{"username1", "username2"}, internal.Spec.SuperUsers)
	assert.Equal(t, "Standard", internal.Spec.Type)
	assert.Equal(t, "kafka:9092", internal.Spec.BootstrapServers)
	assert.Equal(t, expectedClientProperties, internal.Spec.ClientProperties)
	assert.Equal(t, expectedACLs, internal.Spec.Acls)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, types.StringValue("vcluster-full"), tfModel.Name)
	assert.Equal(t, true, tfModel.Spec.AclEnabled.ValueBool())
	assert.Equal(t, false, tfModel.Spec.SuperUsers.IsNull())
	assert.Equal(t, false, tfModel.Spec.SuperUsers.IsUnknown())
	assert.Equal(t, types.StringValue("Standard"), tfModel.Spec.SpecType)
	assert.Equal(t, false, tfModel.Spec.ClientProperties.IsNull())
	assert.Equal(t, false, tfModel.Spec.ClientProperties.IsUnknown())
	assert.Equal(t, false, tfModel.Spec.Acls.IsNull())
	assert.Equal(t, false, tfModel.Spec.Acls.IsUnknown())

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "VirtualCluster", internal2.Kind)
	assert.Equal(t, "gateway/v2", internal2.ApiVersion)
	assert.Equal(t, "vcluster-full", internal2.Metadata.Name)
	assert.Equal(t, true, internal2.Spec.AclEnabled)
	assert.Equal(t, []string{"username1", "username2"}, internal2.Spec.SuperUsers)
	assert.Equal(t, "Standard", internal2.Spec.Type)
	assert.Equal(t, "kafka:9092", internal2.Spec.BootstrapServers)
	assert.Equal(t, expectedClientProperties, internal2.Spec.ClientProperties)
	assert.Equal(t, expectedACLs, internal2.Spec.Acls)

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

func TestMiniVirtualClusterV2KafkaModelMapping(t *testing.T) {
	ctx := context.Background()

	jsonVirtualClusterV2Resource := []byte(test.TestAccTestdata(t, "gateway/virtual_cluster_v2/mini_api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonVirtualClusterV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "VirtualCluster", ctlResource.Kind)
	assert.Equal(t, "gateway/v2", ctlResource.Version)
	assert.Equal(t, "vcluster-mini", ctlResource.Name)
	assert.Equal(t, jsonVirtualClusterV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := gateway.NewVirtualClusterResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, "VirtualCluster", internal.Kind)
	assert.Equal(t, "gateway/v2", internal.ApiVersion)
	assert.Equal(t, "vcluster-mini", internal.Metadata.Name)
	assert.Equal(t, false, internal.Spec.AclEnabled)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, types.StringValue("vcluster-mini"), tfModel.Name)
	assert.Equal(t, types.BoolValue(false), tfModel.Spec.AclEnabled)
	assert.Equal(t, false, tfModel.Spec.SuperUsers.IsNull())
	assert.Equal(t, false, tfModel.Spec.SuperUsers.IsUnknown())
	assert.Equal(t, true, tfModel.Spec.ClientProperties.IsNull())
	assert.Equal(t, false, tfModel.Spec.ClientProperties.IsUnknown())
	assert.Equal(t, true, tfModel.Spec.Acls.IsNull())
	assert.Equal(t, false, tfModel.Spec.Acls.IsUnknown())

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}

	emptyMap := make(map[string]map[string]string)
	assert.Equal(t, "VirtualCluster", internal2.Kind)
	assert.Equal(t, "gateway/v2", internal2.ApiVersion)
	assert.Equal(t, "vcluster-mini", internal2.Metadata.Name)
	assert.Equal(t, false, internal2.Spec.AclEnabled)
	assert.Equal(t, emptyMap, internal2.Spec.ClientProperties)

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
