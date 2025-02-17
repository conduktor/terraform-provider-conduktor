package console_kafka_cluster_v2

import (
	"context"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_kafka_cluster_v2"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *schema.ConsoleKafkaClusterV2Model) (console.KafkaClusterResource, error) {

	labels, diag := schemaUtils.MapValueToStringMap(ctx, r.Labels)
	if diag.HasError() {
		return console.KafkaClusterResource{}, mapper.WrapDiagError(diag, "labels", mapper.FromTerraform)
	}

	spec, err := specTFToInternalModel(ctx, &r.Spec)
	if err != nil {
		return console.KafkaClusterResource{}, err
	}

	return console.NewKafkaClusterResource(
		r.Name.ValueString(),
		labels,
		spec,
	), nil
}

func specTFToInternalModel(ctx context.Context, r *schema.SpecValue) (console.KafkaClusterSpec, error) {
	properties, diag := schemaUtils.MapValueToStringMap(ctx, r.Properties)
	if diag.HasError() {
		return console.KafkaClusterSpec{}, mapper.WrapDiagError(diag, "properties", mapper.FromTerraform)
	}

	kafkaFlavor, err := kafkaFlavorTFToInternalModel(ctx, &r.KafkaFlavor)
	if err != nil {
		return console.KafkaClusterSpec{}, err
	}

	schemaRegistry, err := schemaRegistryTFToInternalModel(ctx, &r.SchemaRegistry)
	if err != nil {
		return console.KafkaClusterSpec{}, err
	}

	return console.KafkaClusterSpec{
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

func kafkaFlavorTFToInternalModel(ctx context.Context, r *basetypes.ObjectValue) (*console.KafkaFlavor, error) {
	if r.IsNull() {
		return nil, nil
	}

	kafkaFlavorValue, diag := schema.NewKafkaFlavorValue(r.AttributeTypes(ctx), r.Attributes())
	if diag.HasError() {
		return nil, mapper.WrapDiagError(diag, "kafka_flavor", mapper.FromTerraform)
	}

	var Aiven *console.Aiven = nil
	if schemaUtils.AttrIsSet(kafkaFlavorValue.Aiven) {
		aivenValue, diag := schema.NewAivenValue(kafkaFlavorValue.Aiven.AttributeTypes(ctx), kafkaFlavorValue.Aiven.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "kafka_flavor.aiven", mapper.FromTerraform)
		}
		Aiven = &console.Aiven{
			Type:        "Aiven",
			ApiToken:    aivenValue.ApiToken.ValueString(),
			Project:     aivenValue.Project.ValueString(),
			ServiceName: aivenValue.ServiceName.ValueString(),
		}
	}

	var Confluent *console.Confluent = nil
	if schemaUtils.AttrIsSet(kafkaFlavorValue.Confluent) {
		confluentValue, diag := schema.NewConfluentValue(kafkaFlavorValue.Confluent.AttributeTypes(ctx), kafkaFlavorValue.Confluent.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "kafka_flavor.confluent", mapper.FromTerraform)
		}
		Confluent = &console.Confluent{
			Type:                   "Confluent",
			Key:                    confluentValue.Key.ValueString(),
			Secret:                 confluentValue.Secret.ValueString(),
			ConfluentEnvironmentId: confluentValue.ConfluentEnvironmentId.ValueString(),
			ConfluentClusterId:     confluentValue.ConfluentClusterId.ValueString(),
		}
	}

	var Gateway *console.Gateway = nil
	if schemaUtils.AttrIsSet(kafkaFlavorValue.Gateway) {
		gatewayValue, diag := schema.NewGatewayValue(kafkaFlavorValue.Gateway.AttributeTypes(ctx), kafkaFlavorValue.Gateway.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "kafka_flavor.gateway", mapper.FromTerraform)
		}
		Gateway = &console.Gateway{
			Type:                       "Gateway",
			Url:                        gatewayValue.Url.ValueString(),
			User:                       gatewayValue.User.ValueString(),
			Password:                   gatewayValue.Password.ValueString(),
			VirtualCluster:             gatewayValue.VirtualCluster.ValueString(),
			IgnoreUntrustedCertificate: gatewayValue.IgnoreUntrustedCertificate.ValueBool(),
		}
	}
	return &console.KafkaFlavor{
		Aiven:     Aiven,
		Confluent: Confluent,
		Gateway:   Gateway,
	}, nil
}

func confluentLikeSchemaRegistrySecurityTFToInternalModel(ctx context.Context, r *schema.ConfluentSecurityValue) (model.ConfluentLikeSchemaRegistrySecurity, error) {
	var BasicAuth *model.BasicAuth = nil
	if schemaUtils.AttrIsSet(r.BasicAuth) {
		basicAuthValue, diag := schema.NewBasicAuthValue(r.BasicAuth.AttributeTypes(ctx), r.BasicAuth.Attributes())
		if diag.HasError() {
			return model.ConfluentLikeSchemaRegistrySecurity{}, mapper.WrapDiagError(diag, "schema_registry.security.basic_auth", mapper.FromTerraform)
		}
		BasicAuth = &model.BasicAuth{
			Type:     "BasicAuth",
			UserName: basicAuthValue.Username.ValueString(),
			Password: basicAuthValue.Password.ValueString(),
		}
	}

	var BearerToken *model.BearerToken = nil
	if schemaUtils.AttrIsSet(r.BearerToken) {
		bearerTokenValue, diag := schema.NewBearerTokenValue(r.BearerToken.AttributeTypes(ctx), r.BearerToken.Attributes())
		if diag.HasError() {
			return model.ConfluentLikeSchemaRegistrySecurity{}, mapper.WrapDiagError(diag, "schema_registry.security.bearer_token", mapper.FromTerraform)
		}
		BearerToken = &model.BearerToken{
			Type:  "BearerToken",
			Token: bearerTokenValue.Token.ValueString(),
		}
	}

	var SSLAuth *model.SSLAuth = nil
	if schemaUtils.AttrIsSet(r.SslAuth) {
		sslAuthValue, diag := schema.NewSslAuthValue(r.SslAuth.AttributeTypes(ctx), r.SslAuth.Attributes())
		if diag.HasError() {
			return model.ConfluentLikeSchemaRegistrySecurity{}, mapper.WrapDiagError(diag, "schema_registry.security.ssl_auth", mapper.FromTerraform)
		}
		SSLAuth = &model.SSLAuth{
			Type:             "SSLAuth",
			Key:              sslAuthValue.Key.ValueString(),
			CertificateChain: sslAuthValue.CertificateChain.ValueString(),
		}
	}

	var NoSecurity *model.NoSecurity = nil
	if BasicAuth == nil && BearerToken == nil && SSLAuth == nil {
		NoSecurity = &model.NoSecurity{
			Type: "NoSecurity",
		}
	}

	return model.ConfluentLikeSchemaRegistrySecurity{
		BasicAuth:   BasicAuth,
		BearerToken: BearerToken,
		SSLAuth:     SSLAuth,
		NoSecurity:  NoSecurity,
	}, nil
}

func amazonSchemaRegistrySecurityTFToInternalModel(ctx context.Context, r *schema.GlueSecurityValue) (model.AmazonSecurity, error) {

	var credentials *model.Credentials = nil
	if schemaUtils.AttrIsSet(r.Credentials) {
		credentialsValue, diag := schema.NewCredentialsValue(r.Credentials.AttributeTypes(ctx), r.Credentials.Attributes())
		if diag.HasError() {
			return model.AmazonSecurity{}, mapper.WrapDiagError(diag, "schema_registry.security.credentials", mapper.FromTerraform)
		}
		credentials = &model.Credentials{
			Type:        "Credentials",
			AccessKeyId: credentialsValue.AccessKeyId.ValueString(),
			SecretKey:   credentialsValue.SecretKey.ValueString(),
		}
	}

	var fromContext *model.FromContext = nil
	if schemaUtils.AttrIsSet(r.FromContext) {
		fromContextValue, diag := schema.NewFromContextValue(r.FromContext.AttributeTypes(ctx), r.FromContext.Attributes())
		if diag.HasError() {
			return model.AmazonSecurity{}, mapper.WrapDiagError(diag, "schema_registry.security.from_context", mapper.FromTerraform)
		}
		fromContext = &model.FromContext{
			Type:    "FromContext",
			Profile: fromContextValue.Profile.ValueString(),
		}
	}

	var fromRole *model.FromRole = nil
	if schemaUtils.AttrIsSet(r.FromRole) {
		fromRoleValue, diag := schema.NewFromRoleValue(r.FromRole.AttributeTypes(ctx), r.FromRole.Attributes())
		if diag.HasError() {
			return model.AmazonSecurity{}, mapper.WrapDiagError(diag, "schema_registry.security.from_role", mapper.FromTerraform)
		}
		fromRole = &model.FromRole{
			Type: "FromRole",
			Role: fromRoleValue.Role.ValueString(),
		}
	}

	var iamAnywhere *model.IAMAnywhere = nil
	if schemaUtils.AttrIsSet(r.IamAnywhere) {
		iamAnywhereValue, diag := schema.NewIamAnywhereValue(r.IamAnywhere.AttributeTypes(ctx), r.IamAnywhere.Attributes())
		if diag.HasError() {
			return model.AmazonSecurity{}, mapper.WrapDiagError(diag, "schema_registry.security.iam_anywhere", mapper.FromTerraform)
		}
		iamAnywhere = &model.IAMAnywhere{
			Type:           "IAMAnywhere",
			TrustAnchorArn: iamAnywhereValue.TrustAnchorArn.ValueString(),
			ProfileArn:     iamAnywhereValue.ProfileArn.ValueString(),
			RoleArn:        iamAnywhereValue.RoleArn.ValueString(),
			Certificate:    iamAnywhereValue.Certificate.ValueString(),
			PrivateKey:     iamAnywhereValue.PrivateKey.ValueString(),
		}
	}

	return model.AmazonSecurity{
		Credentials: credentials,
		FromContext: fromContext,
		FromRole:    fromRole,
		IAMAnywhere: iamAnywhere,
	}, nil
}

func schemaRegistryTFToInternalModel(ctx context.Context, r *basetypes.ObjectValue) (*model.SchemaRegistry, error) {
	if r.IsNull() {
		return nil, nil
	}

	schemaRegistryValue, diag := schema.NewSchemaRegistryValue(r.AttributeTypes(ctx), r.Attributes())
	if diag.HasError() {
		return nil, mapper.WrapDiagError(diag, "schema_registry", mapper.FromTerraform)
	}

	var confluentLike *model.ConfluentLike = nil
	if schemaUtils.AttrIsSet(schemaRegistryValue.ConfluentLike) {
		confluentLikeValue, diag := schema.NewConfluentLikeValue(schemaRegistryValue.ConfluentLike.AttributeTypes(ctx), schemaRegistryValue.ConfluentLike.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "schema_registry.confluent_like", mapper.FromTerraform)
		}

		confluentSecurityValue, diag := schema.NewConfluentSecurityValue(confluentLikeValue.Security.AttributeTypes(ctx), confluentLikeValue.Security.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "schema_registry.confluent_like.security", mapper.FromTerraform)
		}

		security, err := confluentLikeSchemaRegistrySecurityTFToInternalModel(ctx, &confluentSecurityValue)
		if err != nil {
			return nil, err
		}

		confluentLike = &model.ConfluentLike{
			Type:                       "ConfluentLike",
			Url:                        confluentLikeValue.Url.ValueString(),
			Properties:                 confluentLikeValue.Properties.ValueString(),
			IgnoreUntrustedCertificate: confluentLikeValue.IgnoreUntrustedCertificate.ValueBool(),
			Security:                   security,
		}
	}

	var glue *model.Glue = nil
	if schemaUtils.AttrIsSet(schemaRegistryValue.Glue) {
		glueValue, diag := schema.NewGlueValue(schemaRegistryValue.Glue.AttributeTypes(ctx), schemaRegistryValue.Glue.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "schema_registry.glue", mapper.FromTerraform)
		}

		securityValue, diag := schema.NewGlueSecurityValue(glueValue.Security.AttributeTypes(ctx), glueValue.Security.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "schema_registry.glue.security", mapper.FromTerraform)
		}

		security, err := amazonSchemaRegistrySecurityTFToInternalModel(ctx, &securityValue)
		if err != nil {
			return nil, err
		}

		glue = &model.Glue{
			Type:         "Glue",
			RegistryName: glueValue.RegistryName.ValueString(),
			Region:       glueValue.Region.ValueString(),
			Security:     security,
		}
	}

	return &model.SchemaRegistry{
		ConfluentLike: confluentLike,
		Glue:          glue,
	}, nil
}
