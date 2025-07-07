package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/gateway_service_account_v2"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_service_account_v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const gatewayServiceAccountV2ApiPath = "/gateway/v2/service-account"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GatewayServiceAccountV2Resource{}
var _ resource.ResourceWithImportState = &GatewayServiceAccountV2Resource{}

func NewGatewayServiceAccountV2Resource() resource.Resource {
	return &GatewayServiceAccountV2Resource{}
}

// GatewayServiceAccountV2Resource defines the resource implementation.
type GatewayServiceAccountV2Resource struct {
	apiClient *client.Client
}

func (r *GatewayServiceAccountV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway_service_account_v2"
}

func (r *GatewayServiceAccountV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.GatewayServiceAccountV2ResourceSchema(ctx)
}

func (r *GatewayServiceAccountV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GatewayServiceAccountV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.GatewayServiceAccountV2Model
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Create service account named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create service account with TF data: %+v", data))

	gatewayResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create service account, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Service Account to create : %+v", gatewayResource))

	apply, err := r.apiClient.Apply(ctx, gatewayServiceAccountV2ApiPath, gatewayResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service account, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Service account created with result: %s", apply))

	var gatewayRes gateway.GatewayServiceAccountResource
	err = gatewayRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as service account : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New service account state : %+v", gatewayRes))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read service account, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayServiceAccountV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.GatewayServiceAccountV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Only appending vcluster if present
	queryString := "name=" + data.Name.ValueString()
	if data.Vcluster.ValueString() != "" {
		queryString += "&vcluster=" + data.Vcluster.ValueString()
	}

	tflog.Info(ctx, fmt.Sprintf("Read service account named %s on vcluster %s", data.Name.String(), data.Vcluster.String()))
	get, err := r.apiClient.Describe(ctx, fmt.Sprintf("%s?%s", gatewayServiceAccountV2ApiPath, queryString))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service account, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Service account %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var gatewayResult = []gateway.GatewayServiceAccountResource{}
	err = json.Unmarshal(get, &gatewayResult)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service account, got error: %s", err))
		return
	}
	if len(gatewayResult) > 1 {
		tflog.Warn(ctx, fmt.Sprintf("Multiple service accounts found with query name %s and vcluster %s, using best match", data.Name.String(), data.Vcluster.ValueString()))
	}
	var gatewayResource = gateway.GatewayServiceAccountResource{}
	var matchFound = false
	for _, res := range gatewayResult {
		nameMatch := res.Metadata.Name == data.Name.ValueString()
		vclusterMatch := res.Metadata.VCluster == data.Vcluster.ValueString()
		passthroughMatch := res.Metadata.VCluster == "" && data.Vcluster.ValueString() == "passthrough"

		if nameMatch && (vclusterMatch || passthroughMatch) {
			gatewayResource = res
			matchFound = true
			break
		}
	}
	if !matchFound {
		tflog.Debug(ctx, fmt.Sprintf("Service account %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("New service account state : %+v", gatewayResource))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayResource)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read service account, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayServiceAccountV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.GatewayServiceAccountV2Model
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Update service account named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update service account with TF data: %+v", data))

	gatewayResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create service account, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Service account to update : %+v", gatewayResource))

	apply, err := r.apiClient.Apply(ctx, gatewayServiceAccountV2ApiPath, gatewayResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service account, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Service account updated with result: %s", apply))

	var gatewayRes gateway.GatewayServiceAccountResource
	err = gatewayRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as service account : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New service account state : %+v", gatewayRes))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read service account, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayServiceAccountV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.GatewayServiceAccountV2Model
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Delete service account named %s", data.Name.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	deleteRes := gateway.GatewayServiceAccountMetadata{
		Name:     data.Name.ValueString(),
		VCluster: data.Vcluster.ValueString(),
	}

	err := r.apiClient.Delete(ctx, client.GATEWAY, gatewayServiceAccountV2ApiPath, deleteRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service account, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Service account %s deleted", data.Name.String()))
}

// ImportState imports the state of the resource from the given ID.
// The ID is expected to be in the format: <service_account_name>/<vcluster> with backward compatibility as <service_account_name> only.
func (r *GatewayServiceAccountV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <service_account_name>/<vcluster>. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[0])...)

	// optional vcluster part for backward compatibility
	if len(idParts) == 2 && idParts[1] != "" && idParts[1] != "null" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vcluster"), idParts[1])...)
	}
}
