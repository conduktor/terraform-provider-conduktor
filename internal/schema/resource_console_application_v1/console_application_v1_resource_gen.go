// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_console_application_v1

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func ConsoleApplicationV1ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Application name, must be unique, acts as an ID for import",
				MarkdownDescription: "Application name, must be unique, acts as an ID for import",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[0-9a-z\\_\\-]+$"), ""),
				},
			},
			"spec": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"description": schema.StringAttribute{
						Optional:            true,
						Description:         "Application description",
						MarkdownDescription: "Application description",
					},
					"owner": schema.StringAttribute{
						Required:            true,
						Description:         "Application owner, must be a valid Console Group id",
						MarkdownDescription: "Application owner, must be a valid Console Group id",
					},
					"title": schema.StringAttribute{
						Required:            true,
						Description:         "Application title",
						MarkdownDescription: "Application title",
					},
				},
				CustomType: SpecType{
					ObjectType: types.ObjectType{
						AttrTypes: SpecValue{}.AttributeTypes(ctx),
					},
				},
				Required:            true,
				Description:         "Application specification",
				MarkdownDescription: "Application specification",
			},
		},
	}
}

type ConsoleApplicationV1Model struct {
	Name types.String `tfsdk:"name"`
	Spec SpecValue    `tfsdk:"spec"`
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

	descriptionAttribute, ok := attributes["description"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`description is missing from object`)

		return nil, diags
	}

	descriptionVal, ok := descriptionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`description expected to be basetypes.StringValue, was: %T`, descriptionAttribute))
	}

	ownerAttribute, ok := attributes["owner"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`owner is missing from object`)

		return nil, diags
	}

	ownerVal, ok := ownerAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`owner expected to be basetypes.StringValue, was: %T`, ownerAttribute))
	}

	titleAttribute, ok := attributes["title"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`title is missing from object`)

		return nil, diags
	}

	titleVal, ok := titleAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`title expected to be basetypes.StringValue, was: %T`, titleAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return SpecValue{
		Description: descriptionVal,
		Owner:       ownerVal,
		Title:       titleVal,
		state:       attr.ValueStateKnown,
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

	descriptionAttribute, ok := attributes["description"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`description is missing from object`)

		return NewSpecValueUnknown(), diags
	}

	descriptionVal, ok := descriptionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`description expected to be basetypes.StringValue, was: %T`, descriptionAttribute))
	}

	ownerAttribute, ok := attributes["owner"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`owner is missing from object`)

		return NewSpecValueUnknown(), diags
	}

	ownerVal, ok := ownerAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`owner expected to be basetypes.StringValue, was: %T`, ownerAttribute))
	}

	titleAttribute, ok := attributes["title"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`title is missing from object`)

		return NewSpecValueUnknown(), diags
	}

	titleVal, ok := titleAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`title expected to be basetypes.StringValue, was: %T`, titleAttribute))
	}

	if diags.HasError() {
		return NewSpecValueUnknown(), diags
	}

	return SpecValue{
		Description: descriptionVal,
		Owner:       ownerVal,
		Title:       titleVal,
		state:       attr.ValueStateKnown,
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
	Description basetypes.StringValue `tfsdk:"description"`
	Owner       basetypes.StringValue `tfsdk:"owner"`
	Title       basetypes.StringValue `tfsdk:"title"`
	state       attr.ValueState
}

func (v SpecValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 3)

	var val tftypes.Value
	var err error

	attrTypes["description"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["owner"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["title"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 3)

		val, err = v.Description.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["description"] = val

		val, err = v.Owner.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["owner"] = val

		val, err = v.Title.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["title"] = val

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

	attributeTypes := map[string]attr.Type{
		"description": basetypes.StringType{},
		"owner":       basetypes.StringType{},
		"title":       basetypes.StringType{},
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
			"description": v.Description,
			"owner":       v.Owner,
			"title":       v.Title,
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

	if !v.Description.Equal(other.Description) {
		return false
	}

	if !v.Owner.Equal(other.Owner) {
		return false
	}

	if !v.Title.Equal(other.Title) {
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
		"description": basetypes.StringType{},
		"owner":       basetypes.StringType{},
		"title":       basetypes.StringType{},
	}
}
