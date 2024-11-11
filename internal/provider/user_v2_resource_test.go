package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserV2Resource(t *testing.T) {
	resourceRef := "conduktor_user_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + test.TestAccTestdata(t, "user_v2_resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "pam.beesly@dunder.mifflin.com"),
					resource.TestCheckResourceAttr(resourceRef, "spec.firstname", "Pam"),
					resource.TestCheckResourceAttr(resourceRef, "spec.lastname", "Beesly"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.resource_type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.name", "team1.test-topic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.cluster", "*"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.permissions.#", "3"),
				),
			},
			// Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "pam.beesly@dunder.mifflin.com",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfig + test.TestAccTestdata(t, "user_v2_resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "pam.beesly@dunder.mifflin.com"),
					resource.TestCheckResourceAttr(resourceRef, "spec.firstname", "Pam"),
					resource.TestCheckResourceAttr(resourceRef, "spec.lastname", "Halpert"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.#", "2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.resource_type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.name", "team1.test-topic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.cluster", "*"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.1.resource_type", "PLATFORM"),
					resource.TestCheckResourceAttr(resourceRef, "spec.permissions.1.permissions.#", "4"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccUserV2Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfig + test.TestAccTestdata(t, "user_v2_resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_user_v2.minimal", "name", "angela.martin@dunder-mifflin.com"),
					resource.TestCheckResourceAttr("conduktor_user_v2.minimal", "spec.permissions.#", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccUserV2ExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from example
			{
				Config: providerConfig + test.TestAccExample(t, "resources", "conduktor_user_v2", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "name", "bob@company.io"),
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "spec.firstname", "Bob"),
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "spec.lastname", "Smith"),
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "spec.permissions.#", "0"),
				),
			},
			// Create and Read from example
			{
				Config: providerConfig + test.TestAccExample(t, "resources", "conduktor_user_v2", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "name", "bob@company.io"),
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "spec.firstname", "Bob"),
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "spec.lastname", "Smith"),
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "spec.permissions.#", "2"),
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "spec.permissions.0.resource_type", "TOPIC"),
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "spec.permissions.0.name", "test-topic"),
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "spec.permissions.0.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "spec.permissions.0.cluster", "*"),
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "spec.permissions.0.permissions.#", "3"),
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "spec.permissions.1.resource_type", "PLATFORM"),
					resource.TestCheckResourceAttr("conduktor_user_v2.example", "spec.permissions.1.permissions.#", "3"),
				),
			},
		},
	})
}
