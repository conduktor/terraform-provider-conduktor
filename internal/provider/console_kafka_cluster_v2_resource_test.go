package provider

import (
	"regexp"
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKafkaClusterV2Resource(t *testing.T) {
	resourceRef := "conduktor_console_kafka_cluster_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/kafka_cluster_v2/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "labels.%", "1"),
					resource.TestCheckResourceAttr(resourceRef, "labels.env", "test"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "Test Cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.bootstrap_servers", "localhost:9092"),
					resource.TestCheckResourceAttr(resourceRef, "spec.icon", "kafka"),
					resource.TestCheckResourceAttr(resourceRef, "spec.color", "#FF0000"),
					resource.TestCheckResourceAttr(resourceRef, "spec.properties.%", "3"),
					resource.TestCheckResourceAttr(resourceRef, "spec.properties.sasl.jaas.config", "org.apache.kafka.common.security.plain.PlainLoginModule required username=admin password=admin-secret"),
					resource.TestCheckResourceAttr(resourceRef, "spec.properties.security.protocol", "SASL_SSL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.properties.sasl.mechanism", "PLAIN"),
					resource.TestCheckResourceAttr(resourceRef, "spec.ignore_untrusted_certificate", "true"),
					resource.TestCheckResourceAttr(resourceRef, "spec.kafka_flavor.confluent.key", "confluent-key"),
					resource.TestCheckResourceAttr(resourceRef, "spec.kafka_flavor.confluent.secret", "confluent-secret"),
					resource.TestCheckResourceAttr(resourceRef, "spec.kafka_flavor.confluent.confluent_cluster_id", "confluent-cluster-id"),
					resource.TestCheckResourceAttr(resourceRef, "spec.kafka_flavor.confluent.confluent_environment_id", "confluent-environment-id"),
					resource.TestCheckResourceAttr(resourceRef, "spec.schema_registry.confluent_like.url", "http://localhost:8081"),
					resource.TestCheckResourceAttr(resourceRef, "spec.schema_registry.confluent_like.security.bearer_token.token", "auth-token"),
				),
			},
			//Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "test-cluster",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/kafka_cluster_v2/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "labels.%", "2"),
					resource.TestCheckResourceAttr(resourceRef, "labels.env", "test"),
					resource.TestCheckResourceAttr(resourceRef, "labels.sec", "C1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "Test Cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.bootstrap_servers", "cluster.aiven.io:9092"),
					resource.TestCheckResourceAttr(resourceRef, "spec.icon", "kafka"),
					resource.TestCheckResourceAttr(resourceRef, "spec.color", "#FF0000"),
					resource.TestCheckResourceAttr(resourceRef, "spec.properties.%", "3"),
					resource.TestCheckResourceAttr(resourceRef, "spec.properties.sasl.jaas.config", "org.apache.kafka.common.security.plain.PlainLoginModule required username=admin-update password=admin-secret-update"),
					resource.TestCheckResourceAttr(resourceRef, "spec.properties.security.protocol", "SASL_SSL"),
					resource.TestCheckResourceAttr(resourceRef, "spec.properties.sasl.mechanism", "PLAIN"),
					resource.TestCheckResourceAttr(resourceRef, "spec.ignore_untrusted_certificate", "true"),
					resource.TestCheckResourceAttr(resourceRef, "spec.kafka_flavor.aiven.api_token", "aiven-api-token"),
					resource.TestCheckResourceAttr(resourceRef, "spec.kafka_flavor.aiven.project", "aiven-project"),
					resource.TestCheckResourceAttr(resourceRef, "spec.kafka_flavor.aiven.service_name", "aiven-service-name"),
					resource.TestCheckResourceAttr(resourceRef, "spec.schema_registry.confluent_like.url", "http://localhost:8081"),
					resource.TestCheckResourceAttr(resourceRef, "spec.schema_registry.confluent_like.security.basic_auth.username", "user"),
					resource.TestCheckResourceAttr(resourceRef, "spec.schema_registry.confluent_like.security.basic_auth.password", "password"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccKafkaClusterV2Minimal(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resourceRef := "conduktor_console_kafka_cluster_v2.minimal"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/kafka_cluster_v2/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "mini-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "Minimal Cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.bootstrap_servers", "localhost:9092"),
					resource.TestCheckResourceAttr(resourceRef, "spec.ignore_untrusted_certificate", "false"),
					resource.TestCheckResourceAttr(resourceRef, "labels.%", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccKafkaClusterV2Constraints(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Try to create with conflicting flavors
				Config:      providerConfigConsole + test.TestAccTestdata(t, "console/kafka_cluster_v2/resource_not_valid_flavor.tf"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
			{ // Try to create with conflicting schema registry
				Config:      providerConfigConsole + test.TestAccTestdata(t, "console/kafka_cluster_v2/resource_not_valid_sr.tf"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
			{ // Try to create with conflicting confluent security
				Config:      providerConfigConsole + test.TestAccTestdata(t, "console/kafka_cluster_v2/resource_not_valid_confluent_security.tf"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
			{ // Try to create with conflicting glue security
				Config:      providerConfigConsole + test.TestAccTestdata(t, "console/kafka_cluster_v2/resource_not_valid_glue_security.tf"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccKafkaClusterV2ExampleResource tests the kafka_cluster_v2 resource with example configurations.
func TestAccKafkaClusterV2ExampleResource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)

	var simpleResourceRef = "conduktor_console_kafka_cluster_v2.simple"
	var gatewayResourceRef = "conduktor_console_kafka_cluster_v2.gateway"
	var aivenResourceRef = "conduktor_console_kafka_cluster_v2.aiven"
	var awsResourceRef = "conduktor_console_kafka_cluster_v2.aws_msk"
	var confluentResourceRef = "conduktor_console_kafka_cluster_v2.confluent"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from simple example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_kafka_cluster_v2", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(simpleResourceRef, "name", "simple-cluster"),
					resource.TestCheckResourceAttr(simpleResourceRef, "labels.%", "0"),
					resource.TestCheckResourceAttr(simpleResourceRef, "spec.display_name", "Simple kafka Cluster"),
					resource.TestCheckResourceAttr(simpleResourceRef, "spec.bootstrap_servers", "localhost:9092"),
					resource.TestCheckResourceAttr(simpleResourceRef, "spec.ignore_untrusted_certificate", "true"),
					resource.TestCheckResourceAttr(simpleResourceRef, "spec.color", "#000000"),
				),
			},
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_kafka_cluster_v2", "gateway.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(gatewayResourceRef, "name", "gateway-cluster"),
					resource.TestCheckResourceAttr(gatewayResourceRef, "labels.%", "1"),
					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.display_name", "Gateway Cluster"),
					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.bootstrap_servers", "gateway:6969"),
					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.properties.%", "3"),
					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.ignore_untrusted_certificate", "true"),
					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.kafka_flavor.gateway.url", "http://gateway:8888"),
					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.kafka_flavor.gateway.user", "admin"),
					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.kafka_flavor.gateway.password", "admin"),
					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.kafka_flavor.gateway.virtual_cluster", "passthrough"),
					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.kafka_flavor.gateway.ignore_untrusted_certificate", "true"),
					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.schema_registry.confluent_like.url", "http://localhost:8081"),
					resource.TestCheckResourceAttr(gatewayResourceRef, "spec.schema_registry.confluent_like.security.bearer_token.token", "auth-token"),
				),
			},
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_kafka_cluster_v2", "aiven.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(aivenResourceRef, "name", "aiven-cluster"),
					resource.TestCheckResourceAttr(aivenResourceRef, "labels.%", "1"),
					resource.TestCheckResourceAttr(aivenResourceRef, "spec.display_name", "Aiven Cluster"),
					resource.TestCheckResourceAttr(aivenResourceRef, "spec.bootstrap_servers", "cluster.aiven.io:9092"),
					resource.TestCheckResourceAttr(aivenResourceRef, "spec.properties.%", "3"),
					resource.TestCheckResourceAttr(aivenResourceRef, "spec.ignore_untrusted_certificate", "true"),
					resource.TestCheckResourceAttr(aivenResourceRef, "spec.kafka_flavor.aiven.api_token", "a1b2c3d4e5f6g7h8i9j0"),
					resource.TestCheckResourceAttr(aivenResourceRef, "spec.kafka_flavor.aiven.project", "my-kafka-project"),
					resource.TestCheckResourceAttr(aivenResourceRef, "spec.kafka_flavor.aiven.service_name", "my-kafka-service"),
					resource.TestCheckResourceAttr(aivenResourceRef, "spec.schema_registry.confluent_like.url", "https://sr.aiven.io:8081"),
					resource.TestCheckResourceAttr(aivenResourceRef, "spec.schema_registry.confluent_like.ignore_untrusted_certificate", "false"),
					resource.TestCheckResourceAttr(aivenResourceRef, "spec.schema_registry.confluent_like.security.basic_auth.username", "uuuuuuu"),
					resource.TestCheckResourceAttr(aivenResourceRef, "spec.schema_registry.confluent_like.security.basic_auth.password", "ppppppp"),
				),
			},
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_kafka_cluster_v2", "aws_msk.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(awsResourceRef, "name", "aws-cluster"),
					resource.TestCheckResourceAttr(awsResourceRef, "labels.%", "1"),
					resource.TestCheckResourceAttr(awsResourceRef, "spec.display_name", "AWS MSK Cluster"),
					resource.TestCheckResourceAttr(awsResourceRef, "spec.bootstrap_servers", "b-3-public.xxxxx.yyyyy.zz.kafka.eu-west-1.amazonaws.com:9198,b-2-public.xxxxx.yyyyy.zz.kafka.eu-west-1.amazonaws.com:9198,b-1-public.xxxxx.yyyyy.zz.kafka.eu-west-1.amazonaws.com:9198"),
					resource.TestCheckResourceAttr(awsResourceRef, "spec.properties.%", "4"),
					resource.TestCheckResourceAttr(awsResourceRef, "spec.ignore_untrusted_certificate", "true"),
					resource.TestCheckResourceAttr(awsResourceRef, "spec.schema_registry.glue.region", "eu-west-1"),
					resource.TestCheckResourceAttr(awsResourceRef, "spec.schema_registry.glue.registry_name", "default"),
					resource.TestCheckResourceAttr(awsResourceRef, "spec.schema_registry.glue.security.credentials.access_key_id", "accessKey"),
					resource.TestCheckResourceAttr(awsResourceRef, "spec.schema_registry.glue.security.credentials.secret_key", "secretKey"),
				),
			},
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_kafka_cluster_v2", "confluent.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(confluentResourceRef, "name", "confluent-cluster"),
					resource.TestCheckResourceAttr(confluentResourceRef, "labels.%", "1"),
					resource.TestCheckResourceAttr(confluentResourceRef, "spec.display_name", "Confluent Cluster"),
					resource.TestCheckResourceAttr(confluentResourceRef, "spec.bootstrap_servers", "aaa-aaaa.us-west4.gcp.confluent.cloud:9092"),
					resource.TestCheckResourceAttr(confluentResourceRef, "spec.properties.%", "3"),
					resource.TestCheckResourceAttr(confluentResourceRef, "spec.kafka_flavor.confluent.key", "yourApiKey123456"),
					resource.TestCheckResourceAttr(confluentResourceRef, "spec.kafka_flavor.confluent.secret", "yourApiSecret123456"),
					resource.TestCheckResourceAttr(confluentResourceRef, "spec.kafka_flavor.confluent.confluent_environment_id", "env-12345"),
					resource.TestCheckResourceAttr(confluentResourceRef, "spec.kafka_flavor.confluent.confluent_cluster_id", "lkc-67890"),
					resource.TestCheckResourceAttr(confluentResourceRef, "spec.schema_registry.confluent_like.url", "https://bbb-bbbb.us-west4.gcp.confluent.cloud:8081"),
					resource.TestCheckResourceAttr(confluentResourceRef, "spec.schema_registry.confluent_like.ignore_untrusted_certificate", "false"),
					resource.TestCheckResourceAttr(confluentResourceRef, "spec.schema_registry.confluent_like.security.ssl_auth.key", "-----BEGIN PRIVATE KEY-----\nMIIOXzCCDUegAwIBAgIRAPRytMVYJNUgCbhnA+eYumgwDQYJKoZIhvcNAQELBQAw\n...\nIFyCs+xkcgvHFtBjjel4pnIET0agtbGJbGDEQBNxX+i4MDA=\n-----END PRIVATE KEY-----\n"),
					resource.TestCheckResourceAttr(confluentResourceRef, "spec.schema_registry.confluent_like.security.ssl_auth.certificate_chain", "-----BEGIN CERTIFICATE-----\nMIIOXzCCDUegAwIBAgIRAPRytMVYJNUgCbhnA+eYumgwDQYJKoZIhvcNAQELBQAw\n...\nIFyCs+xkcgvHFtBjjel4pnIET0agtbGJbGDEQBNxX+i4MDA=\n-----END CERTIFICATE-----\n"),
				),
			},
		},
	})
}
