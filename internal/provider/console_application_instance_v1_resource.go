package provider

import (
	"context"
	"fmt"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/console_application_instance_v1"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_application_instance_v1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
)

const applicationInstanceV1ApiPath = "/public/self-serve/v1/application-instance"
const appInstanceMininumVersion = "v1.31.0"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ApplicationInstanceV1Resource{}
var _ resource.ResourceWithImportState = &ApplicationInstanceV1Resource{}

func NewApplicationInstanceV1Resource() resource.Resource {
	return &ApplicationInstanceV1Resource{}
}

// ApplicationInstanceV1Resource defines the resource implementation.
type ApplicationInstanceV1Resource struct {
	apiClient *client.Client
}

func (r *ApplicationInstanceV1Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_application_instance_v1"
}

func (r *ApplicationInstanceV1Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.ConsoleApplicationInstanceV1ResourceSchema(ctx)
}

func (r *ApplicationInstanceV1Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationInstanceV1Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.ConsoleApplicationInstanceV1Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating application instance named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create application instance with desired state : %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create application instance, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Application Instance to create : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, applicationInstanceV1ApiPath, consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create application instance, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Application Instance created with result: %s", apply.UpsertResult))

	var consoleRes = console.ApplicationInstanceConsoleResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as application instance : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New application instance state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read application instance, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationInstanceV1Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.ConsoleApplicationInstanceV1Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read application instance named %s", data.Name.String()))
	get, err := r.apiClient.Describe(ctx, fmt.Sprintf("%s/%s", applicationInstanceV1ApiPath, data.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read application instance, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Application Instance %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var consoleRes = console.ApplicationInstanceConsoleResource{}
	err = jsoniter.Unmarshal(get, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Parsing Error", fmt.Sprintf("Unable to read application instance, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New application instance state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read application instance, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationInstanceV1Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.ConsoleApplicationInstanceV1Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating application instance named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update application instance with TF data: %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create application instance, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Application Instance to update : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, applicationInstanceV1ApiPath, consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create application instance, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Application Instance updated with result: %s", apply))

	var consoleRes = console.ApplicationInstanceConsoleResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as application instance : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New application instance state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read application instance, got error: %s", err))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationInstanceV1Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.ConsoleApplicationInstanceV1Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Deleting application instance named %s", data.Name.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	resourcePath := fmt.Sprintf("%s/%s", applicationInstanceV1ApiPath, data.Name.ValueString())
	err := r.apiClient.Delete(ctx, client.CONSOLE, resourcePath, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete application instance, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Application Instance %s deleted", data.Name.String()))
}

func (r *ApplicationInstanceV1Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
