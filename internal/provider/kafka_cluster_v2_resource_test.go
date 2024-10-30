package provider

//func TestAccKafkaClusterV2Resource(t *testing.T) {
//	resourceRef := "conduktor_kafka_cluster_v2.test"
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { test.TestAccPreCheck(t) },
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			// Create and Read testing
//			{
//				Config: providerConfig + test.TestAccTestdata(t, "kafka_cluster_v2_resource_create.tf"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceRef, "name", "test-cluster"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "Test Cluster"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.bootstrap_servers", "localhost:9092"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.icon", "kafka"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.color", "#FF0000"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.properties", "sasl.jaas.config: org.apache.kafka.common.security.plain.PlainLoginModule required username=\"admin\" password=\"admin-secret\";\nsecurity.protocol: SASL_SSL\nsasl.mechanism: PLAIN"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.ignore_untrusted_certificate", "true"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.kafka_flavor.type", "Confluent"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.kafka_flavor.key", "confluent-key"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.kafka_flavor.secret", "confluent-secret"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.kafka_flavor.confluent_cluster_id", "confluent-cluster-id"),
//					resource.TestCheckResourceAttr(resourceRef, "spec.kafka_flavor.confluent_environment_id", "confluent-environment-id"),
//				),
//			},
//			// Importing matches the state of the previous step.
//			//{
//			//	ResourceName:                         resourceRef,
//			//	ImportState:                          true,
//			//	ImportStateVerify:                    true,
//			//	ImportStateId:                        "sales",
//			//	ImportStateVerifyIdentifierAttribute: "name",
//			//},
//			//// Update and Read testing
//			//{
//			//	Config: providerConfig + test.TestAccTestdata(t, "group_v2_resource_update.tf"),
//			//	Check: resource.ComposeAggregateTestCheckFunc(
//			//		resource.TestCheckResourceAttr(resourceRef, "name", "sales"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "New Sales Department"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.description", "New Sales Department Group"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.external_groups.#", "2"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.external_groups.0", "sales"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.external_groups.1", "scranton"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.members.#", "1"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.members.0", "michael.scott@dunder.mifflin.com"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.members_from_external_groups.#", "0"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.permissions.#", "2"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.resource_type", "PLATFORM"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.permissions.0.permissions.#", "4"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.permissions.1.resource_type", "TOPIC"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.permissions.1.name", "test-topic"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.permissions.1.cluster", "*"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.permissions.1.pattern_type", "LITERAL"),
//			//		resource.TestCheckResourceAttr(resourceRef, "spec.permissions.1.permissions.#", "3"),
//			//	),
//			//},
//			// Delete testing automatically occurs in TestCase
//		},
//	})
//}

//func TestAccKafkaClusterV2Minimal(t *testing.T) {
//	test.CheckEnterpriseEnabled(t)
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { test.TestAccPreCheck(t) },
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			// Create and Read from minimal example
//			{
//				Config: providerConfig + test.TestAccTestdata(t, "group_v2_resource_minimal.tf"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr("conduktor_group_v2.minimal", "name", "minimal"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.minimal", "spec.display_name", "Minimal"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.minimal", "spec.external_groups.#", "0"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.minimal", "spec.members.#", "0"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.minimal", "spec.members_from_external_groups.#", "0"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.minimal", "spec.permissions.#", "0"),
//				),
//			},
//			// Delete testing automatically occurs in TestCase
//		},
//	})
//}
//
//func TestAccKafkaClusterV2ExampleResource(t *testing.T) {
//	test.CheckEnterpriseEnabled(t)
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { test.TestAccPreCheck(t) },
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//
//		Steps: []resource.TestStep{
//			// Create and Read from simple example
//			{
//				Config: providerConfig + test.TestAccExample(t, "resources", "conduktor_group_v2", "simple.tf"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "name", "simple-group"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "spec.display_name", "Simple Group"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "spec.description", "Simple group description"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "spec.external_groups.#", "0"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "spec.members.#", "0"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "spec.members_from_external_groups.#", "0"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "spec.permissions.#", "0"),
//				),
//			},
//			// Create and Read from complex example
//			{
//				Config: providerConfig + test.TestAccExample(t, "resources", "conduktor_group_v2", "complex.tf"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "name", "complex-group"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "spec.display_name", "Complex group"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "spec.description", "Complex group description"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "spec.external_groups.#", "1"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "spec.external_groups.0", "sso-group1"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "spec.members.#", "1"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "spec.members.0", "user1@company.com"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "spec.members_from_external_groups.#", "0"),
//					resource.TestCheckResourceAttr("conduktor_group_v2.example", "spec.permissions.#", "2"),
//				),
//			},
//		},
//	})
//}
