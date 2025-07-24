package console_service_account_v1

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	serviceacc "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_service_account_v1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *serviceacc.ConsoleServiceAccountV1Model) (console.ServiceAccountResource, error) {
	labels, diag := schema.MapValueToStringMap(ctx, r.Labels)
	if diag.HasError() {
		return console.ServiceAccountResource{}, mapper.WrapDiagError(diag, "labels", mapper.FromTerraform)
	}

	auth, err := objectValueToAuthorization(ctx, r.Spec.Authorization)
	if err != nil {
		return console.ServiceAccountResource{}, mapper.WrapError(err, "authorization", mapper.FromTerraform)
	}

	return console.NewServiceAccountResource(
		console.ServiceAccountMetadata{
			Name:        r.Name.ValueString(),
			Cluster:     r.Cluster.ValueString(),
			Labels:      labels,
			AppInstance: r.AppInstance.ValueString(),
		},
		console.ServiceAccountSpec{
			Authorization: auth,
		},
	), nil

}

func objectValueToAuthorization(ctx context.Context, r basetypes.ObjectValue) (*console.ServiceAccountAuthorization, error) {
	if r.IsNull() {
		return nil, nil
	}

	authValue, diag := serviceacc.NewAuthorizationValue(r.AttributeTypes(ctx), r.Attributes())
	if diag.HasError() {
		return &console.ServiceAccountAuthorization{}, mapper.WrapDiagError(diag, "authorization", mapper.FromTerraform)
	}

	var aiven *console.ServiceAccountAuthAiven = nil
	if schema.AttrIsSet(authValue.Aiven) {
		aivenValue, diag := serviceacc.NewAivenValue(authValue.Aiven.AttributeTypes(ctx), authValue.Aiven.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "authorization.aiven", mapper.FromTerraform)
		}

		aivenACLs, err := setValueToAivenAclArray(ctx, aivenValue.Acls)
		if err != nil {
			return nil, mapper.WrapError(err, "authorization.aiven", mapper.FromTerraform)
		}
		aiven = &console.ServiceAccountAuthAiven{
			ACLS: aivenACLs,
			Type: "AIVEN_ACL",
		}
	}

	var kafka *console.ServiceAccountAuthKafka = nil
	if schema.AttrIsSet(authValue.Kafka) {
		kafkaValue, diag := serviceacc.NewKafkaValue(authValue.Kafka.AttributeTypes(ctx), authValue.Kafka.Attributes())
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "authorization.kafka", mapper.FromTerraform)
		}

		kafkaACLs, err := setValueToKafkaAclArray(ctx, kafkaValue.Acls)
		if err != nil {
			return nil, mapper.WrapError(err, "authorization.aiven", mapper.FromTerraform)
		}
		kafka = &console.ServiceAccountAuthKafka{
			ACLS: kafkaACLs,
			Type: "KAFKA_ACL",
		}
	}

	return &console.ServiceAccountAuthorization{
		Aiven: aiven,
		Kafka: kafka,
	}, nil
}

func setValueToAivenAclArray(ctx context.Context, set basetypes.SetValue) ([]console.ServiceAccountAuthAivenACL, error) {
	acls := make([]console.ServiceAccountAuthAivenACL, 0)
	var diag diag.Diagnostics

	if !set.IsNull() && !set.IsUnknown() {
		var aivenACLs []serviceacc.AivenAclsValue
		diag = set.ElementsAs(ctx, &aivenACLs, false)
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "acls", mapper.FromTerraform)
		}
		for _, p := range aivenACLs {
			acls = append(acls, console.ServiceAccountAuthAivenACL{
				Name:         p.Name.ValueString(),
				ResourceType: p.ResourceType.ValueString(),
				Permission:   p.Permission.ValueString(),
			})
		}
	}
	return acls, nil
}

func setValueToKafkaAclArray(ctx context.Context, set basetypes.SetValue) ([]console.ServiceAccountAuthKafkaACL, error) {
	acls := make([]console.ServiceAccountAuthKafkaACL, 0)
	var diag diag.Diagnostics

	if !set.IsNull() && !set.IsUnknown() {
		var kafkaACLs []serviceacc.KafkaAclsValue
		diag = set.ElementsAs(ctx, &kafkaACLs, false)
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "acls", mapper.FromTerraform)
		}
		for _, p := range kafkaACLs {
			operations, diag := schema.SetValueToStringArray(ctx, p.Operations)
			if diag.HasError() {
				return nil, mapper.WrapDiagError(diag, "acls.operations", mapper.FromTerraform)
			}

			acls = append(acls, console.ServiceAccountAuthKafkaACL{
				Type:           p.KafkaAclsType.ValueString(),
				Name:           p.Name.ValueString(),
				PatternType:    p.PatternType.ValueString(),
				ConnectCluster: p.ConnectCluster.ValueString(),
				Operations:     operations,
				Host:           p.Host.ValueString(),
				Permission:     p.Permission.ValueString(),
			})
		}
	}
	return acls, nil
}
