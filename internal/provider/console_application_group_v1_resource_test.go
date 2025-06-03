package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationGroupV1Resource(t *testing.T) {
	resourceRef := "conduktor_console_application_group_v1.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/application_group_v1/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "myappgroup"),
					resource.TestCheckResourceAttr(resourceRef, "application", "myapp"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "My Application Group"),
					resource.TestCheckResourceAttr(resourceRef, "spec.description", "test"),
					resource.TestCheckResourceAttr(resourceRef, "spec.external_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.external_groups.0", "mygroup"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.app_instance", "my-app-instance"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.resource_type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.name", "*"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.permissions.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.permissions.0", "topicViewConfig"),
				),
			},
			// Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "myappgroup",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/application_group_v1/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "myappgroup"),
					resource.TestCheckResourceAttr(resourceRef, "application", "myapp"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "My Updated Application Group"),
					resource.TestCheckResourceAttr(resourceRef, "spec.description", "update test"),
					resource.TestCheckResourceAttr(resourceRef, "spec.external_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.external_groups.0", "mygroup"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.#", "2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.app_instance", "my-app-instance"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.resource_type", "SUBJECT"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.name", "*"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.permissions.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.permissions.0", "subjectView"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.1.app_instance", "my-app-instance"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.1.resource_type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.1.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.1.name", "*"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.1.permissions.#", "2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.1.permissions.0", "topicConsume"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.1.permissions.1", "topicViewConfig"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccApplicationGroupV1Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/application_group_v1/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.minimal", "name", "minimal"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.minimal", "application", "myapp"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.minimal", "spec.display_name", "Minimal Application Group"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccApplicationGroupV1ExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from simple example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_application_group_v1", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.simple", "name", "simple"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.simple", "application", "myapp"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.simple", "spec.display_name", "Simple Application Group"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.simple", "spec.description", "Simple Description"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.simple", "spec.permissions.#", "1"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.simple", "spec.permissions.0.app_instance", "my-app-instance"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.simple", "spec.permissions.0.resource_type", "TOPIC"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.simple", "spec.permissions.0.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.simple", "spec.permissions.0.name", "*"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.simple", "spec.permissions.0.permissions.#", "2"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.simple", "spec.permissions.0.permissions.0", "topicConsume"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.simple", "spec.permissions.0.permissions.1", "topicViewConfig"),
				),
			},
			// Create and Read from complex example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_application_group_v1", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "name", "complex"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "application", "myapp"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.display_name", "Complex Application Group"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.description", "Complex Description"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.#", "3"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.0.app_instance", "my-app-instance"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.0.resource_type", "CONNECTOR"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.0.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.0.name", "*"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.0.connect_cluster", "kafka-connect"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.0.permissions.#", "3"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.0.permissions.0", "kafkaConnectRestart"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.0.permissions.1", "kafkaConnectorStatus"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.0.permissions.2", "kafkaConnectorViewConfig"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.1.app_instance", "my-app-instance"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.1.resource_type", "CONSUMER_GROUP"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.1.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.1.name", "*"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.1.permissions.#", "4"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.1.permissions.0", "consumerGroupCreate"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.1.permissions.1", "consumerGroupDelete"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.1.permissions.2", "consumerGroupReset"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.1.permissions.3", "consumerGroupView"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.2.app_instance", "my-app-instance"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.2.resource_type", "TOPIC"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.2.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.2.name", "*"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.2.permissions.#", "2"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.2.permissions.0", "topicConsume"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.permissions.2.permissions.1", "topicViewConfig"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.members.#", "2"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.members.0", "user1@company.org"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.members.1", "user2@company.org"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.external_groups.#", "1"),
					resource.TestCheckResourceAttr("conduktor_console_application_group_v1.complex", "spec.external_groups.0", "GP-COMPANY-CLICKSTREAM-SUPPORT"),
				),
			},
		},
	})
}
