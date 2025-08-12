package console_ksqldb_cluster_v2

import (
	"context"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_ksqldb_cluster_v2"
)

func TFToInternalModel(ctx context.Context, r *schema.ConsoleKsqldbClusterV2Model) (console.KsqlDBClusterResource, error) {

	spec, err := specTFToInternalModel(ctx, &r.Spec)
	if err != nil {
		return console.KsqlDBClusterResource{}, err
	}

	return console.NewKsqlDBClusterResource(
		r.Name.ValueString(),
		r.Cluster.ValueString(),
		spec,
	), nil
}

func specTFToInternalModel(ctx context.Context, r *schema.SpecValue) (console.KsqlDBClusterSpec, error) {
	headers, diag := schemaUtils.MapValueToStringMap(ctx, r.Headers)
	if diag.HasError() {
		return console.KsqlDBClusterSpec{}, mapper.WrapDiagError(diag, "headers", mapper.FromTerraform)
	}

	var securityValue = schema.NewSecurityValueNull()
	if !r.Security.IsNull() {
		securityValue, diag = schema.NewSecurityValue(r.Security.AttributeTypes(ctx), r.Security.Attributes())
		if diag.HasError() {
			return console.KsqlDBClusterSpec{}, mapper.WrapDiagError(diag, "security", mapper.FromTerraform)
		}
	}

	security, err := securityTFToInternalModel(ctx, &securityValue)
	if err != nil {
		return console.KsqlDBClusterSpec{}, err
	}

	return console.KsqlDBClusterSpec{
		DisplayName:                r.DisplayName.ValueString(),
		Url:                        r.Url.ValueString(),
		IgnoreUntrustedCertificate: r.IgnoreUntrustedCertificate.ValueBool(),
		Headers:                    headers,
		Security:                   security,
	}, nil
}

func securityTFToInternalModel(ctx context.Context, r *schema.SecurityValue) (*console.KsqlDBClusterSecurity, error) {
	if r.IsNull() {
		return nil, nil
	}

	var basicAuth *console.KsqlDBClusterBasicAuth = nil
	if schemaUtils.AttrIsSet(r.BasicAuth) {
		basicAuthValue, diag := schema.NewBasicAuthValue(r.BasicAuth.AttributeTypes(ctx), r.BasicAuth.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "security.basic_auth", mapper.FromTerraform)
		}

		basicAuth = &console.KsqlDBClusterBasicAuth{
			Type:     string(console.KSQLDB_BASIC_AUTH),
			Username: basicAuthValue.Username.ValueString(),
			Password: basicAuthValue.Password.ValueString(),
		}
	}

	var bearerToken *console.KsqlDBClusterBearerToken = nil
	if schemaUtils.AttrIsSet(r.BearerToken) {
		bearerTokenValue, diag := schema.NewBearerTokenValue(r.BearerToken.AttributeTypes(ctx), r.BearerToken.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "security.bearer_token", mapper.FromTerraform)
		}

		bearerToken = &console.KsqlDBClusterBearerToken{
			Type:  string(console.KSQLDB_BEARER_TOKEN),
			Token: bearerTokenValue.Token.ValueString(),
		}
	}

	var sslAuth *console.KsqlDBClusterSSLAuth = nil
	if schemaUtils.AttrIsSet(r.SslAuth) {
		sslAUthValue, diag := schema.NewSslAuthValue(r.SslAuth.AttributeTypes(ctx), r.SslAuth.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "security.ssl_auth", mapper.FromTerraform)
		}

		sslAuth = &console.KsqlDBClusterSSLAuth{
			Type:             string(console.KSQLDB_SSL_AUTH),
			CertificateChain: sslAUthValue.CertificateChain.ValueString(),
			Key:              sslAUthValue.Key.ValueString(),
		}
	}

	return &console.KsqlDBClusterSecurity{
		BasicAuth:   basicAuth,
		BearerToken: bearerToken,
		SSLAuth:     sslAuth,
	}, nil
}
