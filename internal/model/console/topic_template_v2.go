package console

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const TopicTemplateV2Kind = "TopicTemplate"
const TopicTemplateV2ApiVersion = "v2"

type TopicTemplateMetadata struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
}

func (r TopicTemplateMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type TopicTemplateDefaultsMetadata struct {
	Name   string            `json:"name,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

type TopicTemplateDefaultsSpec struct {
	Partitions        *int64            `json:"partitions,omitempty"`
	ReplicationFactor *int64            `json:"replicationFactor,omitempty"`
	Configs           map[string]string `json:"configs,omitempty"`
}

type TopicTemplateDefaults struct {
	Metadata TopicTemplateDefaultsMetadata `json:"metadata"`
	Spec     TopicTemplateDefaultsSpec     `json:"spec"`
}

type TopicTemplateSpec struct {
	DisplayName string                `json:"displayName"`
	Description string                `json:"description,omitempty"`
	Defaults    TopicTemplateDefaults `json:"defaults"`
}

type TopicTemplateResource struct {
	ApiVersion string                `json:"apiVersion"`
	Kind       string                `json:"kind"`
	Metadata   TopicTemplateMetadata `json:"metadata"`
	Spec       TopicTemplateSpec     `json:"spec"`
}

func NewTopicTemplateResource(name string, labels map[string]string, spec TopicTemplateSpec) TopicTemplateResource {
	return TopicTemplateResource{
		ApiVersion: TopicTemplateV2ApiVersion,
		Kind:       TopicTemplateV2Kind,
		Metadata: TopicTemplateMetadata{
			Name:   name,
			Labels: labels,
		},
		Spec: spec,
	}
}

func (r *TopicTemplateResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *TopicTemplateResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *TopicTemplateResource) FromRawJsonInterface(jsonInterface any) error {
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

func NewTopicTemplateResourceFromClientResource(cliResource ctlresource.Resource) (TopicTemplateResource, error) {
	var consoleResource TopicTemplateResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return TopicTemplateResource{}, err
	}
	return consoleResource, nil
}
