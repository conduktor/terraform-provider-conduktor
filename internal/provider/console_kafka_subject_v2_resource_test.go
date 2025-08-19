package provider

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var schemaValue = "{\"$id\":\"https://mycompany.com/myrecord\",\"$schema\":\"https://json-schema.org/draft/2019-09/schema\",\"type\":\"object\",\"title\":\"MyRecord\",\"description\":\"Json schema for MyRecord\",\"properties\":{\"id\":{\"type\":\"string\"},\"name\":{\"type\":[\"string\",\"null\"]}},\"required\":[\"id\"],\"additionalProperties\":false}"

func TestAccKafkaSubjectV2Resource(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceRef, "spec.id", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.version", "1"),
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
					resource.TestCheckResourceAttr(resourceRef, "spec.id", "2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.version", "2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.schema", schemaValue),
					resource.TestCheckResourceAttr(resourceRef, "spec.references.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.references.0.name", "example-subject.value"),
					resource.TestCheckResourceAttr(resourceRef, "spec.references.0.subject", "example-subject.value"),
					resource.TestCheckResourceAttr(resourceRef, "spec.references.0.version", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase

		},
	})
}

func TestAccKafkaSubjectV2Minimal(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
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

func TestAccKafkaSubjectV2ResourceFileSchema(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	minimalRef := "conduktor_console_kafka_subject_v2.minimal"
	// fileRef := "conduktor_console_kafka_subject_v2.json_file"
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
					testCheckJSONEquality(minimalRef, "spec.schema", schemaValue),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// testCheckJSONEquality compares two JSON strings for equality by parsing them into structs and comparing the structs
func testCheckJSONEquality(resourceName, attributeName, expectedJSON string) resource.TestCheckFunc {
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
