package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTopicTemplateV2Resource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, topicTemplateV2MinimumVersion)

	resourceRef := "conduktor_console_topic_template_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/topic_template_v2/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "my-template"),
					resource.TestCheckResourceAttr(resourceRef, "labels.env", "prod"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "My Topic Template"),
					resource.TestCheckResourceAttr(resourceRef, "spec.description", "A template for production topics"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.metadata.name", "topic-{{department}}"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.metadata.labels.throughput", "high"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.spec.partitions", "6"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.spec.replication_factor", "3"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.spec.configs.retention.ms", "604800000"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.spec.configs.cleanup.policy", "delete"),
				),
			},
			// Import
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "my-template",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/topic_template_v2/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "my-template"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "My Topic Template Updated"),
					resource.TestCheckResourceAttr(resourceRef, "spec.description", "Updated description"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.spec.partitions", "12"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.spec.configs.retention.ms", "1209600000"),
				),
			},
			// Delete occurs automatically
		},
	})
}

func TestAccTopicTemplateV2Minimal(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, topicTemplateV2MinimumVersion)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/topic_template_v2/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_topic_template_v2.minimal", "name", "minimal-template"),
					resource.TestCheckResourceAttr("conduktor_console_topic_template_v2.minimal", "spec.display_name", "Minimal Template"),
				),
			},
		},
	})
}

func TestAccTopicTemplateV2ExampleResource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, topicTemplateV2MinimumVersion)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_topic_template_v2", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_topic_template_v2.simple", "name", "simple-template"),
					resource.TestCheckResourceAttr("conduktor_console_topic_template_v2.simple", "spec.display_name", "Simple Template"),
					resource.TestCheckResourceAttr("conduktor_console_topic_template_v2.simple", "spec.defaults.spec.partitions", "3"),
					resource.TestCheckResourceAttr("conduktor_console_topic_template_v2.simple", "spec.defaults.spec.replication_factor", "1"),
					resource.TestCheckResourceAttr("conduktor_console_topic_template_v2.simple", "spec.defaults.spec.configs.retention.ms", "604800000"),
				),
			},
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_topic_template_v2", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_topic_template_v2.complex", "name", "complex-template"),
					resource.TestCheckResourceAttr("conduktor_console_topic_template_v2.complex", "spec.display_name", "Complex Production Template"),
					resource.TestCheckResourceAttr("conduktor_console_topic_template_v2.complex", "spec.defaults.metadata.name", "svc-{{department}}-{{env}}"),
					resource.TestCheckResourceAttr("conduktor_console_topic_template_v2.complex", "spec.defaults.spec.partitions", "12"),
					resource.TestCheckResourceAttr("conduktor_console_topic_template_v2.complex", "spec.defaults.spec.replication_factor", "3"),
				),
			},
		},
	})
}
