package console

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	"github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const KafkaSubjectV2Kind = "Subject"
const KafkaSubjectV2ApiVersion = "v2"

type KafkaSubjectMetadata struct {
	Cluster string            `json:"cluster"`
	Name    string            `json:"name"`
	Labels  map[string]string `json:"labels,omitempty"`
}

func (r KafkaSubjectMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type KafkaSubjectSpec struct {
	Schema        string                   `json:"schema"`
	Format        string                   `json:"format"`
	Version       *int                     `json:"version,omitempty"`
	Compatibility string                   `json:"compatibility,omitempty"`
	Id            *int                     `json:"id,omitempty"`
	References    []KafkaSubjectReferences `json:"references,omitempty"`
}

type KafkaSubjectReferences struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Version int    `json:"version"`
}

type KafkaSubjectResource struct {
	ApiVersion string               `json:"apiVersion"`
	Kind       string               `json:"kind"`
	Metadata   KafkaSubjectMetadata `json:"metadata"`
	Spec       KafkaSubjectSpec     `json:"spec"`
}

func NewKafkaSubjectResource(name string, cluster string, labels map[string]string, spec KafkaSubjectSpec) KafkaSubjectResource {
	return KafkaSubjectResource{
		ApiVersion: KafkaSubjectV2ApiVersion,
		Kind:       KafkaSubjectV2Kind,
		Metadata: KafkaSubjectMetadata{
			Name:    name,
			Cluster: cluster,
			Labels:  labels,
		},
		Spec: spec,
	}
}

func (r *KafkaSubjectResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *KafkaSubjectResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *KafkaSubjectResource) FromRawJsonInterface(jsonInterface any) error {
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

func NewKafkaSubjectResourceFromClientResource(cliResource ctlresource.Resource) (KafkaSubjectResource, error) {
	var consoleResource KafkaSubjectResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return KafkaSubjectResource{}, err
	}
	return consoleResource, nil
}
