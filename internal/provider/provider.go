package provider

import (
	"context"
	"strings"
	"sync"

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

// Mutex to make resource operations sequential.
var resourceMutex sync.Mutex

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
	Mode   client.Mode
	Client *client.Client
}

func (p *ConduktorProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "conduktor"
	resp.Version = p.version
}

func (p *ConduktorProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.ConduktorProviderSchema(ctx)
}

func (p *ConduktorProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var input schema.ConduktorModel
	var apiClient *client.Client
	var data ProviderData
	var err error

	resp.Diagnostics.Append(req.Config.Get(ctx, &input)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mode := strings.ToLower(schemaUtils.GetStringConfig(input.Mode, []string{"CDK_PROVIDER_MODE"}))

	// Data will only contain the mode being either CONSOLE or GATEWAY
	apiParameter, data, resp := p.PreFlightChecks(mode, input, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "mode", mode)
	ctx = tflog.SetField(ctx, "api_token", apiParameter.ApiKey)
	ctx = tflog.SetField(ctx, "base_url", apiParameter.BaseUrl)
	ctx = tflog.SetField(ctx, "admin_user", apiParameter.CdkUser)
	ctx = tflog.SetField(ctx, "admin_password", apiParameter.CdkPassword)
	ctx = tflog.SetField(ctx, "cert", apiParameter.TLSParameters.Cert)
	ctx = tflog.SetField(ctx, "cacert", apiParameter.TLSParameters.Cacert)
	ctx = tflog.SetField(ctx, "key", apiParameter.TLSParameters.Key)
	ctx = tflog.SetField(ctx, "insecure", apiParameter.TLSParameters.Insecure)
	// Avoid leaking sensitive information in logs
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "api_token")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "admin_password")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "key")

	tflog.Debug(ctx, "Creating Conduktor client for: "+string(data.Mode))

	apiClient, err = client.Make(ctx, data.Mode, apiParameter, p.version)

	if err != nil {
		resp.Diagnostics.AddError("Could not create the Conduktor "+string(data.Mode)+" API client", err.Error())
		return
	}

	data.Client = apiClient

	tflog.Info(ctx, "Configured Conduktor "+string(data.Mode)+" client", map[string]any{"success": true})

	resp.DataSourceData = &data
	resp.ResourceData = &data
}

func (p *ConduktorProvider) PreFlightChecks(mode string, input schema.ConduktorModel, resp *provider.ConfigureResponse) (client.ApiParameter, ProviderData, *provider.ConfigureResponse) {
	var data ProviderData
	var apiParameter client.ApiParameter

	switch mode {
	case "console":
		{
			data.Mode = client.CONSOLE
			apiParameter.BaseUrl = schemaUtils.GetStringConfig(input.BaseUrl, []string{"CDK_CONSOLE_BASE_URL", "CDK_BASE_URL"})
			apiParameter.ApiKey = schemaUtils.GetStringConfig(input.ApiToken, []string{"CDK_API_TOKEN", "CDK_API_KEY"})
			apiParameter.CdkUser = schemaUtils.GetStringConfig(input.AdminUser, []string{"CDK_CONSOLE_USER", "CDK_ADMIN_EMAIL", "CDK_ADMIN_USER"})
			apiParameter.CdkPassword = schemaUtils.GetStringConfig(input.AdminPassword, []string{"CDK_CONSOLE_PASSWORD", "CDK_ADMIN_PASSWORD"})
			apiParameter.TLSParameters.Key = schemaUtils.GetStringConfig(input.Cert, []string{"CDK_CONSOLE_CERT", "CDK_CERT"})
			apiParameter.TLSParameters.Cacert = schemaUtils.GetStringConfig(input.Cacert, []string{"CDK_CONSOLE_CACERT", "CDK_CACERT"})
			apiParameter.TLSParameters.Key = schemaUtils.GetStringConfig(input.Key, []string{"CDK_CONSOLE_KEY", "CDK_KEY"})
			apiParameter.TLSParameters.Insecure = schemaUtils.GetBooleanConfig(input.Insecure, []string{"CDK_CONSOLE_INSECURE", "CDK_INSECURE"}, false)

			if apiParameter.ApiKey == "" {
				// We only need to check user and password if no apiToken is provided.
				if apiParameter.CdkUser == "" || apiParameter.CdkPassword == "" {
					details := "The provider cannot create the Console API client as there is a missing or empty value for the API Token and missing or empty values for the admin user and password. \n" +
						"Set either : \n" +
						" - the api_token value in the configuration or use the CDK_API_TOKEN environment variable. \n" +
						" - the admin_user and admin_password value in the configuration or use the CDK_ADMIN_EMAIL and CDK_ADMIN_PASSWORD environment variable. \n" +
						"If either is already set, ensure the value is not empty."

					resp.Diagnostics.AddAttributeError(path.Root("api_token"), "Missing API token", details)
					resp.Diagnostics.AddAttributeError(path.Root("admin_user"), "Missing Admin email", details)
					resp.Diagnostics.AddAttributeError(path.Root("admin_password"), "Missing Admin password", details)
				}
			}
		}
	case "gateway":
		{
			data.Mode = client.GATEWAY
			apiParameter.BaseUrl = schemaUtils.GetStringConfig(input.BaseUrl, []string{"CDK_GATEWAY_BASE_URL", "CDK_BASE_URL"})
			apiParameter.CdkUser = schemaUtils.GetStringConfig(input.AdminUser, []string{"CDK_GATEWAY_USER", "CDK_ADMIN_USER"})
			apiParameter.CdkPassword = schemaUtils.GetStringConfig(input.AdminPassword, []string{"CDK_GATEWAY_PASSWORD", "CDK_ADMIN_PASSWORD"})
			apiParameter.TLSParameters.Cacert = schemaUtils.GetStringConfig(input.Cert, []string{"CDK_GATEWAY_CERT", "CDK_CERT"})
			apiParameter.TLSParameters.Cacert = schemaUtils.GetStringConfig(input.Cacert, []string{"CDK_GATEWAY_CACERT", "CDK_CACERT"})
			apiParameter.TLSParameters.Key = schemaUtils.GetStringConfig(input.Key, []string{"CDK_GATEWAY_KEY", "CDK_KEY"})
			apiParameter.TLSParameters.Insecure = schemaUtils.GetBooleanConfig(input.Insecure, []string{"CDK_GATEWAY_INSECURE", "CDK_INSECURE"}, false)

			if apiParameter.CdkUser == "" || apiParameter.CdkPassword == "" {
				details := "The provider cannot create the Gateway API client as there is a missing or empty value for the admin user and password. \n" +
					"Set both : \n" +
					" - the admin_user value in the configuration or use the CDK_GATEWAY_USER environment variable. \n" +
					" - the admin_password value in the configuration or use the CDK_GATEWAY_PASSWORD environment variable. \n" +
					"If either is already set, ensure the value is not empty."

				resp.Diagnostics.AddAttributeError(path.Root("admin_user"), "Missing Gateway Admin login", details)
				resp.Diagnostics.AddAttributeError(path.Root("admin_password"), "Missing Gateway Admin password", details)
			}
		}
	}

	if apiParameter.BaseUrl == "" {
		details := "The provider cannot create any API client as there is a missing or empty value for the Base URL. \n" +
			"Set: \n" +
			" - the base_url value in the configuration or use the following environment variables: CDK_BASE_URL or CDK_CONSOLE_BASE_URL for Console, CDK_GATEWAY_BASE_URL for Gateway. \n" +
			"If either is already set, ensure the value is not empty."

		resp.Diagnostics.AddAttributeError(path.Root("base_url"), "Missing Console URL", details)
	}

	return apiParameter, data, resp
}

func (p *ConduktorProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserV2Resource,
		NewGroupV2Resource,
		NewGenericResource,
		NewKafkaClusterV2Resource,
		NewKafkaConnectV2Resource,
		NewGatewayServiceAccountV2Resource,
		NewGatewayTokenV2Resource,
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
