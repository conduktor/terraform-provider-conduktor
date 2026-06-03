package validation

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func permissionObjType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"resource_type": types.StringType,
			"permissions":   types.SetType{ElemType: types.StringType},
			"name":          types.StringType,
			"pattern_type":  types.StringType,
			"cluster":       types.StringType,
			"kafka_connect": types.StringType,
			"ksqldb":        types.StringType,
		},
	}
}

func makePermissionObj(resourceType string, fields map[string]string) types.Object {
	attrs := map[string]attr.Value{
		"resource_type": types.StringValue(resourceType),
		"permissions":   types.SetValueMust(types.StringType, []attr.Value{types.StringValue("testPerm")}),
		"name":          types.StringNull(),
		"pattern_type":  types.StringNull(),
		"cluster":       types.StringNull(),
		"kafka_connect": types.StringNull(),
		"ksqldb":        types.StringNull(),
	}
	for k, v := range fields {
		attrs[k] = types.StringValue(v)
	}
	return types.ObjectValueMust(permissionObjType().AttrTypes, attrs)
}

func makePermissionSet(perms ...types.Object) basetypes.SetValue {
	elems := make([]attr.Value, len(perms))
	for i, p := range perms {
		elems[i] = p
	}
	return types.SetValueMust(permissionObjType(), elems)
}

func runValidator(t *testing.T, set basetypes.SetValue) validator.SetResponse {
	t.Helper()
	v := PermissionResourceType()
	req := validator.SetRequest{
		ConfigValue: set,
	}
	resp := validator.SetResponse{}
	v.ValidateSet(context.Background(), req, &resp)
	return resp
}

// --- Valid configurations ---

func TestPermissionValidator_PlatformValid(t *testing.T) {
	set := makePermissionSet(makePermissionObj("PLATFORM", nil))
	resp := runValidator(t, set)
	assert.False(t, resp.Diagnostics.HasError(), "PLATFORM with no optional fields should be valid")
}

func TestPermissionValidator_ClusterValid(t *testing.T) {
	set := makePermissionSet(makePermissionObj("CLUSTER", map[string]string{"name": "my-cluster"}))
	resp := runValidator(t, set)
	assert.False(t, resp.Diagnostics.HasError(), "CLUSTER with name should be valid")
}

func TestPermissionValidator_TopicValid(t *testing.T) {
	set := makePermissionSet(makePermissionObj("TOPIC", map[string]string{
		"name": "test-topic", "pattern_type": "LITERAL", "cluster": "my-cluster",
	}))
	resp := runValidator(t, set)
	assert.False(t, resp.Diagnostics.HasError(), "TOPIC with name, pattern_type, cluster should be valid")
}

func TestPermissionValidator_SubjectValid(t *testing.T) {
	set := makePermissionSet(makePermissionObj("SUBJECT", map[string]string{
		"name": "test-subject", "pattern_type": "PREFIXED", "cluster": "my-cluster",
	}))
	resp := runValidator(t, set)
	assert.False(t, resp.Diagnostics.HasError(), "SUBJECT with name, pattern_type, cluster should be valid")
}

func TestPermissionValidator_ConsumerGroupValid(t *testing.T) {
	set := makePermissionSet(makePermissionObj("CONSUMER_GROUP", map[string]string{
		"name": "test-cg", "pattern_type": "LITERAL", "cluster": "my-cluster",
	}))
	resp := runValidator(t, set)
	assert.False(t, resp.Diagnostics.HasError(), "CONSUMER_GROUP with name, pattern_type, cluster should be valid")
}

func TestPermissionValidator_KafkaConnectValid(t *testing.T) {
	set := makePermissionSet(makePermissionObj("KAFKA_CONNECT", map[string]string{
		"name": "*", "pattern_type": "LITERAL", "cluster": "my-cluster", "kafka_connect": "my-connect",
	}))
	resp := runValidator(t, set)
	assert.False(t, resp.Diagnostics.HasError(), "KAFKA_CONNECT with name, pattern_type, cluster, kafka_connect should be valid")
}

func TestPermissionValidator_KsqldbValid(t *testing.T) {
	set := makePermissionSet(makePermissionObj("KSQLDB", map[string]string{
		"cluster": "my-cluster", "ksqldb": "*",
	}))
	resp := runValidator(t, set)
	assert.False(t, resp.Diagnostics.HasError(), "KSQLDB with cluster, ksqldb should be valid")
}

// --- Invalid configurations (the customer's bug scenario) ---

func TestPermissionValidator_KsqldbWithKafkaConnect(t *testing.T) {
	// Customer scenario: KSQLDB with kafka_connect set
	set := makePermissionSet(makePermissionObj("KSQLDB", map[string]string{
		"cluster": "my-cluster", "ksqldb": "*", "kafka_connect": "some-connect",
	}))
	resp := runValidator(t, set)
	assert.True(t, resp.Diagnostics.HasError(), "KSQLDB with kafka_connect should be invalid")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "kafka_connect")
}

func TestPermissionValidator_KsqldbWithNameAndPatternType(t *testing.T) {
	// Customer scenario: KSQLDB with name and pattern_type set (not valid per OpenAPI)
	set := makePermissionSet(makePermissionObj("KSQLDB", map[string]string{
		"name": "*", "cluster": "my-cluster", "pattern_type": "LITERAL", "ksqldb": "*",
	}))
	resp := runValidator(t, set)
	assert.True(t, resp.Diagnostics.HasError(), "KSQLDB with name and pattern_type should be invalid")
	// Should have errors for both name and pattern_type
	assert.GreaterOrEqual(t, len(resp.Diagnostics.Errors()), 2)
}

func TestPermissionValidator_ClusterWithKafkaConnectAndKsqldb(t *testing.T) {
	// Customer scenario: CLUSTER with kafka_connect and ksqldb set
	set := makePermissionSet(makePermissionObj("CLUSTER", map[string]string{
		"name": "my-cluster", "kafka_connect": "some-connect", "ksqldb": "some-ksqldb",
	}))
	resp := runValidator(t, set)
	assert.True(t, resp.Diagnostics.HasError(), "CLUSTER with kafka_connect and ksqldb should be invalid")
	assert.GreaterOrEqual(t, len(resp.Diagnostics.Errors()), 2)
}

func TestPermissionValidator_PlatformWithExtraFields(t *testing.T) {
	set := makePermissionSet(makePermissionObj("PLATFORM", map[string]string{
		"name": "something", "cluster": "something",
	}))
	resp := runValidator(t, set)
	assert.True(t, resp.Diagnostics.HasError(), "PLATFORM with name and cluster should be invalid")
	assert.GreaterOrEqual(t, len(resp.Diagnostics.Errors()), 2)
}

func TestPermissionValidator_TopicWithKafkaConnect(t *testing.T) {
	set := makePermissionSet(makePermissionObj("TOPIC", map[string]string{
		"name": "test", "pattern_type": "LITERAL", "cluster": "my-cluster", "kafka_connect": "some-connect",
	}))
	resp := runValidator(t, set)
	assert.True(t, resp.Diagnostics.HasError(), "TOPIC with kafka_connect should be invalid")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "kafka_connect")
}

func TestPermissionValidator_TopicWithKsqldb(t *testing.T) {
	set := makePermissionSet(makePermissionObj("TOPIC", map[string]string{
		"name": "test", "pattern_type": "LITERAL", "cluster": "my-cluster", "ksqldb": "some-ksqldb",
	}))
	resp := runValidator(t, set)
	assert.True(t, resp.Diagnostics.HasError(), "TOPIC with ksqldb should be invalid")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "ksqldb")
}

func TestPermissionValidator_KafkaConnectWithKsqldb(t *testing.T) {
	set := makePermissionSet(makePermissionObj("KAFKA_CONNECT", map[string]string{
		"name": "*", "pattern_type": "LITERAL", "cluster": "c", "kafka_connect": "kc", "ksqldb": "ks",
	}))
	resp := runValidator(t, set)
	assert.True(t, resp.Diagnostics.HasError(), "KAFKA_CONNECT with ksqldb should be invalid")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "ksqldb")
}

// --- Multiple permissions in same set ---

func TestPermissionValidator_MultiplePermissions_MixedValidity(t *testing.T) {
	set := makePermissionSet(
		// Valid TOPIC
		makePermissionObj("TOPIC", map[string]string{
			"name": "test", "pattern_type": "LITERAL", "cluster": "my-cluster",
		}),
		// Invalid KSQLDB (has kafka_connect)
		makePermissionObj("KSQLDB", map[string]string{
			"cluster": "my-cluster", "ksqldb": "*", "kafka_connect": "bad",
		}),
	)
	resp := runValidator(t, set)
	assert.True(t, resp.Diagnostics.HasError())
	// Only one error (for the KSQLDB permission)
	assert.Equal(t, 1, len(resp.Diagnostics.Errors()))
}

// --- Null and unknown values ---

func TestPermissionValidator_NullSet(t *testing.T) {
	set := types.SetNull(permissionObjType())
	resp := runValidator(t, set)
	assert.False(t, resp.Diagnostics.HasError(), "null set should be valid")
}

func TestPermissionValidator_UnknownSet(t *testing.T) {
	set := types.SetUnknown(permissionObjType())
	resp := runValidator(t, set)
	assert.False(t, resp.Diagnostics.HasError(), "unknown set should be valid")
}

func TestPermissionValidator_EmptySet(t *testing.T) {
	set := types.SetValueMust(permissionObjType(), []attr.Value{})
	resp := runValidator(t, set)
	assert.False(t, resp.Diagnostics.HasError(), "empty set should be valid")
}

// --- Description ---

func TestPermissionValidator_Description(t *testing.T) {
	v := PermissionResourceType()
	desc := v.Description(context.Background())
	assert.NotEmpty(t, desc)
}
