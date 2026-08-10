package console_connector_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	connector "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_connector_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *connector.ConsoleConnectorV2Model) (console.ConnectorConsoleResource, error) {
	userLabels, diag := schema.MapValueToStringMap(ctx, r.Labels)
	if diag.HasError() {
		return console.ConnectorConsoleResource{}, mapper.WrapDiagError(diag, "labels", mapper.FromTerraform)
	}
	managedLabels, diag := schema.MapValueToStringMap(ctx, r.ManagedLabels)
	if diag.HasError() {
		return console.ConnectorConsoleResource{}, mapper.WrapDiagError(diag, "managed_labels", mapper.FromTerraform)
	}

	config, diag := schema.MapValueToStringMap(ctx, r.Spec.Config)
	if diag.HasError() {
		return console.ConnectorConsoleResource{}, mapper.WrapDiagError(diag, "spec.config", mapper.FromTerraform)
	}

	var autoRestart *console.AutoRestart = nil
	if schema.AttrIsSet(r.AutoRestart) {
		autoRestart = &console.AutoRestart{
			Enabled:          r.AutoRestart.Enabled.ValueBool(),
			FrequencySeconds: r.AutoRestart.FrequencySeconds.ValueInt64(),
		}
	}

	return console.NewConnectorConsoleResource(
		console.ConnectorConsoleMetadata{
			Name:           r.Name.ValueString(),
			Cluster:        r.Cluster.ValueString(),
			ConnectCluster: r.ConnectCluster.ValueString(),
			Labels:         mapper.MergeLabels(managedLabels, userLabels),
			Description:    r.Description.ValueString(),
			AutoRestart:    autoRestart,
		},
		console.ConnectorConsoleSpec{
			Config: config,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *console.ConnectorConsoleResource) (connector.ConsoleConnectorV2Model, error) {

	var modelUserLabels, modelManagedLabels = mapper.SplitLabels(r.Metadata.Labels)
	labels, diag := schema.StringMapToMapValue(ctx, modelUserLabels)
	if diag.HasError() {
		return connector.ConsoleConnectorV2Model{}, mapper.WrapDiagError(diag, "labels", mapper.IntoTerraform)
	}

	managedLabels, diag := schema.StringMapToMapValue(ctx, modelManagedLabels)
	if diag.HasError() {
		return connector.ConsoleConnectorV2Model{}, mapper.WrapDiagError(diag, "managed_labels", mapper.IntoTerraform)
	}

	config, diag := schema.StringMapToMapValue(ctx, r.Spec.Config)
	if diag.HasError() {
		return connector.ConsoleConnectorV2Model{}, mapper.WrapDiagError(diag, "spec.config", mapper.IntoTerraform)
	}

	autoRestart, err := autoRestartInternalToTerraform(r.Metadata.AutoRestart)
	if err != nil {
		return connector.ConsoleConnectorV2Model{}, err
	}

	specValue, diag := connector.NewSpecValue(
		map[string]attr.Type{
			"config": config.Type(ctx),
		},
		map[string]attr.Value{
			"config": config,
		},
	)
	if diag.HasError() {
		return connector.ConsoleConnectorV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return connector.ConsoleConnectorV2Model{
		Name:           schema.NewStringValue(r.Metadata.Name),
		Cluster:        schema.NewStringValue(r.Metadata.Cluster),
		ConnectCluster: schema.NewStringValue(r.Metadata.ConnectCluster),
		Labels:         labels,
		ManagedLabels:  managedLabels,
		Description:    schema.NewStringValue(r.Metadata.Description),
		AutoRestart:    autoRestart,
		Spec:           specValue,
	}, nil
}

func autoRestartInternalToTerraform(r *console.AutoRestart) (connector.AutoRestartValue, error) {
	if r == nil {
		return connector.NewAutoRestartValueNull(), nil
	}

	autoRestart, diag := connector.NewAutoRestartValue(
		map[string]attr.Type{
			"enabled":           basetypes.BoolType{},
			"frequency_seconds": basetypes.Int64Type{},
		},
		map[string]attr.Value{
			"enabled":           basetypes.NewBoolValue(r.Enabled),
			"frequency_seconds": schema.NewInt64Value(r.FrequencySeconds),
		},
	)
	if diag.HasError() {
		return connector.AutoRestartValue{}, mapper.WrapDiagError(diag, "auto_restart", mapper.IntoTerraform)
	}

	return autoRestart, nil
}
