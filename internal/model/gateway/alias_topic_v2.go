package gateway

import (
	"encoding/json"
	"fmt"

	model "github.com/conduktor/terraform-provider-conduktor/internal/model"

	ctlresource "github.com/conduktor/ctl/resource"
	jsoniter "github.com/json-iterator/go"
)

const GatewayAliasTopicV2Kind = "AliasTopic"
const GatewayAliasTopicV2ApiVersion = "gateway/v2"

type GatewayAliasTopicMetadata struct {
	Name     string `json:"name"`
	VCluster string `json:"vCluster,omitempty"`
}

func (r GatewayAliasTopicMetadata) String() string {
	return fmt.Sprintf(`name: %s, vCluster: %s`, r.Name, r.VCluster)
}

type GatewayAliasTopicSpec struct {
	PhysicalName    string `json:"physicalName"`
	PhysicalCluster string `json:"physicalCluster,omitempty"`
}

type GatewayAliasTopicResource struct {
	Kind       string                    `json:"kind"`
	ApiVersion string                    `json:"apiVersion"`
	Metadata   GatewayAliasTopicMetadata `json:"metadata"`
	Spec       GatewayAliasTopicSpec     `json:"spec"`
}

func NewGatewayAliasTopicResource(metadata GatewayAliasTopicMetadata, spec GatewayAliasTopicSpec) GatewayAliasTopicResource {
	return GatewayAliasTopicResource{
		Kind:       GatewayAliasTopicV2Kind,
		ApiVersion: GatewayAliasTopicV2ApiVersion,
		Metadata:   metadata,
		Spec:       spec,
	}
}

func (r *GatewayAliasTopicResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *GatewayAliasTopicResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *GatewayAliasTopicResource) FromRawJsonInterface(jsonInterface any) error {
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

func NewGatewayAliasTopicResourceFromClientResource(cliResource ctlresource.Resource) (GatewayAliasTopicResource, error) {
	var gatewayResource GatewayAliasTopicResource
	err := gatewayResource.FromClientResource(cliResource)
	if err != nil {
		return GatewayAliasTopicResource{}, err
	}
	return gatewayResource, nil
}
