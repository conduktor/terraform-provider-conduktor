package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTopicV2Resource(t *testing.T) {
	resourceRef := "conduktor_console_topic_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/topic_v2/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "topic-test"),
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
				ImportStateId:                        "kafka-cluster/topic-test",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/topic_v2/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "topic-test"),
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

// func TestAccTopicV2Minimal(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { test.TestAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Create and Read from minimal example
// 			{
// 				Config: providerConfigConsole + test.TestAccTestdata(t, "console/topic_v2/resource_minimal.tf"),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.minimal", "name", "minimal"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.minimal", "cluster", "kafka-cluster"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.minimal", "spec.partitions", "3"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.minimal", "spec.replication_factor", "1"),
// 				),
// 			},
// 			// Delete testing automatically occurs in TestCase
// 		},
// 	})
// }

//
// func TestAccTopicV2ExampleResource(t *testing.T) {
// 	v, err := fetchClientVersion(client.CONSOLE)
// 	if err != nil {
// 		t.Fatalf("Error fetching current version: %s", err)
// 	}
// 	test.CheckMinimumVersionRequirement(t, v, appInstanceMininumVersion)
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { test.TestAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//
// 		Steps: []resource.TestStep{
// 			// Create and Read from simple example
// 			{
// 				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_topic_v2", "simple.tf"),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "name", "simple"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "application", "myapp"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "spec.cluster", "kafka-cluster"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "spec.resources.#", "1"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "spec.resources.0.type", "TOPIC"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "spec.resources.0.name", "topic"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "spec.resources.0.pattern_type", "PREFIXED"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.simple", "spec.application_managed_service_account", "false"),
// 				),
// 			},
// 			// Create and Read from complex example
// 			{
// 				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_topic_v2", "complex.tf"),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "name", "complex"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "application", "myapp"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.cluster", "kafka-cluster"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.service_account", "my-service-account"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.topic_policy_ref.#", "1"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.topic_policy_ref.0", "topic-policy"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.default_catalog_visibility", "PUBLIC"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.#", "5"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.0.type", "CONNECTOR"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.0.connect_cluster", "kafka-connect"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.0.name", "click."),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.0.pattern_type", "PREFIXED"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.1.type", "CONSUMER_GROUP"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.1.name", "click."),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.1.pattern_type", "PREFIXED"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.2.type", "SUBJECT"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.2.name", "click."),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.2.pattern_type", "PREFIXED"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.3.type", "TOPIC"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.3.name", "click."),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.3.pattern_type", "PREFIXED"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.4.type", "TOPIC"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.4.name", "legacy-click."),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.4.pattern_type", "PREFIXED"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.resources.4.ownership_mode", "LIMITED"),
// 					resource.TestCheckResourceAttr("conduktor_console_topic_v2.complex", "spec.application_managed_service_account", "false"),
// 				),
// 			},
// 		},
// 	})
// }
