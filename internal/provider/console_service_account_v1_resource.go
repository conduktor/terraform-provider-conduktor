package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/console_service_account_v1"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_service_account_v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/mod/semver"
)

func serviceAccountV1ApiPutPath(cluster string) string {
	return fmt.Sprintf("/public/self-serve/v1/cluster/%s/service-account", cluster)
}

func serviceAccountV1ApiGetPath(cluster string, connectServerName string) string {
	return fmt.Sprintf("/public/self-serve/v1/cluster/%s/service-account/%s", cluster, connectServerName)
}

const consoleServiceAccountMininumVersion = "v1.30.0"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ServiceAccountV1Resource{}
var _ resource.ResourceWithImportState = &ServiceAccountV1Resource{}
var _ resource.ResourceWithConfigValidators = &ServiceAccountV1Resource{}

func NewServiceAccountV1Resource() resource.Resource {
	return &ServiceAccountV1Resource{}
}

// ServiceAccountV1Resource defines the resource implementation.
type ServiceAccountV1Resource struct {
	apiClient *client.Client
}

func (r *ServiceAccountV1Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_service_account_v1"
}

func (r *ServiceAccountV1Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.ConsoleServiceAccountV1ResourceSchema(ctx)
}

func (r *ServiceAccountV1Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	if semver.IsValid(consoleVersion) && semver.Compare(consoleVersion, consoleServiceAccountMininumVersion) < 0 {
		resp.Diagnostics.AddError(
			"Minimum version requirement not met",
			"This resource requires Conduktor Console API version "+consoleServiceAccountMininumVersion+" but targeted Conduktor Console API is "+consoleVersion,
		)
		return
	}

	r.apiClient = data.Client
}

func (r *ServiceAccountV1Resource) ConfigValidators(_ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("spec").AtName("authorization").AtName("aiven"),
			path.MatchRoot("spec").AtName("authorization").AtName("kafka"),
		),
	}
}

func (r *ServiceAccountV1Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.ConsoleServiceAccountV1Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating service account named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create service account with desired state : %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create service account, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("service account to create : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, serviceAccountV1ApiPutPath(consoleResource.Metadata.Cluster), consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service account, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Service account created with result: %s", apply.UpsertResult))

	var consoleRes = console.ServiceAccountResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as service account : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New service account state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read service account, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceAccountV1Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.ConsoleServiceAccountV1Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read service account named %s", data.Name.String()))
	get, err := r.apiClient.Describe(ctx, serviceAccountV1ApiGetPath(data.Cluster.ValueString(), data.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service account, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Service account %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var consoleRes = console.ServiceAccountResource{}
	err = jsoniter.Unmarshal(get, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Parsing Error", fmt.Sprintf("Unable to read service account, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New service account state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read service account, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceAccountV1Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.ConsoleServiceAccountV1Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating service account named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update service account with TF data: %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create service account, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Service account to update : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, serviceAccountV1ApiPutPath(consoleResource.Metadata.Cluster), consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service account, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Service account updated with result: %s", apply))

	var consoleRes = console.ServiceAccountResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as service account : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New service account state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read service account, got error: %s", err))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceAccountV1Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.ConsoleServiceAccountV1Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Deleting service account named %s", data.Name.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	resourcePath := serviceAccountV1ApiGetPath(data.Cluster.ValueString(), data.Name.ValueString())
	err := r.apiClient.Delete(ctx, client.CONSOLE, resourcePath, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service account, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Service account %s deleted", data.Name.String()))
}

func (r *ServiceAccountV1Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
