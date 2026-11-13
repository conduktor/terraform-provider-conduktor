package provider

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccGatewayGroupV2Resource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resourceRef := "conduktor_gateway_group_v2.test"

	gwClient, err := testClient(client.GATEWAY)
	if err != nil {
		t.Fatalf("Error creating gateway client: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/group_v2/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-group"),
					resource.TestCheckResourceAttr(resourceRef, "spec.members.#", "0"),
					resource.TestCheckResourceAttr(resourceRef, "spec.external_groups.#", "1"),
				),
			},
			// Import
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "test-group",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Test plan after external deletion
			{
				PreConfig: func() {
					time.Sleep(1 * time.Second)
					t.Log("Externally deleting gateway group test-group")
					err := gwClient.Delete(context.Background(), client.GATEWAY, fmt.Sprintf("%s/%s", gatewayGroupV2ApiPath, "test-group"), nil)
					if err != nil {
						t.Fatalf("Error externally deleting gateway group: %s", err)
					}
				},
				Config:             providerConfigGateway + test.TestAccTestdata(t, "gateway/group_v2/resource_create.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
			// Re-create
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/group_v2/resource_create.tf"),
			},
			// Update and Read
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/group_v2/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-group"),
					resource.TestCheckResourceAttr(resourceRef, "spec.external_groups.#", "2"),
				),
			},
		},
	})
}

func TestAccGatewayGroupV2Minimal(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/group_v2/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_group_v2.minimal", "name", "minimal-group"),
					resource.TestCheckResourceAttr("conduktor_gateway_group_v2.minimal", "spec.members.#", "0"),
					resource.TestCheckResourceAttr("conduktor_gateway_group_v2.minimal", "spec.external_groups.#", "0"),
				),
			},
		},
	})
}
