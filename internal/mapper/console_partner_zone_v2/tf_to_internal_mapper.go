package console_partner_zone_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	partnerZone "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_partner_zone_v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *partnerZone.ConsolePartnerZoneV2Model) (console.PartnerZoneConsoleResource, error) {
	labels, diag := schema.MapValueToStringMap(ctx, r.Labels)
	if diag.HasError() {
		return console.PartnerZoneConsoleResource{}, mapper.WrapDiagError(diag, "labels", mapper.FromTerraform)
	}

	authMode, err := objectValueToAuthenticationMode(ctx, r.Spec.AuthenticationMode)
	if err != nil {
		return console.PartnerZoneConsoleResource{}, err
	}

	topics, err := setValueToTopicsArray(ctx, r.Spec.Topics)
	if err != nil {
		return console.PartnerZoneConsoleResource{}, err
	}

	partner, err := objectValueToPartner(ctx, r.Spec.Partner)
	if err != nil {
		return console.PartnerZoneConsoleResource{}, err
	}

	tfc, err := objectValueToTrafficControlPolicies(ctx, r.Spec.TrafficControlPolicies)
	if err != nil {
		return console.PartnerZoneConsoleResource{}, err
	}

	headers, err := objectValueToHeaders(ctx, r.Spec.Headers)
	if err != nil {
		return console.PartnerZoneConsoleResource{}, err
	}

	return console.NewPartnerZoneConsoleResource(
		console.PartnerZoneConsoleMetadata{
			Name:   r.Name.ValueString(),
			Labels: labels,
		},
		console.PartnerZoneConsoleSpec{
			Cluster:                r.Spec.Cluster.ValueString(),
			DisplayName:            r.Spec.DisplayName.ValueString(),
			Description:            r.Spec.Description.ValueString(),
			Url:                    r.Spec.Url.ValueString(),
			AuthenticationMode:     authMode,
			Topics:                 topics,
			Partner:                partner,
			TrafficControlPolicies: tfc,
			Headers:                headers,
		},
	), nil
}

func objectValueToAuthenticationMode(ctx context.Context, r basetypes.ObjectValue) (console.PartnerZoneAuthenticationMode, error) {
	if r.IsNull() {
		return console.PartnerZoneAuthenticationMode{}, nil
	}

	authModeValue, diag := partnerZone.NewAuthenticationModeValue(r.AttributeTypes(ctx), r.Attributes())
	if diag.HasError() {
		return console.PartnerZoneAuthenticationMode{}, mapper.WrapDiagError(diag, "authentication_mode", mapper.FromTerraform)
	}
	return console.PartnerZoneAuthenticationMode{
		ServiceAccount: authModeValue.ServiceAccount.ValueString(),
		Type:           authModeValue.AuthenticationModeType.ValueString(),
	}, nil
}

func setValueToTopicsArray(ctx context.Context, set basetypes.SetValue) ([]console.PartnerZoneTopic, error) {
	topics := make([]console.PartnerZoneTopic, 0)
	var diag diag.Diagnostics

	if !set.IsNull() && !set.IsUnknown() {
		var tfResources []partnerZone.TopicsValue
		diag = set.ElementsAs(ctx, &tfResources, false)
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "topics", mapper.FromTerraform)
		}

		for _, p := range tfResources {
			topics = append(topics, console.PartnerZoneTopic{
				Name:         p.Name.ValueString(),
				BackingTopic: p.BackingTopic.ValueString(),
				Permission:   p.Permission.ValueString(),
			})
		}

	}
	return topics, nil
}

func objectValueToPartner(ctx context.Context, r basetypes.ObjectValue) (console.PartnerZonePartner, error) {
	if r.IsNull() {
		return console.PartnerZonePartner{}, nil
	}

	partnerValue, diag := partnerZone.NewPartnerValue(r.AttributeTypes(ctx), r.Attributes())
	if diag.HasError() {
		return console.PartnerZonePartner{}, mapper.WrapDiagError(diag, "partner", mapper.FromTerraform)
	}
	return console.PartnerZonePartner{
		Name:  partnerValue.Name.ValueString(),
		Role:  partnerValue.Role.ValueString(),
		Email: partnerValue.Email.ValueString(),
		Phone: partnerValue.Phone.ValueString(),
	}, nil
}

func objectValueToTrafficControlPolicies(ctx context.Context, r basetypes.ObjectValue) (console.PartnerZoneTrafficControlPolicies, error) {
	if r.IsNull() {
		return console.PartnerZoneTrafficControlPolicies{}, nil
	}

	tcpValue, diag := partnerZone.NewTrafficControlPoliciesValue(r.AttributeTypes(ctx), r.Attributes())
	if diag.HasError() {
		return console.PartnerZoneTrafficControlPolicies{}, mapper.WrapDiagError(diag, "traffic_control_policies", mapper.FromTerraform)
	}
	return console.PartnerZoneTrafficControlPolicies{
		MaxProduceRate:    tcpValue.MaxProduceRate.ValueInt64(),
		MaxConsumeRate:    tcpValue.MaxConsumeRate.ValueInt64(),
		LimitCommitOffset: tcpValue.LimitCommitOffset.ValueInt64(),
	}, nil
}

func objectValueToHeaders(ctx context.Context, r basetypes.ObjectValue) (console.PartnerZoneHeaders, error) {
	if r.IsNull() {
		return console.PartnerZoneHeaders{}, nil
	}

	headers, diag := partnerZone.NewHeadersValue(r.AttributeTypes(ctx), r.Attributes())
	if diag.HasError() {
		return console.PartnerZoneHeaders{}, mapper.WrapDiagError(diag, "headers", mapper.FromTerraform)
	}

	toAdd, err := setValueToAddArray(ctx, headers.AddOnProduce)
	if err != nil {
		return console.PartnerZoneHeaders{}, err
	}

	toRemove, err := setValueToRemoveArray(ctx, headers.RemoveOnConsume)
	if err != nil {
		return console.PartnerZoneHeaders{}, err
	}

	return console.PartnerZoneHeaders{
		AddOnProduce:    toAdd,
		RemoveOnConsume: toRemove,
	}, nil
}

func setValueToAddArray(ctx context.Context, set basetypes.SetValue) ([]console.PartnerZoneToAdd, error) {
	toAdd := make([]console.PartnerZoneToAdd, 0)
	var diag diag.Diagnostics

	if !set.IsNull() && !set.IsUnknown() {
		var tfResources []partnerZone.AddOnProduceValue
		diag = set.ElementsAs(ctx, &tfResources, false)
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "add_on_produce", mapper.FromTerraform)
		}

		for _, p := range tfResources {
			toAdd = append(toAdd, console.PartnerZoneToAdd{
				Key:              p.Key.ValueString(),
				Value:            p.Value.ValueString(),
				OverrideIfExists: p.OverrideIfExists.ValueBool(),
			})
		}

	}
	return toAdd, nil
}

func setValueToRemoveArray(ctx context.Context, set basetypes.SetValue) ([]console.PartnerZoneToRemove, error) {
	toRemove := make([]console.PartnerZoneToRemove, 0)
	var diag diag.Diagnostics

	if !set.IsNull() && !set.IsUnknown() {
		var tfResources []partnerZone.RemoveOnConsumeValue
		diag = set.ElementsAs(ctx, &tfResources, false)
		if diag.HasError() {
			return nil, mapper.WrapDiagError(diag, "remove_on_consume", mapper.FromTerraform)
		}

		for _, p := range tfResources {
			toRemove = append(toRemove, console.PartnerZoneToRemove{
				KeyRegex: p.KeyRegex.ValueString(),
			})
		}

	}
	return toRemove, nil
}
