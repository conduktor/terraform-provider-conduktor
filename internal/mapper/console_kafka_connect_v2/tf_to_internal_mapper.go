package console_kafka_connect_v2

import (
	"context"
	"fmt"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_kafka_connect_v2"
	"github.com/conduktor/terraform-provider-conduktor/internal/schema/validation"
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

func securityTFToInternalModel(_ context.Context, r *schema.SecurityValue) (*console.KafkaConnectSecurity, error) {
	if r.IsNull() {
		return nil, nil
	}

	securityType := r.SecurityType.ValueString()
	switch securityType {
	case validation.BasicAuthKafkaConnectSecurity:
		return &console.KafkaConnectSecurity{
			BasicAuth: &console.KafkaConnectBasicAuth{
				Type:     securityType,
				Username: r.Username.ValueString(),
				Password: r.Password.ValueString(),
			},
		}, nil
	case validation.BearerTokenKafkaConnectSecurity:
		return &console.KafkaConnectSecurity{
			BearerToken: &console.KafkaConnectBearerToken{
				Type:  securityType,
				Token: r.Token.ValueString(),
			},
		}, nil
	case validation.SSLAuthKafkaConnectSecurity:
		return &console.KafkaConnectSecurity{
			SSLAuth: &console.KafkaConnectSSLAuth{
				Type:             securityType,
				Key:              r.Key.ValueString(),
				CertificateChain: r.CertificateChain.ValueString(),
			},
		}, nil
	default:
		return &console.KafkaConnectSecurity{}, mapper.WrapError(fmt.Errorf("unsupported SecurityType: %s", securityType), "security", mapper.FromTerraform)
	}
}
