package console

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const ApplicationGroupV1Kind = "ApplicationGroup"
const ApplicationGroupV1ApiVersion = "v1"

type ApplicationGroupMetadata struct {
	Name        string `json:"name"`
	Application string `json:"application"`
}

func (r ApplicationGroupMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type ApplicationGroupPermission struct {
	AppInstance    string   `json:"appInstance"`
	PatternType    string   `json:"patternType"`
	ConnectCluster string   `json:"connectCluster,omitempty"`
	ResourceType   string   `json:"resourceType"`
	Name           string   `json:"name"`
	Permissions    []string `json:"permissions"`
}

type ApplicationGroupSpec struct {
	DisplayName               string                       `json:"displayName"`
	Description               string                       `json:"description"`
	Permissions               []ApplicationGroupPermission `json:"permissions"`
	Members                   []string                     `json:"members"`
	ExternalGroups            []string                     `json:"externalGroups"`
	MembersFromExternalGroups []string                     `json:"membersFromExternalGroups"`
}

type ApplicationGroupConsoleResource struct {
	Kind       string                   `json:"kind"`
	ApiVersion string                   `json:"apiVersion"`
	Metadata   ApplicationGroupMetadata `json:"metadata"`
	Spec       ApplicationGroupSpec     `json:"spec"`
}

func NewApplicationGroupConsoleResource(name string, app string, spec ApplicationGroupSpec) ApplicationGroupConsoleResource {
	return ApplicationGroupConsoleResource{
		Kind:       ApplicationGroupV1Kind,
		ApiVersion: ApplicationGroupV1ApiVersion,
		Metadata: ApplicationGroupMetadata{
			Name:        name,
			Application: app,
		},
		Spec: spec,
	}
}

func (r *ApplicationGroupConsoleResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *ApplicationGroupConsoleResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *ApplicationGroupConsoleResource) FromRawJsonInterface(jsonInterface interface{}) error {
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

func NewApplicationGroupConsoleResourceFromClientResource(cliResource ctlresource.Resource) (ApplicationGroupConsoleResource, error) {
	var consoleResource ApplicationGroupConsoleResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return ApplicationGroupConsoleResource{}, err
	}
	return consoleResource, nil
}
