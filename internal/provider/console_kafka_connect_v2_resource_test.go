package provider

import (
	"regexp"
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKafkaConnectV2Resource(t *testing.T) {
	resourceRef := "conduktor_console_kafka_connect_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console_kafka_connect_v2_resource_create.tf"),
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
					resource.TestCheckResourceAttr(resourceRef, "spec.security.bearer_token.token", "auth-token"),
				),
			},
			//Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "mini-cluster/test-connect",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console_kafka_connect_v2_resource_update.tf"),
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
					resource.TestCheckResourceAttr(resourceRef, "spec.security.basic_auth.username", "user"),
					resource.TestCheckResourceAttr(resourceRef, "spec.security.basic_auth.password", "password"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccKafkaConnectV2Minimal(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resourceRef := "conduktor_console_kafka_connect_v2.minimal"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console_kafka_connect_v2_resource_minimal.tf"),
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

func TestAccKafkaConnectV2Constraints(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Try to create with conflicting security attributes
			{
				Config:      providerConfigConsole + test.TestAccTestdata(t, "console_kafka_connect_v2_resource_not_valid.tf"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

// TestAccKafkaConnectV2ExampleResource tests the kafka_connect_v2 resource with example configurations.
func TestAccKafkaConnectV2ExampleResource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)

	var simpleResourceRef = "conduktor_console_kafka_connect_v2.simple"
	var mtlsResourceRef = "conduktor_console_kafka_connect_v2.mtls"
	var basicResourceRef = "conduktor_console_kafka_connect_v2.basic"
	var bearerResourceRef = "conduktor_console_kafka_connect_v2.bearer"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from simple example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_kafka_connect_v2", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(simpleResourceRef, "name", "simple-connect"),
					resource.TestCheckResourceAttr(simpleResourceRef, "cluster", "mini-cluster"),
					resource.TestCheckResourceAttr(simpleResourceRef, "labels.%", "0"),
					resource.TestCheckResourceAttr(simpleResourceRef, "spec.display_name", "Simple Connect Server"),
					resource.TestCheckResourceAttr(simpleResourceRef, "spec.urls", "http://localhost:8083"),
					resource.TestCheckResourceAttr(simpleResourceRef, "spec.ignore_untrusted_certificate", "false"),
					resource.TestCheckResourceAttr(simpleResourceRef, "spec.headers.%", "0"),
				),
			},
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_kafka_connect_v2", "mtls.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(mtlsResourceRef, "name", "mtls-connect"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "cluster", "mini-cluster"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "labels.%", "3"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "labels.env", "dev"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "labels.description", "This is a complex connect using mTLS authentication"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "labels.documentation", "https://docs.mycompany.com/complex-connect"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.display_name", "mTLS Connect server"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.urls", "http://localhost:8083"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.headers.%", "2"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.headers.X-PROJECT-HEADER", "value"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.headers.Cache-Control", "no-cache"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.ignore_untrusted_certificate", "false"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.security.ssl_auth.key", "-----BEGIN PRIVATE KEY-----\nMIIOXzCCDUegAwIBAgIRAPRytMVYJNUgCbhnA+eYumgwDQYJKoZIhvcNAQELBQAw\n...\nIFyCs+xkcgvHFtBjjel4pnIET0agtbGJbGDEQBNxX+i4MDA=\n-----END PRIVATE KEY-----\n"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.security.ssl_auth.certificate_chain", "-----BEGIN CERTIFICATE-----\nMIIOXzCCDUegAwIBAgIRAPRytMVYJNUgCbhnA+eYumgwDQYJKoZIhvcNAQELBQAw\n...\nIFyCs+xkcgvHFtBjjel4pnIET0agtbGJbGDEQBNxX+i4MDA=\n-----END CERTIFICATE-----\n"),
				),
			},
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_kafka_connect_v2", "basicAuth.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicResourceRef, "name", "basic-connect"),
					resource.TestCheckResourceAttr(basicResourceRef, "cluster", "mini-cluster"),
					resource.TestCheckResourceAttr(basicResourceRef, "labels.%", "3"),
					resource.TestCheckResourceAttr(basicResourceRef, "labels.env", "dev"),
					resource.TestCheckResourceAttr(basicResourceRef, "labels.description", "This is a complex connect using basic authentication"),
					resource.TestCheckResourceAttr(basicResourceRef, "labels.documentation", "https://docs.mycompany.com/complex-connect"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.display_name", "Basic Connect server"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.urls", "http://localhost:8083"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.headers.%", "2"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.headers.X-PROJECT-HEADER", "value"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.headers.Cache-Control", "no-cache"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.ignore_untrusted_certificate", "false"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.security.basic_auth.username", "user"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.security.basic_auth.password", "password"),
				),
			},
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_kafka_connect_v2", "bearerToken.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bearerResourceRef, "name", "bearer-connect"),
					resource.TestCheckResourceAttr(bearerResourceRef, "cluster", "mini-cluster"),
					resource.TestCheckResourceAttr(bearerResourceRef, "labels.%", "3"),
					resource.TestCheckResourceAttr(bearerResourceRef, "labels.env", "dev"),
					resource.TestCheckResourceAttr(bearerResourceRef, "labels.description", "This is a complex connect using bearer token authentication"),
					resource.TestCheckResourceAttr(bearerResourceRef, "labels.documentation", "https://docs.mycompany.com/complex-connect"),
					resource.TestCheckResourceAttr(bearerResourceRef, "spec.display_name", "Bearer Connect server"),
					resource.TestCheckResourceAttr(bearerResourceRef, "spec.urls", "http://localhost:8083"),
					resource.TestCheckResourceAttr(bearerResourceRef, "spec.headers.%", "2"),
					resource.TestCheckResourceAttr(bearerResourceRef, "spec.headers.X-PROJECT-HEADER", "value"),
					resource.TestCheckResourceAttr(bearerResourceRef, "spec.headers.Cache-Control", "no-cache"),
					resource.TestCheckResourceAttr(bearerResourceRef, "spec.ignore_untrusted_certificate", "false"),
					resource.TestCheckResourceAttr(bearerResourceRef, "spec.security.bearer_token.token", "token"),
				),
			},
		},
	})
}
