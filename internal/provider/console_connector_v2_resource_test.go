package provider

import (
	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConnectorV2Resource(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, connectorMininumRecommendedVersion)

	resourceRef := "conduktor_console_connector_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/connector_v2/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "connector-test"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "connect_cluster", "kafka-connect"),
					resource.TestCheckResourceAttr(resourceRef, "labels.key1", "value1"),
					resource.TestCheckResourceAttr(resourceRef, "description", "description"),
					resource.TestCheckResourceAttr(resourceRef, "auto_restart.frequency_seconds", "800"),
					resource.TestCheckResourceAttr(resourceRef, "auto_restart.enabled", "true"),
					resource.TestCheckResourceAttr(resourceRef, "spec.config.connector.class", "org.apache.kafka.connect.tools.MockSourceConnector"),
					resource.TestCheckResourceAttr(resourceRef, "spec.config.tasks.max", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.config.topic", "click.pageviews"),
					resource.TestCheckResourceAttr(resourceRef, "spec.config.file", "/etc/kafka/consumer.properties"),
				),
			},
			// Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "kafka-cluster/kafka-connect/connector-test",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/connector_v2/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "connector-test"),
					resource.TestCheckResourceAttr(resourceRef, "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "connect_cluster", "kafka-connect"),
					resource.TestCheckResourceAttr(resourceRef, "labels.env", "test"),
					resource.TestCheckResourceAttr(resourceRef, "labels.sec", "C1"),
					resource.TestCheckResourceAttr(resourceRef, "description", "description update"),
					resource.TestCheckResourceAttr(resourceRef, "auto_restart.enabled", "false"),
					resource.TestCheckResourceAttr(resourceRef, "auto_restart.frequency_seconds", "600"),
					resource.TestCheckResourceAttr(resourceRef, "spec.config.connector.class", "org.apache.kafka.connect.tools.MockSourceConnector"),
					resource.TestCheckResourceAttr(resourceRef, "spec.config.tasks.max", "2"),
					resource.TestCheckResourceAttr(resourceRef, "spec.config.topic", "click.pageviews.new"),
					resource.TestCheckResourceAttr(resourceRef, "spec.config.file", "/etc/kafka/producer.properties"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccConnectorV2Minimal(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}

	test.CheckMinimumVersionRequirement(t, v, connectorMininumRecommendedVersion)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/connector_v2/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.minimal", "name", "minimal"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.minimal", "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.minimal", "connect_cluster", "kafka-connect"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.minimal", "spec.config.connector.class", "org.apache.kafka.connect.tools.MockSourceConnector"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.minimal", "spec.config.tasks.max", "1"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.minimal", "spec.config.topic", "click.pageviews"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.minimal", "spec.config.file", "/etc/kafka/consumer.properties"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccConnectorV2Labels(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}

	test.CheckMinimumVersionRequirement(t, v, connectorMininumRecommendedVersion)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerConfigConsole + test.TestAccTestdata(t, "console/connector_v2/resource_with_managed_labels.tf"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Managed Label Key"),
			},
			{
				Config:      providerConfigConsole + test.TestAccTestdata(t, "console/connector_v2/resource_with_managed_labels_not_ro.tf"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid Configuration for Read-Only Attribute"),
			},
		},
	})
}

func TestAccConnectorV2ExampleSimpleResource(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}

	test.CheckMinimumVersionRequirement(t, v, connectorMininumRecommendedVersion)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from simple example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_connector_v2", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.simple", "name", "simple"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.simple", "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.simple", "connect_cluster", "kafka-connect"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.simple", "description", "# Simple kafka connector"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.simple", "spec.config.connector.class", "org.apache.kafka.connect.tools.MockSourceConnector"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.simple", "spec.config.tasks.max", "1"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.simple", "spec.config.topic", "click.pageviews"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.simple", "spec.config.file", "/etc/kafka/consumer.properties"),
				),
			},
		},
	})
}

func TestAccConnectorV2ExampleComplexResource(t *testing.T) {
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}

	test.CheckMinimumVersionRequirement(t, v, connectorMininumRecommendedVersion)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from complex example
			{
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_connector_v2", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.complex", "name", "complex"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.complex", "cluster", "kafka-cluster"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.complex", "connect_cluster", "kafka-connect"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.complex", "labels.domain", "clickstream"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.complex", "labels.appcode", "clk"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.complex", "description", "# Complex kafka connector"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.complex", "auto_restart.enabled", "true"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.complex", "auto_restart.frequency_seconds", "800"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.complex", "spec.config.connector.class", "org.apache.kafka.connect.tools.MockSourceConnector"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.complex", "spec.config.tasks.max", "1"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.complex", "spec.config.topic", "click.pageviews"),
					resource.TestCheckResourceAttr("conduktor_console_connector_v2.complex", "spec.config.file", "/etc/kafka/consumer.properties"),
				),
			},
		},
	})
}
