package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourcePolicyV1Resource(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, resourcePolicyMininumVersion)

	resourceRef := "conduktor_console_resource_policy_v1.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/resource_policy_v1/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "resourcepolicy"),
					resource.TestCheckResourceAttr(resourceRef, "labels.label1", "value1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.target_kind", "Topic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.description", "This is a test resource policy"),
					resource.TestCheckResourceAttr(resourceRef, "spec.rules.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.rules.0.condition", "spec.replicationFactor == 3"),
					resource.TestCheckResourceAttr(resourceRef, "spec.rules.0.error_message", "replication factor should be 3"),
				),
			},
			// Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "resourcepolicy",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/resource_policy_v1/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "resourcepolicy"),
					resource.TestCheckResourceAttr(resourceRef, "labels.label1", "value1"),
					resource.TestCheckResourceAttr(resourceRef, "labels.label2", "value2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.target_kind", "Topic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.description", "This is an updated test resource policy"),
					resource.TestCheckResourceAttr(resourceRef, "spec.rules.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.rules.0.condition", "spec.replicationFactor == 3"),
					resource.TestCheckResourceAttr(resourceRef, "spec.rules.0.error_message", "replication factor should be 3"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccResourcePolicyV1Minimal(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, resourcePolicyMininumVersion)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/resource_policy_v1/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.minimal", "name", "minimal"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.minimal", "spec.target_kind", "Connector"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.minimal", "spec.rules.#", "1"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.minimal", "spec.rules.0.condition", "spec.replicationFactor == 3"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.minimal", "spec.rules.0.error_message", "replication factor should be 3"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccResourcePolicyV1ExampleResource(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, resourcePolicyMininumVersion)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from simple example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_resource_policy_v1", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.simple", "name", "simple"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.simple", "labels.business-unit", "delivery"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.simple", "spec.target_kind", "Topic"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.simple", "spec.description", "A policy to check some basic rule for a topic"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.simple", "spec.rules.#", "1"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.simple", "spec.rules.0.condition", "int(string(spec.configs[\"retention.ms\"])) >= 60000 && int(string(spec.configs[\"retention.ms\"])) <= 3600000"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.simple", "spec.rules.0.error_message", "retention should be between 1m and 1h"),
				),
			},
			// Create and Read from complex example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_resource_policy_v1", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.complex", "name", "complex"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.complex", "labels.business-unit", "delivery"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.complex", "spec.target_kind", "Topic"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.complex", "spec.description", "A policy to check some basic rule for a topic"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.complex", "spec.rules.#", "2"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.complex", "spec.rules.0.condition", "metadata.labels[\"data-criticality\"] in [\"C0\", \"C1\", \"C2\"]"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.complex", "spec.rules.0.error_message", "data-criticality should be one of C0, C1, C2"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.complex", "spec.rules.1.condition", "metadata.name.matches(\"^click\\\\.[a-z0-9-]+\\\\.(avro|json)$\")"),
					resource.TestCheckResourceAttr("conduktor_console_resource_policy_v1.complex", "spec.rules.1.error_message", "topic name should match ^click\\.(?<event>[a-z0-9-]+)\\.(avro|json)$"),
				),
			},
		},
	})
}
