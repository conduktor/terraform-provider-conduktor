package provider

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/conduktor/terraform-provider-conduktor/internal/customtypes"

	ctlresource "github.com/conduktor/ctl/resource"
	ctlschema "github.com/conduktor/ctl/schema"
	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	schemaUtils "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_generic"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GenericResource{}
var _ resource.ResourceWithImportState = &GenericResource{}

func NewGenericResource() resource.Resource {
	return &GenericResource{}
}

// GenericResource defines the resource implementation.
type GenericResource struct {
	client *client.Client
}

func (r *GenericResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_generic"
}

func (r *GenericResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.GenericResourceSchema(ctx)
}

func (r *GenericResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = data.Client
}

func (r *GenericResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.GenericModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Create %s kind named %s", data.Kind.String(), data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Create resource with TF data: %+v", data))

	cliResource, err := ctlresource.FromYamlByte([]byte(data.Manifest.ValueString()), true)

	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create Generic, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Resource to create : %+v", cliResource))

	firstResource := cliResource[0]
	apply, err := r.client.ApplyGeneric(ctx, firstResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Generic, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Resource created with result: %s", apply))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GenericResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.GenericModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read %s kind named %s", data.Kind.String(), data.Name.String()))
	resourcePath, err := resourcePath(data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to build Generic api path, got error: \"%s\" from kind:%s name:%s (cluster:%s)", err, data.Kind.ValueString(), data.Name.ValueString(), data.Cluster.ValueString()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Query resource on path %s", resourcePath))
	get, err := r.client.Describe(ctx, resourcePath)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Generic, got error: %s", err))
		return
	}

	if len(get) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("Resource %s not found, removing from state", resourcePath))
		resp.State.RemoveResource(ctx)
		return
	}

	cliResource, err := ctlresource.FromYamlByte(get, true)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read Generic resource, got error: %s", err))
		return
	}
	if len(cliResource) != 1 {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Received more than one resource on response : %v", cliResource))
		return
	}

	firstResource := cliResource[0]

	tflog.Trace(ctx, fmt.Sprintf("New resource JSON state : %s", string(firstResource.Json)))

	outBytes, err := yaml.JSONToYAML(firstResource.Json)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Generic, got error: %s", err))
		return
	}
	yamlString := string(outBytes)

	// Keep only user-declared fields so server-computed ones (promQl, updatedAt...)
	// don't cause a perpetual diff. Best-effort: fall back to the raw response.
	if prior := data.Manifest.ValueString(); strings.TrimSpace(prior) != "" {
		if reconciled, rerr := reconcileManifest(prior, yamlString); rerr != nil {
			tflog.Warn(ctx, fmt.Sprintf("Unable to reconcile Generic manifest, using raw response: %s", rerr))
		} else {
			yamlString = reconciled
		}
	}
	tflog.Trace(ctx, fmt.Sprintf("New resource YAML state : %s", yamlString))

	data.Kind = schemaUtils.NewStringValue(firstResource.Kind)
	data.Name = schemaUtils.NewStringValue(firstResource.Name)
	data.Version = schemaUtils.NewStringValue(firstResource.Version)
	data.Manifest = customtypes.NewNormalizedValue(yamlString)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GenericResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.GenericModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Update %s kind named %s", data.Kind.String(), data.Name.String()))
	tflog.Trace(ctx, fmt.Sprintf("Update resource with TF data: %+v", data))

	cliResource, err := ctlresource.FromYamlByte([]byte(data.Manifest.ValueString()), true)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to create Generic, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Resource to update : %+v", cliResource))

	firstResource := cliResource[0]
	apply, err := r.client.ApplyGeneric(ctx, firstResource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Generic, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Resource updated with result: %s", apply))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GenericResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schema.GenericModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Delete %s kind named %s", data.Kind.String(), data.Name.String()))
	resourcePath, err := resourcePath(data)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read Generic, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Delete resource on path %s", resourcePath))

	err = r.client.Delete(ctx, client.CONSOLE, resourcePath, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Generic, got error: %s", err))
		return
	}
}

func (r *GenericResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// Search for the kind in the CLI default schema.
func getKindFromName(kindName string) (ctlschema.Kind, error) {
	catalog := ctlschema.ConsoleDefaultCatalog() // TODO support gateway kinds and client too
	kinds := catalog.Kind
	kind, ok := kinds[kindName]
	if !ok {
		return ctlschema.Kind{}, fmt.Errorf("kind %s not found", kindName)
	}

	return kind, nil
}

// Generate the resource path for the given kind, cluster and resource name.
func resourcePath(data schema.GenericModel) (string, error) {
	kind, err := getKindFromName(data.Kind.ValueString())
	if err != nil {
		return "", err
	}

	parentPath := []string{}
	cluster := data.Cluster.ValueString()
	if cluster != "" {
		parentPath = append(parentPath, cluster)
	}

	// Kinds scoped by metadata query params instead of a path segment, e.g. Alert
	// v3's appInstance/group/user (#186). Sized-but-empty values keep DescribePath
	// happy; the resolved params are appended to the final path.
	queryParams, err := parentQueryParams(kind, data.Manifest.ValueString())
	if err != nil {
		return "", err
	}

	describe := kind.DescribePath(parentPath, make([]string, len(kind.GetParentQueryFlag())), data.Name.ValueString())
	return appendQueryParams(describe.Path, queryParams), nil
}

// parentQueryParams extracts the kind's parent query params from the manifest
// metadata, reusing ctl's Kind.ApplyPath (the same extraction Create uses).
func parentQueryParams(kind ctlschema.Kind, manifest string) ([]ctlschema.QueryParam, error) {
	if len(kind.GetParentQueryFlag()) == 0 || strings.TrimSpace(manifest) == "" {
		return nil, nil
	}

	cliResources, err := ctlresource.FromYamlByte([]byte(manifest), true)
	if err != nil {
		return nil, err
	}
	if len(cliResources) == 0 {
		return nil, nil
	}

	applyInfo, err := kind.ApplyPath(&cliResources[0])
	if err != nil {
		return nil, err
	}
	return applyInfo.QueryParams, nil
}

// appendQueryParams appends the given query parameters to the path, URL-encoded.
func appendQueryParams(path string, params []ctlschema.QueryParam) string {
	if len(params) == 0 {
		return path
	}
	values := url.Values{}
	for _, param := range params {
		values.Add(param.Name, param.Value)
	}
	return path + "?" + values.Encode()
}

// reconcileManifest overlays the API manifest onto the prior one, returning YAML
// that keeps only the fields the user declared. Both inputs are YAML.
func reconcileManifest(priorManifest, apiManifest string) (string, error) {
	var prior, api map[string]any
	if err := yaml.Unmarshal([]byte(priorManifest), &prior); err != nil {
		return "", err
	}
	if err := yaml.Unmarshal([]byte(apiManifest), &api); err != nil {
		return "", err
	}

	out, err := yaml.Marshal(overlay(prior, api))
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// overlay keeps only the keys present in prior, taking values from api (recursing
// into nested maps to still detect drift on declared fields). Keys missing from api
// keep their prior value; api-only keys (server-computed) are dropped.
func overlay(prior, api map[string]any) map[string]any {
	result := make(map[string]any, len(prior))
	for key, priorVal := range prior {
		apiVal, ok := api[key]
		if !ok {
			result[key] = priorVal
			continue
		}
		priorMap, priorIsMap := priorVal.(map[string]any)
		apiMap, apiIsMap := apiVal.(map[string]any)
		if priorIsMap && apiIsMap {
			result[key] = overlay(priorMap, apiMap)
		} else {
			result[key] = apiVal
		}
	}
	return result
}
