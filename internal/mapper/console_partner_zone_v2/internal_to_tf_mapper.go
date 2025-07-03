package console_partner_zone_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	partnerZone "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_partner_zone_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func InternalModelToTerraform(ctx context.Context, r *console.PartnerZoneConsoleResource) (partnerZone.ConsolePartnerZoneV2Model, error) {
	var diag diag.Diagnostics

	labels, diag := schema.StringMapToMapValue(ctx, r.Metadata.Labels)
	if diag.HasError() {
		return partnerZone.ConsolePartnerZoneV2Model{}, mapper.WrapDiagError(diag, "labels", mapper.IntoTerraform)
	}

	specValue, err := specInternalModelToTerraform(ctx, &r.Spec)
	if err != nil {
		return partnerZone.ConsolePartnerZoneV2Model{}, err
	}

	return partnerZone.ConsolePartnerZoneV2Model{
		Name:   schema.NewStringValue(r.Metadata.Name),
		Labels: labels,
		Spec:   specValue,
	}, nil
}

func specInternalModelToTerraform(ctx context.Context, r *console.PartnerZoneConsoleSpec) (partnerZone.SpecValue, error) {
	unknownSpecObjectValue, diag := partnerZone.NewSpecValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return partnerZone.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	var typesMap = unknownSpecObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	valuesMap["cluster"] = schema.NewStringValue(r.Cluster)
	valuesMap["display_name"] = schema.NewStringValue(r.DisplayName)
	valuesMap["description"] = schema.NewStringValue(r.Description)
	valuesMap["url"] = schema.NewStringValue(r.Url)

	authMode, err := authenticationModeToTerraform(r.AuthenticationMode)
	if err != nil {
		return partnerZone.SpecValue{}, err
	}
	authModeValue, diag := authMode.ToObjectValue(ctx)
	if diag.HasError() {
		return partnerZone.SpecValue{}, mapper.WrapDiagError(diag, "authentication_mode", mapper.IntoTerraform)
	}
	valuesMap["authentication_mode"] = authModeValue

	topicsSet, err := topicsArrayToSetValue(ctx, r.Topics)
	if err != nil {
		return partnerZone.SpecValue{}, err
	}
	valuesMap["topics"] = topicsSet

	partner, err := partnerToTerraform(r.Partner)
	if err != nil {
		return partnerZone.SpecValue{}, err
	}
	partnerValue, diag := partner.ToObjectValue(ctx)
	if diag.HasError() {
		return partnerZone.SpecValue{}, mapper.WrapDiagError(diag, "partner", mapper.IntoTerraform)
	}
	valuesMap["partner"] = partnerValue

	tfc, err := trafficControlPoliciesToTerraform(r.TrafficControlPolicies)
	if err != nil {
		return partnerZone.SpecValue{}, err
	}
	tfcValue, diag := tfc.ToObjectValue(ctx)
	if diag.HasError() {
		return partnerZone.SpecValue{}, mapper.WrapDiagError(diag, "traffic_control_policies", mapper.IntoTerraform)
	}
	valuesMap["traffic_control_policies"] = tfcValue

	headers, err := headersToTerraform(ctx, r.Headers)
	if err != nil {
		return partnerZone.SpecValue{}, err
	}
	headersValue, diag := headers.ToObjectValue(ctx)
	if diag.HasError() {
		return partnerZone.SpecValue{}, mapper.WrapDiagError(diag, "headers", mapper.IntoTerraform)
	}
	valuesMap["headers"] = headersValue

	value, diag := partnerZone.NewSpecValue(typesMap, valuesMap)
	if diag.HasError() {
		return partnerZone.SpecValue{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}
	return value, nil
}

func authenticationModeToTerraform(r *console.PartnerZoneAuthenticationMode) (partnerZone.AuthenticationModeValue, error) {
	if r == nil {
		return partnerZone.NewAuthenticationModeValueNull(), nil
	}

	types := map[string]attr.Type{
		"service_account": basetypes.StringType{},
		"type":            basetypes.StringType{},
	}
	values := map[string]attr.Value{
		"service_account": schema.NewStringValue(r.ServiceAccount),
		"type":            schema.NewStringValue(r.Type),
	}

	authMode, diag := partnerZone.NewAuthenticationModeValue(types, values)
	if diag.HasError() {
		return partnerZone.AuthenticationModeValue{}, mapper.WrapDiagError(diag, "authentication_mode", mapper.FromTerraform)
	}

	return authMode, nil
}

func partnerToTerraform(r *console.PartnerZonePartner) (partnerZone.PartnerValue, error) {
	if r == nil {
		return partnerZone.NewPartnerValueNull(), nil
	}

	types := map[string]attr.Type{
		"name":  basetypes.StringType{},
		"role":  basetypes.StringType{},
		"email": basetypes.StringType{},
		"phone": basetypes.StringType{},
	}
	values := map[string]attr.Value{
		"name":  schema.NewStringValue(r.Name),
		"role":  schema.NewStringValue(r.Role),
		"email": schema.NewStringValue(r.Email),
		"phone": schema.NewStringValue(r.Phone),
	}

	partner, diag := partnerZone.NewPartnerValue(types, values)
	if diag.HasError() {
		return partnerZone.PartnerValue{}, mapper.WrapDiagError(diag, "partner", mapper.FromTerraform)
	}

	return partner, nil
}

func trafficControlPoliciesToTerraform(r *console.PartnerZoneTrafficControlPolicies) (partnerZone.TrafficControlPoliciesValue, error) {
	if r == nil {
		return partnerZone.NewTrafficControlPoliciesValueNull(), nil
	}

	types := map[string]attr.Type{
		"max_produce_rate":    basetypes.Int64Type{},
		"max_consume_rate":    basetypes.Int64Type{},
		"limit_commit_offset": basetypes.Int64Type{},
	}
	values := map[string]attr.Value{
		"max_produce_rate":    schema.NewInt64Value(r.MaxProduceRate),
		"max_consume_rate":    schema.NewInt64Value(r.MaxConsumeRate),
		"limit_commit_offset": schema.NewInt64Value(r.LimitCommitOffset),
	}

	tcp, diag := partnerZone.NewTrafficControlPoliciesValue(types, values)
	if diag.HasError() {
		return partnerZone.TrafficControlPoliciesValue{}, mapper.WrapDiagError(diag, "traffic_control_policies", mapper.FromTerraform)
	}

	return tcp, nil
}

func headersToTerraform(ctx context.Context, r *console.PartnerZoneHeaders) (partnerZone.HeadersValue, error) {
	if r == nil {
		return partnerZone.NewHeadersValueNull(), nil
	}

	unknownHeadersObjectValue, diag := partnerZone.NewHeadersValueUnknown().ToObjectValue(ctx)
	if diag.HasError() {
		return partnerZone.HeadersValue{}, mapper.WrapDiagError(diag, "headers", mapper.IntoTerraform)
	}
	var typesMap = unknownHeadersObjectValue.AttributeTypes(ctx)
	var valuesMap = schema.ValueMapFromTypes(ctx, typesMap)

	addOnProduce, err := toAddArrayToSetValue(ctx, r.AddOnProduce)
	if err != nil {
		return partnerZone.HeadersValue{}, err
	}
	valuesMap["add_on_produce"] = addOnProduce

	removeOnConsume, err := toRemoveArrayToSetValue(ctx, r.RemoveOnConsume)
	if err != nil {
		return partnerZone.HeadersValue{}, err
	}
	valuesMap["remove_on_consume"] = removeOnConsume

	value, diag := partnerZone.NewHeadersValue(typesMap, valuesMap)
	if diag.HasError() {
		return partnerZone.HeadersValue{}, mapper.WrapDiagError(diag, "headers", mapper.IntoTerraform)
	}
	return value, nil
}

func toAddArrayToSetValue(ctx context.Context, arr []console.PartnerZoneToAdd) (basetypes.SetValue, error) {
	var tfResources []attr.Value
	var diag diag.Diagnostics

	for _, p := range arr {
		types := map[string]attr.Type{
			"key":                basetypes.StringType{},
			"value":              basetypes.StringType{},
			"override_if_exists": basetypes.BoolType{},
		}
		values := map[string]attr.Value{
			"key":                schema.NewStringValue(p.Key),
			"value":              schema.NewStringValue(p.Value),
			"override_if_exists": basetypes.NewBoolValue(p.OverrideIfExists),
		}

		addObj, diag := partnerZone.NewAddOnProduceValue(types, values)
		if diag.HasError() {
			return basetypes.SetValue{}, mapper.WrapDiagError(diag, "add_on_produce", mapper.FromTerraform)
		}
		tfResources = append(tfResources, addObj)
	}

	toAddList, diag := types.SetValue(partnerZone.AddOnProduceValue{}.Type(ctx), tfResources)
	if diag.HasError() {
		return basetypes.SetValue{}, mapper.WrapDiagError(diag, "add_on_produce", mapper.FromTerraform)
	}

	return toAddList, nil
}

func toRemoveArrayToSetValue(ctx context.Context, arr []console.PartnerZoneToRemove) (basetypes.SetValue, error) {
	var tfResources []attr.Value
	var diag diag.Diagnostics

	for _, p := range arr {
		types := map[string]attr.Type{
			"key_regex": basetypes.StringType{},
		}
		values := map[string]attr.Value{
			"key_regex": schema.NewStringValue(p.KeyRegex),
		}

		removeObj, diag := partnerZone.NewRemoveOnConsumeValue(types, values)
		if diag.HasError() {
			return basetypes.SetValue{}, mapper.WrapDiagError(diag, "remove_on_consume", mapper.FromTerraform)
		}
		tfResources = append(tfResources, removeObj)
	}

	toRemoveList, diag := types.SetValue(partnerZone.RemoveOnConsumeValue{}.Type(ctx), tfResources)
	if diag.HasError() {
		return basetypes.SetValue{}, mapper.WrapDiagError(diag, "remove_on_consume", mapper.FromTerraform)
	}

	return toRemoveList, nil
}

func topicsArrayToSetValue(ctx context.Context, arr []console.PartnerZoneTopic) (basetypes.SetValue, error) {
	var tfResources []attr.Value
	var diag diag.Diagnostics

	for _, p := range arr {
		types := map[string]attr.Type{
			"name":          basetypes.StringType{},
			"backing_topic": basetypes.StringType{},
			"permission":    basetypes.StringType{},
		}
		values := map[string]attr.Value{
			"name":          schema.NewStringValue(p.Name),
			"backing_topic": schema.NewStringValue(p.BackingTopic),
			"permission":    schema.NewStringValue(p.Permission),
		}

		permObj, diag := partnerZone.NewTopicsValue(types, values)
		if diag.HasError() {
			return basetypes.SetValue{}, mapper.WrapDiagError(diag, "topics", mapper.FromTerraform)
		}
		tfResources = append(tfResources, permObj)

	}

	topicsList, diag := types.SetValue(partnerZone.TopicsValue{}.Type(ctx), tfResources)
	if diag.HasError() {
		return basetypes.SetValue{}, mapper.WrapDiagError(diag, "topics", mapper.FromTerraform)
	}

	return topicsList, nil
}
