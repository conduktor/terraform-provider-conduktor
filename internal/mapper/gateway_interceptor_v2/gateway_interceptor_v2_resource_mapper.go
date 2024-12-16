package gateway_interceptor_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	gwinterceptor "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_interceptor_v2"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *gwinterceptor.GatewayInterceptorV2Model) (gateway.GatewayInterceptorResource, error) {
	config, diag := schema.NormalizedJsonToString(ctx, r.Spec.Config)
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
			Config: config,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorResource) (gwinterceptor.GatewayInterceptorV2Model, error) {
	normalizedJson, diag := schema.StringToNormalizedJson(ctx, r.Spec.Config)
	if diag.HasError() {
		return gwinterceptor.GatewayInterceptorV2Model{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}

	scopeValue, diag := gwinterceptor.NewScopeValue(
		map[string]attr.Type{
			"group":    basetypes.StringType{},
			"vCluster": basetypes.StringType{},
			"username": basetypes.StringType{},
		},
		map[string]attr.Value{
			"group":    schema.NewStringValue(r.Metadata.Scope.Group),
			"vCluster": schema.NewStringValue(r.Metadata.Scope.VCluster),
			"username": schema.NewStringValue(r.Metadata.Scope.Username),
		},
	)

	specValue, diag := gwinterceptor.NewSpecValue(
		map[string]attr.Type{
			"comment":     basetypes.StringType{},
			"pluginClass": basetypes.StringType{},
			"priority":    basetypes.Int32Type{},
			"config":      jsontypes.NormalizedType{},
		},
		map[string]attr.Value{
			"comment":     schema.NewStringValue(r.Spec.Comment),
			"pluginClass": schema.NewStringValue(r.Spec.PluginClass),
			"priority":    schema.NewInt32Value(r.Spec.Priority),
			"config":      normalizedJson,
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
