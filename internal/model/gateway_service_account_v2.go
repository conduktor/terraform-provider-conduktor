package model

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	jsoniter "github.com/json-iterator/go"
)

const GatewayServiceAccountV2Kind = "GatewayServiceAccount"
const GatewayServiceAccountV2ApiVersion = "v2"

type GatewayServiceAccountMetadata struct {
	Name string `json:"name"`
}

func (r GatewayServiceAccountMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
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

func NewGatewayServiceAccountResource(name string, spec GatewayServiceAccountSpec) GatewayServiceAccountResource {
	return GatewayServiceAccountResource{
		Kind:       GatewayServiceAccountV2Kind,
		ApiVersion: GatewayServiceAccountV2ApiVersion,
		Metadata: GatewayServiceAccountMetadata{
			Name: name,
		},
		Spec: spec,
	}
}

func (r *GatewayServiceAccountResource) ToClientResource() (ctlresource.Resource, error) {
	return toClientResource(r)
}

func (r *GatewayServiceAccountResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *GatewayServiceAccountResource) FromRawJsonInterface(jsonInterface interface{}) error {
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
