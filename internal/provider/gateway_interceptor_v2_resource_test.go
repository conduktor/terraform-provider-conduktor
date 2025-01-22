package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

//func TestAccGatewayInterceptorV2Resource(t *testing.T) {
//	test.CheckEnterpriseEnabled(t)
//	resourceRef := "conduktor_gateway_interceptor_v2.test"
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { test.TestAccPreCheck(t) },
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			// Create and Read testing
//			{
//				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway_interceptor_v2_resource_create.tf"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceRef, "name", "enforce-partition-limit"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.safeguard.CreateTopicPolicyPlugin"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.priority", "1"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.config", `{"numPartition":{"action":"INFO","max":5,"min":5},"topic":"myprefix-.*"}`),
//				),
//			},
//			// Importing matches the state of the previous step.
//			// {
//			// 	ResourceName:                         resourceRef,
//			// 	ImportState:                          true,
//			// 	ImportStateVerify:                    true,
//			// 	ImportStateId:                        "user1",
//			// 	ImportStateVerifyIdentifierAttribute: "name",
//			// },
//			// Update and Read testing
//			{
//				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway_interceptor_v2_resource_update.tf"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceRef, "name", "enforce-partition-limit"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.safeguard.CreateTopicPolicyPlugin"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.priority", "100"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.config", `{"numPartition":{"action":"BLOCK","max":10,"min":5},"retentionMs":{"max":100,"min":10},"topic":"updatemyprefix-.*"}`),
//				),
//			},
//			// Delete testing automatically occurs in TestCase
//		},
//	})
//}

func TestAccGatewayInterceptorV2ComplexResource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resourceRef := "conduktor_gateway_interceptor_v2.test-encryption"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway_interceptor_v2_resource_complex_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "myEncryptPlugin"),
					resource.TestCheckResourceAttr(resourceRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.EncryptPlugin"),
					resource.TestCheckResourceAttr(resourceRef, "spec.priority", "1"),
					//resource.TestCheckResourceAttr(resourceRef, "spec.config", `{"numPartition":{"action":"INFO","max":5,"min":5},"topic":"myprefix-.*"}`),
				),
			},
			// Importing matches the state of the previous step.
			// {
			// 	ResourceName:                         resourceRef,
			// 	ImportState:                          true,
			// 	ImportStateVerify:                    true,
			// 	ImportStateId:                        "user1",
			// 	ImportStateVerifyIdentifierAttribute: "name",
			// },
			// Update and Read testing
			//{
			//	Config: providerConfigGateway + test.TestAccTestdata(t, "gateway_interceptor_v2_resource_update.tf"),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr(resourceRef, "name", "enforce-partition-limit"),
			//		resource.TestCheckResourceAttr(resourceRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.safeguard.CreateTopicPolicyPlugin"),
			//		resource.TestCheckResourceAttr(resourceRef, "spec.priority", "100"),
			//		resource.TestCheckResourceAttr(resourceRef, "spec.config", `{"numPartition":{"action":"BLOCK","max":10,"min":5},"retentionMs":{"max":100,"min":10},"topic":"updatemyprefix-.*"}`),
			//	),
			//},
			// Delete testing automatically occurs in TestCase
		},
	})
}
