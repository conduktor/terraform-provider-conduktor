package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationInstanceV1Resource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
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
					resource.TestCheckResourceAttr(resourceRef, "spec.cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.0.type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.0.name", "mytopic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.0.pattern", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.application_managed_service_account", "true"),
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
					resource.TestCheckResourceAttr(resourceRef, "spec.cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.#", "2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.0.type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.0.name", "mytopic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.0.pattern", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.1.type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.1.name", "mytopic2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resources.1.pattern", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.application_managed_service_account", "true"),
					resource.TestCheckResourceAttr(resourceRef, "spec.service_account", "my-service-account"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccApplicationInstanceV1Minimal(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
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

// func TestAccApplicationInstanceV1ExampleResource(t *testing.T) {
// 	test.CheckEnterpriseEnabled(t)
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { test.TestAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//
// 		Steps: []resource.TestStep{
// 			// Create and Read from simple example
// 			{
// 				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_application_instance_v1", "simple.tf"),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "name", "simple-group"),
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "spec.display_name", "Simple Group"),
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "spec.description", "Simple group description"),
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "spec.external_groups.#", "0"),
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "spec.members.#", "0"),
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "spec.members_from_external_groups.#", "0"),
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "spec.permissions.#", "0"),
// 				),
// 			},
// 			// Create and Read from complex example
// 			{
// 				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_application_instance_v1", "complex.tf"),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "name", "complex-group"),
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "spec.display_name", "Complex group"),
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "spec.description", "Complex group description"),
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "spec.external_groups.#", "1"),
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "spec.external_groups.0", "sso-group1"),
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "spec.members.#", "1"),
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "spec.members.0", "user1@company.com"),
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "spec.members_from_external_groups.#", "0"),
// 					resource.TestCheckResourceAttr("conduktor_console_application_instance_v1.example", "spec.permissions.#", "2"),
// 				),
// 			},
// 		},
// 	})
// }
