package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGatewayVirtualClusterV2Resource(t *testing.T) {
	resourceRef := "conduktor_gateway_virtual_cluster_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/virtual_cluster_v2/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-vcluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acl_enabled", "true"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acl_mode", "REST_API"),
					resource.TestCheckResourceAttr(resourceRef, "spec.type", "Standard"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.0.resource_pattern.resource_type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.0.resource_pattern.name", "test-topic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.0.resource_pattern.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.0.principal", "User:username1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.0.host", "*"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.0.operation", "READ"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.0.permission_type", "ALLOW"),
				),
			},
			// Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "test-vcluster",
				ImportStateVerifyIdentifierAttribute: "name",
				// Ignoring state verify on read only attributes
				ImportStateVerifyIgnore: []string{"spec.bootstrap_servers", "spec.client_properties"},
			},
			// Update and Read testing
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/virtual_cluster_v2/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-vcluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acl_enabled", "true"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acl_mode", "REST_API"),
					resource.TestCheckResourceAttr(resourceRef, "spec.type", "Standard"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.#", "2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.0.resource_pattern.resource_type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.0.resource_pattern.name", "test-topic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.0.resource_pattern.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.0.principal", "User:username1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.0.host", "*"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.0.operation", "READ"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.0.permission_type", "ALLOW"),

					resource.TestCheckResourceAttr(resourceRef, "spec.acls.1.resource_pattern.resource_type", "TOPIC"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.1.resource_pattern.name", "another-topic"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.1.resource_pattern.pattern_type", "PREFIXED"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.1.principal", "User:username2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.1.host", "*"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.1.operation", "WRITE"),
					resource.TestCheckResourceAttr(resourceRef, "spec.acls.1.permission_type", "DENY"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVirtualClusterV2Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/virtual_cluster_v2/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.minimal", "name", "minimal"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVirtualClusterV2ExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from simple example
			{
				Config: providerConfigGateway + test.TestAccExample(t, "resources", "conduktor_gateway_virtual_cluster_v2", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.simple", "name", "simple"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.simple", "spec.acl_enabled", "false"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.simple", "spec.type", "Standard"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.simple", "spec.acl_mode", "KAFKA_API"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.simple", "spec.super_users.#", "1"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.simple", "spec.super_users.0", "user1"),
				),
			},
			// Create and Read from complex example
			{
				Config: providerConfigGateway + test.TestAccExample(t, "resources", "conduktor_gateway_virtual_cluster_v2", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "name", "complex"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acl_enabled", "true"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acl_mode", "REST_API"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.type", "Standard"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.#", "2"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.0.resource_pattern.resource_type", "TOPIC"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.0.resource_pattern.name", "customers"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.0.resource_pattern.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.0.principal", "User:username1"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.0.host", "*"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.0.operation", "READ"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.0.permission_type", "ALLOW"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.1.resource_pattern.resource_type", "TOPIC"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.1.resource_pattern.name", "customers"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.1.resource_pattern.pattern_type", "LITERAL"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.1.principal", "User:username1"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.1.host", "*"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.1.operation", "WRITE"),
					resource.TestCheckResourceAttr("conduktor_gateway_virtual_cluster_v2.complex", "spec.acls.1.permission_type", "ALLOW"),
				),
			},
		},
	})
}
