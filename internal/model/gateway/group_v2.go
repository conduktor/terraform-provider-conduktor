package gateway

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const GatewayGroupV2Kind = "GatewayGroup"
const GatewayGroupV2ApiVersion = "gateway/v2"

type GatewayGroupMetadata struct {
	Name string `json:"name"`
}

func (r GatewayGroupMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

// GatewayServiceAccountId identifies a service account member of a group.
// VCluster is optional — omit for passthrough.
type GatewayServiceAccountId struct {
	Name     string `json:"name"`
	VCluster string `json:"vCluster,omitempty"`
}

type GatewayGroupSpec struct {
	Members        []GatewayServiceAccountId `json:"members"`
	ExternalGroups []string                  `json:"externalGroups,omitempty"`
}

type GatewayGroupResource struct {
	Kind       string               `json:"kind"`
	ApiVersion string               `json:"apiVersion"`
	Metadata   GatewayGroupMetadata `json:"metadata"`
	Spec       GatewayGroupSpec     `json:"spec"`
}

func NewGatewayGroupResource(metadata GatewayGroupMetadata, spec GatewayGroupSpec) GatewayGroupResource {
	return GatewayGroupResource{
		Kind:       GatewayGroupV2Kind,
		ApiVersion: GatewayGroupV2ApiVersion,
		Metadata:   metadata,
		Spec:       spec,
	}
}

func (r *GatewayGroupResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *GatewayGroupResource) FromClientResource(cliResource ctlresource.Resource) error {
	return jsoniter.Unmarshal(cliResource.Json, r)
}

func (r *GatewayGroupResource) FromRawJsonInterface(jsonInterface any) error {
	jsonData, err := json.Marshal(jsonInterface)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(jsonData, r)
}

func NewGatewayGroupResourceFromClientResource(cliResource ctlresource.Resource) (GatewayGroupResource, error) {
	var res GatewayGroupResource
	if err := res.FromClientResource(cliResource); err != nil {
		return GatewayGroupResource{}, err
	}
	return res, nil
}
