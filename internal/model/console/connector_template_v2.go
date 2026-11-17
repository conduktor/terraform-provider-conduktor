// internal/model/console/connector_template_v2.go
package console

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const ConnectorTemplateV2Kind = "ConnectorTemplate"
const ConnectorTemplateV2ApiVersion = "v2"

type ConnectorTemplateMetadata struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
}

func (r ConnectorTemplateMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type ConnectorTemplateDefaultsMetadata struct {
	Name   string            `json:"name,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

type ConnectorTemplateDefaultsSpec struct {
	Class  string            `json:"class"`
	Config map[string]string `json:"config"`
}

type ConnectorTemplateDefaults struct {
	Metadata ConnectorTemplateDefaultsMetadata `json:"metadata"`
	Spec     ConnectorTemplateDefaultsSpec     `json:"spec"`
}

type ConnectorTemplateSpec struct {
	DisplayName string                    `json:"displayName"`
	Description string                    `json:"description,omitempty"`
	Defaults    ConnectorTemplateDefaults `json:"defaults"`
}

type ConnectorTemplateResource struct {
	ApiVersion string                    `json:"apiVersion"`
	Kind       string                    `json:"kind"`
	Metadata   ConnectorTemplateMetadata `json:"metadata"`
	Spec       ConnectorTemplateSpec     `json:"spec"`
}

func NewConnectorTemplateResource(name string, labels map[string]string, spec ConnectorTemplateSpec) ConnectorTemplateResource {
	return ConnectorTemplateResource{
		ApiVersion: ConnectorTemplateV2ApiVersion,
		Kind:       ConnectorTemplateV2Kind,
		Metadata: ConnectorTemplateMetadata{
			Name:   name,
			Labels: labels,
		},
		Spec: spec,
	}
}

func (r *ConnectorTemplateResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *ConnectorTemplateResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *ConnectorTemplateResource) FromRawJsonInterface(jsonInterface any) error {
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

func NewConnectorTemplateResourceFromClientResource(cliResource ctlresource.Resource) (ConnectorTemplateResource, error) {
	var consoleResource ConnectorTemplateResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return ConnectorTemplateResource{}, err
	}
	return consoleResource, nil
}
