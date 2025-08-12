package provider

import (
	"regexp"
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKsqlDBClusterV2Resource(t *testing.T) {
	resourceRef := "conduktor_console_ksqldb_cluster_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/ksqldb_cluster_v2/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-ksqldb"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "mini-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "Test KSQLDB"),
					resource.TestCheckResourceAttr(resourceRef, "spec.url", "http://localhost:8088"),
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
				ImportStateId:                        "mini-cluster/test-ksqldb",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/ksqldb_cluster_v2/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "test-ksqldb"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "mini-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "Test KSQLDB updated"),
					resource.TestCheckResourceAttr(resourceRef, "spec.url", "https://localhost:8088"),
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

func TestAccKsqlDBClusterV2Minimal(t *testing.T) {
	resourceRef := "conduktor_console_ksqldb_cluster_v2.minimal"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/ksqldb_cluster_v2/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "minimal-ksqldb"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "mini-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "Minimal KSQLDB"),
					resource.TestCheckResourceAttr(resourceRef, "spec.url", "http://localhost:8088"),
					resource.TestCheckResourceAttr(resourceRef, "spec.ignore_untrusted_certificate", "false"),
					resource.TestCheckResourceAttr(resourceRef, "spec.headers.%", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccKsqlDBClusterV2Constraints(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Try to create with conflicting security attributes
			{
				Config:      providerConfigConsole + test.TestAccTestdata(t, "console/ksqldb_cluster_v2/resource_not_valid.tf"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

// TestAccKsqlDBClusterV2ExampleResource tests the ksqldb_cluster_v2 resource with example configurations.
func TestAccKsqlDBClusterV2ExampleResource(t *testing.T) {
	var simpleResourceRef = "conduktor_console_ksqldb_cluster_v2.simple"
	var mtlsResourceRef = "conduktor_console_ksqldb_cluster_v2.mtls"
	var basicResourceRef = "conduktor_console_ksqldb_cluster_v2.basic"
	var bearerResourceRef = "conduktor_console_ksqldb_cluster_v2.bearer"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from simple example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_ksqldb_cluster_v2", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(simpleResourceRef, "name", "simple-ksqldb"),
					resource.TestCheckResourceAttr(simpleResourceRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(simpleResourceRef, "spec.display_name", "Simple KSQLDB cluster"),
					resource.TestCheckResourceAttr(simpleResourceRef, "spec.url", "http://localhost:8088"),
					resource.TestCheckResourceAttr(simpleResourceRef, "spec.ignore_untrusted_certificate", "false"),
					resource.TestCheckResourceAttr(simpleResourceRef, "spec.headers.%", "0"),
				),
			},
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_ksqldb_cluster_v2", "mtls.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(mtlsResourceRef, "name", "mtls-ksqldb"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.display_name", "mTLS KSQLDB cluster"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.url", "https://localhost:8088"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.headers.%", "2"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.headers.X-PROJECT-HEADER", "value"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.headers.Cache-Control", "no-cache"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.ignore_untrusted_certificate", "false"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.security.ssl_auth.key", "-----BEGIN PRIVATE KEY-----\nMIIOXzCCDUegAwIBAgIRAPRytMVYJNUgCbhnA+eYumgwDQYJKoZIhvcNAQELBQAw\n...\nIFyCs+xkcgvHFtBjjel4pnIET0agtbGJbGDEQBNxX+i4MDA=\n-----END PRIVATE KEY-----\n"),
					resource.TestCheckResourceAttr(mtlsResourceRef, "spec.security.ssl_auth.certificate_chain", "-----BEGIN CERTIFICATE-----\nMIIOXzCCDUegAwIBAgIRAPRytMVYJNUgCbhnA+eYumgwDQYJKoZIhvcNAQELBQAw\n...\nIFyCs+xkcgvHFtBjjel4pnIET0agtbGJbGDEQBNxX+i4MDA=\n-----END CERTIFICATE-----\n"),
				),
			},
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_ksqldb_cluster_v2", "basicAuth.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(basicResourceRef, "name", "basic-ksqldb"),
					resource.TestCheckResourceAttr(basicResourceRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.display_name", "Basic KSQLDB cluster"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.url", "http://localhost:8088"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.headers.%", "2"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.headers.X-PROJECT-HEADER", "value"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.headers.Cache-Control", "no-cache"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.ignore_untrusted_certificate", "false"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.security.basic_auth.username", "user"),
					resource.TestCheckResourceAttr(basicResourceRef, "spec.security.basic_auth.password", "password"),
				),
			},
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_ksqldb_cluster_v2", "bearerToken.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bearerResourceRef, "name", "bearer-ksqldb"),
					resource.TestCheckResourceAttr(bearerResourceRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(bearerResourceRef, "spec.display_name", "Bearer KSQLDB cluster"),
					resource.TestCheckResourceAttr(bearerResourceRef, "spec.url", "http://localhost:8088"),
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
