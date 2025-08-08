package provider

import (
	"regexp"
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTopicV2Resource(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, topicMininumVersion)

	resourceRef := "conduktor_console_topic_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/topic_v2/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "Kafka-1st-topic-test"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "labels.key1", "value1"),
					resource.TestCheckResourceAttr(resourceRef, "catalog_visibility", "PUBLIC"),
					resource.TestCheckResourceAttr(resourceRef, "description_is_editable", "true"),
					resource.TestCheckResourceAttr(resourceRef, "description", "description"),
					resource.TestCheckResourceAttr(resourceRef, "sql_storage.retention_time_in_second", "86400"),
					resource.TestCheckResourceAttr(resourceRef, "sql_storage.enabled", "true"),
					resource.TestCheckResourceAttr(resourceRef, "spec.partitions", "10"),
					resource.TestCheckResourceAttr(resourceRef, "spec.replication_factor", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.configs.cleanup.policy", "delete"),
				),
			},
			// Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "kafka-cluster/Kafka-1st-topic-test",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/topic_v2/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "Kafka-1st-topic-test"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "labels.key1", "value1"),
					resource.TestCheckResourceAttr(resourceRef, "labels.key2", "value2"),
					resource.TestCheckResourceAttr(resourceRef, "catalog_visibility", "PRIVATE"),
					resource.TestCheckResourceAttr(resourceRef, "description_is_editable", "false"),
					resource.TestCheckResourceAttr(resourceRef, "description", "new description"),
					resource.TestCheckResourceAttr(resourceRef, "sql_storage.retention_time_in_second", "86400"),
					resource.TestCheckResourceAttr(resourceRef, "sql_storage.enabled", "true"),
					resource.TestCheckResourceAttr(resourceRef, "spec.partitions", "10"),
					resource.TestCheckResourceAttr(resourceRef, "spec.replication_factor", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.configs.cleanup.policy", "delete"),
					resource.TestCheckResourceAttr(resourceRef, "spec.configs.retention.ms", "60000"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTopicV2Minimal(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}

	test.CheckMinimumVersionRequirement(t, v, topicMininumVersion)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/topic_v2/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.minimal", "name", "minimal"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.minimal", "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.minimal", "spec.partitions", "3"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.minimal", "spec.replication_factor", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTopicV2Labels(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}

	test.CheckMinimumVersionRequirement(t, v, topicMininumVersion)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerConfigConsole + test.TestAccTestdata(t, "console/topic_v2/resource_with_managed_labels.tf"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Managed Label Key"),
			},
			{
				Config:      providerConfigConsole + test.TestAccTestdata(t, "console/topic_v2/resource_with_managed_labels_not_ro.tf"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid Configuration for Read-Only Attribute"),
			},
		},
	})
}

func TestAccTopicV2ExampleResource(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}

	test.CheckMinimumVersionRequirement(t, v, topicMininumVersion)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from simple example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_topic_v2", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "name", "simple"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "labels.domain", "clickstream"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "description", "# Simple kafka topic"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "spec.partitions", "3"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "spec.replication_factor", "1"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "spec.configs.cleanup.policy", "delete"),
				),
			},
			// Create and Read from complex example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_topic_v2", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "name", "complex"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "labels.domain", "clickstream"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "labels.appcode", "clk"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "catalog_visibility", "PRIVATE"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "description_is_editable", "false"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "description", "# Complex kafka topic"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "sql_storage.retention_time_in_second", "60000"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "sql_storage.enabled", "true"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.partitions", "3"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.replication_factor", "1"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.configs.cleanup.policy", "delete"),
					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.configs.retention.ms", "60000"),
				),
			},
		},
	})
}
