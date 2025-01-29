package provider

import (
	"context"
	"fmt"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/gateway_token_v2"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_token_v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const gatewayTokenV2ApiPath = "/gateway/v2/token"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GatewayTokenV2Resource{}
var _ resource.ResourceWithImportState = &GatewayTokenV2Resource{}

func NewGatewayTokenV2Resource() resource.Resource {
	return &GatewayTokenV2Resource{}
}

// GatewayTokenV2Resource defines the resource implementation.
type GatewayTokenV2Resource struct {
	apiClient *client.Client
}

func (r *GatewayTokenV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway_token_v2"
}

func (r *GatewayTokenV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.GatewayTokenV2ResourceSchema(ctx)
}

func (r *GatewayTokenV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *GatewayTokenV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.GatewayTokenV2Model

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var gatewayRes gateway.GatewayTokenResource
	if data.Token.IsNull() {
		tflog.Info(ctx, fmt.Sprintf("Create token for username %s", data.Username.String()))
		tflog.Trace(ctx, fmt.Sprintf("Create token with TF data: %+v", data))

		gatewayResource, err := mapper.TFToInternalModel(ctx, &data)
		if err != nil {
			resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create token, got error: %s", err))
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Token to create : %+v", gatewayResource))

		apply, err := d.apiClient.ApplyGatewayToken(ctx, gatewayTokenV2ApiPath, gatewayResource)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create token, got error: %s", err))
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Token created with result: %s", apply))

		err = gatewayRes.FromRawJsonInterface(apply.Resource)
		if err != nil {
			resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as token : %v, got error: %s", apply.Resource, err))
			return
		}
		gatewayRes.VCluster = data.Vcluster.ValueString()
		gatewayRes.Username = data.Username.ValueString()
		gatewayRes.LifetimeSeconds = data.LifetimeSeconds.ValueInt64()
		tflog.Debug(ctx, fmt.Sprintf("New token state : %+v", gatewayRes))
	}

	data, err := mapper.InternalModelToTerraform(ctx, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read token, got error: %s", err))
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayTokenV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.GatewayTokenV2Model
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Create token for service account %s", data.Username.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create token with TF data: %+v", data))

	gatewayResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create token, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Token to create : %+v", gatewayResource))

	apply, err := r.apiClient.ApplyGatewayToken(ctx, gatewayTokenV2ApiPath, gatewayResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create token, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("token created with result: %s", apply))

	var gatewayRes gateway.GatewayTokenResource
	err = gatewayRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as token : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New token state : %+v", gatewayRes))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read token, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// func (r *GatewayTokenV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// 	var data schema.GatewayTokenV2Model
//
// 	// Read Terraform prior state data into the model
// 	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
//
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
//
// 	// Only appending vcluster if present
// 	queryString := "name=" + data.Name.ValueString()
// 	if data.Vcluster.ValueString() != "" {
// 		queryString += "&vcluster=" + data.Vcluster.ValueString()
// 	}
//
// 	tflog.Info(ctx, fmt.Sprintf("Read token named %s", data.Name.String()))
// 	get, err := r.apiClient.Describe(ctx, fmt.Sprintf("%s?%s", gatewayTokenV2ApiPath, queryString))
// 	if err != nil {
// 		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read token, got error: %s", err))
// 		return
// 	}
//
// 	if len(get) == 0 {
// 		tflog.Debug(ctx, fmt.Sprintf("Token %s not found, removing from state", data.Name.String()))
// 		resp.State.RemoveResource(ctx)
// 		return
// 	}
//
// 	var gatewayRes = []gateway.GatewayTokenResource{}
// 	err = json.Unmarshal(get, &gatewayRes)
// 	if err != nil || len(gatewayRes) < 1 {
// 		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read token, got error: %s", err))
// 		return
// 	}
// 	tflog.Debug(ctx, fmt.Sprintf("New token state : %+v", gatewayRes))
//
// 	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes[0])
// 	if err != nil {
// 		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read token, got error: %s", err))
// 		return
// 	}
//
// 	// Save updated data into Terraform state
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
// }

func (r *GatewayTokenV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.GatewayTokenV2Model
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Update token for service account %s", data.Username.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update token with TF data: %+v", data))

	gatewayResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to update token, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Token to update : %+v", gatewayResource))

	apply, err := r.apiClient.ApplyGatewayToken(ctx, gatewayTokenV2ApiPath, gatewayResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create token, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Token updated with result: %s", apply))

	var gatewayRes gateway.GatewayTokenResource
	err = gatewayRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as token : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New token state : %+v", gatewayRes))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read token, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayTokenV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.GatewayTokenV2Model
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Delete token for service account: %s", data.Username.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	// deleteRes := gateway.GatewayTokenMetadata{
	// 	Name:     data.Name.ValueString(),
	// 	VCluster: data.Vcluster.ValueString(),
	// }

	// err := r.apiClient.Delete(ctx, client.GATEWAY, gatewayTokenV2ApiPath, deleteRes)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete token account, got error: %s", err))
	// 	return
	// }
	tflog.Debug(ctx, fmt.Sprintf("Token deleted for service account: %s", data.Username.String()))
}

func (r *GatewayTokenV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
