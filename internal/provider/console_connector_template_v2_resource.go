package provider

import (
	"context"
	"fmt"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/console_connector_template_v2"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_connector_template_v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/mod/semver"
)

const connectorTemplateV2ApiPath = "/public/console/v2/connector-template"
const connectorTemplateV2MinimumVersion = "v1.46.0"

var _ resource.Resource = &ConnectorTemplateV2Resource{}
var _ resource.ResourceWithImportState = &ConnectorTemplateV2Resource{}

func NewConnectorTemplateV2Resource() resource.Resource {
	return &ConnectorTemplateV2Resource{}
}

type ConnectorTemplateV2Resource struct {
	apiClient *client.Client
}

func (r *ConnectorTemplateV2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_connector_template_v2"
}

func (r *ConnectorTemplateV2Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.ConsoleConnectorTemplateV2ResourceSchema(ctx)
}

func (r *ConnectorTemplateV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	if semver.IsValid(consoleVersion) && semver.Compare(consoleVersion, connectorTemplateV2MinimumVersion) < 0 {
		resp.Diagnostics.AddError(
			"Minimum version requirement not met",
			"This resource requires Conduktor Console API version "+connectorTemplateV2MinimumVersion+" but targeted Conduktor Console API is "+consoleVersion,
		)
		return
	}

	checkEnterprisePlanRequirement(ctx, data.Client, consoleVersion, connectorTemplateV2MinimumVersion, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	r.apiClient = data.Client
}

func (r *ConnectorTemplateV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.ConsoleConnectorTemplateV2Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating connector template named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create connector template with desired state: %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create connector template, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Connector template to create: %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, connectorTemplateV2ApiPath, consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create connector template, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Connector template created with result: %s", apply.UpsertResult))

	var consoleRes console.ConnectorTemplateResource
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as connector template: %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New connector template state: %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read connector template, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConnectorTemplateV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.ConsoleConnectorTemplateV2Model

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read connector template named %s", data.Name.String()))
	get, err := r.apiClient.Describe(ctx, fmt.Sprintf("%s/%s", connectorTemplateV2ApiPath, data.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read connector template, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Connector template %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var consoleRes console.ConnectorTemplateResource
	err = jsoniter.Unmarshal(get, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Parsing Error", fmt.Sprintf("Unable to read connector template, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New connector template state: %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read connector template, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConnectorTemplateV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.ConsoleConnectorTemplateV2Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating connector template named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update connector template with TF data: %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to update connector template, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Connector template to update: %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, connectorTemplateV2ApiPath, consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update connector template, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Connector template updated with result: %s", apply.UpsertResult))

	var consoleRes console.ConnectorTemplateResource
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as connector template: %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New connector template state: %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read connector template, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConnectorTemplateV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.ConsoleConnectorTemplateV2Model

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting connector template named %s", data.Name.String()))

	resourcePath := fmt.Sprintf("%s/%s", connectorTemplateV2ApiPath, data.Name.ValueString())
	err := r.apiClient.Delete(ctx, client.CONSOLE, resourcePath, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete connector template, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Connector template %s deleted", data.Name.String()))
}

func (r *ConnectorTemplateV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
