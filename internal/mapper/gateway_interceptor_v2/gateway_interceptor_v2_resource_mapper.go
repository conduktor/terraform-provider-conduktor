package gateway_interceptor_v2

import (
	"context"
	"encoding/json"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	gwinterceptor "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_interceptor_v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TFToInternalModel(ctx context.Context, r *gwinterceptor.GatewayInterceptorV2Model) (gateway.GatewayInterceptorResource, error) {
	scope := gateway.GatewayInterceptorScope{
		Group:    r.Scope.Group.ValueString(),
		VCluster: r.Scope.Vcluster.ValueString(),
		Username: r.Scope.Username.ValueString(),
	}

	var config interface{}
	configStr := r.Spec.Config.ValueString()
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return gateway.GatewayInterceptorResource{}, err
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
			Config:      &config,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorResource) (gwinterceptor.GatewayInterceptorV2Model, error) {
	specValue, err := specInternalModelToTerraform(ctx, r.Spec)
	if err != nil {
		return gwinterceptor.GatewayInterceptorV2Model{}, err
	}

	return gwinterceptor.GatewayInterceptorV2Model{
		Name: types.StringValue(r.Metadata.Name),
		Scope: gwinterceptor.ScopeValue{
			Group:    schema.NewStringValue(r.Metadata.Scope.Group),
			Vcluster: schema.NewStringValue(r.Metadata.Scope.VCluster),
			Username: schema.NewStringValue(r.Metadata.Scope.Username),
		},
		Spec: specValue,
	}, nil
}

func specInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorSpec) (gwinterceptor.SpecValue, error) {
	unknownSpecObjectValue, diag := gwinterceptor.NewSpecValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return gwinterceptor.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["comment"] = schema.NewStringValue(r.Comment)
	valuesMap["plugin_class"] = schema.NewStringValue(r.PluginClass)
	valuesMap["priority"] = schema.NewInt64Value(r.Priority)

	config, err := json.Marshal(r.Config)
	if err != nil {
		return gwinterceptor.SpecValue{}, mapper.WrapDiagError(diag, "spec.config", mapper.IntoTerraform)
	}
	// need a patch in the generated code to use custom type instead of primitive string
	// https://github.com/hashicorp/terraform-plugin-codegen-framework/issues/147
	valuesMap["config"] = schema.NewStringValue(string(config))

	value, diag := gwinterceptor.NewSpecValue(typesMap, valuesMap)
	if diag.HasError() {
		return gwinterceptor.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	return value, nil
}
