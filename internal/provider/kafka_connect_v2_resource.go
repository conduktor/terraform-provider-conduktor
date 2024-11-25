package provider

import (
	"context"
	"fmt"
	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/kafka_connect_v2"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_kafka_connect_v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

func kafkaConnectV2ApiPutPath(cluster string) string {
	return fmt.Sprintf("/public/console/v2/cluster/%s/kafka-connect", cluster)
}

func kafkaConnectV2ApiGetPath(cluster string, connectServerName string) string {
	return fmt.Sprintf("/public/console/v2/cluster/%s/kafka-connect/%s", cluster, connectServerName)
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &KafkaConnectV2Resource{}
var _ resource.ResourceWithImportState = &KafkaConnectV2Resource{}

func NewKafkaConnectV2Resource() resource.Resource {
	return &KafkaConnectV2Resource{}
}

// KafkaConnectV2Resource defines the resource implementation.
type KafkaConnectV2Resource struct {
	apiClient *client.ConsoleClient
}

func (r *KafkaConnectV2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kafka_connect_v2"
}

func (r *KafkaConnectV2Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.KafkaConnectV2ResourceSchema(ctx)
}

func (r *KafkaConnectV2Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	if data.ConsoleClient == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Console Client not configured. Please provide client configuration details.",
		)
		return
	}

	r.apiClient = data.ConsoleClient
}

func (r *KafkaConnectV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.KafkaConnectV2Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating kafka connect server named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create kafka connect server with desired state : %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create kafka connect server, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Kafka connect server to create : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, kafkaConnectV2ApiPutPath(consoleResource.Metadata.Cluster), consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create kafka connect server, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Kafka connect server created with result: %s", apply.UpsertResult))

	var consoleRes = model.KafkaConnectResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as kafka connect server : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New kafka connect server state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read kafka connect server, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KafkaConnectV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.KafkaConnectV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read kafka connect server named %s", data.Name.String()))
	get, err := r.apiClient.Describe(ctx, kafkaConnectV2ApiGetPath(data.Cluster.ValueString(), data.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read kafka connect server, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Kafka connect server %s not found, removing from state", data.Name.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	var consoleRes = model.KafkaConnectResource{}
	err = jsoniter.Unmarshal(get, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Parsing Error", fmt.Sprintf("Unable to read kafka connect server, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New kafka connect server state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read kafka connect server, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KafkaConnectV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.KafkaConnectV2Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating kafka connect server named %s", data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update kafka connect server with TF data: %+v", data))

	consoleResource, err := mapper.TFToInternalModel(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create kafka connect server, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Kafka connect server to update : %+v", consoleResource))

	apply, err := r.apiClient.Apply(ctx, kafkaConnectV2ApiPutPath(consoleResource.Metadata.Cluster), consoleResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create kafka connect server, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Kafka connect server updated with result: %s", apply))

	var consoleRes = model.KafkaConnectResource{}
	err = consoleRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as kafka connect server : %v, got error: %s", apply.Resource, err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("New kafka connect server state : %+v", consoleRes))

	data, err = mapper.InternalModelToTerraform(ctx, &consoleRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read kafka connect server, got error: %s", err))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KafkaConnectV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.KafkaConnectV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Deleting kafka connect server named %s", data.Name.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.apiClient.Delete(ctx, kafkaConnectV2ApiGetPath(data.Cluster.ValueString(), data.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete kafka connect server, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Kafka connect server %s deleted", data.Name.String()))
}

func (r *KafkaConnectV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
