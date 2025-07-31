package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
