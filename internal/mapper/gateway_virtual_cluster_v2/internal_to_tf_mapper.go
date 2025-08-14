package gateway_virtual_cluster_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_virtual_cluster_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func InternalModelToTerraform(ctx context.Context, r *gateway.VirtualClusterResource) (schema.GatewayVirtualClusterV2Model, error) {
	specValue, err := specInternalModelToTerraform(ctx, &r.Spec)
	if err != nil {
		return schema.GatewayVirtualClusterV2Model{}, err
	}

	return schema.GatewayVirtualClusterV2Model{
		Name: schemaUtils.NewStringValue(r.Metadata.Name),
		Spec: specValue,
	}, nil
}

func specInternalModelToTerraform(ctx context.Context, r *gateway.VirtualClusterSpec) (schema.SpecValue, error) {
	unknownSpecObjectValue, diag := schema.NewSpecValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schemaUtils.ValueMapFromTypes(ctx, typesMap)

	valuesMap["acl_enabled"] = basetypes.NewBoolValue(r.AclEnabled)
	valuesMap["acl_mode"] = schemaUtils.NewStringValue(r.AclMode)
	valuesMap["type"] = schemaUtils.NewStringValue(r.Type)
	valuesMap["bootstrap_servers"] = schemaUtils.NewStringValue(r.BootstrapServers)

	superUsers, diag := schemaUtils.StringArrayToSetValue(r.SuperUsers)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "super_users", mapper.IntoTerraform)
	}
	valuesMap["super_users"] = superUsers

	clientProperties, err := clientPropertiesMapToMapValue(ctx, r.ClientProperties)
	if err != nil {
		return schema.SpecValue{}, err
	}
	valuesMap["client_properties"] = clientProperties

	acls, err := aclsArrayToSetValue(ctx, r.Acls)
	if err != nil {
		return schema.SpecValue{}, err
	}
	valuesMap["acls"] = acls

	value, diag := schema.NewSpecValue(typesMap, valuesMap)
	if diag.HasError() {
		return schema.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	return value, nil
}

func clientPropertiesMapToMapValue(ctx context.Context, m map[string]map[string]string) (basetypes.MapValue, error) {
	var tfProperties = make(map[string]attr.Value)
	var diag diag.Diagnostics

	for k, v := range m {
		var tfValue = make(map[string]attr.Value)

		for key, value := range v {
			tfValue[key] = schemaUtils.NewStringValue(value)
		}

		permObj, diag := types.MapValue(basetypes.StringType{}, tfValue)
		if diag.HasError() {
			return basetypes.MapValue{}, mapper.WrapDiagError(diag, "client_properties.inside", mapper.FromTerraform)
		}

		tfProperties[k] = permObj
	}

	clientType := types.MapUnknown(basetypes.StringType{}).Type(ctx)

	mapValue, diag := types.MapValue(clientType, tfProperties)
	if diag.HasError() {
		return basetypes.MapValue{}, mapper.WrapDiagError(diag, "client_properties.outside", mapper.FromTerraform)
	}

	return mapValue, nil
}

func aclsArrayToSetValue(ctx context.Context, arr []gateway.VirtualClusterACL) (basetypes.SetValue, error) {
	var aclsList basetypes.SetValue
	var tfACLs []attr.Value
	var diag diag.Diagnostics

	for _, p := range arr {
		resourcePattern, err := resourcePatternInternalModelToTerraform(p.ResourcePattern)
		if err != nil {
			return basetypes.SetValue{}, mapper.WrapError(err, "resource_pattern", mapper.IntoTerraform)
		}

		resourcePatternObj, diag := resourcePattern.ToObjectValue(ctx)
		if diag.HasError() {
			return basetypes.SetValue{}, mapper.WrapDiagError(diag, "resource_pattern", mapper.IntoTerraform)
		}

		types := map[string]attr.Type{
			"resource_pattern": resourcePatternObj.Type(ctx),
			"principal":        basetypes.StringType{},
			"host":             basetypes.StringType{},
			"operation":        basetypes.StringType{},
			"permission_type":  basetypes.StringType{},
		}
		values := map[string]attr.Value{
			"resource_pattern": resourcePatternObj,
			"principal":        schemaUtils.NewStringValue(p.Principal),
			"host":             schemaUtils.NewStringValue(p.Host),
			"operation":        schemaUtils.NewStringValue(p.Operation),
			"permission_type":  schemaUtils.NewStringValue(p.PermissionType),
		}

		permObj, diag := schema.NewAclsValue(types, values)
		if diag.HasError() {
			return basetypes.SetValue{}, mapper.WrapDiagError(diag, "acls", mapper.FromTerraform)
		}
		tfACLs = append(tfACLs, permObj)

	}

	aclsList, diag = types.SetValue(schema.AclsValue{}.Type(ctx), tfACLs)
	if diag.HasError() {
		return basetypes.SetValue{}, mapper.WrapDiagError(diag, "acls", mapper.FromTerraform)
	}

	return aclsList, nil
}

func resourcePatternInternalModelToTerraform(r gateway.VirtualClusterACLResourcePattern) (schema.ResourcePatternValue, error) {
	types := map[string]attr.Type{
		"resource_type": basetypes.StringType{},
		"name":          basetypes.StringType{},
		"pattern_type":  basetypes.StringType{},
	}

	values := map[string]attr.Value{
		"resource_type": schemaUtils.NewStringValue(r.ResourceType),
		"name":          schemaUtils.NewStringValue(r.Name),
		"pattern_type":  schemaUtils.NewStringValue(r.PatternType),
	}

	value, diag := schema.NewResourcePatternValue(types, values)
	if diag.HasError() {
		return schema.ResourcePatternValue{}, mapper.WrapDiagError(diag, "resource_pattern", mapper.IntoTerraform)
	}
	return value, nil
}
