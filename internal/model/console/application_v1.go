package console

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const ApplicationV1Kind = "Application"
const ApplicationV1ApiVersion = "v1"

type ApplicationConsoleMetadata struct {
	Name string `json:"name"`
}

func (r ApplicationConsoleMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type ApplicationConsoleSpec struct {
	Description string `json:"description,omitempty"`
	Title       string `json:"title"`
	Owner       string `json:"owner"`
}

type ApplicationConsoleResource struct {
	Kind       string                     `json:"kind"`
	ApiVersion string                     `json:"apiVersion"`
	Metadata   ApplicationConsoleMetadata `json:"metadata"`
	Spec       ApplicationConsoleSpec     `json:"spec"`
}

func NewApplicationConsoleResource(name string, spec ApplicationConsoleSpec) ApplicationConsoleResource {
	return ApplicationConsoleResource{
		Kind:       ApplicationV1Kind,
		ApiVersion: ApplicationV1ApiVersion,
		Metadata: ApplicationConsoleMetadata{
			Name: name,
		},
		Spec: spec,
	}
}

func (r *ApplicationConsoleResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *ApplicationConsoleResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *ApplicationConsoleResource) FromRawJsonInterface(jsonInterface any) error {
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

func NewApplicationConsoleResourceFromClientResource(cliResource ctlresource.Resource) (ApplicationConsoleResource, error) {
	var consoleResource ApplicationConsoleResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return ApplicationConsoleResource{}, err
	}
	return consoleResource, nil
}
