package console_topic_v2

import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	topic "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_topic_v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFToInternalModel(ctx context.Context, r *topic.ConsoleTopicV2Model) (console.TopicConsoleResource, error) {
	labels, diag := schema.MapValueToStringMap(ctx, r.Labels)
	if diag.HasError() {
		return console.TopicConsoleResource{}, mapper.WrapDiagError(diag, "labels", mapper.FromTerraform)
	}

	configs, diag := schema.MapValueToStringMap(ctx, r.Spec.Configs)
	if diag.HasError() {
		return console.TopicConsoleResource{}, mapper.WrapDiagError(diag, "spec.configs", mapper.FromTerraform)
	}

	return console.NewTopicConsoleResource(
		console.TopicConsoleMetadata{
			Name:                  r.Name.ValueString(),
			Cluster:               r.Cluster.ValueString(),
			Labels:                labels,
			CatalogVisibility:     r.CatalogVisibility.ValueString(),
			DescriptionIsEditable: r.DescriptionIsEditable.ValueBool(),
			Description:           r.Description.ValueString(),
			SqlStorage: console.TopicSqlStorage{
				RetentionTimeInSecond: r.SqlStorage.RetentionTimeInSecond.ValueInt64(),
				Enabled:               r.SqlStorage.Enabled.ValueBool(),
			},
		},
		console.TopicConsoleSpec{
			Partitions:        r.Spec.Partitions.ValueInt64(),
			ReplicationFactor: r.Spec.ReplicationFactor.ValueInt64(),
			Configs:           configs,
		},
	), nil
}

func InternalModelToTerraform(ctx context.Context, r *console.TopicConsoleResource) (topic.ConsoleTopicV2Model, error) {
	labels, diag := schema.StringMapToMapValue(ctx, r.Metadata.Labels)
	if diag.HasError() {
		return topic.ConsoleTopicV2Model{}, mapper.WrapDiagError(diag, "labels", mapper.IntoTerraform)
	}

	configs, diag := schema.StringMapToMapValue(ctx, r.Spec.Configs)
	if diag.HasError() {
		return topic.ConsoleTopicV2Model{}, mapper.WrapDiagError(diag, "spec.configs", mapper.IntoTerraform)
	}

	sqlStorage, diag := topic.NewSqlStorageValue(
		map[string]attr.Type{
			"retention_time_in_second": basetypes.Int64Type{},
			"enabled":                  basetypes.BoolType{},
		},
		map[string]attr.Value{
			"retention_time_in_second": schema.NewInt64Value(r.Metadata.SqlStorage.RetentionTimeInSecond),
			"enabled":                  basetypes.NewBoolValue(r.Metadata.SqlStorage.Enabled),
		},
	)

	specValue, diag := topic.NewSpecValue(
		map[string]attr.Type{
			"partitions":         basetypes.Int64Type{},
			"replication_factor": basetypes.Int64Type{},
			"configs":            configs.Type(ctx),
		},
		map[string]attr.Value{
			"partitions":         schema.NewInt64Value(r.Spec.Partitions),
			"replication_factor": schema.NewInt64Value(r.Spec.ReplicationFactor),
			"configs":            configs,
		},
	)
	if diag.HasError() {
		return topic.ConsoleTopicV2Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return topic.ConsoleTopicV2Model{
		Name:                  schema.NewStringValue(r.Metadata.Name),
		Cluster:               schema.NewStringValue(r.Metadata.Cluster),
		Labels:                labels,
		CatalogVisibility:     schema.NewStringValue(r.Metadata.CatalogVisibility),
		DescriptionIsEditable: basetypes.NewBoolValue(r.Metadata.DescriptionIsEditable),
		Description:           schema.NewStringValue(r.Metadata.Description),
		SqlStorage:            sqlStorage,
		Spec:                  specValue,
	}, nil
}
