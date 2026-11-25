package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/gateway_alias_topic_v2"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_alias_topic_v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/mod/semver"
)

const gatewayAliasTopicV2ApiPath = "/gateway/v2/alias-topic"
const gatewayAliasTopicV2MinimumVersion = "v3.3.0"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GatewayAliasTopicV2Resource{}
var _ resource.ResourceWithImportState = &GatewayAliasTopicV2Resource{}

func NewGatewayAliasTopicV2Resource() resource.Resource {
	return &GatewayAliasTopicV2Resource{}
}

// GatewayAliasTopicV2Resource defines the resource implementation.
type GatewayAliasTopicV2Resource struct {
	apiClient *client.Client
}

func (r *GatewayAliasTopicV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway_alias_topic_v2"
}

func (r *GatewayAliasTopicV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.GatewayAliasTopicV2ResourceSchema(ctx)
}

func (r *GatewayAliasTopicV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	if semver.IsValid(gatewayVersion) && semver.Compare(gatewayVersion, gatewayAliasTopicV2MinimumVersion) < 0 {
		resp.Diagnostics.AddError(
			"Minimum version requirement not met",
			fmt.Sprintf("This resource requires Conduktor Gateway API version %s or higher, but targeted Conduktor Gateway API is %s", gatewayAliasTopicV2MinimumVersion, gatewayVersion),
		)
		return
	}

	r.apiClient = data.Client
}

func (r *GatewayAliasTopicV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.GatewayAliasTopicV2Model
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating gateway alias topic named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create gateway alias topic with TF data: %+v", data))

	gatewayResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create gateway alias topic, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Gateway alias topic to create: %+v", gatewayResource))

	apply, err := r.apiClient.Apply(ctx, gatewayAliasTopicV2ApiPath, gatewayResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create gateway alias topic, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Gateway alias topic created with result: %s", apply.UpsertResult))

	// The upsert response does not include the full resource; read back to get server-populated fields.
	data, err = r.readAliasTopicState(ctx, data.Name.ValueString(), data.Vcluster.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read gateway alias topic after create, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayAliasTopicV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.GatewayAliasTopicV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	queryString := "name=" + data.Name.ValueString()
	if data.Vcluster.ValueString() != "" {
		queryString += "&vcluster=" + data.Vcluster.ValueString()
	}

	tflog.Info(ctx, fmt.Sprintf("Read gateway alias topic named %s on vcluster %s", data.Name.String(), data.Vcluster.String()))
	get, err := r.apiClient.Describe(ctx, fmt.Sprintf("%s?%s", gatewayAliasTopicV2ApiPath, queryString))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read gateway alias topic, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Gateway alias topic %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var gatewayResult = []gateway.GatewayAliasTopicResource{}
	err = json.Unmarshal(get, &gatewayResult)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read gateway alias topic, got error: %s", err))
		return
	}
	if len(gatewayResult) > 1 {
		tflog.Warn(ctx, fmt.Sprintf("Multiple alias topics found with query name %s and vcluster %s, using best match", data.Name.String(), data.Vcluster.ValueString()))
	}
	var gatewayResource = gateway.GatewayAliasTopicResource{}
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
		tflog.Debug(ctx, fmt.Sprintf("Gateway alias topic %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Gateway alias topic state: %+v", gatewayResource))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayResource)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read gateway alias topic, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayAliasTopicV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.GatewayAliasTopicV2Model
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating gateway alias topic named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update gateway alias topic with TF data: %+v", data))

	gatewayResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to update gateway alias topic, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Gateway alias topic to update: %+v", gatewayResource))

	apply, err := r.apiClient.Apply(ctx, gatewayAliasTopicV2ApiPath, gatewayResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update gateway alias topic, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Gateway alias topic updated with result: %s", apply.UpsertResult))

	// The upsert response does not include the full resource; read back to get server-populated fields.
	data, err = r.readAliasTopicState(ctx, data.Name.ValueString(), data.Vcluster.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read gateway alias topic after update, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// readAliasTopicState fetches the current alias topic from the API and returns the TF model.
// Used after create/update since the upsert response does not include the full resource body.
func (r *GatewayAliasTopicV2Resource) readAliasTopicState(ctx context.Context, name, vcluster string) (schema.GatewayAliasTopicV2Model, error) {
	queryString := "name=" + name + "&vcluster=" + vcluster
	get, err := r.apiClient.Describe(ctx, fmt.Sprintf("%s?%s", gatewayAliasTopicV2ApiPath, queryString))
	if err != nil {
		return schema.GatewayAliasTopicV2Model{}, err
	}

	var gatewayResult []gateway.GatewayAliasTopicResource
	if err := json.Unmarshal(get, &gatewayResult); err != nil {
		return schema.GatewayAliasTopicV2Model{}, err
	}

	for _, res := range gatewayResult {
		nameMatch := res.Metadata.Name == name
		vclusterMatch := res.Metadata.VCluster == vcluster
		passthroughMatch := res.Metadata.VCluster == "" && vcluster == "passthrough"
		if nameMatch && (vclusterMatch || passthroughMatch) {
			return mapper.InternalModelToTerraform(ctx, &res)
		}
	}

	return schema.GatewayAliasTopicV2Model{}, fmt.Errorf("alias topic %s/%s not found after upsert", name, vcluster)
}

func (r *GatewayAliasTopicV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.GatewayAliasTopicV2Model
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Deleting gateway alias topic named %s", data.Name.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	deleteRes := gateway.GatewayAliasTopicMetadata{
		Name:     data.Name.ValueString(),
		VCluster: data.Vcluster.ValueString(),
	}

	err := r.apiClient.Delete(ctx, client.GATEWAY, gatewayAliasTopicV2ApiPath, deleteRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete gateway alias topic, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Gateway alias topic %s deleted", data.Name.String()))
}

// ImportState imports the state of the resource from the given ID.
// The ID is expected to be in the format: <alias_topic_name>/<vcluster>.
func (r *GatewayAliasTopicV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <alias_topic_name>/<vcluster>. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vcluster"), idParts[1])...)
}
