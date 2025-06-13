package console

import (
	"encoding/json"
	"fmt"
	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const ResourcePolicyConsoleV1Kind = "ResourcePolicy"
const ResourcePolicyConsoleV1ApiVersion = "v1"

type ResourcePolicyConsoleMetadata struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
}

func (r ResourcePolicyConsoleMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type ResourcePolicyConsoleRule struct {
	Condition    string `json:"condition"`
	ErrorMessage string `json:"errorMessage"`
}

type ResourcePolicyConsoleSpec struct {
	TargetKind  string                      `json:"targetKind"`
	Description string                      `json:"description,omitempty"`
	Rules       []ResourcePolicyConsoleRule `json:"rules"`
}

type ResourcePolicyConsoleResource struct {
	ApiVersion string                        `json:"apiVersion"`
	Kind       string                        `json:"kind"`
	Metadata   ResourcePolicyConsoleMetadata `json:"metadata"`
	Spec       ResourcePolicyConsoleSpec     `json:"spec"`
}

func NewResourcePolicyConsoleResource(meta ResourcePolicyConsoleMetadata, spec ResourcePolicyConsoleSpec) ResourcePolicyConsoleResource {
	return ResourcePolicyConsoleResource{
		ApiVersion: ResourcePolicyConsoleV1ApiVersion,
		Kind:       ResourcePolicyConsoleV1Kind,
		Metadata:   meta,
		Spec:       spec,
	}
}

func (r *ResourcePolicyConsoleResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *ResourcePolicyConsoleResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *ResourcePolicyConsoleResource) FromRawJsonInterface(jsonInterface any) error {
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

func NewResourcePolicyConsoleResourceFromClientResource(cliResource ctlresource.Resource) (ResourcePolicyConsoleResource, error) {
	var consoleResource ResourcePolicyConsoleResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return ResourcePolicyConsoleResource{}, err
	}
	return consoleResource, nil
}
