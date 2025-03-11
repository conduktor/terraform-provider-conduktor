package provider

import (
	"context"
	"fmt"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/console_topic_policy_v1"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_topic_policy_v1"
	topicPolicy "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_topic_policy_v1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
)

const topicPolicyV1ApiPath = "/public/self-serve/v1/topic-policy"
const topicPolicyMininumVersion = "v1.30.0"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TopicPolicyV1Resource{}
var _ resource.ResourceWithImportState = &TopicPolicyV1Resource{}

func NewTopicPolicyV1Resource() resource.Resource {
	return &TopicPolicyV1Resource{}
}

// TopicPolicyV1Resource defines the resource implementation.
type TopicPolicyV1Resource struct {
	apiClient *client.Client
}

func (r *TopicPolicyV1Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_topic_policy_v1"
}

func (r *TopicPolicyV1Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.ConsoleTopicPolicyV1ResourceSchema(ctx)
}

func (r *TopicPolicyV1Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*ProviderData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	if data.Client == nil || data.Mode != client.CONSOLE {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Console Client not configured. Please provide client configuration details for Console API and ensure you have set the right provider mode for this resource. \n"+
				"More info here: \n"+
				" - https://registry.terraform.io/providers/conduktor/conduktor/latest/docs",
		)
		return
	}

	r.apiClient = data.Client
}

func (r *TopicPolicyV1Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data schema.ConsoleTopicPolicyV1Model

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var tfPolicies map[string]topicPolicy.PoliciesValue
	diag := data.Spec.Policies.ElementsAs(ctx, &tfPolicies, false)
	if diag.HasError() {
		return
	}

	for k, v := range tfPolicies {
		if !schemaUtils.AttrIsSet(v.AllowedKeys) && !schemaUtils.AttrIsSet(v.Match) && !schemaUtils.AttrIsSet(v.OneOf) && !schemaUtils.AttrIsSet(v.NoneOf) && !schemaUtils.AttrIsSet(v.Range) {
			resp.Diagnostics.AddError(
				"Invalid Attribute Configuration",
				"Policy '"+k+"' must have one of the following constraints: allowed_keys, match, one_of, none_of, range",
			)
		}
		if schemaUtils.AttrIsSet(v.AllowedKeys) && (schemaUtils.AttrIsSet(v.Match) || schemaUtils.AttrIsSet(v.OneOf) || schemaUtils.AttrIsSet(v.NoneOf) || schemaUtils.AttrIsSet(v.Range)) {
			resp.Diagnostics.AddError(
				"Invalid Attribute Combination",
				"Policy '"+k+"' can only have one of the following constraints: allowed_keys, match, one_of, none_of, range",
			)
		}
		if schemaUtils.AttrIsSet(v.Match) && (schemaUtils.AttrIsSet(v.AllowedKeys) || schemaUtils.AttrIsSet(v.OneOf) || schemaUtils.AttrIsSet(v.NoneOf) || schemaUtils.AttrIsSet(v.Range)) {
			resp.Diagnostics.AddError(
				"Invalid Attribute Combination",
				"Policy '"+k+"' can only have one of the following constraints: allowed_keys, match, one_of, none_of, range",
			)
		}
		if schemaUtils.AttrIsSet(v.OneOf) && (schemaUtils.AttrIsSet(v.AllowedKeys) || schemaUtils.AttrIsSet(v.Match) || schemaUtils.AttrIsSet(v.NoneOf) || schemaUtils.AttrIsSet(v.Range)) {
			resp.Diagnostics.AddError(
				"Invalid Attribute Combination",
				"Policy '"+k+"' can only have one of the following constraints: allowed_keys, match, one_of, none_of, range",
			)
		}
		if schemaUtils.AttrIsSet(v.NoneOf) && (schemaUtils.AttrIsSet(v.AllowedKeys) || schemaUtils.AttrIsSet(v.Match) || schemaUtils.AttrIsSet(v.OneOf) || schemaUtils.AttrIsSet(v.Range)) {
			resp.Diagnostics.AddError(
				"Invalid Attribute Combination",
				"Policy '"+k+"' can only have one of the following constraints: allowed_keys, match, one_of, none_of, range",
			)
		}
		if schemaUtils.AttrIsSet(v.Range) && (schemaUtils.AttrIsSet(v.AllowedKeys) || schemaUtils.AttrIsSet(v.Match) || schemaUtils.AttrIsSet(v.OneOf) || schemaUtils.AttrIsSet(v.NoneOf)) {
			resp.Diagnostics.AddError(
				"Invalid Attribute Combination",
				"Policy '"+k+"' can only have one of the following constraints: allowed_keys, match, one_of, none_of, range",
			)
		}
	}
}

func (r *TopicPolicyV1Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.ConsoleTopicPolicyV1Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating topic policy named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create topic policy with desired state : %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create topic policy, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Topic Policy to create : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, topicPolicyV1ApiPath, consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create topic policy, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Topic Policy created with result: %s", apply.UpsertResult))

	var consoleRes = console.TopicPolicyResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as topic policy : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New topic policy state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read topic policy, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TopicPolicyV1Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.ConsoleTopicPolicyV1Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read topic policy named %s", data.Name.String()))
	get, err := r.apiClient.Describe(ctx, fmt.Sprintf("%s/%s", topicPolicyV1ApiPath, data.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read topic policy, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Topic Policy %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var consoleRes = console.TopicPolicyResource{}
	err = jsoniter.Unmarshal(get, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Parsing Error", fmt.Sprintf("Unable to read topic policy, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New topic policy state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read topic policy, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TopicPolicyV1Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.ConsoleTopicPolicyV1Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating topic policy named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update topic policy with TF data: %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create topic policy, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Topic Policy to update : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, topicPolicyV1ApiPath, consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create topic policy, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Topic Policy updated with result: %s", apply))

	var consoleRes = console.TopicPolicyResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as topic policy : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New topic policy state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read topic policy, got error: %s", err))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TopicPolicyV1Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.ConsoleTopicPolicyV1Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Deleting topic policy named %s", data.Name.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	resourcePath := fmt.Sprintf("%s/%s", topicPolicyV1ApiPath, data.Name.ValueString())
	err := r.apiClient.Delete(ctx, client.CONSOLE, resourcePath, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete topic policy, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Topic Policy %s deleted", data.Name.String()))
}

func (r *TopicPolicyV1Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
