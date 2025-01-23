package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGenericResource(t *testing.T) {
	resourceRef := "conduktor_generic.embedded"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create embedded and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "generic_resource_create_embedded.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "jim.halpert@dunder.mifflin.com"),
					resource.TestCheckResourceAttr(resourceRef, "kind", "User"),
					resource.TestCheckResourceAttr(resourceRef, "version", "v2"),
					resource.TestCheckResourceAttrWith(resourceRef, "manifest",
						test.TestCheckResourceAttrContainsStringsFunc("\"name\": \"jim.halpert@dunder.mifflin.com\"", "\"firstName\": \"Jim\"", "\"lastName\": \"Halpert\"")),
				),
			},
			// Don't support import on generic resources yet
			// ImportState testing
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "generic_resource_update_embedded.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "jim.halpert@dunder.mifflin.com"),
					resource.TestCheckResourceAttr(resourceRef, "kind", "User"),
					resource.TestCheckResourceAttr(resourceRef, "version", "v2"),
					resource.TestCheckResourceAttrWith(resourceRef, "manifest",
						test.TestCheckResourceAttrContainsStringsFunc("\"name\": \"jim.halpert@dunder.mifflin.com\"", "\"firstName\": \"Tim\"", "\"lastName\": \"Canterbury\"")),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccGenericExample2Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_generic", "embedded.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_generic.example", "name", "bob@company.io"),
					resource.TestCheckResourceAttrWith("conduktor_generic.example", "manifest",
						test.TestCheckResourceAttrContainsStringsFunc(
							"\"name\": \"bob@company.io\"",
							"\"firstName\": \"Bob\"",
							"\"lastName\": \"Smith\"",
						)),
				),
			},
		},
	})
}
