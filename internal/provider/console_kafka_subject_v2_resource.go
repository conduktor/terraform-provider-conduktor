package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/console_kafka_subject_v2"
	"github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_kafka_subject_v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/mod/semver"
)

func kafkaSubjectV2ApiPutPath(cluster string) string {
	return fmt.Sprintf("/public/kafka/v2/cluster/%s/subject", cluster)
}

func kafkaSubjectV2ApiGetPath(cluster string, subjectName string) string {
	return fmt.Sprintf("/public/kafka/v2/cluster/%s/subject/%s", cluster, subjectName)
}

// Version when this resource was introduced.
const kafkaSubjectMininumVersion = "v1.29.0"

var _ resource.Resource = &KafkaSubjectV2Resource{}
var _ resource.ResourceWithImportState = &KafkaSubjectV2Resource{}

func NewKafkaSubjectV2Resource() resource.Resource {
	return &KafkaSubjectV2Resource{}
}

type KafkaSubjectV2Resource struct {
	apiClient *client.Client
}

func (r *KafkaSubjectV2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_kafka_subject_v2"
}

func (r *KafkaSubjectV2Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.ConsoleKafkaSubjectV2ResourceSchema(ctx)
}

func (r *KafkaSubjectV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
		resp.Diagnostics.AddError("Error fetching Console version", err.Error())
		return
	}
	if semver.IsValid(consoleVersion) && semver.Compare(consoleVersion, kafkaSubjectMininumVersion) < 0 {
		resp.Diagnostics.AddError(
			"Minimum version requirement not met",
			"This resource requires Conduktor Console API version "+kafkaSubjectMininumVersion+" or higher, but targeted Conduktor Console API is "+consoleVersion,
		)
		return
	}

	r.apiClient = data.Client
}

func (r *KafkaSubjectV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.ConsoleKafkaSubjectV2Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating kafka subject named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create kafka subject with desired state : %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create kafka subject, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Kafka subject to create : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, kafkaSubjectV2ApiPutPath(consoleResource.Metadata.Cluster), consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create kafka subject, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Kafka subject created with result: %s", apply.UpsertResult))

	// Get back the resource to read fields populated by server like ID and Version
	tflog.Debug(ctx, "Get created Kafka subject to update state")
	newState, err := r.pollSubjectState(ctx, data.Cluster.ValueString(), data.Name.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read kafka subject after update, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *KafkaSubjectV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.ConsoleKafkaSubjectV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read kafka subject named %s", data.Name.String()))

	newState, err := r.getSubjectState(ctx, data.Cluster.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read kafka subject, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *KafkaSubjectV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schema.ConsoleKafkaSubjectV2Model
	var state schema.ConsoleKafkaSubjectV2Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check if schema is changed to determine if we need to poll for a new version
	isSchemaSame, diag := state.Spec.Schema.StringSemanticEquals(ctx, plan.Spec.Schema)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating kafka subject named %s", plan.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update kafka subject with TF data: %+v", plan))

	consoleResource, err := mapper.TFToInternalModel(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create kafka subject, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Kafka subject to update : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, kafkaSubjectV2ApiPutPath(consoleResource.Metadata.Cluster), consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create kafka subject, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Kafka subject updated with result: %s", apply))

	// Get back the resource to read fields populated by server like ID and Version
	tflog.Debug(ctx, "Get created Kafka subject to update state")
	previousVersion := state.Spec.Version.ValueInt64Pointer()
	if isSchemaSame {
		// If schema is unchanged no new version will be created skip polling for a new version
		previousVersion = nil
	}
	newState, err := r.pollSubjectState(ctx, plan.Cluster.ValueString(), plan.Name.ValueString(), previousVersion)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read kafka subject after update, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *KafkaSubjectV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.ConsoleKafkaSubjectV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Deleting kafka subject named %s", data.Name.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	resourcePath := kafkaSubjectV2ApiGetPath(data.Cluster.ValueString(), data.Name.ValueString())
	err := r.apiClient.Delete(ctx, client.CONSOLE, resourcePath, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete kafka subject, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Kafka subject %s deleted", data.Name.String()))
}

func (r *KafkaSubjectV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

// Poll until the subject exists with a version greater than previousVersion.
func (r *KafkaSubjectV2Resource) pollSubjectState(ctx context.Context, clusterName, subjectName string, previousVersion *int64) (schema.ConsoleKafkaSubjectV2Model, error) {
	const maxRetries = 20
	sleepTime := 100 * time.Millisecond
	tflog.Debug(ctx, fmt.Sprintf("Polling kafka subject %s for version update every %dms for %d retries", subjectName, sleepTime.Milliseconds(), maxRetries))
	for i := 0; i < maxRetries; i++ {
		newState, err := r.getSubjectState(ctx, clusterName, subjectName)
		if err != nil {
			return schema.ConsoleKafkaSubjectV2Model{}, err
		}

		// If no previous version, return immediately (creation case)
		if previousVersion == nil {
			tflog.Debug(ctx, fmt.Sprintf("Kafka subject %s created", subjectName))
			return newState, nil
		}

		// Check if version has been updated
		if !newState.Spec.Version.IsNull() && !newState.Spec.Version.IsUnknown() {
			currentVersion := *newState.Spec.Version.ValueInt64Pointer()
			if currentVersion > *previousVersion {
				tflog.Debug(ctx, fmt.Sprintf("Kafka subject %s updated, current version: %v, previous version: %v", subjectName, currentVersion, *previousVersion))
				return newState, nil
			}
		}

		tflog.Debug(ctx, fmt.Sprintf("Kafka subject %s not updated yet (attempt %d/%d), current version: %v, previous version: %v, sleep %s", subjectName, i+1, maxRetries, newState.Spec.Version, *previousVersion, sleepTime))
		time.Sleep(sleepTime)
	}

	return schema.ConsoleKafkaSubjectV2Model{}, fmt.Errorf("timeout waiting for kafka subject %s to update after %d retries", subjectName, maxRetries)
}

// Get the current state of the subject from the API and convert it to the Terraform model.
func (r *KafkaSubjectV2Resource) getSubjectState(ctx context.Context, clusterName, subjectName string) (schema.ConsoleKafkaSubjectV2Model, error) {
	get, err := r.apiClient.Describe(ctx, kafkaSubjectV2ApiGetPath(clusterName, subjectName))
	if err != nil {
		return schema.ConsoleKafkaSubjectV2Model{}, fmt.Errorf("unable to read kafka subject after update, got error: %s", err)
	}
	if len(get) == 0 {
		return schema.ConsoleKafkaSubjectV2Model{}, fmt.Errorf("unable to read kafka subject after update, got empty response")
	}

	var consoleRes = console.KafkaSubjectResource{}
	err = jsoniter.Unmarshal(get, &consoleRes)
	if err != nil {
		return schema.ConsoleKafkaSubjectV2Model{}, fmt.Errorf("response resource can't be cast as kafka subject : %v, got error: %s", get, err)
	}
	tflog.Debug(ctx, fmt.Sprintf("New kafka subject state : %+v", consoleRes))

	var newState schema.ConsoleKafkaSubjectV2Model
	newState, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		return schema.ConsoleKafkaSubjectV2Model{}, fmt.Errorf("unable to read kafka subject, got error: %s", err)
	}
	return newState, nil
}
