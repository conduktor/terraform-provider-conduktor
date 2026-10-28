package provider

import (
	"context"
	"fmt"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/console_topic_template_v2"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_topic_template_v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/mod/semver"
)

const topicTemplateV2ApiPath = "/public/console/v2/topic-template"
const topicTemplateV2MinimumVersion = "v1.26.0"

var _ resource.Resource = &TopicTemplateV2Resource{}
var _ resource.ResourceWithImportState = &TopicTemplateV2Resource{}

func NewTopicTemplateV2Resource() resource.Resource {
	return &TopicTemplateV2Resource{}
}

type TopicTemplateV2Resource struct {
	apiClient *client.Client
}

func (r *TopicTemplateV2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_topic_template_v2"
}

func (r *TopicTemplateV2Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.ConsoleTopicTemplateV2ResourceSchema(ctx)
}

func (r *TopicTemplateV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	consoleVersion, err := data.Client.GetAPIVersion(ctx, client.CONSOLE)
	if err != nil {
		resp.Diagnostics.AddError("Error fetching Console version", err.Error())
		return
	}
	if semver.IsValid(consoleVersion) && semver.Compare(consoleVersion, topicTemplateV2MinimumVersion) < 0 {
		resp.Diagnostics.AddError(
			"Minimum version requirement not met",
			"This resource requires Conduktor Console API version "+topicTemplateV2MinimumVersion+" but targeted Conduktor Console API is "+consoleVersion,
		)
		return
	}

	checkEnterprisePlanRequirement(ctx, data.Client, consoleVersion, topicTemplateV2MinimumVersion, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	r.apiClient = data.Client
}

func (r *TopicTemplateV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.ConsoleTopicTemplateV2Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating topic template named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create topic template with desired state: %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create topic template, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Topic template to create: %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, topicTemplateV2ApiPath, consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create topic template, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Topic template created with result: %s", apply.UpsertResult))

	var consoleRes console.TopicTemplateResource
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as topic template: %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New topic template state: %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read topic template, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TopicTemplateV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.ConsoleTopicTemplateV2Model

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read topic template named %s", data.Name.String()))
	get, err := r.apiClient.Describe(ctx, fmt.Sprintf("%s/%s", topicTemplateV2ApiPath, data.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read topic template, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Topic template %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var consoleRes console.TopicTemplateResource
	err = jsoniter.Unmarshal(get, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Parsing Error", fmt.Sprintf("Unable to read topic template, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New topic template state: %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read topic template, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TopicTemplateV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.ConsoleTopicTemplateV2Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating topic template named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update topic template with TF data: %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to update topic template, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Topic template to update: %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, topicTemplateV2ApiPath, consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update topic template, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Topic template updated with result: %s", apply.UpsertResult))

	var consoleRes console.TopicTemplateResource
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as topic template: %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New topic template state: %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read topic template, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TopicTemplateV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.ConsoleTopicTemplateV2Model

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting topic template named %s", data.Name.String()))

	resourcePath := fmt.Sprintf("%s/%s", topicTemplateV2ApiPath, data.Name.ValueString())
	err := r.apiClient.Delete(ctx, client.CONSOLE, resourcePath, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete topic template, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Topic template %s deleted", data.Name.String()))
}

func (r *TopicTemplateV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
