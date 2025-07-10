package console

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const GroupV2Kind = "Group"
const GroupV2ApiVersion = "v2"

type GroupConsoleMetadata struct {
	Name string `json:"name"`
}

func (r GroupConsoleMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type GroupConsoleSpec struct {
	Description               string             `json:"description,omitempty"`
	DisplayName               string             `json:"displayName"`
	ExternalGroups            []string           `json:"externalGroups"`
	Members                   []string           `json:"members"`
	MembersFromExternalGroups []string           `json:"membersFromExternalGroups"`
	Permissions               []model.Permission `json:"permissions"`
}

type GroupConsoleResource struct {
	Kind       string               `json:"kind"`
	ApiVersion string               `json:"apiVersion"`
	Metadata   GroupConsoleMetadata `json:"metadata"`
	Spec       GroupConsoleSpec     `json:"spec"`
}

func NewGroupConsoleResource(name string, spec GroupConsoleSpec) GroupConsoleResource {
	return GroupConsoleResource{
		Kind:       GroupV2Kind,
		ApiVersion: GroupV2ApiVersion,
		Metadata: GroupConsoleMetadata{
			Name: name,
		},
		Spec: spec,
	}
}

func (r *GroupConsoleResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *GroupConsoleResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *GroupConsoleResource) FromRawJsonInterface(jsonInterface any) error {
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

func NewGroupConsoleResourceFromClientResource(cliResource ctlresource.Resource) (GroupConsoleResource, error) {
	var consoleResource GroupConsoleResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return GroupConsoleResource{}, err
	}
	return consoleResource, nil
}
