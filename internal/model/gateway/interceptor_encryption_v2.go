package gateway

import (
	"encoding/json"
	"fmt"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"

	ctlresource "github.com/conduktor/ctl/resource"
	jsoniter "github.com/json-iterator/go"
)

const GatewayInterceptorV2Kind = "GatewayInterceptor"
const GatewayInterceptorEncryptionV2ApiVersion = "gateway/v2"

type GatewayInterceptorEncryptionScope struct {
	Group    string `json:"group,omitempty"`
	VCluster string `json:"vCluster,omitempty"`
	Username string `json:"username,omitempty"`
}

type GatewayInterceptorEncryptionMetadata struct {
	Name  string                            `json:"name"`
	Scope GatewayInterceptorEncryptionScope `json:"scope"`
}

func (r GatewayInterceptorEncryptionScope) String() string {
	return fmt.Sprintf(`group: %s, vCluster: %s, username: %s`, r.Group, r.VCluster, r.Username)
}

func (r GatewayInterceptorEncryptionMetadata) String() string {
	return fmt.Sprintf(`name: %s, Scope: %s`, r.Name, r.Scope)
}

type GatewayInterceptorEncryptionRetryPolicyConfig struct {
	MaxRetries int64 `json:"maxRetries"`
	DelayMs    int64 `json:"delayMs"`
	MaxDelayMs int64 `json:"maxDelayMs"`
}

type GatewayInterceptorEncryptionAzureTokenCredential struct {
	ClientId     string `json:"clientId"`
	TenantId     string `json:"tenantId"`
	ClientSecret string `json:"clientSecret"`
}

type GatewayInterceptorEncryptionAzureUsernamePasswordCredential struct {
	ClientId string `json:"clientId"`
	TenantId string `json:"tenantId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type GatewayInterceptorEncryptionBasicAWSCredentialsConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type GatewayInterceptorEncryptionBasicSessionCredentialsConfig struct {
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	SessionToken string `json:"sessionToken"`
}

type GatewayInterceptorEncryptionAzureKMSConfig struct {
	RetryPolicy                *GatewayInterceptorEncryptionRetryPolicyConfig               `json:"retryPolicy"`
	TokenCredential            *GatewayInterceptorEncryptionAzureTokenCredential            `json:"tokenCredential"`
	UsernamePasswordCredential *GatewayInterceptorEncryptionAzureUsernamePasswordCredential `json:"usernamePasswordCredential"`
}

type GatewayInterceptorEncryptionAWSKMSConfig struct {
	BasicCredentials   GatewayInterceptorEncryptionBasicAWSCredentialsConfig     `json:"basicCredentials"`
	SessionCredentials GatewayInterceptorEncryptionBasicSessionCredentialsConfig `json:"sessionCredentials"`
}

type GatewayInterceptorEncryptionGCPKMSConfig struct {
	ServiceAccountCredentialsFilePath string `json:"serviceAccountCredentialsFilePath"`
}

type GatewayInterceptorEncryptionKMSConfig struct {
	KeyTtlMs int64                                      `json:"keyTtlMs"`
	Azure    GatewayInterceptorEncryptionAzureKMSConfig `json:"azure"`
	AWS      GatewayInterceptorEncryptionAWSKMSConfig   `json:"aws"`
	GCP      GatewayInterceptorEncryptionGCPKMSConfig   `json:"gcp"`
}

type GatewayInterceptorEncryptionSchemaRegistryConfig struct {
	Host              string            `json:"host"`
	CacheSize         int64             `json:"cacheSize"`
	AdditionalConfigs map[string]string `json:"additionalConfigs"`
}

type GatewayInterceptorEncryptionConfig struct {
	EnableAuditLogOnError bool                                              `json:"enableAuditLogOnError"`
	ExternalStorage       bool                                              `json:"externalStorage"`
	SchemaDataMode        string                                            `json:"schemaDataMode"`
	SchemaRegistryConfig  *GatewayInterceptorEncryptionSchemaRegistryConfig `json:"schemaRegistryConfig"`
	KmsConfig             *GatewayInterceptorEncryptionKMSConfig            `json:"kmsConfig"`
	Topic                 string                                            `json:"topic"`
}

type GatewayInterceptorEncryptionSpec struct {
	Comment     string                              `json:"comment,omitempty"`
	PluginClass string                              `json:"pluginClass"`
	Priority    int64                               `json:"priority"` // API accepts int32 but terraform doesn't support that.
	Config      *GatewayInterceptorEncryptionConfig `json:"config"`
}

type GatewayInterceptorEncryptionResource struct {
	Kind       string                                `json:"kind"`
	ApiVersion string                                `json:"apiVersion"`
	Metadata   *GatewayInterceptorEncryptionMetadata `json:"metadata"`
	Spec       *GatewayInterceptorEncryptionSpec     `json:"spec"`
}

func NewGatewayInterceptorEncryptionResource(metadata GatewayInterceptorEncryptionMetadata, spec GatewayInterceptorEncryptionSpec) GatewayInterceptorEncryptionResource {
	return GatewayInterceptorEncryptionResource{
		Kind:       GatewayInterceptorV2Kind,
		ApiVersion: GatewayInterceptorEncryptionV2ApiVersion,
		Metadata:   &metadata,
		Spec:       &spec,
	}
}

func (r *GatewayInterceptorEncryptionResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *GatewayInterceptorEncryptionResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *GatewayInterceptorEncryptionResource) FromRawJsonInterface(jsonInterface interface{}) error {
	jsonData, err := json.Marshal(jsonInterface)
	if err != nil {
		return err
	}
	err = jsoniter.Unmarshal(jsonData, r)
	if err != nil {
		return err
	}
	return nil
}

func NewGatewayInterceptorEncryptionResourceFromClientResource(cliResource ctlresource.Resource) (GatewayInterceptorEncryptionResource, error) {
	var gatewaynresource GatewayInterceptorEncryptionResource
	err := gatewaynresource.FromClientResource(cliResource)
	if err != nil {
		return GatewayInterceptorEncryptionResource{}, err
	}
	return gatewaynresource, nil
}
