package provider

import (
	"encoding/json"
	"fmt"
	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"reflect"
	"regexp"
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var schemaValue = "{\"$id\":\"https://mycompany.com/myrecord\",\"$schema\":\"https://json-schema.org/draft/2019-09/schema\",\"additionalProperties\":false,\"description\":\"Json schema for MyRecord\",\"properties\":{\"id\":{\"type\":\"string\"},\"name\":{\"type\":[\"string\",\"null\"]}},\"required\":[\"id\"],\"title\":\"MyRecord\",\"type\":\"object\"}"
var schemaValueUpdate = "{\"$id\":\"https://mycompany.com/myrecord\",\"$schema\":\"https://json-schema.org/draft/2019-09/schema\",\"additionalProperties\":false,\"description\":\"Json schema for MyRecord\",\"properties\":{\"ext_ref\":{\"$ref\":\"https://mycompany.com/example.json\"},\"id\":{\"type\":\"string\"},\"name\":{\"type\":[\"string\",\"null\"]}},\"required\":[\"id\"],\"title\":\"MyRecord\",\"type\":\"object\"}"
var schemaValuePretty = `{
  "$id": "https://mycompany.com/myrecord",
  "$schema": "https://json-schema.org/draft/2019-09/schema",
  "type": "object",
  "title": "MyRecord",
  "description": "Json schema for MyRecord",
  "properties": {
    "id": {
      "type": "string"
    },
    "name": {
      "type": ["string", "null"]
    }
  },
  "required": ["id"],
  "additionalProperties": false
}`

var schemaAvroValuePretty = `{
  "type": "record",
  "name": "MyRecord",
  "namespace": "com.mycompany",
  "fields": [
    {
      "name": "id",
      "type": "long"
    }
  ]
}`

var schemaProtobufValue = `syntax = "proto3";

message MyRecord {
  int32 id = 1;
  string createdAt = 2;
  string name = 3;
}
`

func TestAccKafkaSubjectV2Resource(t *testing.T) {
	checkMinimalVersion(t)
	resourceRef := "conduktor_console_kafka_subject_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/kafka_subject_v2/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "api-json-example-subject.value"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "labels.%", "2"),
					resource.TestCheckResourceAttr(resourceRef, "labels.team", "test"),
					resource.TestCheckResourceAttr(resourceRef, "labels.environment", "test"),
					resource.TestCheckResourceAttr(resourceRef, "spec.format", "JSON"),
					resource.TestCheckResourceAttr(resourceRef, "spec.compatibility", "BACKWARD"),
					resource.TestCheckResourceAttr(resourceRef, "spec.schema", schemaValue),
				),
			},
			//Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "kafka-cluster/api-json-example-subject.value",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/kafka_subject_v2/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "api-json-example-subject.value"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "labels.%", "2"),
					resource.TestCheckResourceAttr(resourceRef, "labels.team", "test"),
					resource.TestCheckResourceAttr(resourceRef, "labels.environment", "test"),
					resource.TestCheckResourceAttr(resourceRef, "spec.format", "JSON"),
					resource.TestCheckResourceAttr(resourceRef, "spec.compatibility", "BACKWARD"),
					resource.TestCheckResourceAttr(resourceRef, "spec.schema", schemaValueUpdate),
					resource.TestCheckResourceAttr(resourceRef, "spec.references.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.references.0.name", "https://mycompany.com/example.json"),
					resource.TestCheckResourceAttr(resourceRef, "spec.references.0.subject", "example-subject.value"),
					resource.TestCheckResourceAttr(resourceRef, "spec.references.0.version", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase

		},
	})
}

func TestAccKafkaSubjectV2Minimal(t *testing.T) {
	checkMinimalVersion(t)
	resourceRef := "conduktor_console_kafka_subject_v2.minimal"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/kafka_subject_v2/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "api-json-example-subject.value"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.format", "JSON"),
					resource.TestCheckResourceAttr(resourceRef, "spec.schema", schemaValue),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccKafkaSubjectV2Errors(t *testing.T) {
	checkMinimalVersion(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerConfigConsole + test.TestAccTestdata(t, "console/kafka_subject_v2/resource_unknown_schema.tf"),
				ExpectError: regexp.MustCompile(`Unknown Schema Format`),
			},
		},
	})
}

func TestAccKafkaSubjectV2ExampleResource(t *testing.T) {
	checkMinimalVersion(t)
	minimalRef := "conduktor_console_kafka_subject_v2.minimal"
	complexRef := "conduktor_console_kafka_subject_v2.complex"
	avroRef := "conduktor_console_kafka_subject_v2.avro"
	protobufRef := "conduktor_console_kafka_subject_v2.protobuf"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_kafka_subject_v2", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(minimalRef, "name", "minimal.value"),
					resource.TestCheckResourceAttr(minimalRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(minimalRef, "spec.format", "JSON"),
					testCheckJSONEquality(minimalRef, "spec.schema", schemaValuePretty),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from complex example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_kafka_subject_v2", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(complexRef, "name", "complex.value"),
					resource.TestCheckResourceAttr(complexRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(complexRef, "labels.%", "2"),
					resource.TestCheckResourceAttr(complexRef, "labels.team", "test"),
					resource.TestCheckResourceAttr(complexRef, "labels.environment", "test"),
					resource.TestCheckResourceAttr(complexRef, "spec.format", "JSON"),
					resource.TestCheckResourceAttr(complexRef, "spec.compatibility", "BACKWARD"),
					testCheckJSONEquality(complexRef, "spec.schema", schemaValuePretty),
					resource.TestCheckResourceAttr(complexRef, "spec.references.#", "1"),
					resource.TestCheckResourceAttr(complexRef, "spec.references.0.name", "example-reference"),
					resource.TestCheckResourceAttr(complexRef, "spec.references.0.subject", "minimal_subject"),
					resource.TestCheckResourceAttr(complexRef, "spec.references.0.version", "1"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from protobuf example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_kafka_subject_v2", "protobuf.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(protobufRef, "name", "protobuf.value"),
					resource.TestCheckResourceAttr(protobufRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(protobufRef, "labels.%", "2"),
					resource.TestCheckResourceAttr(protobufRef, "labels.team", "test"),
					resource.TestCheckResourceAttr(protobufRef, "labels.environment", "test"),
					resource.TestCheckResourceAttr(protobufRef, "spec.format", "PROTOBUF"),
					resource.TestCheckResourceAttr(protobufRef, "spec.compatibility", "BACKWARD"),
					resource.TestCheckResourceAttr(protobufRef, "spec.schema", schemaProtobufValue),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from avro example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_kafka_subject_v2", "avro.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(avroRef, "name", "avro.value"),
					resource.TestCheckResourceAttr(avroRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(avroRef, "labels.%", "2"),
					resource.TestCheckResourceAttr(avroRef, "labels.team", "test"),
					resource.TestCheckResourceAttr(avroRef, "labels.environment", "test"),
					resource.TestCheckResourceAttr(avroRef, "spec.format", "AVRO"),
					resource.TestCheckResourceAttr(avroRef, "spec.compatibility", "FORWARD_TRANSITIVE"),
					testCheckJSONEquality(avroRef, "spec.schema", schemaAvroValuePretty),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

	// TODO basic AVRO of {"name": "id", "type": "long"}
}

func testCheckJSONEquality(resourceName, attributeName, expectedJSON string) resource.TestCheckFunc {
	// This function returns a resource.TestCheckFunc that checks if the value of an attribute named attributeName
	// on a resource named resourceName is equal to the expectedJSON in terms of JSON equality.
	// This function is useful if you want to compare JSON strings in a resource.TestStep.

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		actualJSON := rs.Primary.Attributes[attributeName]

		var expected, actual interface{}

		if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
			return fmt.Errorf("failed to unmarshal expected JSON: %w", err)
		}

		if err := json.Unmarshal([]byte(actualJSON), &actual); err != nil {
			return fmt.Errorf("failed to unmarshal actual JSON: %w", err)
		}

		if !reflect.DeepEqual(expected, actual) {
			return fmt.Errorf("%s: expected %s, got %s", attributeName, expectedJSON, actualJSON)
		}

		return nil
	}
}

func checkMinimalVersion(t *testing.T) {

	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, kafkaSubjectMininumVersion)
}
