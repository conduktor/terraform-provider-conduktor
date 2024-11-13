package model

import (
	"encoding/json"
	"fmt"
	ctlresource "github.com/conduktor/ctl/resource"
	jsoniter "github.com/json-iterator/go"
)

const KafkaConnectV2Kind = "KafkaConnectCluster"
const KafkaConnectV2ApiVersion = "v2"

type KafkaConnectMetadata struct {
	Labels  map[string]string `json:"labels,omitempty"`
	Name    string            `json:"name"`
	Cluster string            `json:"cluster"`
}

func (r KafkaConnectMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type KafkaConnectSpec struct {
	Urls                       string                `json:"urls"`
	DisplayName                string                `json:"displayName"`
	IgnoreUntrustedCertificate bool                  `json:"ignoreUntrustedCertificate"`
	Headers                    map[string]string     `json:"headers,omitempty"`
	Security                   *KafkaConnectSecurity `json:"security,omitempty"`
}

type KafkaConnectSecurity struct {
	BasicAuth   *BasicAuth
	BearerToken *BearerToken
	SSLAuth     *SSLAuth
}

func (s *KafkaConnectSecurity) UnmarshalJSON(bytes []byte) error {
	var disc Discriminable
	err := json.Unmarshal(bytes, &disc)
	if err != nil {
		return err
	}
	switch disc.Type {
	case "BasicAuth":
		var basic BasicAuth
		err = json.Unmarshal(bytes, &basic)
		if err != nil {
			return err
		}
		s.BasicAuth = &basic
	case "BearerToken":
		var bearertoken BearerToken
		err = json.Unmarshal(bytes, &bearertoken)
		if err != nil {
			return err
		}
		s.BearerToken = &bearertoken
	case "SSLAuth":
		var sslauth SSLAuth
		err = json.Unmarshal(bytes, &sslauth)
		if err != nil {
			return err
		}
		s.SSLAuth = &sslauth
	default:
		return fmt.Errorf("unknown confluentLikeSchemaRegistrySecurity type %s", disc.Type)
	}
	return nil
}

func (s KafkaConnectSecurity) MarshalJSON() ([]byte, error) {
	if s.BasicAuth != nil {
		return json.Marshal(s.BasicAuth)
	} else if s.BearerToken != nil {
		return json.Marshal(s.BearerToken)
	} else if s.SSLAuth != nil {
		return json.Marshal(s.SSLAuth)
	} else {
		return nil, fmt.Errorf("unknown confluentLikeSchemaRegistrySecurity type")
	}
}

type SSLAuth struct {
	Key              string `json:"key"`
	CertificateChain string `json:"certificateChain"`
	Type             string `json:"type"`
}

type BearerToken struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

type BasicAuth struct {
	Type     string `json:"type"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type KafkaConnectResource struct {
	ApiVersion string               `json:"apiVersion"`
	Kind       string               `json:"kind"`
	Metadata   KafkaConnectMetadata `json:"metadata"`
	Spec       KafkaConnectSpec     `json:"spec"`
}

func NewKafkaConnectResource(name string, cluster string, labels map[string]string, spec KafkaConnectSpec) KafkaConnectResource {
	return KafkaConnectResource{
		ApiVersion: KafkaConnectV2ApiVersion,
		Kind:       KafkaConnectV2Kind,
		Metadata: KafkaConnectMetadata{
			Name:    name,
			Cluster: cluster,
			Labels:  labels,
		},
		Spec: spec,
	}
}

func (r *KafkaConnectResource) ToClientResource() (ctlresource.Resource, error) {
	return toClientResource(r)
}

func (r *KafkaConnectResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *KafkaConnectResource) FromRawJsonInterface(jsonInterface interface{}) error {
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

func NewKafkaConnectResourceFromClientResource(cliResource ctlresource.Resource) (KafkaConnectResource, error) {
	var consoleResource KafkaConnectResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return KafkaConnectResource{}, err
	}
	return consoleResource, nil
}
