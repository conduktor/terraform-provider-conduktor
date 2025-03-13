package console

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const TopicV2Kind = "Topic"
const TopicV2ApiVersion = "v2"

type TopicSqlStorage struct {
	RetentionTimeInSecond int64 `json:"retentionTimeInSecond"`
	Enabled               bool  `json:"enabled,omitempty"`
}

type TopicConsoleMetadata struct {
	Name                  string            `json:"name"`
	Cluster               string            `json:"cluster"`
	Labels                map[string]string `json:"labels,omitempty"`
	CatalogVisibility     string            `json:"catalogVisibility,omitempty"`
	DescriptionIsEditable bool              `json:"descriptionIsEditable,omitempty"`
	Description           string            `json:"description,omitempty"`
	SqlStorage            *TopicSqlStorage  `json:"sqlStorage,omitempty"`
}

func (r TopicConsoleMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type TopicConsoleSpec struct {
	Partitions        int64             `json:"partitions"`
	ReplicationFactor int64             `json:"replicationFactor"`
	Configs           map[string]string `json:"configs,omitempty"`
}

type TopicConsoleResource struct {
	Kind       string               `json:"kind"`
	ApiVersion string               `json:"apiVersion"`
	Metadata   TopicConsoleMetadata `json:"metadata"`
	Spec       TopicConsoleSpec     `json:"spec"`
}

func NewTopicConsoleResource(meta TopicConsoleMetadata, spec TopicConsoleSpec) TopicConsoleResource {
	return TopicConsoleResource{
		Kind:       TopicV2Kind,
		ApiVersion: TopicV2ApiVersion,
		Metadata:   meta,
		Spec:       spec,
	}
}

//	func NewTopicConsoleMetadata(name string, cluster string, labels map[string]string, cVisibility string, dEditable bool, desc string, sql TopicSqlStorage) TopicConsoleMetadata {
//		return TopicConsoleMetadata{
//			Name:                  name,
//			Cluster:               cluster,
//			Labels:                labels,
//			CatalogVisibility:     cVisibility,
//			DescriptionIsEditable: dEditable,
//			Description:           desc,
//			SqlStorage:            sql,
//		}
//	}

func (r *TopicConsoleResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *TopicConsoleResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *TopicConsoleResource) FromRawJsonInterface(jsonInterface interface{}) error {
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

func NewTopicConsoleResourceFromClientResource(cliResource ctlresource.Resource) (TopicConsoleResource, error) {
	var consoleResource TopicConsoleResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return TopicConsoleResource{}, err
	}
	return consoleResource, nil
}
