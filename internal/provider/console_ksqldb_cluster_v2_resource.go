package provider

import (
	"context"
	"fmt"
	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/console_ksqldb_cluster_v2"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_ksqldb_cluster_v2"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

func ksqldbClusterV2ApiPutPath(cluster string) string {
	return fmt.Sprintf("/public/console/v2/cluster/%s/ksqldb", cluster)
}

func ksqldbClusterV2ApiGetPath(cluster string, ksqldbClusterName string) string {
	return fmt.Sprintf("/public/console/v2/cluster/%s/ksqldb/%s", cluster, ksqldbClusterName)
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &KsqlDBClusterV2Resource{}
var _ resource.ResourceWithImportState = &KsqlDBClusterV2Resource{}
var _ resource.ResourceWithConfigValidators = &KsqlDBClusterV2Resource{}

func NewKsqlDBClusterV2Resource() resource.Resource {
	return &KsqlDBClusterV2Resource{}
}

// KsqlDBClusterV2Resource defines the resource implementation.
type KsqlDBClusterV2Resource struct {
	apiClient *client.Client
}

func (r *KsqlDBClusterV2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_ksqldb_cluster_v2"
}

func (r *KsqlDBClusterV2Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.ConsoleKsqldbClusterV2ResourceSchema(ctx)
}

func (r *KsqlDBClusterV2Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *KsqlDBClusterV2Resource) ConfigValidators(_ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("spec").AtName("security").AtName("basic_auth"),
			path.MatchRoot("spec").AtName("security").AtName("bearer_token"),
			path.MatchRoot("spec").AtName("security").AtName("ssl_auth"),
		),
	}
}

func (r *KsqlDBClusterV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.ConsoleKsqldbClusterV2Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating KsqlDB cluster server named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create KsqlDB cluster server with desired state : %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create KsqlDB cluster server, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("KsqlDB cluster server to create : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, ksqldbClusterV2ApiPutPath(consoleResource.Metadata.Cluster), consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create KsqlDB cluster server, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("KsqlDB cluster server created with result: %s", apply.UpsertResult))

	var consoleRes = console.KsqlDBClusterResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as KsqlDB cluster server : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New KsqlDB cluster server state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read KsqlDB cluster server, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KsqlDBClusterV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.ConsoleKsqldbClusterV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read KsqlDB cluster server named %s", data.Name.String()))
	get, err := r.apiClient.Describe(ctx, ksqldbClusterV2ApiGetPath(data.Cluster.ValueString(), data.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read KsqlDB cluster server, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("KsqlDB cluster server %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var consoleRes = console.KsqlDBClusterResource{}
	err = jsoniter.Unmarshal(get, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Parsing Error", fmt.Sprintf("Unable to read KsqlDB cluster server, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New KsqlDB cluster server state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read KsqlDB cluster server, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KsqlDBClusterV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.ConsoleKsqldbClusterV2Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating KsqlDB cluster server named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update KsqlDB cluster server with TF data: %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create KsqlDB cluster server, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("KsqlDB cluster server to update : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, ksqldbClusterV2ApiPutPath(consoleResource.Metadata.Cluster), consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create KsqlDB cluster server, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("KsqlDB cluster server updated with result: %s", apply))

	var consoleRes = console.KsqlDBClusterResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as KsqlDB cluster server : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New KsqlDB cluster server state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read KsqlDB cluster server, got error: %s", err))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KsqlDBClusterV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.ConsoleKsqldbClusterV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Deleting KsqlDB cluster server named %s", data.Name.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	resourcePath := ksqldbClusterV2ApiGetPath(data.Cluster.ValueString(), data.Name.ValueString())
	err := r.apiClient.Delete(ctx, client.CONSOLE, resourcePath, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete KsqlDB cluster server, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("KsqlDB cluster server %s deleted", data.Name.String()))
}

func (r *KsqlDBClusterV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: cluster/name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cluster"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
}
