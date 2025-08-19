package provider

import (
	"context"
	"fmt"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/gateway_virtual_cluster_v2"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_virtual_cluster_v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
)

const virtualClusterV2ApiPath = "/gateway/v2/virtual-cluster"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &VirtualClusterV2Resource{}
var _ resource.ResourceWithImportState = &VirtualClusterV2Resource{}

func NewVirtualClusterV2Resource() resource.Resource {
	return &VirtualClusterV2Resource{}
}

// VirtualClusterV2Resource defines the resource implementation.
type VirtualClusterV2Resource struct {
	apiClient *client.Client
}

func (r *VirtualClusterV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway_virtual_cluster_v2"
}

func (r *VirtualClusterV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.GatewayVirtualClusterV2ResourceSchema(ctx)
}

func (r *VirtualClusterV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.apiClient = data.Client
}

func (r *VirtualClusterV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.GatewayVirtualClusterV2Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating virtual cluster named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create virtual cluster with desired state : %+v", data))

	gatewayResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create virtual cluster, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Virtual Cluster to create : %+v", gatewayResource))
	apply, err := r.apiClient.Apply(ctx, virtualClusterV2ApiPath, gatewayResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create virtual cluster, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Virtual Cluster created with result: %s", apply.UpsertResult))

	var gatewayRes = gateway.VirtualClusterResource{}
	err = gatewayRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as virtual cluster : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New virtual cluster state : %+v", gatewayRes))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read virtual cluster, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtualClusterV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.GatewayVirtualClusterV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read virtual cluster named %s", data.Name.String()))
	get, err := r.apiClient.Describe(ctx, fmt.Sprintf("%s/%s", virtualClusterV2ApiPath, data.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read virtual cluster, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Virtual Cluster %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var gatewayRes = gateway.VirtualClusterResource{}
	err = jsoniter.Unmarshal(get, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Parsing Error", fmt.Sprintf("Unable to read virtual cluster, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New virtual cluster state : %+v", gatewayRes))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read virtual cluster, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtualClusterV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.GatewayVirtualClusterV2Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating virtual cluster named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update virtual cluster with TF data: %+v", data))

	gatewayResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create virtual cluster, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Virtual Cluster to update : %+v", gatewayResource))

	apply, err := r.apiClient.Apply(ctx, virtualClusterV2ApiPath, gatewayResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create virtual cluster, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Virtual Cluster updated with result: %s", apply))

	var gatewayRes = gateway.VirtualClusterResource{}
	err = gatewayRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as virtual cluster : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New virtual cluster state : %+v", gatewayRes))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read virtual cluster, got error: %s", err))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtualClusterV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.GatewayVirtualClusterV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Deleting virtual cluster named %s", data.Name.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	resourcePath := fmt.Sprintf("%s/%s", virtualClusterV2ApiPath, data.Name.ValueString())
	// Although this is a Gateway resource, it uses the same mode as the Console API, so we use the CONSOLE mode here.
	// i.e. ID of the resource is expected in the URL path.
	err := r.apiClient.Delete(ctx, client.CONSOLE, resourcePath, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete virtual cluster, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Virtual Cluster %s deleted", data.Name.String()))
}

func (r *VirtualClusterV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
