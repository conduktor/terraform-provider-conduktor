package kafka_cluster_v2

import (
	"context"
	"testing"

	ctlresource "github.com/conduktor/ctl/resource"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestKafkaClusterV2ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonKafkaClusterV2Resource := []byte(test.TestAccTestdata(t, "kafka_cluster_v2_confluent_api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonKafkaClusterV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "KafkaCluster", ctlResource.Kind)
	assert.Equal(t, "v2", ctlResource.Version)
	assert.Equal(t, "cluster-name", ctlResource.Name)
	assert.Equal(t, map[string]interface{}{"name": "cluster-name", "labels": map[string]interface{}{"key1": "value1"}}, ctlResource.Metadata)
	assert.Equal(t, jsonKafkaClusterV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := model.NewKafkaClusterResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "KafkaCluster", internal.Kind)
	assert.Equal(t, "v2", internal.ApiVersion)
	assert.Equal(t, "cluster-name", internal.Metadata.Name)
	assert.Equal(t, "Cluster display name", internal.Spec.DisplayName)
	assert.Equal(t, "localhost:9092", internal.Spec.BootstrapServers)
	assert.Equal(t, "#000000", internal.Spec.Color)
	assert.Equal(t, "kafka", internal.Spec.Icon)
	assert.Equal(t, false, internal.Spec.IgnoreUntrustedCertificate)
	assert.Equal(t, map[string]string{
		"security.protocol": "SASL_SSL",
		"sasl.mechanism":    "PLAIN",
		"sasl.jaas.config":  "org.apache.kafka.common.security.plain.PlainLoginModule required username=\"admin\" password=\"admin-secret\";",
	}, internal.Spec.Properties)
	assert.Equal(t, "http://localhost:8080", internal.Spec.SchemaRegistry.ConfluentLike.Url)
	assert.Equal(t, true, internal.Spec.SchemaRegistry.ConfluentLike.IgnoreUntrustedCertificate)
	assert.Equal(t, "some_user", internal.Spec.SchemaRegistry.ConfluentLike.Security.BasicAuth.UserName)
	assert.Equal(t, "some_password", internal.Spec.SchemaRegistry.ConfluentLike.Security.BasicAuth.Password)
	assert.Equal(t, "key", internal.Spec.KafkaFlavor.Confluent.Key)
	assert.Equal(t, "secret", internal.Spec.KafkaFlavor.Confluent.Secret)
	assert.Equal(t, "env", internal.Spec.KafkaFlavor.Confluent.ConfluentEnvironmentId)
	assert.Equal(t, "cluster", internal.Spec.KafkaFlavor.Confluent.ConfluentClusterId)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("cluster-name"), tfModel.Name)
	assert.Equal(t, types.StringValue("Cluster display name"), tfModel.Spec.DisplayName)
	assert.Equal(t, false, tfModel.Spec.IsNull())
	assert.Equal(t, false, tfModel.Spec.IsUnknown())
	assert.Equal(t, types.StringValue("localhost:9092"), tfModel.Spec.BootstrapServers)
	assert.Equal(t, types.StringValue("#000000"), tfModel.Spec.Color)
	assert.Equal(t, types.StringValue("kafka"), tfModel.Spec.Icon)
	assert.Equal(t, types.BoolValue(false), tfModel.Spec.IgnoreUntrustedCertificate)
	assert.Equal(t, false, tfModel.Spec.KafkaFlavor.IsNull())
	assert.Equal(t, false, tfModel.Spec.KafkaFlavor.IsUnknown())
	assert.Equal(t, false, tfModel.Spec.SchemaRegistry.IsNull())
	assert.Equal(t, false, tfModel.Spec.SchemaRegistry.IsUnknown())

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "KafkaCluster", internal2.Kind)
	assert.Equal(t, "v2", internal2.ApiVersion)
	assert.Equal(t, "cluster-name", internal2.Metadata.Name)
	assert.Equal(t, "Cluster display name", internal2.Spec.DisplayName)
	assert.Equal(t, "localhost:9092", internal2.Spec.BootstrapServers)
	assert.Equal(t, "#000000", internal2.Spec.Color)
	assert.Equal(t, "kafka", internal2.Spec.Icon)
	assert.Equal(t, false, internal2.Spec.IgnoreUntrustedCertificate)
	assert.Equal(t, map[string]string{
		"security.protocol": "SASL_SSL",
		"sasl.mechanism":    "PLAIN",
		"sasl.jaas.config":  "org.apache.kafka.common.security.plain.PlainLoginModule required username=\"admin\" password=\"admin-secret\";",
	}, internal2.Spec.Properties)
	assert.Equal(t, "http://localhost:8080", internal2.Spec.SchemaRegistry.ConfluentLike.Url)
	assert.Equal(t, true, internal2.Spec.SchemaRegistry.ConfluentLike.IgnoreUntrustedCertificate)
	assert.Equal(t, "some_user", internal2.Spec.SchemaRegistry.ConfluentLike.Security.BasicAuth.UserName)
	assert.Equal(t, "some_password", internal2.Spec.SchemaRegistry.ConfluentLike.Security.BasicAuth.Password)
	assert.Equal(t, "key", internal2.Spec.KafkaFlavor.Confluent.Key)
	assert.Equal(t, "secret", internal2.Spec.KafkaFlavor.Confluent.Secret)
	assert.Equal(t, "env", internal2.Spec.KafkaFlavor.Confluent.ConfluentEnvironmentId)
	assert.Equal(t, "cluster", internal2.Spec.KafkaFlavor.Confluent.ConfluentClusterId)
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

func TestAWSKafkaClusterV2ModelMapping(t *testing.T) {
	ctx := context.Background()

	jsonKafkaClusterV2Resource := []byte(test.TestAccTestdata(t, "kafka_cluster_v2_aws_api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonKafkaClusterV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "KafkaCluster", ctlResource.Kind)
	assert.Equal(t, "v2", ctlResource.Version)
	assert.Equal(t, "aws-cluster", ctlResource.Name)
	assert.Equal(t, jsonKafkaClusterV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := model.NewKafkaClusterResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "KafkaCluster", internal.Kind)
	assert.Equal(t, "v2", internal.ApiVersion)
	assert.Equal(t, "aws-cluster", internal.Metadata.Name)
	assert.Equal(t, "MSK Cluster display name", internal.Spec.DisplayName)
	assert.Equal(t, "a-3-public.xxxxx.xxxxx.a1.kafka.eu-west-1.amazonaws.com:9198", internal.Spec.BootstrapServers)
	assert.Equal(t, map[string]string{
		"security.protocol":                  "SASL_SSL",
		"sasl.mechanism":                     "AWS_MSK_IAM",
		"sasl.jaas.config":                   "software.amazon.msk.auth.iam.IAMLoginModule required;",
		"sasl.client.callback.handler.class": "io.conduktor.aws.IAMClientCallbackHandler",
		"aws_access_key_id":                  "XXXXXXXXXX",
		"aws_secret_access_key":              "YYYYYYYYYY",
	}, internal.Spec.Properties)
	assert.Equal(t, "eu-west-1", internal.Spec.SchemaRegistry.Glue.Region)
	assert.Equal(t, "default", internal.Spec.SchemaRegistry.Glue.RegistryName)
	assert.Equal(t, "XXXXXXXXXX", internal.Spec.SchemaRegistry.Glue.Security.Credentials.AccessKeyId)
	assert.Equal(t, "YYYYYYYYYY", internal.Spec.SchemaRegistry.Glue.Security.Credentials.SecretKey)
	assert.Equal(t, (*model.KafkaFlavor)(nil), internal.Spec.KafkaFlavor)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("aws-cluster"), tfModel.Name)
	assert.Equal(t, false, tfModel.Spec.IsNull())
	assert.Equal(t, false, tfModel.Spec.IsUnknown())
	assert.Equal(t, types.StringValue("MSK Cluster display name"), tfModel.Spec.DisplayName)
	assert.Equal(t, types.StringValue("a-3-public.xxxxx.xxxxx.a1.kafka.eu-west-1.amazonaws.com:9198"), tfModel.Spec.BootstrapServers)
	assert.Equal(t, types.BoolValue(false), tfModel.Spec.IgnoreUntrustedCertificate)
	assert.Equal(t, true, tfModel.Spec.KafkaFlavor.IsNull())
	assert.Equal(t, false, tfModel.Spec.KafkaFlavor.IsUnknown())
	assert.Equal(t, false, tfModel.Spec.SchemaRegistry.IsNull())
	assert.Equal(t, false, tfModel.Spec.SchemaRegistry.IsUnknown())

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
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

func TestMinimalKafkaClusterV2ModelMapping(t *testing.T) {
	ctx := context.Background()

	jsonKafkaClusterV2Resource := []byte(test.TestAccTestdata(t, "kafka_cluster_v2_minimal_api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonKafkaClusterV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "KafkaCluster", ctlResource.Kind)
	assert.Equal(t, "v2", ctlResource.Version)
	assert.Equal(t, "cluster-minimal", ctlResource.Name)
	assert.Equal(t, jsonKafkaClusterV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := model.NewKafkaClusterResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "KafkaCluster", internal.Kind)
	assert.Equal(t, "v2", internal.ApiVersion)
	assert.Equal(t, "cluster-minimal", internal.Metadata.Name)
	assert.Equal(t, "Minimal Cluster display name", internal.Spec.DisplayName)
	assert.Equal(t, "localhost:9092", internal.Spec.BootstrapServers)
	assert.Equal(t, "", internal.Spec.Color)
	assert.Equal(t, "", internal.Spec.Icon)
	assert.Equal(t, map[string]string(nil), internal.Spec.Properties)
	assert.Equal(t, (*model.SchemaRegistry)(nil), internal.Spec.SchemaRegistry)
	assert.Equal(t, (*model.KafkaFlavor)(nil), internal.Spec.KafkaFlavor)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, types.StringValue("cluster-minimal"), tfModel.Name)
	assert.Equal(t, false, tfModel.Spec.IsNull())
	assert.Equal(t, false, tfModel.Spec.IsUnknown())
	assert.Equal(t, types.StringValue("Minimal Cluster display name"), tfModel.Spec.DisplayName)
	assert.Equal(t, types.StringValue("localhost:9092"), tfModel.Spec.BootstrapServers)
	assert.Equal(t, types.StringNull(), tfModel.Spec.Color)
	assert.Equal(t, types.StringNull(), tfModel.Spec.Icon)
	assert.Equal(t, types.BoolValue(false), tfModel.Spec.IgnoreUntrustedCertificate)
	assert.Equal(t, true, tfModel.Spec.KafkaFlavor.IsNull())
	assert.Equal(t, true, tfModel.Spec.SchemaRegistry.IsNull())

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
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
