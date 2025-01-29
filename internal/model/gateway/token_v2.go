package gateway

import (
	"encoding/json"

	model "github.com/conduktor/terraform-provider-conduktor/internal/model"

	ctlresource "github.com/conduktor/ctl/resource"
	jsoniter "github.com/json-iterator/go"
)

type GatewayTokenResource struct {
	VCluster        string `json:"vCluster"`
	Username        string `json:"username"`
	LifetimeSeconds int64  `json:"lifeTimeSeconds"`
	Token           string `json:"token"`
}

func NewGatewayTokenResource(vCluster string, username string, lifetimeSeconds int64) GatewayTokenResource {
	return GatewayTokenResource{
		VCluster:        vCluster,
		Username:        username,
		LifetimeSeconds: lifetimeSeconds,
	}
}

func (r *GatewayTokenResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *GatewayTokenResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *GatewayTokenResource) FromRawJson(jsonData []byte) error {
	err := jsoniter.Unmarshal(jsonData, &r)
	if err != nil {
		return err
	}
	return nil
}

func (r *GatewayTokenResource) FromRawJsonInterface(jsonInterface interface{}) error {
	jsonData, err := json.Marshal(jsonInterface)
	if err != nil {
		return err
	}
	err = r.FromRawJson(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func NewGatewayTokenResourceFromClientResource(cliResource ctlresource.Resource) (GatewayTokenResource, error) {
	var gatewaynresource GatewayTokenResource
	err := gatewaynresource.FromClientResource(cliResource)
	if err != nil {
		return GatewayTokenResource{}, err
	}
	return gatewaynresource, nil
}
