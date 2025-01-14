package gateway_interceptor_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	gwinterceptor "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_interceptor_encryption_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *gwinterceptor.GatewayInterceptorEncryptionV2Model) (gateway.GatewayInterceptorEncryptionResource, error) {
	scope := gateway.GatewayInterceptorEncryptionScope{
		Group:    r.Scope.Group.ValueString(),
		VCluster: r.Scope.Vcluster.ValueString(),
		Username: r.Scope.Username.ValueString(),
	}

	config, err := ObjectValueToInterceptorEncryptionConfig(ctx, &r.Spec.Config)
	if err != nil {
		return gateway.GatewayInterceptorEncryptionResource{}, err
	}

	return gateway.NewGatewayInterceptorEncryptionResource(
		gateway.GatewayInterceptorEncryptionMetadata{
			Name:  r.Name.ValueString(),
			Scope: scope,
		},
		gateway.GatewayInterceptorEncryptionSpec{
			Comment:     r.Spec.Comment.ValueString(),
			PluginClass: r.Spec.PluginClass.ValueString(),
			Priority:    r.Spec.Priority.ValueInt64(),
			Config:      config,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *gateway.GatewayInterceptorEncryptionResource) (gwinterceptor.GatewayInterceptorEncryptionV2Model, error) {
	config, err := schema.InterceptorConfigToObjectValue(ctx, *r.Spec.Config)
	if err != nil {
		return gwinterceptor.GatewayInterceptorEncryptionV2Model{}, err
	}

	specValue, diag := gwinterceptor.NewSpecValue(
		map[string]attr.Type{
			"comment":      basetypes.StringType{},
			"plugin_class": basetypes.StringType{},
			"priority":     basetypes.Int64Type{},
			"config":       config.Type(ctx),
		},
		map[string]attr.Value{
			"comment":      schema.NewStringValue(r.Spec.Comment),
			"plugin_class": schema.NewStringValue(r.Spec.PluginClass),
			"priority":     schema.NewInt64Value(r.Spec.Priority),
			"config":       config,
		},
	)
	if diag.HasError() {
		return gwinterceptor.GatewayInterceptorEncryptionV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return gwinterceptor.GatewayInterceptorEncryptionV2Model{
		Name: types.StringValue(r.Metadata.Name),
		Scope: gwinterceptor.ScopeValue{
			Group:    schema.NewStringValue(r.Metadata.Scope.Group),
			Vcluster: schema.NewStringValue(r.Metadata.Scope.VCluster),
			Username: schema.NewStringValue(r.Metadata.Scope.Username),
		},
		Spec: specValue,
	}, nil
}

func ObjectValueToInterceptorEncryptionConfig(ctx context.Context, r *basetypes.ObjectValue) (*gateway.GatewayInterceptorEncryptionConfig, error) {
	if r.IsNull() {
		return nil, nil
	}

	configValue, diag := gwinterceptor.NewConfigValue(r.AttributeTypes(ctx), r.Attributes())
	if diag.HasError() {
		return nil, mapper.WrapDiagError(diag, "config", mapper.FromTerraform)
	}

	return &gateway.GatewayInterceptorEncryptionConfig{
		VirtualTopic: configValue.VirtualTopic.ValueString(),
		Statement:    configValue.Statement.ValueString(),
	}, nil
}
