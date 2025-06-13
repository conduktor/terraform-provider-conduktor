package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationInstanceV1Resource(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, appInstanceMininumVersion)

	resourceRef := "conduktor_console_application_instance_v1.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/application_instance_v1/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "appinstance"),
					resource.TestCheckResourceAttr(resourceRef, "application", "myapp"),
					resource.TestCheckResourceAttr(resourceRef, "spec.cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.0.type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.0.name", "mytopic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.0.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.application_managed_service_account", "false"),
					resource.TestCheckResourceAttr(resourceRef, "spec.service_account", "my-service-account"),
				),
			},
			// Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "appinstance",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/application_instance_v1/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "appinstance"),
					resource.TestCheckResourceAttr(resourceRef, "application", "myapp"),
					resource.TestCheckResourceAttr(resourceRef, "spec.cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.#", "2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.0.type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.0.name", "mytopic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.0.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.1.type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.1.name", "mytopic2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.1.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.application_managed_service_account", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// This test contains a resource definition for creating an application instance targeting specifically Conduktor Console v1.34.
// Currently used to test the new `spec.policy_ref` field.
func TestAccApplicationInstanceV1Resource2(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, "v1.34.0")

	resourceRef := "conduktor_console_application_instance_v1.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/application_instance_v1/resource_create_2.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "appinstance"),
					resource.TestCheckResourceAttr(resourceRef, "application", "myapp"),
					resource.TestCheckResourceAttr(resourceRef, "spec.cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.policy_ref.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.policy_ref.0", "resource-policy"),
					resource.TestCheckResourceAttr(resourceRef, "spec.application_managed_service_account", "false"),
					resource.TestCheckResourceAttr(resourceRef, "spec.service_account", "my-service-account"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccApplicationInstanceV1Minimal(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, appInstanceMininumVersion)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/application_instance_v1/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.minimal", "name", "minimal"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.minimal", "application", "myapp"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.minimal", "spec.cluster", "kafka-cluster"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccApplicationInstanceV1ExampleResource(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, appInstanceMininumVersion)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from simple example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_application_instance_v1", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.simple", "name", "simple"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.simple", "application", "myapp"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.simple", "spec.cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.simple", "spec.resources.#", "1"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.simple", "spec.resources.0.type", "TOPIC"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.simple", "spec.resources.0.name", "topic"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.simple", "spec.resources.0.pattern_type", "PREFIXED"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.simple", "spec.application_managed_service_account", "false"),
				),
			},
			// Create and Read from complex example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_application_instance_v1", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "name", "complex"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "application", "myapp"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.service_account", "my-service-account"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.topic_policy_ref.#", "1"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.topic_policy_ref.0", "topic-policy"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.default_catalog_visibility", "PUBLIC"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.#", "5"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.0.type", "CONNECTOR"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.0.connect_cluster", "kafka-connect"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.0.name", "click."),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.0.pattern_type", "PREFIXED"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.1.type", "CONSUMER_GROUP"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.1.name", "click."),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.1.pattern_type", "PREFIXED"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.2.type", "SUBJECT"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.2.name", "click."),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.2.pattern_type", "PREFIXED"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.3.type", "TOPIC"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.3.name", "click."),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.3.pattern_type", "PREFIXED"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.4.type", "TOPIC"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.4.name", "legacy-click."),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.4.pattern_type", "PREFIXED"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.resources.4.ownership_mode", "LIMITED"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.complex", "spec.application_managed_service_account", "false"),
				),
			},
		},
	})
}
