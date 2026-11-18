package provider

import (
	"context"
	"fmt"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/gateway_group_v2"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_group_v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/mod/semver"
)

const gatewayGroupV2ApiPath = "/gateway/v2/group"
const gatewayGroupV2MinimumVersion = "v3.3.0"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GatewayGroupV2Resource{}
var _ resource.ResourceWithImportState = &GatewayGroupV2Resource{}

func NewGatewayGroupV2Resource() resource.Resource {
	return &GatewayGroupV2Resource{}
}

// GatewayGroupV2Resource defines the resource implementation.
type GatewayGroupV2Resource struct {
	apiClient *client.Client
}

func (r *GatewayGroupV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway_group_v2"
}

func (r *GatewayGroupV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.GatewayGroupV2ResourceSchema(ctx)
}

func (r *GatewayGroupV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	if data.Client == nil || data.Mode != client.GATEWAY {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Gateway Client not configured. Please provide client configuration details for Gateway API and ensure you have set the right provider mode for this resource. \n"+
				"More info here: \n"+
				" - https://registry.terraform.io/providers/conduktor/conduktor/latest/docs",
		)
		return
	}

	gatewayVersion, err := data.Client.GetAPIVersion(ctx, client.GATEWAY)
	if err != nil {
		resp.Diagnostics.AddError("Error fetching Gateway version", err.Error())
		return
	}
	if semver.IsValid(gatewayVersion) && semver.Compare(gatewayVersion, gatewayGroupV2MinimumVersion) < 0 {
		resp.Diagnostics.AddError(
			"Minimum version requirement not met",
			fmt.Sprintf("This resource requires Conduktor Gateway API version %s or higher, but targeted Conduktor Gateway API is %s", gatewayGroupV2MinimumVersion, gatewayVersion),
		)
		return
	}

	r.apiClient = data.Client
}

func (r *GatewayGroupV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.GatewayGroupV2Model
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating gateway group named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create gateway group with TF data: %+v", data))

	gatewayResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create gateway group, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Gateway group to create: %+v", gatewayResource))

	apply, err := r.apiClient.Apply(ctx, gatewayGroupV2ApiPath, gatewayResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create gateway group, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Gateway group created with result: %s", apply.UpsertResult))

	var gatewayRes gateway.GatewayGroupResource
	err = gatewayRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as gateway group: %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New gateway group state: %+v", gatewayRes))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read gateway group, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayGroupV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.GatewayGroupV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read gateway group named %s", data.Name.String()))
	get, err := r.apiClient.Describe(ctx, fmt.Sprintf("%s/%s", gatewayGroupV2ApiPath, data.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read gateway group, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Gateway group %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var gatewayRes gateway.GatewayGroupResource
	err = jsoniter.Unmarshal(get, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Parsing Error", fmt.Sprintf("Unable to read gateway group, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Gateway group state: %+v", gatewayRes))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read gateway group, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayGroupV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.GatewayGroupV2Model
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating gateway group named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update gateway group with TF data: %+v", data))

	gatewayResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to update gateway group, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Gateway group to update: %+v", gatewayResource))

	apply, err := r.apiClient.Apply(ctx, gatewayGroupV2ApiPath, gatewayResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update gateway group, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Gateway group updated with result: %s", apply))

	var gatewayRes gateway.GatewayGroupResource
	err = gatewayRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as gateway group: %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New gateway group state: %+v", gatewayRes))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read gateway group, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayGroupV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.GatewayGroupV2Model
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting gateway group named %s", data.Name.String()))

	resourcePath := fmt.Sprintf("%s/%s", gatewayGroupV2ApiPath, data.Name.ValueString())
	err := r.apiClient.Delete(ctx, client.GATEWAY, resourcePath, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete gateway group, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Gateway group %s deleted", data.Name.String()))
}

func (r *GatewayGroupV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
