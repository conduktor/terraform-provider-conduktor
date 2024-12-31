package gateway_interceptor_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	gwinterceptor "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_interceptor_v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *gwinterceptor.GatewayInterceptorV2Model) (gateway.GatewayInterceptorResource, error) {
	scope := gateway.GatewayInterceptorScope{
		Group:    r.Scope.Group.ValueString(),
		VCluster: r.Scope.Vcluster.ValueString(),
		Username: r.Scope.Username.ValueString(),
	}

	config, err := configTFToInternalModel(ctx, &r.Spec.Config)
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
			Config:      config,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorResource) (gwinterceptor.GatewayInterceptorV2Model, error) {
	// scopeValue := gwinterceptor.ScopeValue{
	// 	Group:    schema.NewStringValue(r.Metadata.Scope.Group),
	// 	Vcluster: schema.NewStringValue(r.Metadata.Scope.VCluster),
	// 	Username: schema.NewStringValue(r.Metadata.Scope.Username),
	// }

	unknownScopeObjectValue, diag := gwinterceptor.NewScopeValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return gwinterceptor.GatewayInterceptorV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	var typesMapScope = unknownScopeObjectValue.AttributeTypes(ctx)
	var valuesMapScope = schema.ValueMapFromTypes(ctx, typesMapScope)

	valuesMapScope["group"] = schema.NewStringValue(r.Metadata.Scope.Group)
	valuesMapScope["vcluster"] = schema.NewStringValue(r.Metadata.Scope.VCluster)
	valuesMapScope["username"] = schema.NewStringValue(r.Metadata.Scope.Username)

	scopeValue, diag := gwinterceptor.NewScopeValue(typesMapScope, valuesMapScope)
	if diag.HasError() {
		return gwinterceptor.GatewayInterceptorV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	unknownSpecObjectValue, diag := gwinterceptor.NewSpecValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return gwinterceptor.GatewayInterceptorV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["comment"] = schema.NewStringValue(r.Spec.Comment)
	valuesMap["plugin_class"] = schema.NewStringValue(r.Spec.PluginClass)
	valuesMap["priority"] = schema.NewInt64Value(r.Spec.Priority)

	config, err := configInternalModelToTerraform(ctx, r.Spec.Config)
	if err != nil {
		return gwinterceptor.GatewayInterceptorV2Model{}, err
	}
	configValue, diag := config.ToObjectValue(ctx)
	if diag.HasError() {
		return gwinterceptor.GatewayInterceptorV2Model{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}
	valuesMap["config"] = configValue

	specValue, diag := gwinterceptor.NewSpecValue(typesMap, valuesMap)
	if diag.HasError() {
		return gwinterceptor.GatewayInterceptorV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return gwinterceptor.GatewayInterceptorV2Model{
		Name:  types.StringValue(r.Metadata.Name),
		Scope: scopeValue,
		Spec:  specValue,
	}, nil
}

func configTFToInternalModel(ctx context.Context, r *basetypes.ObjectValue) (*gateway.GatewayInterceptorConfig, error) {
	if r.IsNull() {
		return nil, nil
	}

	configValue, diag := gwinterceptor.NewConfigValue(r.AttributeTypes(ctx), r.Attributes())
	if diag.HasError() {
		return nil, mapper.WrapDiagError(diag, "config", mapper.FromTerraform)
	}

	return &gateway.GatewayInterceptorConfig{
		VirtualTopic: configValue.VirtualTopic.ValueString(),
		Statement:    configValue.Statement.ValueString(),
	}, nil
}

func scopeInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorScope) (gwinterceptor.ScopeValue, error) {
	if r == nil {
		return gwinterceptor.NewScopeValueNull(), nil
	}

	var unknownSR = gwinterceptor.NewConfigValueUnknown()
	unknownSRObjectValue, diag := unknownSR.ToObjectValue(ctx)
	if diag.HasError() {
		return gwinterceptor.ScopeValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}
	var typesMap = unknownSRObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["vcluster"] = schema.NewStringValue(r.VCluster)
	valuesMap["username"] = schema.NewStringValue(r.Username)
	valuesMap["group"] = schema.NewStringValue(r.Group)

	value, diag := gwinterceptor.NewScopeValue(typesMap, valuesMap)
	if diag.HasError() {
		return gwinterceptor.ScopeValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}
	return value, nil
}

func configInternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorConfig) (gwinterceptor.ConfigValue, error) {
	if r == nil {
		return gwinterceptor.NewConfigValueNull(), nil
	}

	var unknownSR = gwinterceptor.NewConfigValueUnknown()
	unknownSRObjectValue, diag := unknownSR.ToObjectValue(ctx)
	if diag.HasError() {
		return gwinterceptor.ConfigValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}
	var typesMap = unknownSRObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["virtual_topic"] = schema.NewStringValue(r.VirtualTopic)
	valuesMap["statement"] = schema.NewStringValue(r.Statement)

	value, diag := gwinterceptor.NewConfigValue(typesMap, valuesMap)
	if diag.HasError() {
		return gwinterceptor.ConfigValue{}, mapper.WrapDiagError(diag, "config", mapper.IntoTerraform)
	}
	return value, nil
}
