package console

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const ApplicationInstanceV1Kind = "ApplicationInstance"
const ApplicationInstanceV1ApiVersion = "v1"

type ApplicationInstanceConsoleMetadata struct {
	Name        string `json:"name"`
	Application string `json:"application"`
}

func (r ApplicationInstanceConsoleMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type ApplicationInstanceConsoleSpec struct {
	Cluster                          string                        `json:"cluster"`
	TopicPolicyRef                   []string                      `json:"topicPolicyRef,omitempty"`
	Resources                        []model.ResourceWithOwnership `json:"resources,omitempty"`
	ApplicationManagedServiceAccount bool                          `json:"applicationManagedServiceAccount"`
	ServiceAccount                   string                        `json:"serviceAccount,omitempty"`
	DefaultCatalogVisibility         string                        `json:"defaultCatalogVisibility,omitempty"`
}

type ApplicationInstanceConsoleResource struct {
	Kind       string                             `json:"kind"`
	ApiVersion string                             `json:"apiVersion"`
	Metadata   ApplicationInstanceConsoleMetadata `json:"metadata"`
	Spec       ApplicationInstanceConsoleSpec     `json:"spec"`
}

func NewApplicationInstanceConsoleResource(name string, app string, spec ApplicationInstanceConsoleSpec) ApplicationInstanceConsoleResource {
	return ApplicationInstanceConsoleResource{
		Kind:       ApplicationInstanceV1Kind,
		ApiVersion: ApplicationInstanceV1ApiVersion,
		Metadata: ApplicationInstanceConsoleMetadata{
			Name:        name,
			Application: app,
		},
		Spec: spec,
	}
}

func (r *ApplicationInstanceConsoleResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *ApplicationInstanceConsoleResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *ApplicationInstanceConsoleResource) FromRawJsonInterface(jsonInterface interface{}) error {
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

func NewApplicationInstanceConsoleResourceFromClientResource(cliResource ctlresource.Resource) (ApplicationInstanceConsoleResource, error) {
	var consoleResource ApplicationInstanceConsoleResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return ApplicationInstanceConsoleResource{}, err
	}
	return consoleResource, nil
}
