package provider

import (
	"context"
	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGatewayInterceptorV2Resource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	policyDefaultRef := "conduktor_gateway_interceptor_v2.topic-policy-default"
	policyVCSARef := "conduktor_gateway_interceptor_v2.topic-policy-vcluster_sa"
	schemaEncRef := "conduktor_gateway_interceptor_v2.schema-encryption"
	fullEncRef := "conduktor_gateway_interceptor_v2.full-encryption"
	datamaskingRef := "conduktor_gateway_interceptor_v2.datamasking"

	gwClient, err := testClient(client.GATEWAY)
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway_interceptor_v2/resource_create.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					// check topic policy interceptor
					resource.TestCheckResourceAttr(policyDefaultRef, "name", "enforce-partition-limit-test"),
					resource.TestCheckResourceAttr(policyDefaultRef, "scope.vcluster", "passthrough"),
					resource.TestCheckNoResourceAttr(policyDefaultRef, "scope.username"),
					resource.TestCheckNoResourceAttr(policyDefaultRef, "scope.group"),
					resource.TestCheckResourceAttr(policyDefaultRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.safeguard.CreateTopicPolicyPlugin"),
					resource.TestCheckResourceAttr(policyDefaultRef, "spec.priority", "1"),
					resource.TestCheckResourceAttr(policyDefaultRef, "spec.config", `{"numPartition":{"action":"INFO","max":5,"min":5},"topic":"myprefix-.*"}`),
					// check topic policy interceptor
					resource.TestCheckResourceAttr(policyVCSARef, "name", "enforce-partition-limit-test"),
					resource.TestCheckResourceAttr(policyVCSARef, "scope.vcluster", "vcluster_sa"),
					resource.TestCheckResourceAttr(policyVCSARef, "scope.username", "my.user"),
					resource.TestCheckNoResourceAttr(policyVCSARef, "scope.group"),
					resource.TestCheckResourceAttr(policyVCSARef, "spec.plugin_class", "io.conduktor.gateway.interceptor.safeguard.CreateTopicPolicyPlugin"),
					resource.TestCheckResourceAttr(policyVCSARef, "spec.priority", "4"),
					resource.TestCheckResourceAttr(policyVCSARef, "spec.config", `{"numPartition":{"action":"INFO","max":6,"min":3},"topic":"other-.*"}`),
					// check schema encryption interceptor
					resource.TestCheckResourceAttr(schemaEncRef, "name", "schema-encryption"),
					resource.TestCheckResourceAttr(schemaEncRef, "scope.vcluster", "vcluster_sa"),
					resource.TestCheckResourceAttr(schemaEncRef, "scope.group", "group-a"),
					resource.TestCheckNoResourceAttr(schemaEncRef, "scope.username"),
					resource.TestCheckResourceAttr(schemaEncRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.EncryptSchemaBasedPlugin"),
					resource.TestCheckResourceAttr(schemaEncRef, "spec.priority", "2"),
					resource.TestCheckResourceAttr(schemaEncRef, "spec.config", `{"defaultAlgorithm":"AES128_EAX","defaultKeySecretId":"in-memory-kms://myDefaultKeySecret","kmsConfig":{},"namespace":"conduktor.","schemaDataMode":"convert_json","tags":["PII","ENCRYPTION"]}`),
					// check full encryption interceptor
					resource.TestCheckResourceAttr(fullEncRef, "name", "full-encryption"),
					resource.TestCheckResourceAttr(fullEncRef, "scope.vcluster", "passthrough"),
					resource.TestCheckNoResourceAttr(fullEncRef, "scope.group"),
					resource.TestCheckNoResourceAttr(fullEncRef, "scope.username"),
					resource.TestCheckResourceAttr(fullEncRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.EncryptPlugin"),
					resource.TestCheckResourceAttr(fullEncRef, "spec.priority", "3"),
					resource.TestCheckResourceAttr(fullEncRef, "spec.config", `{"kmsConfig":{"aws":{"basicCredentials":{"accessKey":"test","secretKey":"test"}}},"recordValue":{"payload":{"algorithm":"AES128_GCM","keySecretId":"aws-kms://test-arn"}},"topic":"full-encrypt.*"}`),
					// check data masking interceptor
					resource.TestCheckResourceAttr(datamaskingRef, "name", "mask-sensitive-fields"),
					resource.TestCheckResourceAttr(datamaskingRef, "scope.vcluster", "passthrough"),
					resource.TestCheckNoResourceAttr(datamaskingRef, "scope.group"),
					resource.TestCheckNoResourceAttr(datamaskingRef, "scope.username"),
					resource.TestCheckResourceAttr(datamaskingRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.FieldLevelDataMaskingPlugin"),
					resource.TestCheckResourceAttr(datamaskingRef, "spec.priority", "100"),
					resource.TestCheckResourceAttr(datamaskingRef, "spec.config", `{"policies":[{"fields":["profile.creditCardNumber","contact.email"],"name":"Mask credit card","rule":{"type":"MASK_ALL"}},{"fields":["contact.phone"],"name":"Partial mask phone","rule":{"maskingChar":"*","numberOfChars":9,"type":"MASK_FIRST_N"}}],"topic":"^[A-Za-z]*_masked$"}`),
				),
			},
			//Importing matches the state of the previous step.
			{
				ResourceName:                         schemaEncRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "schema-encryption/vcluster_sa/group-a/null", // username is empty
				ImportStateVerifyIdentifierAttribute: "name",
			},
			//Importing matches the state of the previous step.
			{
				ResourceName:                         datamaskingRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "mask-sensitive-fields/passthrough//", // group and username are empty
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Test plan changes if externally deleted resource
			{
				PreConfig: func() {
					// wait a bit to ensure the interceptor is created
					time.Sleep(1 * time.Second)
					deleteScope := gateway.GatewayInterceptorScope{
						VCluster: "passthrough",
						Username: "",
						Group:    "",
					}
					deletePath := gatewayInterceptorV2ApiPath + "/" + "mask-sensitive-fields"
					err := gwClient.Delete(context.Background(), client.GATEWAY, deletePath, deleteScope)
					if err != nil {
						t.Fatalf("Error externally deleting interceptor: %s", err)
					}
				},
				Config:             providerConfigGateway + test.TestAccTestdata(t, "gateway_interceptor_v2/resource_create.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
			//Update and Read testing
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway_interceptor_v2/resource_update.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyDefaultRef, "name", "enforce-partition-limit-default"),
					resource.TestCheckResourceAttr(policyDefaultRef, "scope.vcluster", "passthrough"),
					resource.TestCheckResourceAttr(policyDefaultRef, "scope.username", "my.user2"),
					resource.TestCheckNoResourceAttr(policyDefaultRef, "scope.group"),
					resource.TestCheckResourceAttr(policyDefaultRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.safeguard.CreateTopicPolicyPlugin"),
					resource.TestCheckResourceAttr(policyDefaultRef, "spec.priority", "100"),
					resource.TestCheckResourceAttr(policyDefaultRef, "spec.config", `{"numPartition":{"action":"BLOCK","max":10,"min":5},"retentionMs":{"max":100,"min":10},"topic":"updatemyprefix-.*"}`),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccGatewayInterceptorV2ExampleResource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	fieldEncRef := "conduktor_gateway_interceptor_v2.field-encryption"
	headerRemoveRef := "conduktor_gateway_interceptor_v2.header-removal"
	policyRef := "conduktor_gateway_interceptor_v2.topic-policy"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from example field-encrypt.tf
			{
				Config: providerConfigGateway + test.TestAccExample(t, "resources", "conduktor_gateway_interceptor_v2", "field-encrypt.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fieldEncRef, "name", "field-encryption"),
					resource.TestCheckResourceAttr(fieldEncRef, "scope.vcluster", "passthrough"),
					resource.TestCheckResourceAttr(fieldEncRef, "scope.username", "my.user"),
					resource.TestCheckNoResourceAttr(fieldEncRef, "scope.group"),
					resource.TestCheckResourceAttr(fieldEncRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.EncryptPlugin"),
					resource.TestCheckResourceAttr(fieldEncRef, "spec.priority", "1"),
					resource.TestCheckResourceAttr(fieldEncRef, "spec.config", `{"kmsConfig":{"vault":{"token":"test","uri":"http://vault:8200","version":1}},"recordValue":{"fields":[{"algorithm":"AES128_GCM","fieldName":"password","keySecretId":"vault-kms://vault:8200/transit/keys/password-secret"},{"algorithm":"AES128_GCM","fieldName":"visa","keySecretId":"vault-kms://vault:8200/transit/keys/{{record.header.test-header}}-visa-secret-{{record.key}}-{{record.value.username}}-{{record.value.education.account.accountId}}"},{"algorithm":"AES128_GCM","fieldName":"education.account.username","keySecretId":"in-memory-kms://myDefaultKeySecret"}]},"topic":"encrypt.*"}`),
				),
			},
			// Create and Read from example header-removal.tf
			{
				Config: providerConfigGateway + test.TestAccExample(t, "resources", "conduktor_gateway_interceptor_v2", "header-removal.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(headerRemoveRef, "name", "remove-headers"),
					resource.TestCheckResourceAttr(headerRemoveRef, "scope.vcluster", "passthrough"),
					resource.TestCheckNoResourceAttr(headerRemoveRef, "scope.group"),
					resource.TestCheckNoResourceAttr(headerRemoveRef, "scope.username"),
					resource.TestCheckResourceAttr(headerRemoveRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.safeguard.MessageHeaderRemovalPlugin"),
					resource.TestCheckResourceAttr(headerRemoveRef, "spec.priority", "100"),
					resource.TestCheckResourceAttr(headerRemoveRef, "spec.config", `{"headerKeyRegex":"headerKey.*","topic":"topic-.*"}`),
				),
			},
			// Create and Read from example topic-policy.tf
			{
				Config: providerConfigGateway + test.TestAccExample(t, "resources", "conduktor_gateway_interceptor_v2", "topic-policy.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyRef, "name", "enforce-partition-limit"),
					resource.TestCheckResourceAttr(policyRef, "scope.vcluster", "passthrough"),
					resource.TestCheckNoResourceAttr(policyRef, "scope.group"),
					resource.TestCheckNoResourceAttr(policyRef, "scope.username"),
					resource.TestCheckResourceAttr(policyRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.safeguard.CreateTopicPolicyPlugin"),
					resource.TestCheckResourceAttr(policyRef, "spec.priority", "1"),
					resource.TestCheckResourceAttr(policyRef, "spec.config", `{"numPartition":{"action":"INFO","max":5,"min":5},"topic":"myprefix-.*"}`),
				),
			},
		},
	})
}
