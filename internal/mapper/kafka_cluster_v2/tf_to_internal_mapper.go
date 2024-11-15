package kafka_cluster_v2

import (
	"context"
	"fmt"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_kafka_cluster_v2"
	"github.com/conduktor/terraform-provider-conduktor/internal/schema/validation"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *schema.KafkaClusterV2Model) (model.KafkaClusterResource, error) {

	labels, diag := schemaUtils.MapValueToStringMap(ctx, r.Labels)
	if diag.HasError() {
		return model.KafkaClusterResource{}, mapper.WrapDiagError(diag, "labels", mapper.FromTerraform)
	}

	spec, err := specTFToInternalModel(ctx, &r.Spec)
	if err != nil {
		return model.KafkaClusterResource{}, err
	}

	return model.NewKafkaClusterResource(
		r.Name.ValueString(),
		labels,
		spec,
	), nil
}

func specTFToInternalModel(ctx context.Context, r *schema.SpecValue) (model.KafkaClusterSpec, error) {
	properties, diag := schemaUtils.MapValueToStringMap(ctx, r.Properties)
	if diag.HasError() {
		return model.KafkaClusterSpec{}, mapper.WrapDiagError(diag, "properties", mapper.FromTerraform)
	}

	kafkaFlavor, err := kafkaFlavorTFToInternalModel(ctx, &r.KafkaFlavor)
	if err != nil {
		return model.KafkaClusterSpec{}, err
	}

	schemaRegistry, err := schemaRegistryTFToInternalModel(ctx, &r.SchemaRegistry)
	if err != nil {
		return model.KafkaClusterSpec{}, err
	}

	return model.KafkaClusterSpec{
		DisplayName:                r.DisplayName.ValueString(),
		BootstrapServers:           r.BootstrapServers.ValueString(),
		Color:                      r.Color.ValueString(),
		Icon:                       r.Icon.ValueString(),
		IgnoreUntrustedCertificate: r.IgnoreUntrustedCertificate.ValueBool(),
		Properties:                 properties,
		KafkaFlavor:                kafkaFlavor,
		SchemaRegistry:             schemaRegistry,
	}, nil
}

func kafkaFlavorTFToInternalModel(ctx context.Context, r *basetypes.ObjectValue) (*model.KafkaFlavor, error) {
	if r.IsNull() {
		return nil, nil
	}

	kafkaFlavorValue, diag := schema.NewKafkaFlavorValue(r.AttributeTypes(ctx), r.Attributes())
	if diag.HasError() {
		return nil, mapper.WrapDiagError(diag, "kafka_flavor", mapper.FromTerraform)
	}

	flavorType := kafkaFlavorValue.KafkaFlavorType.ValueString()
	switch flavorType {
	case validation.ConfluentKafkaFlavor:
		return &model.KafkaFlavor{
			Confluent: &model.Confluent{
				Type:                   flavorType,
				Key:                    kafkaFlavorValue.Key.ValueString(),
				Secret:                 kafkaFlavorValue.Secret.ValueString(),
				ConfluentEnvironmentId: kafkaFlavorValue.ConfluentEnvironmentId.ValueString(),
				ConfluentClusterId:     kafkaFlavorValue.ConfluentClusterId.ValueString(),
			},
		}, nil
	case validation.AivenKafkaFlavor:
		return &model.KafkaFlavor{
			Aiven: &model.Aiven{
				Type:        flavorType,
				ApiToken:    kafkaFlavorValue.ApiToken.ValueString(),
				Project:     kafkaFlavorValue.Project.ValueString(),
				ServiceName: kafkaFlavorValue.ServiceName.ValueString(),
			},
		}, nil
	case validation.GatewayKafkaFlavor:
		return &model.KafkaFlavor{
			Gateway: &model.Gateway{
				Type:                       flavorType,
				Url:                        kafkaFlavorValue.Url.ValueString(),
				User:                       kafkaFlavorValue.User.ValueString(),
				Password:                   kafkaFlavorValue.Password.ValueString(),
				VirtualCluster:             kafkaFlavorValue.VirtualCluster.ValueString(),
				IgnoreUntrustedCertificate: kafkaFlavorValue.IgnoreUntrustedCertificate.ValueBool(),
			},
		}, nil
	default:
		return nil, mapper.WrapError(fmt.Errorf("unsupported KafkaFlavorType: %s", flavorType), "kafka_flavor", mapper.FromTerraform)
	}
}

func confluentLikeSchemaRegistrySecurityTFToInternalModel(_ context.Context, r *schema.SecurityValue) (model.ConfluentLikeSchemaRegistrySecurity, error) {
	securityType := r.SecurityType.ValueString()
	switch securityType {
	case validation.NoSecuritySchemaRegistrySecurity:
		return model.ConfluentLikeSchemaRegistrySecurity{
			NoSecurity: &model.NoSecurity{
				Type: securityType,
			},
		}, nil
	case validation.BasicAuthSchemaRegistrySecurity:
		return model.ConfluentLikeSchemaRegistrySecurity{
			BasicAuth: &model.BasicAuth{
				Type:     securityType,
				UserName: r.Username.ValueString(),
				Password: r.Password.ValueString(),
			},
		}, nil
	case validation.BearerTokenSchemaRegistrySecurity:
		return model.ConfluentLikeSchemaRegistrySecurity{
			BearerToken: &model.BearerToken{
				Type:  securityType,
				Token: r.Token.ValueString(),
			},
		}, nil
	case validation.SSLAuthSchemaRegistrySecurity:
		return model.ConfluentLikeSchemaRegistrySecurity{
			SSLAuth: &model.SSLAuth{
				Type:             securityType,
				Key:              r.Key.ValueString(),
				CertificateChain: r.CertificateChain.ValueString(),
			},
		}, nil
	default:
		return model.ConfluentLikeSchemaRegistrySecurity{}, mapper.WrapError(fmt.Errorf("unsupported ConfluentLike SecurityType: %s", securityType), "schema_registry.security", mapper.FromTerraform)
	}
}

func amazonSchemaRegistrySecurityTFToInternalModel(_ context.Context, r *schema.SecurityValue) (model.AmazonSecurity, error) {
	securityType := r.SecurityType.ValueString()
	switch securityType {
	case validation.CredentialsSchemaRegistrySecurity:
		return model.AmazonSecurity{
			Credentials: &model.Credentials{
				Type:        securityType,
				AccessKeyId: r.AccessKeyId.ValueString(),
				SecretKey:   r.SecretKey.ValueString(),
			},
		}, nil
	case validation.FromContextSchemaRegistrySecurity:
		return model.AmazonSecurity{
			FromContext: &model.FromContext{
				Type:    securityType,
				Profile: r.Profile.ValueString(),
			},
		}, nil
	case validation.FromRoleSchemaRegistrySecurity:
		return model.AmazonSecurity{
			FromRole: &model.FromRole{
				Type: securityType,
				Role: r.Role.ValueString(),
			},
		}, nil
	case validation.IAMAnywhereSchemaRegistrySecurity:
		return model.AmazonSecurity{
			IAMAnywhere: &model.IAMAnywhere{
				Type:           securityType,
				TrustAnchorArn: r.TrustAnchorArn.ValueString(),
				ProfileArn:     r.ProfileArn.ValueString(),
				RoleArn:        r.RoleArn.ValueString(),
				Certificate:    r.Certificate.ValueString(),
				PrivateKey:     r.PrivateKey.ValueString(),
			},
		}, nil
	default:
		return model.AmazonSecurity{}, mapper.WrapError(fmt.Errorf("unsupported Amazon SecurityType: %s", securityType), "schema_registry.security", mapper.FromTerraform)
	}
}

func schemaRegistryTFToInternalModel(ctx context.Context, r *basetypes.ObjectValue) (*model.SchemaRegistry, error) {
	if r.IsNull() {
		return nil, nil
	}

	schemaRegistryValue, diag := schema.NewSchemaRegistryValue(r.AttributeTypes(ctx), r.Attributes())
	if diag.HasError() {
		return nil, mapper.WrapDiagError(diag, "schema_registry", mapper.FromTerraform)
	}
	schemaRegistryType := schemaRegistryValue.SchemaRegistryType.ValueString()

	securityValue, diag := schema.NewSecurityValue(schemaRegistryValue.Security.AttributeTypes(ctx), schemaRegistryValue.Security.Attributes())
	if diag.HasError() {
		return nil, mapper.WrapDiagError(diag, "schema_registry.security", mapper.FromTerraform)
	}

	switch schemaRegistryType {
	case validation.ConfluentLikeSchemaRegistry:
		security, err := confluentLikeSchemaRegistrySecurityTFToInternalModel(ctx, &securityValue)
		if err != nil {
			return nil, err
		}

		return &model.SchemaRegistry{
			ConfluentLike: &model.ConfluentLike{
				Type:                       schemaRegistryType,
				Url:                        schemaRegistryValue.Url.ValueString(),
				Properties:                 schemaRegistryValue.Properties.ValueString(),
				IgnoreUntrustedCertificate: schemaRegistryValue.IgnoreUntrustedCertificate.ValueBool(),
				Security:                   security,
			},
		}, nil
	case validation.GlueSchemaRegistry:
		security, err := amazonSchemaRegistrySecurityTFToInternalModel(ctx, &securityValue)
		if err != nil {
			return nil, err
		}

		return &model.SchemaRegistry{
			Glue: &model.Glue{
				Type:         schemaRegistryType,
				RegistryName: schemaRegistryValue.RegistryName.ValueString(),
				Region:       schemaRegistryValue.Region.ValueString(),
				Security:     security,
			},
		}, nil
	default:
		return nil, mapper.WrapError(fmt.Errorf("unsupported SchemaRegistryType: %s", schemaRegistryType), "schema_registry", mapper.FromTerraform)
	}
}
