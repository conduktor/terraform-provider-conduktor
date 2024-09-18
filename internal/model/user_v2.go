package model

import (
	"bytes"
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	jsoniter "github.com/json-iterator/go"
)

const UserV2Kind = "User"
const UserV2ApiVersion = "v2"

type UserConsoleMetadata struct {
	Name string `json:"name"`
}

func (r UserConsoleMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type UserConsoleSpec struct {
	FirstName   string       `json:"firstName"`
	LastName    string       `json:"lastName"`
	Permissions []Permission `json:"permissions"`
}

type UserConsoleResource struct {
	Kind       string              `json:"kind"`
	ApiVersion string              `json:"apiVersion"`
	Metadata   UserConsoleMetadata `json:"metadata"`
	Spec       UserConsoleSpec     `json:"spec"`
}

func NewUserConsoleResource(name string, spec UserConsoleSpec) UserConsoleResource {
	return UserConsoleResource{
		Kind:       UserV2Kind,
		ApiVersion: UserV2ApiVersion,
		Metadata: UserConsoleMetadata{
			Name: name,
		},
		Spec: spec,
	}
}

func (r *UserConsoleResource) ToClientResource() (ctlresource.Resource, error) {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return ctlresource.Resource{}, err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, jsonData, "", "\t")
	if err != nil {
		return ctlresource.Resource{}, err
	}

	ctlResource := ctlresource.Resource{}
	err = ctlResource.UnmarshalJSON(prettyJSON.Bytes())
	if err != nil {
		return ctlresource.Resource{}, err
	}
	return ctlResource, nil
}

func (r *UserConsoleResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func NewUserConsoleResourceFromClientResource(cliResource ctlresource.Resource) (UserConsoleResource, error) {
	var consoleResource UserConsoleResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return UserConsoleResource{}, err
	}
	return consoleResource, nil
}
