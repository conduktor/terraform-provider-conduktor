package gateway_interceptor_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	gwinterceptor "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_interceptor_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *gwinterceptor.GatewayInterceptorV2Model) (gateway.GatewayInterceptorResource, error) {
	config, diag := schema.JsonToNormalizedString(ctx, r.Spec.Config)
	if diag.HasError() {
		return gateway.GatewayInterceptorResource{}, mapper.WrapDiagError(diag, "external_names", mapper.FromTerraform)
	}

	scope := gateway.GatewayInterceptorScope{
		Group:    r.Scope.Group.ValueString(),
		VCluster: r.Scope.Vcluster.ValueString(),
		Username: r.Scope.Username.ValueString(),
	}

	return gateway.NewGatewayInterceptorResource(
		gateway.GatewayInterceptorMetadata{
			Name:  r.Name.ValueString(),
			Scope: scope,
		},
		gateway.GatewayInterceptorSpec{
			Comment:     r.Spec.Comment.ValueString(),
			PluginClass: r.Spec.PluginClass.ValueString(),
			Priority:    r.Spec.Priority.ValueInt64(),
			Config:      config,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorResource) (gwinterceptor.GatewayInterceptorV2Model, error) {
	// normalizedJson, diag := schema.StringToNormalizedJson(ctx, r.Spec.Config)
	// if diag.HasError() {
	// 	return gwinterceptor.GatewayInterceptorV2Model{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	// }

	scopeValue := gwinterceptor.ScopeValue{
		Group:    schema.NewStringValue(r.Metadata.Scope.Group),
		Vcluster: schema.NewStringValue(r.Metadata.Scope.VCluster),
		Username: schema.NewStringValue(r.Metadata.Scope.Username),
	}

	specValue, diag := gwinterceptor.NewSpecValue(
		map[string]attr.Type{
			"comment":      basetypes.StringType{},
			"plugin_class": basetypes.StringType{},
			"priority":     basetypes.Int64Type{},
			"config":       basetypes.StringType{},
		},
		map[string]attr.Value{
			"comment":      schema.NewStringValue(r.Spec.Comment),
			"plugin_class": schema.NewStringValue(r.Spec.PluginClass),
			"priority":     schema.NewInt64Value(r.Spec.Priority),
			"config":       schema.NewStringValue(r.Spec.Config),
		},
	)
	if diag.HasError() {
		return gwinterceptor.GatewayInterceptorV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return gwinterceptor.GatewayInterceptorV2Model{
		Name:  types.StringValue(r.Metadata.Name),
		Scope: scopeValue,
		Spec:  specValue,
	}, nil
}
