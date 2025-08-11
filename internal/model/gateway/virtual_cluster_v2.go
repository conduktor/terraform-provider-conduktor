package gateway

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const VirtualClusterV2Kind = "VirtualCluster"
const VirtualClusterV2ApiVersion = "v2"

type VirtualClusterMetadata struct {
	Name string `json:"name"`
}

func (r VirtualClusterMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type VirtualClusterSpec struct {
	AclEnabled       bool                `json:"aclEnabled,omitempty"`
	AclMode          string              `json:"aclMode,omitempty"`
	SuperUsers       []string            `json:"superUsers,omitempty"`
	Type             string              `json:"type,omitempty"`
	BootstrapServers string              `json:"bootstrapServers,omitempty"`
	ClientProperties map[string]string   `json:"clientProperties,omitempty"`
	Acls             []VirtualClusterACL `json:"acls,omitempty"`
}

type VirtualClusterACL struct {
	ResourcePattern VirtualClusterACLResourcePattern `json:"resourcePattern"`
	Principal       string                           `json:"principal"`
	Host            string                           `json:"host"`
	Operation       string                           `json:"operation"`
	PermissionType  string                           `json:"permissionType"`
}

type VirtualClusterACLResourcePattern struct {
	ResourceType string `json:"resourceType"`
	Name         string `json:"name"`
	PatternType  string `json:"patternType"`
}

type VirtualClusterResource struct {
	ApiVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   VirtualClusterMetadata `json:"metadata"`
	Spec       VirtualClusterSpec     `json:"spec"`
}

func NewVirtualClusterResource(meta VirtualClusterMetadata, spec VirtualClusterSpec) VirtualClusterResource {
	return VirtualClusterResource{
		ApiVersion: VirtualClusterV2ApiVersion,
		Kind:       VirtualClusterV2Kind,
		Metadata:   meta,
		Spec:       spec,
	}
}

func (r *VirtualClusterResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *VirtualClusterResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *VirtualClusterResource) FromRawJsonInterface(jsonInterface any) error {
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

func NewVirtualClusterResourceFromClientResource(cliResource ctlresource.Resource) (VirtualClusterResource, error) {
	var consoleResource VirtualClusterResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return VirtualClusterResource{}, err
	}
	return consoleResource, nil
}
