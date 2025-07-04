package provider

import (
	"context"
	"fmt"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/console_resource_policy_v1"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_resource_policy_v1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/mod/semver"
)

const resourcePolicyV1ApiPath = "/public/self-serve/v1/resource-policy"
const resourcePolicyMininumVersion = "v1.34.0"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ResourcePolicyV1Resource{}
var _ resource.ResourceWithImportState = &ResourcePolicyV1Resource{}

func NewResourcePolicyV1Resource() resource.Resource {
	return &ResourcePolicyV1Resource{}
}

// ResourcePolicyV1Resource defines the resource implementation.
type ResourcePolicyV1Resource struct {
	apiClient *client.Client
}

func (r *ResourcePolicyV1Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_resource_policy_v1"
}

func (r *ResourcePolicyV1Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.ConsoleResourcePolicyV1ResourceSchema(ctx)
}

func (r *ResourcePolicyV1Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	consoleVersion, err := data.Client.GetAPIVersion(ctx, client.CONSOLE)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching Console version",
			err.Error(),
		)
		return
	}
	if semver.IsValid(consoleVersion) && semver.Compare(consoleVersion, resourcePolicyMininumVersion) < 0 {
		resp.Diagnostics.AddError(
			"Minimum version requirement not met",
			"This resource requires Conduktor Console API version "+resourcePolicyMininumVersion+" but targeted Conduktor Console API is "+consoleVersion,
		)
		return
	}

	r.apiClient = data.Client
}

func (r *ResourcePolicyV1Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.ConsoleResourcePolicyV1Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating resource policy named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create resource policy with desired state : %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create resource policy, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Resource Policy to create : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, resourcePolicyV1ApiPath, consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create resource policy, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Resource Policy created with result: %s", apply.UpsertResult))

	var consoleRes = console.ResourcePolicyConsoleResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as resource policy : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New resource policy state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read resource policy, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourcePolicyV1Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.ConsoleResourcePolicyV1Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read resource policy named %s", data.Name.String()))
	get, err := r.apiClient.Describe(ctx, fmt.Sprintf("%s/%s", resourcePolicyV1ApiPath, data.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read resource policy, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Resource Policy %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var consoleRes = console.ResourcePolicyConsoleResource{}
	err = jsoniter.Unmarshal(get, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Parsing Error", fmt.Sprintf("Unable to read resource policy, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New resource policy state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read resource policy, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourcePolicyV1Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.ConsoleResourcePolicyV1Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating resource policy named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update resource policy with TF data: %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create resource policy, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Resource Policy to update : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, resourcePolicyV1ApiPath, consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create resource policy, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Resource Policy updated with result: %s", apply))

	var consoleRes = console.ResourcePolicyConsoleResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as resource policy : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New resource policy state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read resource policy, got error: %s", err))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourcePolicyV1Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.ConsoleResourcePolicyV1Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Deleting resource policy named %s", data.Name.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	resourcePath := fmt.Sprintf("%s/%s", resourcePolicyV1ApiPath, data.Name.ValueString())
	err := r.apiClient.Delete(ctx, client.CONSOLE, resourcePath, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete resource policy, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Resource Policy %s deleted", data.Name.String()))
}

func (r *ResourcePolicyV1Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
