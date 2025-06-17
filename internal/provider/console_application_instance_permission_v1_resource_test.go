package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationInstancePermissionV1Resource(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, applicationInstancePermissionMininumVersion)

	resourceRef := "conduktor_console_application_instance_permission_v1.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/application_instance_permission_v1/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "appinstance-permission"),
					resource.TestCheckResourceAttr(resourceRef, "application", "myapp"),
					resource.TestCheckResourceAttr(resourceRef, "app_instance", "my-app-instance"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resource.type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resource.name", "my-topic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resource.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.user_permission", "READ"),
					resource.TestCheckResourceAttr(resourceRef, "spec.service_account_permission", "WRITE"),
					resource.TestCheckResourceAttr(resourceRef, "spec.granted_to", "my-app-instance"),
				),
			},
			// Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "appinstance-permission",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/application_instance_permission_v1/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "appinstance-permission"),
					resource.TestCheckResourceAttr(resourceRef, "application", "myapp"),
					resource.TestCheckResourceAttr(resourceRef, "app_instance", "my-app-instance"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resource.type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resource.name", "my-topic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.resource.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.user_permission", "READ"),
					resource.TestCheckResourceAttr(resourceRef, "spec.service_account_permission", "NONE"),
					resource.TestCheckResourceAttr(resourceRef, "spec.granted_to", "my-app-instance"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccApplicationInstancePermissionV1ExampleResource(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, applicationInstancePermissionMininumVersion)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from complex example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_application_instance_permission_v1", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_application_instance_permission_v1.complex", "name", "complex"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_permission_v1.complex", "application", "myapp"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_permission_v1.complex", "app_instance", "my-app-instance"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_permission_v1.complex", "spec.resource.type", "TOPIC"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_permission_v1.complex", "spec.resource.name", "my-topic"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_permission_v1.complex", "spec.resource.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_permission_v1.complex", "spec.user_permission", "WRITE"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_permission_v1.complex", "spec.service_account_permission", "NONE"),
					resource.TestCheckResourceAttr("conduktor_console_application_instance_permission_v1.complex", "spec.granted_to", "my-app-instance"),
				),
			},
		},
	})
}
