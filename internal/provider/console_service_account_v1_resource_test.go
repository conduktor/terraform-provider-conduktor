package provider

import (
	"regexp"
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceAccountV1Resource(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, consoleServiceAccountMininumVersion)

	resourceRef := "conduktor_console_service_account_v1.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/service_account_v1/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-service-account"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.0.name", "test-topic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.0.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.0.operations.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.0.operations.0", "Write"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.0.type", "TOPIC"),
				),
			},
			//Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "kafka-cluster/test-service-account",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/service_account_v1/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-service-account"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.#", "2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.0.name", "test-topic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.0.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.0.operations.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.0.operations.0", "Write"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.0.type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.0.host", "*"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.0.permission", "Deny"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.1.name", "test-topic-2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.1.pattern_type", "PREFIXED"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.1.operations.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.1.operations.0", "Write"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.1.type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.1.host", "*"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authorization.kafka.acls.1.permission", "Allow"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccServiceAccountV1Constraints(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, consoleServiceAccountMininumVersion)

	test.CheckEnterpriseEnabled(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Try to create with conflicting security attributes
			{
				Config:      providerConfigConsole + test.TestAccTestdata(t, "console/service_account_v1/resource_not_valid.tf"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestAccServiceAccountV1ExampleResource(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, consoleServiceAccountMininumVersion)

	test.CheckEnterpriseEnabled(t)

	var kafkaResourceRef = "conduktor_console_service_account_v1.kafka_sa"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_service_account_v1", "kafka.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(kafkaResourceRef, "name", "kafka-service-account"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "labels.%", "2"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "labels.domain", "clickstream"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "labels.appcode", "clk"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.#", "4"),

					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.0.name", "*"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.0.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.0.operations.#", "1"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.0.operations.0", "Describe"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.0.type", "TOPIC"),

					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.1.name", "click.event-stream."),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.1.pattern_type", "PREFIXED"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.1.operations.#", "1"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.1.operations.0", "Read"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.1.type", "CONSUMER_GROUP"),

					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.2.name", "click.event-stream.avro"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.2.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.2.operations.#", "2"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.2.operations.0", "Read"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.2.operations.1", "Write"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.2.type", "TOPIC"),

					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.3.name", "public_"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.3.pattern_type", "PREFIXED"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.3.operations.#", "1"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.3.operations.0", "Read"),
					resource.TestCheckResourceAttr(kafkaResourceRef, "spec.authorization.kafka.acls.3.type", "TOPIC"),
				),
			},
		},
	})
}
