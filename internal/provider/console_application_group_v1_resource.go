package provider

import (
	"context"
	"fmt"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/console_application_group_v1"
	"github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_application_group_v1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
)

const applicationGroupV1ApiPath = "/public/self-serve/v1/application-group"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ApplicationGroupV1Resource{}
var _ resource.ResourceWithImportState = &ApplicationGroupV1Resource{}

func NewApplicationGroupV1Resource() resource.Resource {
	return &GroupV2Resource{}
}

// GroupV2Resource defines the resource implementation.
type ApplicationGroupV1Resource struct {
	apiClient *client.Client
}

func (r *ApplicationGroupV1Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_application_group_v1"
}
func (r *ApplicationGroupV1Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.ConsoleApplicationGroupV1ResourceSchema(ctx)
}

func (r *ApplicationGroupV1Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
				"Please refer to the documentation for more information.",
		)
		return
	}

	r.apiClient = data.Client
}

func (r *ApplicationGroupV1Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.ConsoleApplicationGroupV1Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating Application Group named %s", data.Name.String()))
	tflog.Debug(ctx, fmt.Sprintf("Creating Application Group with desired state : %v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error,", fmt.Sprintf("Unable to create group, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Creating Application Group with internal model : %v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, applicationGroupV1ApiPath, consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create group, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Created Application Group with response : %v", apply))

	var consoleRes = console.ApplicationGroupConsoleResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshalling Error", fmt.Sprintf("Response resource can't be cast as application group : %v got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New Application Group state : %v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create group, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationGroupV1Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.ConsoleApplicationGroupV1Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Reading Application Group named %s", data.Name.String()))

	get, err := r.apiClient.Describe(ctx, fmt.Sprintf("%s/%s", applicationGroupV1ApiPath, data.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Application Group, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Application Group %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var consoleRes = console.ApplicationGroupConsoleResource{}
	err = jsoniter.Unmarshal(get, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Parsing Error", fmt.Sprintf("Unable to read Application Group, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Application Group state : %v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read Application Group, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationGroupV1Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.ConsoleApplicationGroupV1Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating Application Group named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update Application Group with TF data: %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error,", fmt.Sprintf("Unable to update application group, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Updating Application Group with internal model : %v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, applicationGroupV1ApiPath, consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update group, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Application Group updated with response : %v", apply))

	var consoleRes = console.ApplicationGroupConsoleResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshalling Error", fmt.Sprintf("Response resource can't be cast as application group : %v got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New Application Group state : %v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to update group, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationGroupV1Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.ConsoleApplicationGroupV1Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Deleting Application Group named %s", data.Name.String()))

	err := r.apiClient.Delete(ctx, client.CONSOLE, fmt.Sprintf("%s/%s", applicationGroupV1ApiPath, data.Name.ValueString()), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete group, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Application Group %s deleted", data.Name.String()))
}

func (r *ApplicationGroupV1Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
