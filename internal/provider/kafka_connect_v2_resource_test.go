package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKafkaConnectV2Resource(t *testing.T) {
	resourceRef := "conduktor_kafka_connect_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + test.TestAccTestdata(t, "kafka_connect_v2_resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-connect"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "mini-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "labels.%", "1"),
					resource.TestCheckResourceAttr(resourceRef, "labels.env", "test"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "Test Connect"),
					resource.TestCheckResourceAttr(resourceRef, "spec.urls", "http://localhost:8083"),
					resource.TestCheckResourceAttr(resourceRef, "spec.headers.%", "2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.headers.X-PROJECT-HEADER", "value"),
					resource.TestCheckResourceAttr(resourceRef, "spec.headers.AnotherHeader", "test"),
					resource.TestCheckResourceAttr(resourceRef, "spec.ignore_untrusted_certificate", "true"),
					resource.TestCheckResourceAttr(resourceRef, "spec.security.type", "BearerToken"),
					resource.TestCheckResourceAttr(resourceRef, "spec.security.token", "auth-token"),
				),
			},
			//Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "test-connect",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfig + test.TestAccTestdata(t, "kafka_connect_v2_resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-connect"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "mini-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "labels.%", "2"),
					resource.TestCheckResourceAttr(resourceRef, "labels.env", "test"),
					resource.TestCheckResourceAttr(resourceRef, "labels.security", "C1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "Test Connect updated"),
					resource.TestCheckResourceAttr(resourceRef, "spec.urls", "http://localhost:8083"),
					resource.TestCheckResourceAttr(resourceRef, "spec.headers.%", "3"),
					resource.TestCheckResourceAttr(resourceRef, "spec.headers.X-PROJECT-HEADER", "value"),
					resource.TestCheckResourceAttr(resourceRef, "spec.headers.AnotherHeader", "test"),
					resource.TestCheckResourceAttr(resourceRef, "spec.headers.Cache-Control", "no-store"),
					resource.TestCheckResourceAttr(resourceRef, "spec.ignore_untrusted_certificate", "false"),
					resource.TestCheckResourceAttr(resourceRef, "spec.security.type", "BasicAuth"),
					resource.TestCheckResourceAttr(resourceRef, "spec.security.username", "user"),
					resource.TestCheckResourceAttr(resourceRef, "spec.security.password", "password"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccKafkaConnectV2Minimal(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resourceRef := "conduktor_kafka_connect_v2.minimal"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfig + test.TestAccTestdata(t, "kafka_connect_v2_resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "minimal-connect"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "mini-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "labels.%", "0"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "Minimal Connect"),
					resource.TestCheckResourceAttr(resourceRef, "spec.urls", "http://localhost:8083"),
					resource.TestCheckResourceAttr(resourceRef, "spec.ignore_untrusted_certificate", "false"),
					resource.TestCheckResourceAttr(resourceRef, "spec.headers.%", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

//// TestAccKafkaConnectV2ExampleResource tests the kafka_connect_v2 resource with example configurations.
//func TestAccKafkaConnectV2ExampleResource(t *testing.T) {
//	test.CheckEnterpriseEnabled(t)
//
//	var simpleResourceRef = "conduktor_kafka_connect_v2.simple"
//	var gatewayResourceRef = "conduktor_kafka_connect_v2.gateway"
//	var aivenResourceRef = "conduktor_kafka_connect_v2.aiven"
//	var awsResourceRef = "conduktor_kafka_connect_v2.aws_msk"
//	var confluentResourceRef = "conduktor_kafka_connect_v2.confluent"
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { test.TestAccPreCheck(t) },
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//
//		Steps: []resource.TestStep{
//			// Create and Read from simple example
//			{
//				Config: providerConfig + test.TestAccExample(t, "resources", "conduktor_kafka_connect_v2", "simple.tf"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(simpleResourceRef, "name", "simple-cluster"),
//					resource.TestCheckResourceAttr(simpleResourceRef, "labels.%", "0"),
//					resource.TestCheckResourceAttr(simpleResourceRef, "spec.display_name", "Simple kafka Cluster"),
//					resource.TestCheckResourceAttr(simpleResourceRef, "spec.bootstrap_servers", "localhost:9092"),
//					resource.TestCheckResourceAttr(simpleResourceRef, "spec.ignore_untrusted_certificate", "true"),
//					resource.TestCheckResourceAttr(simpleResourceRef, "spec.color", "#000000"),
//				),
//			},
//			{
//				Config: providerConfig + test.TestAccExample(t, "resources", "conduktor_kafka_connect_v2", "gateway.tf"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(gatewayResourceRef, "name", "gateway-cluster"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "labels.%", "1"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.display_name", "Gateway Cluster"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.bootstrap_servers", "gateway:6969"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.properties.%", "3"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.ignore_untrusted_certificate", "true"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.kafka_flavor.type", "Gateway"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.kafka_flavor.url", "http://gateway:8888"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.kafka_flavor.user", "admin"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.kafka_flavor.password", "admin"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.kafka_flavor.virtual_cluster", "passthrough"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.kafka_flavor.ignore_untrusted_certificate", "true"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.schema_registry.type", "ConfluentLike"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.schema_registry.url", "http://localhost:8081"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.schema_registry.security.type", "BearerToken"),
//					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.schema_registry.security.token", "auth-token"),
//				),
//			},
//			{
//				Config: providerConfig + test.TestAccExample(t, "resources", "conduktor_kafka_connect_v2", "aiven.tf"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(aivenResourceRef, "name", "aiven-cluster"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "labels.%", "1"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "spec.display_name", "Aiven Cluster"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "spec.bootstrap_servers", "cluster.aiven.io:9092"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "spec.properties.%", "3"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "spec.ignore_untrusted_certificate", "true"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "spec.kafka_flavor.type", "Aiven"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "spec.kafka_flavor.api_token", "a1b2c3d4e5f6g7h8i9j0"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "spec.kafka_flavor.project", "my-kafka-project"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "spec.kafka_flavor.service_name", "my-kafka-service"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "spec.schema_registry.type", "ConfluentLike"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "spec.schema_registry.url", "https://sr.aiven.io:8081"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "spec.schema_registry.ignore_untrusted_certificate", "false"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "spec.schema_registry.security.type", "BasicAuth"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "spec.schema_registry.security.username", "uuuuuuu"),
//					resource.TestCheckResourceAttr(aivenResourceRef, "spec.schema_registry.security.password", "ppppppp"),
//				),
//			},
//			{
//				Config: providerConfig + test.TestAccExample(t, "resources", "conduktor_kafka_connect_v2", "aws_msk.tf"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(awsResourceRef, "name", "aws-cluster"),
//					resource.TestCheckResourceAttr(awsResourceRef, "labels.%", "1"),
//					resource.TestCheckResourceAttr(awsResourceRef, "spec.display_name", "AWS MSK Cluster"),
//					resource.TestCheckResourceAttr(awsResourceRef, "spec.bootstrap_servers", "b-3-public.xxxxx.yyyyy.zz.kafka.eu-west-1.amazonaws.com:9198,b-2-public.xxxxx.yyyyy.zz.kafka.eu-west-1.amazonaws.com:9198,b-1-public.xxxxx.yyyyy.zz.kafka.eu-west-1.amazonaws.com:9198"),
//					resource.TestCheckResourceAttr(awsResourceRef, "spec.properties.%", "4"),
//					resource.TestCheckResourceAttr(awsResourceRef, "spec.ignore_untrusted_certificate", "true"),
//					resource.TestCheckResourceAttr(awsResourceRef, "spec.schema_registry.type", "Glue"),
//					resource.TestCheckResourceAttr(awsResourceRef, "spec.schema_registry.region", "eu-west-1"),
//					resource.TestCheckResourceAttr(awsResourceRef, "spec.schema_registry.registry_name", "default"),
//					resource.TestCheckResourceAttr(awsResourceRef, "spec.schema_registry.security.type", "Credentials"),
//					resource.TestCheckResourceAttr(awsResourceRef, "spec.schema_registry.security.access_key_id", "accessKey"),
//					resource.TestCheckResourceAttr(awsResourceRef, "spec.schema_registry.security.secret_key", "secretKey"),
//				),
//			},
//			{
//				Config: providerConfig + test.TestAccExample(t, "resources", "conduktor_kafka_connect_v2", "confluent.tf"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(confluentResourceRef, "name", "confluent-cluster"),
//					resource.TestCheckResourceAttr(confluentResourceRef, "labels.%", "1"),
//					resource.TestCheckResourceAttr(confluentResourceRef, "spec.display_name", "Confluent Cluster"),
//					resource.TestCheckResourceAttr(confluentResourceRef, "spec.bootstrap_servers", "aaa-aaaa.us-west4.gcp.confluent.cloud:9092"),
//					resource.TestCheckResourceAttr(confluentResourceRef, "spec.properties.%", "3"),
//					resource.TestCheckResourceAttr(confluentResourceRef, "spec.kafka_flavor.type", "Confluent"),
//					resource.TestCheckResourceAttr(confluentResourceRef, "spec.kafka_flavor.key", "yourApiKey123456"),
//					resource.TestCheckResourceAttr(confluentResourceRef, "spec.kafka_flavor.secret", "yourApiSecret123456"),
//					resource.TestCheckResourceAttr(confluentResourceRef, "spec.kafka_flavor.confluent_environment_id", "env-12345"),
//					resource.TestCheckResourceAttr(confluentResourceRef, "spec.kafka_flavor.confluent_cluster_id", "lkc-67890"),
//					resource.TestCheckResourceAttr(confluentResourceRef, "spec.schema_registry.type", "ConfluentLike"),
//					resource.TestCheckResourceAttr(confluentResourceRef, "spec.schema_registry.url", "https://bbb-bbbb.us-west4.gcp.confluent.cloud:8081"),
//					resource.TestCheckResourceAttr(confluentResourceRef, "spec.schema_registry.ignore_untrusted_certificate", "false"),
//					resource.TestCheckResourceAttr(confluentResourceRef, "spec.schema_registry.security.type", "SSLAuth"),
//				),
//			},
//		},
//	})
//}
