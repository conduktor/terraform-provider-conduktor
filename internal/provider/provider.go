package provider

import (
	"context"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/provider_conduktor"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure ConduktorProvider satisfies various provider interfaces.
var _ provider.Provider = &ConduktorProvider{}
var _ provider.ProviderWithFunctions = &ConduktorProvider{}

// ConduktorProvider defines the provider implementation.
type ConduktorProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
	// commit is set to the provider commit on release, "none" otherwise
	commit string
	// date is set to the provider release date, "unknown" otherwise
	date string
}

type ProviderData struct {
	ConsoleClient *client.ConsoleClient
	GatewayClient *client.GatewayClient
}

func (p *ConduktorProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "conduktor"
	resp.Version = p.version
}

func (p *ConduktorProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.ConduktorProviderSchema(ctx)
}

func (p *ConduktorProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config schema.ConduktorModel
	var consoleApiClient *client.ConsoleClient
	var gatewayApiClient *client.GatewayClient
	var err error

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	consoleUrl := schemaUtils.GetStringConfig(config.ConsoleUrl, []string{"CDK_BASE_URL", "CDK_CONSOLE_URL"})
	apiToken := schemaUtils.GetStringConfig(config.ApiToken, []string{"CDK_API_TOKEN", "CDK_API_KEY"})
	adminEmail := schemaUtils.GetStringConfig(config.AdminEmail, []string{"CDK_ADMIN_EMAIL"})
	adminPassword := schemaUtils.GetStringConfig(config.AdminPassword, []string{"CDK_ADMIN_PASSWORD"})
	gatewayUrl := schemaUtils.GetStringConfig(config.GatewayUrl, []string{"CDK_GATEWAY_BASE_URL"})
	gatewayUser := schemaUtils.GetStringConfig(config.GatewayUser, []string{"CDK_GATEWAY_USER"})
	gatewayPassword := schemaUtils.GetStringConfig(config.GatewayPassword, []string{"CDK_GATEWAY_PASSWORD"})
	cert := schemaUtils.GetStringConfig(config.Cert, []string{"CDK_CERT"})
	cacert := schemaUtils.GetStringConfig(config.Cacert, []string{"CDK_CACERT"})
	key := schemaUtils.GetStringConfig(config.Key, []string{"CDK_KEY"})
	insecure := schemaUtils.GetBooleanConfig(config.Insecure, []string{"CDK_INSECURE"}, false)
	gatewayCert := schemaUtils.GetStringConfig(config.GatewayCert, []string{"CDK_GATEWAY_CERT"})
	gatewayCacert := schemaUtils.GetStringConfig(config.GatewayCacert, []string{"CDK_GATEWAY_CACERT"})
	gatewayKey := schemaUtils.GetStringConfig(config.GatewayKey, []string{"CDK_GATEWAY_KEY"})
	gatewayInsecure := schemaUtils.GetBooleanConfig(config.GatewayInsecure, []string{"CDK_GATEWAY_INSECURE"}, false)

	// Validate mandatory configurations
	if consoleUrl == "" && gatewayUrl == "" {
		details := "The provider cannot create any API client as there is a missing or empty value for both the Console and Gateway URL. " +
			"Set either : " +
			" - the console_url value in the configuration or use the CDK_BASE_URL or CDK_CONSOLE_URL environment variable. " +
			" - the gateway_url value in the configuration or use the CDK_GATEWAY_BASE_URL environment variable. " +
			"If either is already set, ensure the value is not empty."

		resp.Diagnostics.AddAttributeError(path.Root("console_url"), "Missing Console URL", details)
		resp.Diagnostics.AddAttributeError(path.Root("gateway_url"), "Missing Gateway URL", details)
	}

	if consoleUrl != "" {
		if apiToken == "" && adminEmail == "" && adminPassword == "" {
			details := "The provider cannot create the Console API client as there is a missing or empty value for the API Token and missing or empty values for the admin email/password. " +
				"Set either : " +
				" - the api_token value in the configuration or use the CDK_API_TOKEN environment variable. " +
				" - the admin_email and admin_password value in the configuration or use the CDK_ADMIN_EMAIL/CDK_ADMIN_PASSWORD environment variable. " +
				"If either is already set, ensure the value is not empty."

			resp.Diagnostics.AddAttributeError(path.Root("api_token"), "Missing API token", details)
			resp.Diagnostics.AddAttributeError(path.Root("admin_email"), "Missing Admin email", details)
			resp.Diagnostics.AddAttributeError(path.Root("admin_password"), "Missing Admin password", details)
		}
	}

	if gatewayUrl != "" && (gatewayUser == "" || gatewayPassword == "") {
		details := "The provider cannot create the Gateway API client as there is a missing or empty value for the admin email/password. " +
			"Set both : " +
			" - the gateway_user value in the configuration or use the CDK_GATEWAY_USER environment variable. " +
			" - the gateway_password value in the configuration or use the CDK_GATEWAY_PASSWORD environment variable. " +
			"If either is already set, ensure the value is not empty."

		resp.Diagnostics.AddAttributeError(path.Root("gateway_user"), "Missing Gateway Admin login", details)
		resp.Diagnostics.AddAttributeError(path.Root("gateway_password"), "Missing Gateway Admin password", details)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "console_url", consoleUrl)
	ctx = tflog.SetField(ctx, "api_token", apiToken)
	ctx = tflog.SetField(ctx, "admin_email", adminEmail)
	ctx = tflog.SetField(ctx, "admin_password", adminPassword)
	ctx = tflog.SetField(ctx, "gateway_url", gatewayUrl)
	ctx = tflog.SetField(ctx, "gateway_user", gatewayUser)
	ctx = tflog.SetField(ctx, "gateway_password", gatewayPassword)
	ctx = tflog.SetField(ctx, "cert", cert)
	ctx = tflog.SetField(ctx, "cacert", cacert)
	ctx = tflog.SetField(ctx, "key", key)
	ctx = tflog.SetField(ctx, "insecure", insecure)
	ctx = tflog.SetField(ctx, "gateway_cert", gatewayCert)
	ctx = tflog.SetField(ctx, "gateway_cacert", gatewayCacert)
	ctx = tflog.SetField(ctx, "gateway_key", gatewayKey)
	ctx = tflog.SetField(ctx, "gateway_insecure", gatewayInsecure)
	// Avoid leaking sensitive information in logs
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "api_token")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "admin_password")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "gateway_password")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "key")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "gatewayKey")
	tflog.Debug(ctx, "Creating Conduktor Console client")

	// Create Console or Gateway clients only if the respective URL is provided
	if consoleUrl != "" {
		consoleApiClient, err = client.Make(ctx,
			client.ApiParameter{
				ApiKey:      apiToken,
				BaseUrl:     consoleUrl,
				CdkUser:     adminEmail,
				CdkPassword: adminPassword,
				TLSParameters: client.TLSParameters{
					Key:      key,
					Cert:     cert,
					Cacert:   cacert,
					Insecure: insecure,
				},
			},
			p.version,
		)
		if err != nil {
			resp.Diagnostics.AddError("Could not create the Conduktor Console API client", err.Error())
			return
		}
	}

	if gatewayUrl != "" {
		gatewayApiClient, err = client.MakeGateway(ctx,
			client.GatewayApiParameters{
				BaseUrl:         gatewayUrl,
				GatewayUser:     gatewayUser,
				GatewayPassword: gatewayPassword,
				TLSParameters: client.TLSParameters{
					Key:      gatewayKey,
					Cert:     gatewayCert,
					Cacert:   gatewayCacert,
					Insecure: gatewayInsecure,
				},
			},
			p.version,
		)
		if err != nil {
			resp.Diagnostics.AddError("Could not create the Conduktor Gateway API client", err.Error())
			return
		}
	}

	data := &ProviderData{
		ConsoleClient: consoleApiClient,
		GatewayClient: gatewayApiClient,
	}
	resp.DataSourceData = data
	resp.ResourceData = data

	if consoleApiClient != nil {
		tflog.Info(ctx, "Configured Conduktor Console client", map[string]any{"success": true})
	}
	if gatewayApiClient != nil {
		tflog.Info(ctx, "Configured Conduktor Gateway client", map[string]any{"success": true})
	}
}

func (p *ConduktorProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserV2Resource,
		NewGroupV2Resource,
		NewGenericResource,
		NewKafkaClusterV2Resource,
		NewKafkaConnectV2Resource,
		NewGatewayServiceAccountV2Resource,
	}
}

func (p *ConduktorProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *ConduktorProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version, commit, date string) func() provider.Provider {
	return func() provider.Provider {
		return &ConduktorProvider{
			version: version,
			commit:  commit,
			date:    date,
		}
	}
}
