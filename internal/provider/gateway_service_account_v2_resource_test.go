package provider

import (
	"context"
	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"testing"
	"time"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGatewayServiceAccountV2Resource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resourceRef := "conduktor_gateway_service_account_v2.test"

	gwClient, err := testClient(client.GATEWAY)
	if err != nil {
		t.Fatalf("Error creating gateway client: %s", err)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/service_account_v2/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-sa"),
					resource.TestCheckResourceAttr(resourceRef, "vcluster", "vcluster_sa"),
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
				ImportStateId:                        "test-sa/vcluster_sa",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Test plan changes if externally deleted resource
			{
				PreConfig: func() {
					// wait a bit to ensure the service account is created
					time.Sleep(1 * time.Second)
					deleteRes := gateway.GatewayServiceAccountMetadata{
						Name:     "test-sa",
						VCluster: "vcluster_sa",
					}
					t.Logf("Deleting service account %s in vcluster %s", deleteRes.Name, deleteRes.VCluster)
					err := gwClient.Delete(context.Background(), client.GATEWAY, gatewayServiceAccountV2ApiPath, deleteRes)
					if err != nil {
						t.Fatalf("Error externally deleting interceptor: %s", err)
					}
				},
				Config:             providerConfigGateway + test.TestAccTestdata(t, "gateway/service_account_v2/resource_create.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
			// Re-create and Read testing for update test
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/service_account_v2/resource_create.tf"),
			},
			// Update and Read testing
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/service_account_v2/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-sa"),
					resource.TestCheckResourceAttr(resourceRef, "vcluster", "vcluster_sa"),
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
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/service_account_v2/resource_minimal.tf"),
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
				Config: providerConfigGateway + test.TestAccExample(t, "resources", "conduktor_gateway_service_account_v2", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.local_sa", "name", "simple-service-account"),
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.local_sa", "vcluster", "passthrough"),
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.local_sa", "spec.type", "LOCAL"),
				),
			},
			// Create and Read from complex example
			{
				Config: providerConfigGateway + test.TestAccExample(t, "resources", "conduktor_gateway_service_account_v2", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.external_sa", "name", "complex-service-account"),
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.external_sa", "vcluster", "vcluster_sa"),
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.external_sa", "spec.type", "EXTERNAL"),
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.external_sa", "spec.external_names.#", "1"),
					resource.TestCheckResourceAttr("conduktor_gateway_service_account_v2.external_sa", "spec.external_names.0", "externalName"),
				),
			},
		},
	})
}
