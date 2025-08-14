package gateway_virtual_cluster_v2

import (
	"context"
	"fmt"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_virtual_cluster_v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *schema.GatewayVirtualClusterV2Model) (model.VirtualClusterResource, error) {
	superUsers, diag := schemaUtils.SetValueToStringArray(ctx, r.Spec.SuperUsers)
	if diag.HasError() {
		return model.VirtualClusterResource{}, mapper.WrapDiagError(diag, "super_users", mapper.FromTerraform)
	}

	clientProperties, err := mapValueToClientPropertiesMap(ctx, r.Spec.ClientProperties)
	if err != nil {
		return model.VirtualClusterResource{}, err
	}

	acls, err := setValueToACLs(ctx, r.Spec.Acls)
	if err != nil {
		return model.VirtualClusterResource{}, err
	}

	return model.NewVirtualClusterResource(
		model.VirtualClusterMetadata{
			Name: r.Name.ValueString(),
		},
		model.VirtualClusterSpec{
			AclEnabled:       r.Spec.AclEnabled.ValueBool(),
			AclMode:          r.Spec.AclMode.ValueString(),
			SuperUsers:       superUsers,
			Type:             r.Spec.SpecType.ValueString(),
			BootstrapServers: r.Spec.BootstrapServers.ValueString(),
			ClientProperties: clientProperties,
			Acls:             acls,
		},
	), nil

}

func mapValueToClientPropertiesMap(ctx context.Context, m basetypes.MapValue) (map[string]map[string]string, error) {
	clientProperties := make(map[string]map[string]string)

	if !m.IsNull() && !m.IsUnknown() {
		for k, v := range m.Elements() {
			value, ok := v.(basetypes.MapValue)
			if !ok {
				return nil, fmt.Errorf("client_properties.%s is not a map", k)
			}
			properties, diag := schemaUtils.MapValueToStringMap(ctx, value)
			if diag.HasError() {
				return nil, mapper.WrapDiagError(diag, "client_properties", mapper.FromTerraform)
			}
			clientProperties[k] = properties
		}
	}

	return clientProperties, nil
}

func setValueToACLs(ctx context.Context, set basetypes.SetValue) ([]model.VirtualClusterACL, error) {
	acls := make([]model.VirtualClusterACL, 0)
	var diag diag.Diagnostics

	if !set.IsNull() && !set.IsUnknown() {
		var aclsValue []schema.AclsValue
		diag = set.ElementsAs(ctx, &aclsValue, false)
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "acls", mapper.FromTerraform)
		}
		for _, p := range aclsValue {

			resourcePattern, err := objectValueToResourcePattern(ctx, p.ResourcePattern)
			if err != nil {
				return nil, err
			}

			acls = append(acls, model.VirtualClusterACL{
				ResourcePattern: *resourcePattern,
				Principal:       p.Principal.ValueString(),
				Host:            p.Host.ValueString(),
				Operation:       p.Operation.ValueString(),
				PermissionType:  p.PermissionType.ValueString(),
			})
		}
	}
	return acls, nil
}

func objectValueToResourcePattern(ctx context.Context, r basetypes.ObjectValue) (*model.VirtualClusterACLResourcePattern, error) {
	if r.IsNull() {
		return nil, nil
	}

	resPattern, diag := schema.NewResourcePatternValue(r.AttributeTypes(ctx), r.Attributes())
	if diag.HasError() {
		return &model.VirtualClusterACLResourcePattern{}, mapper.WrapDiagError(diag, "resource_pattern", mapper.FromTerraform)
	}

	return &model.VirtualClusterACLResourcePattern{
		ResourceType: resPattern.ResourceType.ValueString(),
		Name:         resPattern.Name.ValueString(),
		PatternType:  resPattern.PatternType.ValueString(),
	}, nil
}
