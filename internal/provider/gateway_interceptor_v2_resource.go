package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/gateway_interceptor_v2"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_interceptor_v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

const gatewayInterceptorV2ApiPath = "/gateway/v2/interceptor"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GatewayInterceptorV2Resource{}
var _ resource.ResourceWithImportState = &GatewayInterceptorV2Resource{}

func NewGatewayInterceptorV2Resource() resource.Resource {
	return &GatewayInterceptorV2Resource{}
}

// GatewayInterceptorV2Resource defines the resource implementation.
type GatewayInterceptorV2Resource struct {
	apiClient *client.Client
}

func (r *GatewayInterceptorV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway_interceptor_v2"
}

func (r *GatewayInterceptorV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.GatewayInterceptorV2ResourceSchema(ctx)
}

func (r *GatewayInterceptorV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GatewayInterceptorV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.GatewayInterceptorV2Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Create interceptor named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create interceptor with TF data: %+v", data))

	gatewayResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create interceptor, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Interceptor to create : %+v", gatewayResource))

	apply, err := r.apiClient.Apply(ctx, gatewayInterceptorV2ApiPath, gatewayResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create interceptor, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Interceptor created with result: %s", apply))

	var gatewayRes gateway.GatewayInterceptorResource
	err = gatewayRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as interceptor : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New interceptor state : %+v", gatewayRes))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read interceptor, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayInterceptorV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.GatewayInterceptorV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	queryString := "name=" + data.Name.ValueString()
	queryString += "&global=false"
	if data.Scope.Vcluster.ValueString() != "" {
		queryString += "&vcluster=" + data.Scope.Vcluster.ValueString()
	}
	if data.Scope.Username.ValueString() != "" {
		queryString += "&username=" + data.Scope.Username.ValueString()
	}
	if data.Scope.Group.ValueString() != "" {
		queryString += "&group=" + data.Scope.Group.ValueString()
	}

	tflog.Info(ctx, fmt.Sprintf("Read interceptor named %s", data.Name.String()))
	get, err := r.apiClient.Describe(ctx, fmt.Sprintf("%s?%s", gatewayInterceptorV2ApiPath, queryString))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read interceptor, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Interceptor %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var gatewayRes = []gateway.GatewayInterceptorResource{}
	err = json.Unmarshal(get, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read interceptor, got error: %s", err))
		return
	}
	if len(gatewayRes) < 1 {
		tflog.Debug(ctx, fmt.Sprintf("Interceptor %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New interceptor state : %+v", gatewayRes))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes[0])
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read interceptor, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayInterceptorV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.GatewayInterceptorV2Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Update interceptor named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update interceptor with TF data: %+v", data))

	gatewayResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create interceptor, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Interceptor to update : %+v", gatewayResource))

	apply, err := r.apiClient.Apply(ctx, gatewayInterceptorV2ApiPath, gatewayResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create interceptor, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Interceptor updated with result: %s", apply))

	var gatewayRes gateway.GatewayInterceptorResource
	err = gatewayRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as interceptor : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New interceptor state : %+v", gatewayRes))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read interceptor, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayInterceptorV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.GatewayInterceptorV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Delete interceptor named %s", data.Name.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	deleteScope := gateway.GatewayInterceptorScope{
		VCluster: data.Scope.Vcluster.ValueString(),
		Username: data.Scope.Username.ValueString(),
		Group:    data.Scope.Group.ValueString(),
	}
	deletePath := gatewayInterceptorV2ApiPath + "/" + data.Name.ValueString()

	err := r.apiClient.Delete(ctx, client.GATEWAY, deletePath, deleteScope)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete interceptor, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Interceptor %s deleted", data.Name.String()))
}

// ImportState imports the state of the resource from the given ID.
// The ID is expected to be in the format: <interceptor_name>/<vcluster>/<group>/<username>.
func (r *GatewayInterceptorV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 4 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <interceptor_name>/<vcluster>/<group>/<username>. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[0])...)
	if idParts[1] != "" && idParts[1] != "null" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope").AtName("vcluster"), idParts[1])...)
	}
	if idParts[2] != "" && idParts[2] != "null" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope").AtName("group"), idParts[2])...)
	}
	if idParts[3] != "" && idParts[3] != "null" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope").AtName("username"), idParts[3])...)
	}
}
