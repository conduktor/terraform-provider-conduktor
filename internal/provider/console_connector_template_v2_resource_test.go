package provider

import (
	"testing"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConnectorTemplateV2Resource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, connectorTemplateV2MinimumVersion)

	resourceRef := "conduktor_console_connector_template_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/connector_template_v2/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "my-connector-template"),
					resource.TestCheckResourceAttr(resourceRef, "labels.env", "prod"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "My Connector Template"),
					resource.TestCheckResourceAttr(resourceRef, "spec.description", "A template for S3 sink connectors"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.metadata.name", "connector-{{team}}"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.metadata.labels.managed-by", "terraform"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.spec.class", "io.confluent.connect.s3.S3SinkConnector"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.spec.config.tasks.max", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.spec.config.s3.bucket.name", "my-bucket"),
				),
			},
			// Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "my-connector-template",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/connector_template_v2/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "my-connector-template"),
					resource.TestCheckResourceAttr(resourceRef, "labels.env", "prod"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "My Connector Template Updated"),
					resource.TestCheckResourceAttr(resourceRef, "spec.description", "Updated description"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.metadata.name", "connector-{{team}}"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.metadata.labels.managed-by", "terraform"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.spec.class", "io.confluent.connect.s3.S3SinkConnector"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.spec.config.tasks.max", "2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.defaults.spec.config.s3.bucket.name", "my-bucket-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccConnectorTemplateV2Minimal(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, connectorTemplateV2MinimumVersion)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/connector_template_v2/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_connector_template_v2.minimal", "name", "minimal-connector-template"),
					resource.TestCheckResourceAttr("conduktor_console_connector_template_v2.minimal", "spec.display_name", "Minimal Connector Template"),
					resource.TestCheckResourceAttr("conduktor_console_connector_template_v2.minimal", "spec.defaults.spec.class", "org.apache.kafka.connect.mirror.MirrorSourceConnector"),
					resource.TestCheckResourceAttr("conduktor_console_connector_template_v2.minimal", "spec.defaults.spec.config.tasks.max", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
