// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package provider_conduktor

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

func ConduktorProviderSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"admin_email": schema.StringAttribute{
				Optional:            true,
				Description:         "The email of the admin user. May be set using environment variable `CDK_ADMIN_EMAIL`. Required if admin_password is set. If not provided, the API token will be used to authenticate.",
				MarkdownDescription: "The email of the admin user. May be set using environment variable `CDK_ADMIN_EMAIL`. Required if admin_password is set. If not provided, the API token will be used to authenticate.",
			},
			"admin_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         "The password of the admin user. May be set using environment variable `CDK_ADMIN_PASSWORD`. Required if admin_email is set. If not provided, the API token will be used to authenticater.",
				MarkdownDescription: "The password of the admin user. May be set using environment variable `CDK_ADMIN_PASSWORD`. Required if admin_email is set. If not provided, the API token will be used to authenticater.",
			},
			"api_token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         "The API token to authenticate with the Conduktor API. May be set using environment variable `CDK_API_TOKEN` or `CDK_API_KEY`. If not provided, admin_email and admin_password will be used to authenticate. See [documentation](https://docs.conduktor.io/platform/reference/api-reference/#generate-an-api-key) for more information.",
				MarkdownDescription: "The API token to authenticate with the Conduktor API. May be set using environment variable `CDK_API_TOKEN` or `CDK_API_KEY`. If not provided, admin_email and admin_password will be used to authenticate. See [documentation](https://docs.conduktor.io/platform/reference/api-reference/#generate-an-api-key) for more information.",
			},
			"cacert": schema.StringAttribute{
				Optional:            true,
				Description:         "Root CA certificate in PEM format to verify the Conduktor Console certificate. May be set using environment variable `CDK_CACERT`. If not provided, the system's root CA certificates will be used.",
				MarkdownDescription: "Root CA certificate in PEM format to verify the Conduktor Console certificate. May be set using environment variable `CDK_CACERT`. If not provided, the system's root CA certificates will be used.",
			},
			"cert": schema.StringAttribute{
				Optional:            true,
				Description:         "Cert in PEM format to authenticate using client certificates. May be set using environment variable `CDK_CERT`. Must be used with key. If key is provided, cert is required. Useful when Console behind a reverse proxy with client certificate authentication.",
				MarkdownDescription: "Cert in PEM format to authenticate using client certificates. May be set using environment variable `CDK_CERT`. Must be used with key. If key is provided, cert is required. Useful when Console behind a reverse proxy with client certificate authentication.",
			},
			"console_url": schema.StringAttribute{
				Optional:            true,
				Description:         "The URL of the Conduktor Console. May be set using environment variable `CDK_BASE_URL` or `CDK_CONSOLE_URL`. Required either here or in the environment.",
				MarkdownDescription: "The URL of the Conduktor Console. May be set using environment variable `CDK_BASE_URL` or `CDK_CONSOLE_URL`. Required either here or in the environment.",
			},
			"insecure": schema.BoolAttribute{
				Optional:            true,
				Description:         "Skip TLS verification flag. May be set using environment variable `CDK_INSECURE`.",
				MarkdownDescription: "Skip TLS verification flag. May be set using environment variable `CDK_INSECURE`.",
			},
			"key": schema.StringAttribute{
				Optional:            true,
				Description:         "Key in PEM format to authenticate using client certificates. May be set using environment variable `CDK_KEY`. Must be used with cert. If cert is provided, key is required. Useful when Console behind a reverse proxy with client certificate authentication.",
				MarkdownDescription: "Key in PEM format to authenticate using client certificates. May be set using environment variable `CDK_KEY`. Must be used with cert. If cert is provided, key is required. Useful when Console behind a reverse proxy with client certificate authentication.",
			},
		},
	}
}

type ConduktorModel struct {
	AdminEmail    types.String `tfsdk:"admin_email"`
	AdminPassword types.String `tfsdk:"admin_password"`
	ApiToken      types.String `tfsdk:"api_token"`
	Cacert        types.String `tfsdk:"cacert"`
	Cert          types.String `tfsdk:"cert"`
	ConsoleUrl    types.String `tfsdk:"console_url"`
	Insecure      types.Bool   `tfsdk:"insecure"`
	Key           types.String `tfsdk:"key"`
}