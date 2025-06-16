package console

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const ApplicationInstancePermissionV1Kind = "ApplicationInstancePermission"
const ApplicationInstancePermissionV1ApiVersion = "v1"

type ApplicationInstancePermissionConsoleMetadata struct {
	Name        string `json:"name"`
	Application string `json:"application"`
	AppInstance string `json:"appInstance"`
}

func (r ApplicationInstancePermissionConsoleMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type AppInstancePermissionResource struct {
	Type           string `json:"type"`
	Name           string `json:"name"`
	PatternType    string `json:"patternType"`
	ConnectCluster string `json:"connectCluster,omitempty"`
}

type ApplicationInstancePermissionConsoleSpec struct {
	Resource                 AppInstancePermissionResource `json:"resource"`
	UserPermission           string                        `json:"userPermission"`
	ServiceAccountPermission string                        `json:"serviceAccountPermission"`
	GrantedTo                string                        `json:"grantedTo"`
}

type ApplicationInstancePermissionConsoleResource struct {
	Kind       string                                       `json:"kind"`
	ApiVersion string                                       `json:"apiVersion"`
	Metadata   ApplicationInstancePermissionConsoleMetadata `json:"metadata"`
	Spec       ApplicationInstancePermissionConsoleSpec     `json:"spec"`
}

func NewApplicationInstancePermissionConsoleResource(meta ApplicationInstancePermissionConsoleMetadata, spec ApplicationInstancePermissionConsoleSpec) ApplicationInstancePermissionConsoleResource {
	return ApplicationInstancePermissionConsoleResource{
		Kind:       ApplicationInstancePermissionV1Kind,
		ApiVersion: ApplicationInstancePermissionV1ApiVersion,
		Metadata:   meta,
		Spec:       spec,
	}
}

func (r *ApplicationInstancePermissionConsoleResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *ApplicationInstancePermissionConsoleResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *ApplicationInstancePermissionConsoleResource) FromRawJsonInterface(jsonInterface interface{}) error {
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

func NewApplicationInstancePermissionConsoleResourceFromClientResource(cliResource ctlresource.Resource) (ApplicationInstancePermissionConsoleResource, error) {
	var consoleResource ApplicationInstancePermissionConsoleResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return ApplicationInstancePermissionConsoleResource{}, err
	}
	return consoleResource, nil
}
