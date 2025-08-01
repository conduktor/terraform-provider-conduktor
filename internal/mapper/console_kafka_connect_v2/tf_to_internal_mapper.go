package console_kafka_connect_v2

import (
	"context"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_kafka_connect_v2"
)

func TFToInternalModel(ctx context.Context, r *schema.ConsoleKafkaConnectV2Model) (console.KafkaConnectResource, error) {

	labels, diag := schemaUtils.MapValueToStringMap(ctx, r.Labels)
	if diag.HasError() {
		return console.KafkaConnectResource{}, mapper.WrapDiagError(diag, "labels", mapper.FromTerraform)
	}

	spec, err := specTFToInternalModel(ctx, &r.Spec)
	if err != nil {
		return console.KafkaConnectResource{}, err
	}

	return console.NewKafkaConnectResource(
		r.Name.ValueString(),
		r.Cluster.ValueString(),
		labels,
		spec,
	), nil
}

func specTFToInternalModel(ctx context.Context, r *schema.SpecValue) (console.KafkaConnectSpec, error) {
	headers, diag := schemaUtils.MapValueToStringMap(ctx, r.Headers)
	if diag.HasError() {
		return console.KafkaConnectSpec{}, mapper.WrapDiagError(diag, "headers", mapper.FromTerraform)
	}

	var securityValue = schema.NewSecurityValueNull()
	if !r.Security.IsNull() {
		securityValue, diag = schema.NewSecurityValue(r.Security.AttributeTypes(ctx), r.Security.Attributes())
		if diag.HasError() {
			return console.KafkaConnectSpec{}, mapper.WrapDiagError(diag, "security", mapper.FromTerraform)
		}
	}

	security, err := securityTFToInternalModel(ctx, &securityValue)
	if err != nil {
		return console.KafkaConnectSpec{}, err
	}

	return console.KafkaConnectSpec{
		DisplayName:                r.DisplayName.ValueString(),
		Urls:                       r.Urls.ValueString(),
		IgnoreUntrustedCertificate: r.IgnoreUntrustedCertificate.ValueBool(),
		Headers:                    headers,
		Security:                   security,
	}, nil
}

func securityTFToInternalModel(ctx context.Context, r *schema.SecurityValue) (*console.KafkaConnectSecurity, error) {
	if r.IsNull() {
		return nil, nil
	}

	var basicAuth *console.KafkaConnectBasicAuth = nil
	if schemaUtils.AttrIsSet(r.BasicAuth) {
		basicAuthValue, diag := schema.NewBasicAuthValue(r.BasicAuth.AttributeTypes(ctx), r.BasicAuth.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "security.basic_auth", mapper.FromTerraform)
		}

		basicAuth = &console.KafkaConnectBasicAuth{
			Type:     string(console.BASIC_AUTH),
			Username: basicAuthValue.Username.ValueString(),
			Password: basicAuthValue.Password.ValueString(),
		}
	}

	var bearerToken *console.KafkaConnectBearerToken = nil
	if schemaUtils.AttrIsSet(r.BearerToken) {
		bearerTokenValue, diag := schema.NewBearerTokenValue(r.BearerToken.AttributeTypes(ctx), r.BearerToken.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "security.bearer_token", mapper.FromTerraform)
		}

		bearerToken = &console.KafkaConnectBearerToken{
			Type:  string(console.BEARER_TOKEN),
			Token: bearerTokenValue.Token.ValueString(),
		}
	}

	var sslAuth *console.KafkaConnectSSLAuth = nil
	if schemaUtils.AttrIsSet(r.SslAuth) {
		sslAUthValue, diag := schema.NewSslAuthValue(r.SslAuth.AttributeTypes(ctx), r.SslAuth.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "security.ssl_auth", mapper.FromTerraform)
		}

		sslAuth = &console.KafkaConnectSSLAuth{
			Type:             string(console.SSL_AUTH),
			CertificateChain: sslAUthValue.CertificateChain.ValueString(),
			Key:              sslAUthValue.Key.ValueString(),
		}
	}

	return &console.KafkaConnectSecurity{
		BasicAuth:   basicAuth,
		BearerToken: bearerToken,
		SSLAuth:     sslAuth,
	}, nil
}
