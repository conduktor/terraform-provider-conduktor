// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_console_topic_v2

import (
	"context"
	"fmt"
	"github.com/conduktor/terraform-provider-conduktor/internal/schema/validation"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ConsoleTopicV2ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"catalog_visibility": schema.StringAttribute{
				Optional:            true,
				Description:         "Catalog visibility for the topic, valid values are: PRIVATE, PUBLIC",
				MarkdownDescription: "Catalog visibility for the topic, valid values are: PRIVATE, PUBLIC",
				Validators: []validator.String{
					stringvalidator.OneOf(validation.ValidCatalogVisibilities...),
				},
			},
			"cluster": schema.StringAttribute{
				Required:            true,
				Description:         "Kafka cluster name linked with Kafka topic. Must exist in Conduktor Console",
				MarkdownDescription: "Kafka cluster name linked with Kafka topic. Must exist in Conduktor Console",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[0-9a-z\\_\\-.]+$"), ""),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Description:         "Topic description",
				MarkdownDescription: "Topic description",
			},
			"description_is_editable": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "is optional (defaults 'true'). Defines whether the description can be updated in the UI",
				MarkdownDescription: "is optional (defaults 'true'). Defines whether the description can be updated in the UI",
				Default:             booldefault.StaticBool(true),
			},
			"labels": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "Topic labels",
				MarkdownDescription: "Topic labels",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Topic name, must be unique, acts as an ID for import",
				MarkdownDescription: "Topic name, must be unique, acts as an ID for import",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[0-9a-z\\_\\-.]+$"), ""),
				},
			},
			"spec": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"configs": schema.MapAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Description:         "Must be valid Kafka Topic configs",
						MarkdownDescription: "Must be valid Kafka Topic configs",
					},
					"partitions": schema.Int64Attribute{
						Required:            true,
						Description:         "Immutable field. Any change will require the Topic to be destroyed and re created",
						MarkdownDescription: "Immutable field. Any change will require the Topic to be destroyed and re created",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.RequiresReplace(),
						},
						Validators: []validator.Int64{
							int64validator.Between(1, 2147483647),
						},
					},
					"replication_factor": schema.Int64Attribute{
						Required:            true,
						Description:         "Immutable field. Any change will require the Topic to be destroyed and re created",
						MarkdownDescription: "Immutable field. Any change will require the Topic to be destroyed and re created",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.RequiresReplace(),
						},
						Validators: []validator.Int64{
							int64validator.Between(1, 2147483647),
						},
					},
				},
				CustomType: SpecType{
					ObjectType: types.ObjectType{
						AttrTypes: SpecValue{}.AttributeTypes(ctx),
					},
				},
				Required:            true,
				Description:         "Topic specification",
				MarkdownDescription: "Topic specification",
			},
			"sql_storage": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"retention_time_in_second": schema.Int64Attribute{
						Required: true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
				},
				CustomType: SqlStorageType{
					ObjectType: types.ObjectType{
						AttrTypes: SqlStorageValue{}.AttributeTypes(ctx),
					},
				},
				Optional: true,
			},
		},
	}
}

type ConsoleTopicV2Model struct {
	CatalogVisibility     types.String    `tfsdk:"catalog_visibility"`
	Cluster               types.String    `tfsdk:"cluster"`
	Description           types.String    `tfsdk:"description"`
	DescriptionIsEditable types.Bool      `tfsdk:"description_is_editable"`
	Labels                types.Map       `tfsdk:"labels"`
	Name                  types.String    `tfsdk:"name"`
	Spec                  SpecValue       `tfsdk:"spec"`
	SqlStorage            SqlStorageValue `tfsdk:"sql_storage"`
}

var _ basetypes.ObjectTypable = SpecType{}

type SpecType struct {
	basetypes.ObjectType
}

func (t SpecType) Equal(o attr.Type) bool {
	other, ok := o.(SpecType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t SpecType) String() string {
	return "SpecType"
}

func (t SpecType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	configsAttribute, ok := attributes["configs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`configs is missing from object`)

		return nil, diags
	}

	configsVal, ok := configsAttribute.(basetypes.MapValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`configs expected to be basetypes.MapValue, was: %T`, configsAttribute))
	}

	partitionsAttribute, ok := attributes["partitions"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`partitions is missing from object`)

		return nil, diags
	}

	partitionsVal, ok := partitionsAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`partitions expected to be basetypes.Int64Value, was: %T`, partitionsAttribute))
	}

	replicationFactorAttribute, ok := attributes["replication_factor"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`replication_factor is missing from object`)

		return nil, diags
	}

	replicationFactorVal, ok := replicationFactorAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`replication_factor expected to be basetypes.Int64Value, was: %T`, replicationFactorAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return SpecValue{
		Configs:           configsVal,
		Partitions:        partitionsVal,
		ReplicationFactor: replicationFactorVal,
		state:             attr.ValueStateKnown,
	}, diags
}

func NewSpecValueNull() SpecValue {
	return SpecValue{
		state: attr.ValueStateNull,
	}
}

func NewSpecValueUnknown() SpecValue {
	return SpecValue{
		state: attr.ValueStateUnknown,
	}
}

func NewSpecValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (SpecValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing SpecValue Attribute Value",
				"While creating a SpecValue value, a missing attribute value was detected. "+
					"A SpecValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("SpecValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid SpecValue Attribute Type",
				"While creating a SpecValue value, an invalid attribute value was detected. "+
					"A SpecValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("SpecValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("SpecValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra SpecValue Attribute Value",
				"While creating a SpecValue value, an extra attribute value was detected. "+
					"A SpecValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra SpecValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewSpecValueUnknown(), diags
	}

	configsAttribute, ok := attributes["configs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`configs is missing from object`)

		return NewSpecValueUnknown(), diags
	}

	configsVal, ok := configsAttribute.(basetypes.MapValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`configs expected to be basetypes.MapValue, was: %T`, configsAttribute))
	}

	partitionsAttribute, ok := attributes["partitions"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`partitions is missing from object`)

		return NewSpecValueUnknown(), diags
	}

	partitionsVal, ok := partitionsAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`partitions expected to be basetypes.Int64Value, was: %T`, partitionsAttribute))
	}

	replicationFactorAttribute, ok := attributes["replication_factor"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`replication_factor is missing from object`)

		return NewSpecValueUnknown(), diags
	}

	replicationFactorVal, ok := replicationFactorAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`replication_factor expected to be basetypes.Int64Value, was: %T`, replicationFactorAttribute))
	}

	if diags.HasError() {
		return NewSpecValueUnknown(), diags
	}

	return SpecValue{
		Configs:           configsVal,
		Partitions:        partitionsVal,
		ReplicationFactor: replicationFactorVal,
		state:             attr.ValueStateKnown,
	}, diags
}

func NewSpecValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) SpecValue {
	object, diags := NewSpecValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewSpecValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t SpecType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewSpecValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewSpecValueUnknown(), nil
	}

	if in.IsNull() {
		return NewSpecValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewSpecValueMust(SpecValue{}.AttributeTypes(ctx), attributes), nil
}

func (t SpecType) ValueType(ctx context.Context) attr.Value {
	return SpecValue{}
}

var _ basetypes.ObjectValuable = SpecValue{}

type SpecValue struct {
	Configs           basetypes.MapValue   `tfsdk:"configs"`
	Partitions        basetypes.Int64Value `tfsdk:"partitions"`
	ReplicationFactor basetypes.Int64Value `tfsdk:"replication_factor"`
	state             attr.ValueState
}

func (v SpecValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 3)

	var val tftypes.Value
	var err error

	attrTypes["configs"] = basetypes.MapType{
		ElemType: types.StringType,
	}.TerraformType(ctx)
	attrTypes["partitions"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["replication_factor"] = basetypes.Int64Type{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 3)

		val, err = v.Configs.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["configs"] = val

		val, err = v.Partitions.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["partitions"] = val

		val, err = v.ReplicationFactor.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["replication_factor"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v SpecValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v SpecValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v SpecValue) String() string {
	return "SpecValue"
}

func (v SpecValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	var configsVal basetypes.MapValue
	switch {
	case v.Configs.IsUnknown():
		configsVal = types.MapUnknown(types.StringType)
	case v.Configs.IsNull():
		configsVal = types.MapNull(types.StringType)
	default:
		var d diag.Diagnostics
		configsVal, d = types.MapValue(types.StringType, v.Configs.Elements())
		diags.Append(d...)
	}

	if diags.HasError() {
		return types.ObjectUnknown(map[string]attr.Type{
			"configs": basetypes.MapType{
				ElemType: types.StringType,
			},
			"partitions":         basetypes.Int64Type{},
			"replication_factor": basetypes.Int64Type{},
		}), diags
	}

	attributeTypes := map[string]attr.Type{
		"configs": basetypes.MapType{
			ElemType: types.StringType,
		},
		"partitions":         basetypes.Int64Type{},
		"replication_factor": basetypes.Int64Type{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"configs":            configsVal,
			"partitions":         v.Partitions,
			"replication_factor": v.ReplicationFactor,
		})

	return objVal, diags
}

func (v SpecValue) Equal(o attr.Value) bool {
	other, ok := o.(SpecValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Configs.Equal(other.Configs) {
		return false
	}

	if !v.Partitions.Equal(other.Partitions) {
		return false
	}

	if !v.ReplicationFactor.Equal(other.ReplicationFactor) {
		return false
	}

	return true
}

func (v SpecValue) Type(ctx context.Context) attr.Type {
	return SpecType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v SpecValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"configs": basetypes.MapType{
			ElemType: types.StringType,
		},
		"partitions":         basetypes.Int64Type{},
		"replication_factor": basetypes.Int64Type{},
	}
}

var _ basetypes.ObjectTypable = SqlStorageType{}

type SqlStorageType struct {
	basetypes.ObjectType
}

func (t SqlStorageType) Equal(o attr.Type) bool {
	other, ok := o.(SqlStorageType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t SqlStorageType) String() string {
	return "SqlStorageType"
}

func (t SqlStorageType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	enabledAttribute, ok := attributes["enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)

		return nil, diags
	}

	enabledVal, ok := enabledAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`enabled expected to be basetypes.BoolValue, was: %T`, enabledAttribute))
	}

	retentionTimeInSecondAttribute, ok := attributes["retention_time_in_second"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`retention_time_in_second is missing from object`)

		return nil, diags
	}

	retentionTimeInSecondVal, ok := retentionTimeInSecondAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`retention_time_in_second expected to be basetypes.Int64Value, was: %T`, retentionTimeInSecondAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return SqlStorageValue{
		Enabled:               enabledVal,
		RetentionTimeInSecond: retentionTimeInSecondVal,
		state:                 attr.ValueStateKnown,
	}, diags
}

func NewSqlStorageValueNull() SqlStorageValue {
	return SqlStorageValue{
		state: attr.ValueStateNull,
	}
}

func NewSqlStorageValueUnknown() SqlStorageValue {
	return SqlStorageValue{
		state: attr.ValueStateUnknown,
	}
}

func NewSqlStorageValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (SqlStorageValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing SqlStorageValue Attribute Value",
				"While creating a SqlStorageValue value, a missing attribute value was detected. "+
					"A SqlStorageValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("SqlStorageValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid SqlStorageValue Attribute Type",
				"While creating a SqlStorageValue value, an invalid attribute value was detected. "+
					"A SqlStorageValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("SqlStorageValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("SqlStorageValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra SqlStorageValue Attribute Value",
				"While creating a SqlStorageValue value, an extra attribute value was detected. "+
					"A SqlStorageValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra SqlStorageValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewSqlStorageValueUnknown(), diags
	}

	enabledAttribute, ok := attributes["enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)

		return NewSqlStorageValueUnknown(), diags
	}

	enabledVal, ok := enabledAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`enabled expected to be basetypes.BoolValue, was: %T`, enabledAttribute))
	}

	retentionTimeInSecondAttribute, ok := attributes["retention_time_in_second"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`retention_time_in_second is missing from object`)

		return NewSqlStorageValueUnknown(), diags
	}

	retentionTimeInSecondVal, ok := retentionTimeInSecondAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`retention_time_in_second expected to be basetypes.Int64Value, was: %T`, retentionTimeInSecondAttribute))
	}

	if diags.HasError() {
		return NewSqlStorageValueUnknown(), diags
	}

	return SqlStorageValue{
		Enabled:               enabledVal,
		RetentionTimeInSecond: retentionTimeInSecondVal,
		state:                 attr.ValueStateKnown,
	}, diags
}

func NewSqlStorageValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) SqlStorageValue {
	object, diags := NewSqlStorageValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewSqlStorageValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t SqlStorageType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewSqlStorageValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewSqlStorageValueUnknown(), nil
	}

	if in.IsNull() {
		return NewSqlStorageValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewSqlStorageValueMust(SqlStorageValue{}.AttributeTypes(ctx), attributes), nil
}

func (t SqlStorageType) ValueType(ctx context.Context) attr.Value {
	return SqlStorageValue{}
}

var _ basetypes.ObjectValuable = SqlStorageValue{}

type SqlStorageValue struct {
	Enabled               basetypes.BoolValue  `tfsdk:"enabled"`
	RetentionTimeInSecond basetypes.Int64Value `tfsdk:"retention_time_in_second"`
	state                 attr.ValueState
}

func (v SqlStorageValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["enabled"] = basetypes.BoolType{}.TerraformType(ctx)
	attrTypes["retention_time_in_second"] = basetypes.Int64Type{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Enabled.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["enabled"] = val

		val, err = v.RetentionTimeInSecond.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["retention_time_in_second"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v SqlStorageValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v SqlStorageValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v SqlStorageValue) String() string {
	return "SqlStorageValue"
}

func (v SqlStorageValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"enabled":                  basetypes.BoolType{},
		"retention_time_in_second": basetypes.Int64Type{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"enabled":                  v.Enabled,
			"retention_time_in_second": v.RetentionTimeInSecond,
		})

	return objVal, diags
}

func (v SqlStorageValue) Equal(o attr.Value) bool {
	other, ok := o.(SqlStorageValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Enabled.Equal(other.Enabled) {
		return false
	}

	if !v.RetentionTimeInSecond.Equal(other.RetentionTimeInSecond) {
		return false
	}

	return true
}

func (v SqlStorageValue) Type(ctx context.Context) attr.Type {
	return SqlStorageType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v SqlStorageValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":                  basetypes.BoolType{},
		"retention_time_in_second": basetypes.Int64Type{},
	}
}
