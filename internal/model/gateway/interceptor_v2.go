package gateway

import (
	"encoding/json"
	"fmt"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"

	ctlresource "github.com/conduktor/ctl/resource"
	jsoniter "github.com/json-iterator/go"
)

const GatewayInterceptorV2Kind = "GatewayInterceptor"
const GatewayInterceptorV2ApiVersion = "gateway/v2"

type GatewayInterceptorScope struct {
	Group    string `json:"group,omitempty"`
	VCluster string `json:"vCluster,omitempty"`
	Username string `json:"username,omitempty"`
}

type GatewayInterceptorMetadata struct {
	Name  string                  `json:"name"`
	Scope GatewayInterceptorScope `json:"scope"`
}

func (r GatewayInterceptorScope) String() string {
	return fmt.Sprintf(`group: %s, vCluster: %s, username: %s`, r.Group, r.VCluster, r.Username)
}

func (r GatewayInterceptorMetadata) String() string {
	return fmt.Sprintf(`name: %s, Scope: %s`, r.Name, r.Scope)
}

type GatewayInterceptorSpec struct {
	Comment     string       `json:"comment,omitempty"`
	PluginClass string       `json:"pluginClass"`
	Priority    int64        `json:"priority"` // API accepts int32 but terraform doesn't support that.
	Config      *interface{} `json:"config"`
}

type GatewayInterceptorResource struct {
	Kind       string                      `json:"kind"`
	ApiVersion string                      `json:"apiVersion"`
	Metadata   *GatewayInterceptorMetadata `json:"metadata"`
	Spec       *GatewayInterceptorSpec     `json:"spec"`
}

func NewGatewayInterceptorResource(metadata GatewayInterceptorMetadata, spec GatewayInterceptorSpec) GatewayInterceptorResource {
	return GatewayInterceptorResource{
		Kind:       GatewayInterceptorV2Kind,
		ApiVersion: GatewayInterceptorV2ApiVersion,
		Metadata:   &metadata,
		Spec:       &spec,
	}
}

func (r *GatewayInterceptorResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *GatewayInterceptorResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *GatewayInterceptorResource) FromRawJsonInterface(jsonInterface interface{}) error {
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

func NewGatewayInterceptorResourceFromClientResource(cliResource ctlresource.Resource) (GatewayInterceptorResource, error) {
	var gatewaynresource GatewayInterceptorResource
	err := gatewaynresource.FromClientResource(cliResource)
	if err != nil {
		return GatewayInterceptorResource{}, err
	}
	return gatewaynresource, nil
}
