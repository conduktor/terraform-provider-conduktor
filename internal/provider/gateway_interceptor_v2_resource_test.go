package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGatewayInterceptorV2Resource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	policyRef := "conduktor_gateway_interceptor_v2.topic-policy"
	schemaEncRef := "conduktor_gateway_interceptor_v2.schema-encryption"
	fullEncRef := "conduktor_gateway_interceptor_v2.full-encryption"
	datamaskingRef := "conduktor_gateway_interceptor_v2.datamasking"

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
					resource.TestCheckResourceAttr(policyRef, "name", "enforce-partition-limit"),
					resource.TestCheckResourceAttr(policyRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.safeguard.CreateTopicPolicyPlugin"),
					resource.TestCheckResourceAttr(policyRef, "spec.priority", "1"),
					resource.TestCheckResourceAttr(policyRef, "spec.config", `{"numPartition":{"action":"INFO","max":5,"min":5},"topic":"myprefix-.*"}`),
					// check schema encryption interceptor
					resource.TestCheckResourceAttr(schemaEncRef, "name", "schema-encryption"),
					resource.TestCheckResourceAttr(schemaEncRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.EncryptSchemaBasedPlugin"),
					resource.TestCheckResourceAttr(schemaEncRef, "spec.priority", "2"),
					resource.TestCheckResourceAttr(schemaEncRef, "spec.config", `{"defaultAlgorithm":"AES128_EAX","defaultKeySecretId":"in-memory-kms://myDefaultKeySecret","kmsConfig":{},"namespace":"conduktor.","schemaDataMode":"convert_json","tags":["PII","ENCRYPTION"]}`),
					// check full encryption interceptor
					resource.TestCheckResourceAttr(fullEncRef, "name", "full-encryption"),
					resource.TestCheckResourceAttr(fullEncRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.EncryptPlugin"),
					resource.TestCheckResourceAttr(fullEncRef, "spec.priority", "3"),
					resource.TestCheckResourceAttr(fullEncRef, "spec.config", `{"kmsConfig":{"aws":{"basicCredentials":{"accessKey":"test","secretKey":"test"}}},"recordValue":{"payload":{"algorithm":"AES128_GCM","keySecretId":"aws-kms://test-arn"}},"topic":"full-encrypt.*"}`),
					// check data masking interceptor
					resource.TestCheckResourceAttr(datamaskingRef, "name", "mask-sensitive-fields"),
					resource.TestCheckResourceAttr(datamaskingRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.FieldLevelDataMaskingPlugin"),
					resource.TestCheckResourceAttr(datamaskingRef, "spec.priority", "100"),
					resource.TestCheckResourceAttr(datamaskingRef, "spec.config", `{"policies":[{"fields":["profile.creditCardNumber","contact.email"],"name":"Mask credit card","rule":{"type":"MASK_ALL"}},{"fields":["contact.phone"],"name":"Partial mask phone","rule":{"maskingChar":"*","numberOfChars":9,"type":"MASK_FIRST_N"}}],"topic":"^[A-Za-z]*_masked$"}`),
				),
			},
			//Importing matches the state of the previous step.
			{
				ResourceName:                         policyRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "enforce-partition-limit",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway_interceptor_v2/resource_update.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyRef, "name", "enforce-partition-limit"),
					resource.TestCheckResourceAttr(policyRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.safeguard.CreateTopicPolicyPlugin"),
					resource.TestCheckResourceAttr(policyRef, "spec.priority", "100"),
					resource.TestCheckResourceAttr(policyRef, "spec.config", `{"numPartition":{"action":"BLOCK","max":10,"min":5},"retentionMs":{"max":100,"min":10},"topic":"updatemyprefix-.*"}`),
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
			//// Create and Read from example field-encrypt.tf
			{
				Config: providerConfigGateway + test.TestAccExample(t, "resources", "conduktor_gateway_interceptor_v2", "field-encrypt.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fieldEncRef, "name", "field-encryption"),
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
					resource.TestCheckResourceAttr(policyRef, "spec.plugin_class", "io.conduktor.gateway.interceptor.safeguard.CreateTopicPolicyPlugin"),
					resource.TestCheckResourceAttr(policyRef, "spec.priority", "1"),
					resource.TestCheckResourceAttr(policyRef, "spec.config", `{"numPartition":{"action":"INFO","max":5,"min":5},"topic":"myprefix-.*"}`),
				),
			},
		},
	})
}
