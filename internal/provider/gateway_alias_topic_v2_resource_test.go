package provider

import (
	"context"
	"testing"
	"time"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccGatewayAliasTopicV2Resource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resourceRef := "conduktor_gateway_alias_topic_v2.test"

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
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/alias_topic_v2/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-alias"),
					resource.TestCheckResourceAttr(resourceRef, "vcluster", "vcluster_sa"),
					resource.TestCheckResourceAttr(resourceRef, "spec.physical_name", "my-topic"),
				),
			},
			// Import
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "test-alias/vcluster_sa",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Test plan after external deletion
			{
				PreConfig: func() {
					time.Sleep(1 * time.Second)
					t.Log("Externally deleting gateway alias topic test-alias")
					deleteRes := gateway.GatewayAliasTopicMetadata{
						Name:     "test-alias",
						VCluster: "vcluster_sa",
					}
					err := gwClient.Delete(context.Background(), client.GATEWAY, gatewayAliasTopicV2ApiPath, deleteRes)
					if err != nil {
						t.Fatalf("Error externally deleting gateway alias topic: %s", err)
					}
				},
				Config:             providerConfigGateway + test.TestAccTestdata(t, "gateway/alias_topic_v2/resource_create.tf"),
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
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/alias_topic_v2/resource_create.tf"),
			},
			// Update and Read
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/alias_topic_v2/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-alias"),
					resource.TestCheckResourceAttr(resourceRef, "vcluster", "vcluster_sa"),
					resource.TestCheckResourceAttr(resourceRef, "spec.physical_name", "my-other-topic"),
				),
			},
		},
	})
}

func TestAccGatewayAliasTopicV2Minimal(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/alias_topic_v2/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_alias_topic_v2.minimal", "name", "minimal-alias"),
					resource.TestCheckResourceAttr("conduktor_gateway_alias_topic_v2.minimal", "vcluster", "passthrough"),
					resource.TestCheckResourceAttr("conduktor_gateway_alias_topic_v2.minimal", "spec.physical_name", "my-topic"),
				),
			},
		},
	})
}

func TestAccGatewayAliasTopicV2ExampleResource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigGateway + test.TestAccExample(t, "resources", "conduktor_gateway_alias_topic_v2", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_alias_topic_v2.simple_alias", "name", "simple-alias"),
					resource.TestCheckResourceAttr("conduktor_gateway_alias_topic_v2.simple_alias", "spec.physical_name", "my-topic"),
				),
			},
			{
				Config: providerConfigGateway + test.TestAccExample(t, "resources", "conduktor_gateway_alias_topic_v2", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_alias_topic_v2.complex_alias", "name", "complex-alias"),
					resource.TestCheckResourceAttr("conduktor_gateway_alias_topic_v2.complex_alias", "vcluster", "vcluster_sa"),
					resource.TestCheckResourceAttr("conduktor_gateway_alias_topic_v2.complex_alias", "spec.physical_name", "my-topic"),
					resource.TestCheckResourceAttr("conduktor_gateway_alias_topic_v2.complex_alias", "spec.physical_cluster", "main"),
				),
			},
		},
	})
}
