package provider

import (
	"context"
	"fmt"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/console_group_v2"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_group_v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
)

const groupV2ApiPath = "/public/iam/v2/group"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GroupV2Resource{}
var _ resource.ResourceWithImportState = &GroupV2Resource{}

func NewGroupV2Resource() resource.Resource {
	return &GroupV2Resource{}
}

// GroupV2Resource defines the resource implementation.
type GroupV2Resource struct {
	apiClient *client.Client
}

func (r *GroupV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_group_v2"
}

func (r *GroupV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.ConsoleGroupV2ResourceSchema(ctx)
}

func (r *GroupV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GroupV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.ConsoleGroupV2Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating group named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create group with desired state : %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create group, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Group to create : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, groupV2ApiPath, consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create group, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Group created with result: %s", apply.UpsertResult))

	var consoleRes = console.GroupConsoleResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as group : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New group state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read group, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.ConsoleGroupV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read group named %s", data.Name.String()))
	get, err := r.apiClient.Describe(ctx, fmt.Sprintf("%s/%s", groupV2ApiPath, data.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read group, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Group %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var consoleRes = console.GroupConsoleResource{}
	err = jsoniter.Unmarshal(get, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Parsing Error", fmt.Sprintf("Unable to read group, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New group state : %+v", consoleRes))

	// Store the current value of members_from_external_groups
	currentMembersFromExternalGroups := data.Spec.MembersFromExternalGroups

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read group, got error: %s", err))
		return
	}

	// Use the current value for members_from_external_groups to prevent unnecessary state changes
	// Only use the backend value if this is the first read (currentMembersFromExternalGroups is empty)
	if !currentMembersFromExternalGroups.IsNull() && !currentMembersFromExternalGroups.IsUnknown() {
		data.Spec.MembersFromExternalGroups = currentMembersFromExternalGroups
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.ConsoleGroupV2Model
	var state schema.ConsoleGroupV2Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Read the current state to get the current value of members_from_external_groups
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Use the state value for members_from_external_groups instead of the plan value
	// This effectively ignores any changes to this field during planning
	data.Spec.MembersFromExternalGroups = state.Spec.MembersFromExternalGroups

	tflog.Info(ctx, fmt.Sprintf("Updating group named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update group with TF data: %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create group, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Group to update : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, groupV2ApiPath, consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create group, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Group updated with result: %s", apply))

	var consoleRes = console.GroupConsoleResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as group : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New group state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read group, got error: %s", err))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.ConsoleGroupV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Deleting group named %s", data.Name.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	resourcePath := fmt.Sprintf("%s/%s", groupV2ApiPath, data.Name.ValueString())
	err := r.apiClient.Delete(ctx, client.CONSOLE, resourcePath, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete group, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Group %s deleted", data.Name.String()))
}

func (r *GroupV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
