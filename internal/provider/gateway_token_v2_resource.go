package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/conduktor/terraform-provider-conduktor/internal/client"
	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper/gateway_token_v2"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_gateway_token_v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
)

const gatewayTokenV2ApiPath = "/gateway/v2/token"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GatewayTokenV2Resource{}

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

func (r *GatewayTokenV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schema.GatewayTokenV2Model

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read token for service account %s", data.Username.String()))

	if !data.Token.IsNull() && data.Token.ValueString() != "" {
		expired, err := isTokenExpired(data.Token.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to validate token, got error: %s", err))
			return
		}

		if expired {
			tflog.Trace(ctx, "Token expired, removing resource from state")
			resp.State.RemoveResource(ctx)
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayTokenV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schema.GatewayTokenV2Model

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

	apply, err := applyGatewayToken(ctx, r.apiClient, gatewayTokenV2ApiPath, gatewayResource)
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

	err = gatewayRes.FromRawJsonInterface(apply.Resource)
	if err != nil {
		resp.Diagnostics.AddError("Unmarshall Error", fmt.Sprintf("Response resource can't be cast as token : %v, got error: %s", apply.Resource, err))
		return
	}
	gatewayRes.VCluster = data.Vcluster.ValueString()
	gatewayRes.Username = data.Username.ValueString()
	gatewayRes.LifetimeSeconds = data.LifetimeSeconds.ValueInt64()
	tflog.Debug(ctx, fmt.Sprintf("New token state : %+v", gatewayRes))

	data, err = mapper.InternalModelToTerraform(ctx, &gatewayRes)
	if err != nil {
		resp.Diagnostics.AddError("Model Error", fmt.Sprintf("Unable to read token, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayTokenV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schema.GatewayTokenV2Model

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

	apply, err := applyGatewayToken(ctx, r.apiClient, gatewayTokenV2ApiPath, gatewayResource)
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
	gatewayRes.VCluster = data.Vcluster.ValueString()
	gatewayRes.Username = data.Username.ValueString()
	gatewayRes.LifetimeSeconds = data.LifetimeSeconds.ValueInt64()
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

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Info(ctx, fmt.Sprintf("Delete token for service account: %s", data.Username.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Token deleted for service account: %s", data.Username.String()))
}

// Function to check if a JWT token is expired.
func isTokenExpired(tokenString string) (bool, error) {
	// Parse the token
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			if time.Now().After(expirationTime) {
				return true, nil
			}
		} else {
			return false, fmt.Errorf("Expiration time not found in token")
		}
	} else {
		return false, fmt.Errorf("Unable to parse claims")
	}

	return false, nil
}

// Helper function to issue a new token.
func applyGatewayToken(ctx context.Context, cli *client.Client, path string, resource any) (client.ApplyResult, error) {
	url := cli.BaseUrl + path
	jsonData, err := jsoniter.Marshal(resource)
	if err != nil {
		return client.ApplyResult{}, fmt.Errorf("Error marshalling resource: %s", err)
	}

	tflog.Trace(ctx, fmt.Sprintf("POST %s request body : %s", path, string(jsonData)))

	resp, err := cli.Client.R().SetBody(jsonData).Post(url)
	if err != nil {
		return client.ApplyResult{}, err
	} else if resp.IsError() {
		return client.ApplyResult{}, fmt.Errorf("%s", client.ExtractApiError(resp))
	}

	bodyBytes := resp.Body()
	tflog.Trace(ctx, fmt.Sprintf("POST %s response body : %s", path, string(bodyBytes)))

	var upsertResponse gateway.GatewayTokenResource
	err = jsoniter.Unmarshal(bodyBytes, &upsertResponse)
	if err != nil {
		return client.ApplyResult{}, fmt.Errorf("Error unmarshalling response: %s", err)
	}
	return client.ApplyResult{Resource: upsertResponse}, nil
}
