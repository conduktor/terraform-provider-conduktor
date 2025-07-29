package console

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const ConnectorV2Kind = "Connector"
const ConnectorV2ApiVersion = "v2"

type AutoRestart struct {
	Enabled          bool  `json:"enabled"`
	FrequencySeconds int64 `json:"frequencySeconds,omitempty"`
}

type ConnectorConsoleMetadata struct {
	Name           string            `json:"name"`
	Cluster        string            `json:"cluster"`
	ConnectCluster string            `json:"connectCluster"`
	Labels         map[string]string `json:"labels,omitempty"`
	AutoRestart    *AutoRestart      `json:"autoRestart,omitempty"`
	Description    string            `json:"description,omitempty"`
}

func (r ConnectorConsoleMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type ConnectorConsoleSpec struct {
	Config map[string]string `json:"config"`
}

type ConnectorConsoleResource struct {
	Kind       string                   `json:"kind"`
	ApiVersion string                   `json:"apiVersion"`
	Metadata   ConnectorConsoleMetadata `json:"metadata"`
	Spec       ConnectorConsoleSpec     `json:"spec"`
}

func NewConnectorConsoleResource(meta ConnectorConsoleMetadata, spec ConnectorConsoleSpec) ConnectorConsoleResource {
	return ConnectorConsoleResource{
		Kind:       ConnectorV2Kind,
		ApiVersion: ConnectorV2ApiVersion,
		Metadata:   meta,
		Spec:       spec,
	}
}

func (r *ConnectorConsoleResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *ConnectorConsoleResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *ConnectorConsoleResource) FromRawJsonInterface(jsonInterface any) error {
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

func NewConnectorConsoleResourceFromClientResource(cliResource ctlresource.Resource) (ConnectorConsoleResource, error) {
	var consoleResource ConnectorConsoleResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return ConnectorConsoleResource{}, err
	}
	return consoleResource, nil
}
