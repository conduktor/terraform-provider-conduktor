package provider

import (
	"context"
	"fmt"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/group_v2"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_group_v2"
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
	resp.TypeName = req.ProviderTypeName + "_group_v2"
}

func (r *GroupV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.GroupV2ResourceSchema(ctx)
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

	r.apiClient = data.Client
}

func (r *GroupV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.GroupV2Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Create group named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create group with TF data: %+v", data))

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

	tflog.Debug(ctx, fmt.Sprintf("Group created with result: %s", apply))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.GroupV2Model

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

	var consoleRes = model.GroupConsoleResource{}
	err = jsoniter.Unmarshal(get, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Parsing Error", fmt.Sprintf("Unable to read group, got error: %s", err))
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

func (r *GroupV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.GroupV2Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Update group named %s", data.Name.String()))
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.GroupV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Delete group named %s", data.Name.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.apiClient.Delete(ctx, fmt.Sprintf("%s/%s", groupV2ApiPath, data.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete group, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Group %s deleted", data.Name.String()))
}

func (r *GroupV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
