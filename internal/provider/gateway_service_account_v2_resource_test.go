package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGatewayServiceAccountV2Resource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resourceRef := "conduktor_gateway_service_account_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + test.TestAccTestdata(t, "gateway_service_account_v2_resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "user1"),
					resource.TestCheckResourceAttr(resourceRef, "vcluster", "passthrough"),
					resource.TestCheckResourceAttr(resourceRef, "spec.type", "EXTERNAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.external_names.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.external_names.0", "externalName"),
				),
			},
			// Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "user1",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfig + test.TestAccTestdata(t, "gateway_service_account_v2_resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "user1"),
					resource.TestCheckResourceAttr(resourceRef, "vcluster", "passthrough"),
					resource.TestCheckResourceAttr(resourceRef, "spec.type", "EXTERNAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.external_names.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.external_names.0", "newExternalName"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccGatewayServiceAccountV2Minimal(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfig + test.TestAccTestdata(t, "gateway_service_account_v2_resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.minimal", "name", "minimal"),
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.minimal", "vcluster", "passthrough"),
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.minimal", "spec.type", "LOCAL"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccGatewayServiceAccountV2ExampleResource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from simple example
			{
				Config: providerConfig + test.TestAccExample(t, "resources", "conduktor_gateway_service_account_v2", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.example", "name", "simple-service-account"),
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.example", "vcluster", "passthrough"),
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.example", "spec.type", "LOCAL"),
				),
			},
			// Create and Read from complex example
			{
				Config: providerConfig + test.TestAccExample(t, "resources", "conduktor_gateway_service_account_v2", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.example", "name", "complex-service-account"),
					// TODO: Add vcluster tests - Needs gateway_vclusters to be deployed by terraform
					// resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.example", "vcluster", "vcluster1"),
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.example", "spec.type", "EXTERNAL"),
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.example", "spec.external_names.#", "1"),
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.example", "spec.external_names.0", "externalName"),
				),
			},
		},
	})
}
