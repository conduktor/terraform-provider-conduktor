package provider

import (
	"context"
	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/provider_conduktor"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Conduktor Console client is properly configured.
	// It is also possible to use the CDK_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfigConsole = `
provider "conduktor" {
  mode = "console"
}
`
	providerConfigGateway = `
provider "conduktor" {
  mode = "gateway"
}
`
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"conduktor": providerserver.NewProtocol6WithError(New("test", "none", "unknown")()),
}

func testClient(mode client.Mode) (*client.Client, error) {
	var apiParameter = client.LoadConfig(schema.ConduktorModel{}, mode)
	return client.Make(context.Background(), mode, apiParameter, "test")
}

// Temporary fix for linteger error: `fetchClientVersion` - `mode` always receives `client.CONSOLE` (`"Console"`).
// TODO: to remove once the fetchClientVersion will be used by a getaway resource.
var _, _ = fetchClientVersion(client.GATEWAY)

// Fetch current client version based on the mode.
// Used for version checks in acceptance tests.
func fetchClientVersion(mode client.Mode) (string, error) {
	var version string

	testClient, err := testClient(mode)
	if err != nil {
		return "", err
	}

	if mode == client.CONSOLE {
		version, err = testClient.GetConsoleVersion(context.Background())
		if err != nil {
			return "", err
		}
	} else if mode == client.GATEWAY {
		// TODO
		return "", nil
	}

	return version, nil
}
