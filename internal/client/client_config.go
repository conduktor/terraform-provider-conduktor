package client

import (
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/provider_conduktor"
)

type ApiParameter struct {
	ApiKey        string
	BaseUrl       string
	CdkUser       string
	CdkPassword   string
	TLSParameters TLSParameters
}

type TLSParameters struct {
	Key      string
	Cert     string
	Cacert   string
	Insecure bool
}

func LoadConfig(providerInputConfig schema.ConduktorModel, mode Mode) ApiParameter {
	var apiParameter ApiParameter

	switch mode {
	case CONSOLE:
		apiParameter.BaseUrl = schemaUtils.GetStringConfig(providerInputConfig.BaseUrl, []string{"CDK_CONSOLE_BASE_URL", "CDK_BASE_URL"})
		apiParameter.ApiKey = schemaUtils.GetStringConfig(providerInputConfig.ApiToken, []string{"CDK_API_TOKEN", "CDK_API_KEY"})
		apiParameter.CdkUser = schemaUtils.GetStringConfig(providerInputConfig.AdminUser, []string{"CDK_CONSOLE_USER", "CDK_ADMIN_EMAIL", "CDK_ADMIN_USER"})
		apiParameter.CdkPassword = schemaUtils.GetStringConfig(providerInputConfig.AdminPassword, []string{"CDK_CONSOLE_PASSWORD", "CDK_ADMIN_PASSWORD"})
		apiParameter.TLSParameters.Key = schemaUtils.GetStringConfig(providerInputConfig.Cert, []string{"CDK_CONSOLE_CERT", "CDK_CERT"})
		apiParameter.TLSParameters.Cacert = schemaUtils.GetStringConfig(providerInputConfig.Cacert, []string{"CDK_CONSOLE_CACERT", "CDK_CACERT"})
		apiParameter.TLSParameters.Key = schemaUtils.GetStringConfig(providerInputConfig.Key, []string{"CDK_CONSOLE_KEY", "CDK_KEY"})
		apiParameter.TLSParameters.Insecure = schemaUtils.GetBooleanConfig(providerInputConfig.Insecure, []string{"CDK_CONSOLE_INSECURE", "CDK_INSECURE"}, false)

	case GATEWAY:
		apiParameter.BaseUrl = schemaUtils.GetStringConfig(providerInputConfig.BaseUrl, []string{"CDK_GATEWAY_BASE_URL", "CDK_BASE_URL"})
		apiParameter.CdkUser = schemaUtils.GetStringConfig(providerInputConfig.AdminUser, []string{"CDK_GATEWAY_USER", "CDK_ADMIN_USER"})
		apiParameter.CdkPassword = schemaUtils.GetStringConfig(providerInputConfig.AdminPassword, []string{"CDK_GATEWAY_PASSWORD", "CDK_ADMIN_PASSWORD"})
		apiParameter.TLSParameters.Cacert = schemaUtils.GetStringConfig(providerInputConfig.Cert, []string{"CDK_GATEWAY_CERT", "CDK_CERT"})
		apiParameter.TLSParameters.Cacert = schemaUtils.GetStringConfig(providerInputConfig.Cacert, []string{"CDK_GATEWAY_CACERT", "CDK_CACERT"})
		apiParameter.TLSParameters.Key = schemaUtils.GetStringConfig(providerInputConfig.Key, []string{"CDK_GATEWAY_KEY", "CDK_KEY"})
		apiParameter.TLSParameters.Insecure = schemaUtils.GetBooleanConfig(providerInputConfig.Insecure, []string{"CDK_GATEWAY_INSECURE", "CDK_INSECURE"}, false)

	}

	return apiParameter
}
