package provider

import (
	"context"
	"fmt"
	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	"github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPartnerZoneV2Resource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, partnerZoneMininumConsoleVersion)
	v, err = fetchClientVersion(client.GATEWAY)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, partnerZoneMininumGatewayVersion)

	resourceRef := "conduktor_console_partner_zone_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/partner_zone_v2/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "partner-zone"),
					resource.TestCheckResourceAttr(resourceRef, "labels.label1", "value1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "Partner Zone"),
					resource.TestCheckResourceAttr(resourceRef, "spec.description", "This is a test partner zone"),
					resource.TestCheckResourceAttr(resourceRef, "spec.cluster", "gw-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authentication_mode.type", "PLAIN"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authentication_mode.service_account", "service-account-123"),
					resource.TestCheckResourceAttr(resourceRef, "spec.topics.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.topics.0.name", "topic-a"),
					resource.TestCheckResourceAttr(resourceRef, "spec.topics.0.backing_topic", "kafka-topic-a"),
					resource.TestCheckResourceAttr(resourceRef, "spec.topics.0.permission", "WRITE"),
				),
			},
			// Importing matches the state of the previous step.
			{
				ResourceName:                         resourceRef,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "partner-zone",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/partner_zone_v2/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", "partner-zone"),
					resource.TestCheckResourceAttr(resourceRef, "labels.label1", "new-value1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.display_name", "Updated Partner Zone"),
					resource.TestCheckResourceAttr(resourceRef, "spec.description", "This is an updated test partner zone"),
					resource.TestCheckResourceAttr(resourceRef, "spec.cluster", "gw-cluster"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authentication_mode.type", "PLAIN"),
					resource.TestCheckResourceAttr(resourceRef, "spec.authentication_mode.service_account", "service-account-234"),
					resource.TestCheckResourceAttr(resourceRef, "spec.topics.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "spec.topics.0.name", "topic-b"),
					resource.TestCheckResourceAttr(resourceRef, "spec.topics.0.backing_topic", "kafka-topic-b"),
					resource.TestCheckResourceAttr(resourceRef, "spec.topics.0.permission", "READ"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccPartnerZoneV2Minimal(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, partnerZoneMininumConsoleVersion)
	v, err = fetchClientVersion(client.GATEWAY)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, partnerZoneMininumGatewayVersion)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigConsole + test.TestAccTestdata(t, "console/partner_zone_v2/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.minimal", "name", "minimal"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.minimal", "spec.cluster", "gw-cluster"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.minimal", "spec.authentication_mode.type", "PLAIN"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.minimal", "spec.authentication_mode.service_account", "service-account-123"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.minimal", "spec.topics.#", "1"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.minimal", "spec.topics.0.name", "topic-a"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.minimal", "spec.topics.0.backing_topic", "kafka-topic-a"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.minimal", "spec.topics.0.permission", "WRITE"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccPartnerZoneV2ExampleResource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	v, err := fetchClientVersion(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, partnerZoneMininumConsoleVersion)
	v, err = fetchClientVersion(client.GATEWAY)
	if err != nil {
		t.Fatalf("Error fetching current version: %s", err)
	}
	test.CheckMinimumVersionRequirement(t, v, partnerZoneMininumGatewayVersion)

	consoleClient, err := testClient(client.CONSOLE)
	if err != nil {
		t.Fatalf("Error creating console client: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from simple example
			{
				PreConfig: func() {
					createBackingTopic(t, consoleClient, "kafka-topic", "gw-cluster")
				},
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_partner_zone_v2", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.simple", "name", "simple"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.simple", "spec.cluster", "gw-cluster"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.simple", "spec.display_name", "Simple Partner Zone"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.simple", "spec.url", "https://partner1.com"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.simple", "spec.authentication_mode.type", "PLAIN"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.simple", "spec.authentication_mode.service_account", "simple-partner"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.simple", "spec.topics.#", "1"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.simple", "spec.topics.0.name", "topic"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.simple", "spec.topics.0.backing_topic", "kafka-topic"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.simple", "spec.topics.0.permission", "WRITE"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.simple", "spec.traffic_control_policies.max_produce_rate", "1000000"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.simple", "spec.traffic_control_policies.max_consume_rate", "1000000"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.simple", "spec.traffic_control_policies.limit_commit_offset", "30"),
				),
			},
			// Create and Read from complex example
			{
				PreConfig: func() {
					createBackingTopic(t, consoleClient, "kafka-topic", "gw-cluster")
					createBackingTopic(t, consoleClient, "kafka-topic", "gw-cluster")
				},
				Config: providerConfigConsole + test.TestAccExample(t, "resources", "conduktor_console_partner_zone_v2", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "name", "complex"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.cluster", "gw-cluster"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.display_name", "Complex Partner Zone"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.description", "An external partner to exchange data with."),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.url", "https://partner1.com"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.partner.name", "John Doe"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.partner.role", "Data analyst"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.partner.email", "johndoe@partner.io"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.partner.phone", "07827 837 177"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.authentication_mode.type", "PLAIN"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.authentication_mode.service_account", "external-partner"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.topics.#", "2"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.topics.0.name", "ext-analytics"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.topics.0.backing_topic", "internal-analytics"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.topics.0.permission", "WRITE"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.topics.1.name", "ext-customers"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.topics.1.backing_topic", "internal-customers"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.topics.1.permission", "READ"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.traffic_control_policies.max_produce_rate", "1000000"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.traffic_control_policies.max_consume_rate", "1000000"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.traffic_control_policies.limit_commit_offset", "30"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.headers.add_on_produce.#", "1"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.headers.add_on_produce.0.key", "key"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.headers.add_on_produce.0.value", "value"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.headers.add_on_produce.0.override_if_exists", "false"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.headers.remove_on_consume.#", "1"),
					resource.TestCheckResourceAttr("conduktor_console_partner_zone_v2.complex", "spec.headers.remove_on_consume.0.key_regex", "my_org_prefix.*"),
				),
			},
		},
	})
}

func createBackingTopic(t *testing.T, client *client.Client, name, cluster string) {
	topicResource := console.TopicConsoleResource{
		Kind:       console.TopicV2Kind,
		ApiVersion: console.TopicV2ApiVersion,
		Metadata: console.TopicConsoleMetadata{
			Name:    name,
			Cluster: cluster,
		},
		Spec: console.TopicConsoleSpec{
			Partitions:        1,
			ReplicationFactor: 1,
		},
	}
	apiPath := fmt.Sprintf("/public/kafka/v2/cluster/%s/topic", cluster)
	_, err := client.Apply(context.Background(), apiPath, topicResource)
	if err != nil {
		t.Fatalf("Error creating backing topic: %s", err)
	}
}
