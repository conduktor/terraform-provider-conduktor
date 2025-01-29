package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGatewayTokenV2Resource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resourceRef := "conduktor_gateway_token_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway_token_v2_resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "vcluster", "vcluster_sa"),
					resource.TestCheckResourceAttr(resourceRef, "username", "user11"),
				),
			},
			// Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "user10",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			// {
			// 	Config: providerConfigGateway + test.TestAccTestdata(t, "gateway_token_v2_resource_update.tf"),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr(resourceRef, "name", "user1"),
			// 		resource.TestCheckResourceAttr(resourceRef, "vcluster", "vcluster_sa"),
			// 		resource.TestCheckResourceAttr(resourceRef, "spec.type", "EXTERNAL"),
			// 		resource.TestCheckResourceAttr(resourceRef, "spec.external_names.#", "1"),
			// 		resource.TestCheckResourceAttr(resourceRef, "spec.external_names.0", "newExternalName"),
			// 	),
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}

// func TestAccGatewayTokenV2Minimal(t *testing.T) {
// 	test.CheckEnterpriseEnabled(t)
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { test.TestAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Create and Read from minimal example
// 			{
// 				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway_token_v2_resource_minimal.tf"),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.minimal", "name", "minimal"),
// 					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.minimal", "vcluster", "passthrough"),
// 					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.minimal", "spec.type", "LOCAL"),
// 				),
// 			},
// 			// Delete testing automatically occurs in TestCase
// 		},
// 	})
// }
//
// func TestAccGatewayTokenV2ExampleResource(t *testing.T) {
// 	test.CheckEnterpriseEnabled(t)
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { test.TestAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//
// 		Steps: []resource.TestStep{
// 			// Create and Read from simple example
// 			{
// 				Config: providerConfigGateway + test.TestAccExample(t, "resources", "conduktor_gateway_token_v2", "simple.tf"),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.local_sa", "name", "simple-service-account"),
// 					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.local_sa", "vcluster", "passthrough"),
// 					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.local_sa", "spec.type", "LOCAL"),
// 				),
// 			},
// 			// Create and Read from complex example
// 			{
// 				Config: providerConfigGateway + test.TestAccExample(t, "resources", "conduktor_gateway_token_v2", "complex.tf"),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.external_sa", "name", "complex-service-account"),
// 					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.external_sa", "vcluster", "vcluster_sa"),
// 					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.external_sa", "spec.type", "EXTERNAL"),
// 					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.external_sa", "spec.external_names.#", "1"),
// 					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.external_sa", "spec.external_names.0", "externalName"),
// 				),
// 			},
// 		},
// 	})
// }
