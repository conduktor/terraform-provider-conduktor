package gateway_service_account_v2

import (
	"context"
	"fmt"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	gwserviceaccounts "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_service_account_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *gwserviceaccounts.GatewayServiceAccountV2Model) (model.GatewayServiceAccountResource, error) {
	externalNames, diag := schema.SetValueToStringArray(ctx, r.Spec.ExternalNames)
	if diag.HasError() {
		return model.GatewayServiceAccountResource{}, mapper.WrapDiagError(diag, "external_names", mapper.FromTerraform)
	}
	if len(externalNames) > 0 {
		if r.Spec.SpecType.ValueString() != "EXTERNAL" {
			return model.GatewayServiceAccountResource{}, fmt.Errorf("external_names only configurable when spec.type = EXTERNAL")
		}

	}

	return model.NewGatewayServiceAccountResource(
		model.GatewayServiceAccountMetadata{
			Name:     r.Name.ValueString(),
			VCluster: r.Vcluster.ValueString(),
		},
		model.GatewayServiceAccountSpec{
			Type:          r.Spec.SpecType.ValueString(),
			ExternalNames: externalNames,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *model.GatewayServiceAccountResource) (gwserviceaccounts.GatewayServiceAccountV2Model, error) {
	// Configuring default value for vcluster
	if r.Metadata.VCluster == "" {
		r.Metadata.VCluster = "passthrough"
	}

	externalNamesList, diag := schema.StringArrayToSetValue(r.Spec.ExternalNames)
	if diag.HasError() {
		return gwserviceaccounts.GatewayServiceAccountV2Model{}, mapper.WrapDiagError(diag, "external_names", mapper.IntoTerraform)
	}

	specValue, diag := gwserviceaccounts.NewSpecValue(
		map[string]attr.Type{
			"type":           basetypes.StringType{},
			"external_names": externalNamesList.Type(ctx),
		},
		map[string]attr.Value{
			"type":           schema.NewStringValue(r.Spec.Type),
			"external_names": externalNamesList,
		},
	)
	if diag.HasError() {
		return gwserviceaccounts.GatewayServiceAccountV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return gwserviceaccounts.GatewayServiceAccountV2Model{
		Name:     types.StringValue(r.Metadata.Name),
		Vcluster: types.StringValue(r.Metadata.VCluster),
		Spec:     specValue,
	}, nil
}
