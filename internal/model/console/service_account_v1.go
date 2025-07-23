package console

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const ServiceAccountV1Kind = "ServiceAccount"
const ServiceAccountV1ApiVersion = "v1"

type ServiceAccountMetadata struct {
	AppInstance string            `json:"appInstance,omitempty"`
	Cluster     string            `json:"cluster"`
	Labels      map[string]string `json:"labels,omitempty"`
	Name        string            `json:"name"`
}

func (r ServiceAccountMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type ServiceAccountSpec struct {
	Authorization *ServiceAccountAuthorization `json:"authorization"`
}

type ServiceAccountAuthorization struct {
	Aiven *ServiceAccountAuthAiven
	Kafka *ServiceAccountAuthKafka
}

func (dst *ServiceAccountAuthorization) UnmarshalJSON(data []byte) error {
	var disc model.Discriminable
	err := json.Unmarshal(data, &disc)
	if err != nil {
		return err
	}

	switch disc.Type {
	case "AIVEN_ACL":
		var aiven ServiceAccountAuthAiven
		err = json.Unmarshal(data, &aiven)
		if err != nil {
			return err
		}
		dst.Aiven = &aiven
	case "KAFKA_ACL":
		var kafka ServiceAccountAuthKafka
		err = json.Unmarshal(data, &kafka)
		if err != nil {
			return err
		}
		dst.Kafka = &kafka
	default:
		return fmt.Errorf("unknown authorization type %s", disc.Type)
	}
	return nil
}

func (src *ServiceAccountAuthorization) MarshalJSON() ([]byte, error) {
	if src.Aiven != nil {
		return json.Marshal(src.Aiven)
	} else if src.Kafka != nil {
		return json.Marshal(src.Kafka)
	} else {
		return nil, fmt.Errorf("unknown authorization type")
	}
}

type ServiceAccountAuthAiven struct {
	ACLS []ServiceAccountAuthAivenACL `json:"acls"`
	Type string                       `json:"type"`
}
type ServiceAccountAuthAivenACL struct {
	ResourceType string `json:"resourceType"`
	Name         string `json:"name"`
	Permission   string `json:"permission"`
}

type ServiceAccountAuthKafka struct {
	ACLS []ServiceAccountAuthKafkaACL `json:"acls"`
	Type string                       `json:"type"`
}

type ServiceAccountAuthKafkaACL struct {
	Type           string   `json:"type"`
	Name           string   `json:"name"`
	PatternType    string   `json:"patternType"`
	ConnectCluster string   `json:"connectCluster,omitempty"`
	Operations     []string `json:"operations"`
	Host           string   `json:"host,omitempty"`
	Permission     string   `json:"permission,omitempty"`
}

type ServiceAccountResource struct {
	ApiVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   ServiceAccountMetadata `json:"metadata"`
	Spec       ServiceAccountSpec     `json:"spec"`
}

func NewServiceAccountResource(meta ServiceAccountMetadata, spec ServiceAccountSpec) ServiceAccountResource {
	return ServiceAccountResource{
		ApiVersion: ServiceAccountV1ApiVersion,
		Kind:       ServiceAccountV1Kind,
		Metadata:   meta,
		Spec:       spec,
	}
}

func (r *ServiceAccountResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *ServiceAccountResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *ServiceAccountResource) FromRawJsonInterface(jsonInterface any) error {
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

func NewServiceAccountResourceFromClientResource(cliResource ctlresource.Resource) (ServiceAccountResource, error) {
	var consoleResource ServiceAccountResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return ServiceAccountResource{}, err
	}
	return consoleResource, nil
}
