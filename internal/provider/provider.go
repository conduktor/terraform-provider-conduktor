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
	var config schema.ConduktorModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	consoleUrl := schemaUtils.GetStringConfig(config.ConsoleUrl, []string{"CDK_BASE_URL", "CDK_CONSOLE_URL"})
	apiToken := schemaUtils.GetStringConfig(config.ApiToken, []string{"CDK_API_TOKEN", "CDK_API_KEY"})
	adminEmail := schemaUtils.GetStringConfig(config.AdminEmail, []string{"CDK_ADMIN_EMAIL"})
	adminPassword := schemaUtils.GetStringConfig(config.AdminPassword, []string{"CDK_ADMIN_PASSWORD"})
	cert := schemaUtils.GetStringConfig(config.Cert, []string{"CDK_CERT"})
	cacert := schemaUtils.GetStringConfig(config.Cacert, []string{"CDK_CACERT"})
	key := schemaUtils.GetStringConfig(config.Key, []string{"CDK_KEY"})
	insecure := schemaUtils.GetBooleanConfig(config.Insecure, []string{"CDK_INSECURE"}, false)

	// Validate mandatory configurations
	if consoleUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("console_url"),
			"Missing Console URL",
			"The provider cannot create the Console API client as there is a missing or empty value for the Console URL. "+
				"Set the host value in the configuration or use the CDK_BASE_URL or CDK_CONSOLE_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

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

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "console_url", consoleUrl)
	ctx = tflog.SetField(ctx, "api_token", apiToken)
	ctx = tflog.SetField(ctx, "admin_email", adminEmail)
	ctx = tflog.SetField(ctx, "admin_password", adminPassword)
	ctx = tflog.SetField(ctx, "cert", cert)
	ctx = tflog.SetField(ctx, "cacert", cacert)
	ctx = tflog.SetField(ctx, "key", key)
	ctx = tflog.SetField(ctx, "insecure", insecure)
	// Avoid leaking sensitive information in logs
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "api_token")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "admin_password")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "key")
	tflog.Debug(ctx, "Creating Conduktor Console client")

	// Create client
	apiClient, err := client.Make(ctx,
		client.ApiParameter{
			ApiKey:      apiToken,
			BaseUrl:     consoleUrl,
			Key:         key,
			Cert:        cert,
			Cacert:      cacert,
			CdkUser:     adminEmail,
			CdkPassword: adminPassword,
			Insecure:    insecure,
		},
		p.version,
	)
	if err != nil {
		resp.Diagnostics.AddError("Could not create the Conduktor Console API client", err.Error())
		return
	}

	data := &ProviderData{
		Client: apiClient,
	}
	resp.DataSourceData = data
	resp.ResourceData = data

	tflog.Info(ctx, "Configured Conduktor Console client", map[string]any{"success": true})
}

func (p *ConduktorProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserV2Resource,
		NewGroupV2Resource,
		NewGenericResource,
		NewKafkaClusterV2Resource,
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
