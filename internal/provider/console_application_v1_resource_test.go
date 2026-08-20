package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const applicationV1PolicyRefMinimumVersion = "v1.35.0"

func TestAccApplicationV1Resource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resourceRef := "conduktor_console_application_v1.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/application_v1/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "my-application"),
					resource.TestCheckResourceAttr(resourceRef, "spec.title", "My Application"),
					// resource.TestCheckResourceAttr(resourceRef, "spec.description", "My Application description"),
					resource.TestCheckResourceAttr(resourceRef, "spec.owner", "admin"),
				),
			},
			// Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "my-application",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/application_v1/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "my-application"),
					resource.TestCheckResourceAttr(resourceRef, "spec.title", "My Application"),
					resource.TestCheckResourceAttr(resourceRef, "spec.description", "My Application description"),
					resource.TestCheckResourceAttr(resourceRef, "spec.owner", "admin"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccApplicationV1PolicyRef(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, applicationV1PolicyRefMinimumVersion)

	resourceRef := "conduktor_console_application_v1.test_policy_ref"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with policy_ref and verify it is stored and returned
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/application_v1/resource_policy_ref.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "my-application-with-policy"),
					resource.TestCheckResourceAttr(resourceRef, "spec.title", "My Application With Policy"),
					resource.TestCheckResourceAttr(resourceRef, "spec.owner", "admin"),
					resource.TestCheckResourceAttr(resourceRef, "spec.policy_ref.#", "1"),
				),
			},
			// Import verifies policy_ref round-trips correctly
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "my-application-with-policy",
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAccApplicationV1ExampleResource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from complex example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_application_v1", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_application_v1.example", "name", "complex-app"),
					resource.TestCheckResourceAttr("conduktor_console_application_v1.example", "spec.title", "Complex Application"),
					resource.TestCheckResourceAttr("conduktor_console_application_v1.example", "spec.description", "Complex Application description"),
					resource.TestCheckResourceAttr("conduktor_console_application_v1.example", "spec.owner", "admin"),
				),
			},
		},
	})
}
