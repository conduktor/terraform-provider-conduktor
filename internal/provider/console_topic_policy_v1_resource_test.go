package provider

import (
	"regexp"
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTopicPolicyV1Resource(t *testing.T) {

	resourceRef := "conduktor_console_topic_policy_v1.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/topic_policy_v1/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "topicpolicy"),
					resource.TestCheckResourceAttr(resourceRef, "spec.policies.my-policy.one_of.values.#", "3"),
					resource.TestCheckResourceAttr(resourceRef, "spec.policies.my-policy.one_of.values.0", "C0"),
					resource.TestCheckResourceAttr(resourceRef, "spec.policies.my-policy.one_of.values.1", "C1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.policies.my-policy.one_of.values.2", "C2"),
				),
			},
			// Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "topicpolicy",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/topic_policy_v1/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "topicpolicy"),
					resource.TestCheckResourceAttr(resourceRef, "spec.policies.my-policy.range.max", "3600000"),
					resource.TestCheckResourceAttr(resourceRef, "spec.policies.my-policy.range.min", "60000"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTopicPolicyV1Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/topic_policy_v1/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.minimal", "name", "minimal"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.minimal", "spec.policies.my-policy.one_of.values.#", "1"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.minimal", "spec.policies.my-policy.one_of.values.0", "value"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTopicPolicyV1Constraints(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerConfigConsole + test.TestAccTestdata(t, "console/topic_policy_v1/resource_not_valid.tf"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTopicPolicyV1ExampleResource(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}

	test.CheckMinimumVersionRequirement(t, v, topicPolicyMininumVersion)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from simple example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_topic_policy_v1", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.simple", "name", "simple"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.simple", "spec.policies.spec.configs.retention.ms.range.optional", "true"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.simple", "spec.policies.spec.configs.retention.ms.range.max", "3600000"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.simple", "spec.policies.spec.configs.retention.ms.range.min", "60000"),
				),
			},
			// Create and Read from complex example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_topic_policy_v1", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "name", "complex"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.metadata.labels.data-criticality.one_of.optional", "false"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.metadata.labels.data-criticality.one_of.values.#", "3"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.metadata.labels.data-criticality.one_of.values.0", "C0"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.metadata.labels.data-criticality.one_of.values.1", "C1"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.metadata.labels.data-criticality.one_of.values.2", "C2"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.spec.configs.retention.ms.range.optional", "false"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.spec.configs.retention.ms.range.max", "604800000"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.spec.configs.retention.ms.range.min", "3600000"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.spec.replicationFactor.none_of.optional", "true"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.spec.replicationFactor.none_of.values.#", "2"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.spec.replicationFactor.none_of.values.0", "1"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.spec.replicationFactor.none_of.values.1", "2"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.metadata.name.match.optional", "false"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.metadata.name.match.pattern", "^website-analytics.(?<event>[a-z0-9-]+).(avro|json)$"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.spec.configs.allowed_keys.optional", "false"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.spec.configs.allowed_keys.keys.#", "2"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.spec.configs.allowed_keys.keys.0", "cleanup.policy"),
					resource.TestCheckResourceAttr("conduktor_console_topic_policy_v1.complex", "spec.policies.spec.configs.allowed_keys.keys.1", "retention.ms"),
				),
			},
		},
	})
}
