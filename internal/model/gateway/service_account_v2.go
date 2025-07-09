package gateway

import (
	"encoding/json"
	"fmt"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"

	ctlresource "github.com/conduktor/ctl/resource"
	jsoniter "github.com/json-iterator/go"
)

const GatewayServiceAccountV2Kind = "GatewayServiceAccount"
const GatewayServiceAccountV2ApiVersion = "gateway/v2"

type GatewayServiceAccountMetadata struct {
	Name     string `json:"name"`
	VCluster string `json:"vCluster,omitempty"`
}

func (r GatewayServiceAccountMetadata) String() string {
	return fmt.Sprintf(`name: %s, vCluster: %s`, r.Name, r.VCluster)
}

type GatewayServiceAccountSpec struct {
	Type          string   `json:"type"`
	ExternalNames []string `json:"externalNames"`
}

type GatewayServiceAccountResource struct {
	Kind       string                        `json:"kind"`
	ApiVersion string                        `json:"apiVersion"`
	Metadata   GatewayServiceAccountMetadata `json:"metadata"`
	Spec       GatewayServiceAccountSpec     `json:"spec"`
}

func NewGatewayServiceAccountResource(metadata GatewayServiceAccountMetadata, spec GatewayServiceAccountSpec) GatewayServiceAccountResource {
	return GatewayServiceAccountResource{
		Kind:       GatewayServiceAccountV2Kind,
		ApiVersion: GatewayServiceAccountV2ApiVersion,
		Metadata:   metadata,
		Spec:       spec,
	}
}

func (r *GatewayServiceAccountResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *GatewayServiceAccountResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *GatewayServiceAccountResource) FromRawJsonInterface(jsonInterface any) error {
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

func NewGatewayServiceAccountResourceFromClientResource(cliResource ctlresource.Resource) (GatewayServiceAccountResource, error) {
	var gatewaynresource GatewayServiceAccountResource
	err := gatewaynresource.FromClientResource(cliResource)
	if err != nil {
		return GatewayServiceAccountResource{}, err
	}
	return gatewaynresource, nil
}
