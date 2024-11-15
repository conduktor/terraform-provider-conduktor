package kafka_connect_v2

import (
	"context"
	"fmt"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_kafka_connect_v2"
	"github.com/conduktor/terraform-provider-conduktor/internal/schema/validation"
)

func TFToInternalModel(ctx context.Context, r *schema.KafkaConnectV2Model) (model.KafkaConnectResource, error) {

	labels, diag := schemaUtils.MapValueToStringMap(ctx, r.Labels)
	if diag.HasError() {
		return model.KafkaConnectResource{}, mapper.WrapDiagError(diag, "labels", mapper.FromTerraform)
	}

	spec, err := specTFToInternalModel(ctx, &r.Spec)
	if err != nil {
		return model.KafkaConnectResource{}, err
	}

	return model.NewKafkaConnectResource(
		r.Name.ValueString(),
		r.Cluster.ValueString(),
		labels,
		spec,
	), nil
}

func specTFToInternalModel(ctx context.Context, r *schema.SpecValue) (model.KafkaConnectSpec, error) {
	headers, diag := schemaUtils.MapValueToStringMap(ctx, r.Headers)
	if diag.HasError() {
		return model.KafkaConnectSpec{}, mapper.WrapDiagError(diag, "headers", mapper.FromTerraform)
	}

	var securityValue = schema.NewSecurityValueNull()
	if !r.Security.IsNull() {
		securityValue, diag = schema.NewSecurityValue(r.Security.AttributeTypes(ctx), r.Security.Attributes())
		if diag.HasError() {
			return model.KafkaConnectSpec{}, mapper.WrapDiagError(diag, "security", mapper.FromTerraform)
		}
	}

	security, err := securityTFToInternalModel(ctx, &securityValue)
	if err != nil {
		return model.KafkaConnectSpec{}, err
	}

	return model.KafkaConnectSpec{
		DisplayName:                r.DisplayName.ValueString(),
		Urls:                       r.Urls.ValueString(),
		IgnoreUntrustedCertificate: r.IgnoreUntrustedCertificate.ValueBool(),
		Headers:                    headers,
		Security:                   security,
	}, nil
}

func securityTFToInternalModel(_ context.Context, r *schema.SecurityValue) (*model.KafkaConnectSecurity, error) {
	if r.IsNull() {
		return nil, nil
	}

	securityType := r.SecurityType.ValueString()
	switch securityType {
	case validation.BasicAuthKafkaConnectSecurity:
		return &model.KafkaConnectSecurity{
			BasicAuth: &model.KafkaConnectBasicAuth{
				Type:     securityType,
				Username: r.Username.ValueString(),
				Password: r.Password.ValueString(),
			},
		}, nil
	case validation.BearerTokenKafkaConnectSecurity:
		return &model.KafkaConnectSecurity{
			BearerToken: &model.KafkaConnectBearerToken{
				Type:  securityType,
				Token: r.Token.ValueString(),
			},
		}, nil
	case validation.SSLAuthKafkaConnectSecurity:
		return &model.KafkaConnectSecurity{
			SSLAuth: &model.KafkaConnectSSLAuth{
				Type:             securityType,
				Key:              r.Key.ValueString(),
				CertificateChain: r.CertificateChain.ValueString(),
			},
		}, nil
	default:
		return &model.KafkaConnectSecurity{}, mapper.WrapError(fmt.Errorf("unsupported SecurityType: %s", securityType), "security", mapper.FromTerraform)
	}
}
