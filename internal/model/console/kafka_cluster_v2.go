package console

import (
	"encoding/json"
	"fmt"
	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const KafkaClusterV2Kind = "KafkaCluster"
const KafkaClusterV2ApiVersion = "v2"

type KafkaClusterMetadata struct {
	Labels map[string]string `json:"labels,omitempty"`
	Name   string            `json:"name"`
}

func (r KafkaClusterMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type KafkaClusterSpec struct {
	BootstrapServers           string                   `json:"bootstrapServers"`
	ClientMetricsCollectors    []ClientMetricsCollector `json:"clientMetricsCollectors,omitempty"`
	Color                      string                   `json:"color,omitempty"`
	DisplayName                string                   `json:"displayName"`
	Icon                       string                   `json:"icon,omitempty"`
	IgnoreUntrustedCertificate bool                     `json:"ignoreUntrustedCertificate"`
	KafkaFlavor                *KafkaFlavor             `json:"kafkaFlavor,omitempty"`
	PoliciesRef                []string                 `json:"policiesRef,omitempty"`
	Properties                 map[string]string        `json:"properties,omitempty"`
	SchemaRegistry             *model.SchemaRegistry    `json:"schemaRegistry,omitempty"`
}

type ClientMetricsCollectorSecurityType string

const (
	CMC_BASIC_AUTH   ClientMetricsCollectorSecurityType = "BasicAuth"
	CMC_BEARER_TOKEN ClientMetricsCollectorSecurityType = "BearerToken"
	CMC_NO_SECURITY  ClientMetricsCollectorSecurityType = "NoSecurity"
)

type ClientMetricsCollectorBasicAuth struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ClientMetricsCollectorBearerToken struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type ClientMetricsCollectorNoSecurity struct {
	Type string `json:"type"`
}

type ClientMetricsCollectorSecurity struct {
	BasicAuth   *ClientMetricsCollectorBasicAuth
	BearerToken *ClientMetricsCollectorBearerToken
	NoSecurity  *ClientMetricsCollectorNoSecurity
}

func (dst *ClientMetricsCollectorSecurity) UnmarshalJSON(data []byte) error {
	var disc model.Discriminable
	if err := json.Unmarshal(data, &disc); err != nil {
		return err
	}
	switch ClientMetricsCollectorSecurityType(disc.Type) {
	case CMC_BASIC_AUTH:
		var v ClientMetricsCollectorBasicAuth
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		dst.BasicAuth = &v
	case CMC_BEARER_TOKEN:
		var v ClientMetricsCollectorBearerToken
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		dst.BearerToken = &v
	case CMC_NO_SECURITY:
		var v ClientMetricsCollectorNoSecurity
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		dst.NoSecurity = &v
	default:
		return fmt.Errorf("unknown ClientMetricsCollectorSecurity type %s", disc.Type)
	}
	return nil
}

func (src *ClientMetricsCollectorSecurity) MarshalJSON() ([]byte, error) {
	if src.BasicAuth != nil {
		return json.Marshal(src.BasicAuth)
	} else if src.BearerToken != nil {
		return json.Marshal(src.BearerToken)
	} else if src.NoSecurity != nil {
		return json.Marshal(src.NoSecurity)
	}
	return nil, fmt.Errorf("unknown ClientMetricsCollectorSecurity type")
}

type ClientMetricsCollector struct {
	Url                        string                          `json:"url"`
	IgnoreUntrustedCertificate bool                            `json:"ignoreUntrustedCertificate,omitempty"`
	Vcluster                   string                          `json:"vcluster,omitempty"`
	Security                   *ClientMetricsCollectorSecurity `json:"security,omitempty"`
}

type KafkaFlavorType string

const (
	AIVEN     KafkaFlavorType = "Aiven"
	CONFLUENT KafkaFlavorType = "Confluent"
	GATEWAY   KafkaFlavorType = "Gateway"
)

type KafkaFlavor struct {
	Aiven     *Aiven
	Confluent *Confluent
	Gateway   *Gateway
}

func (dst *KafkaFlavor) UnmarshalJSON(data []byte) error {
	var disc model.Discriminable
	err := json.Unmarshal(data, &disc)
	if err != nil {
		return err
	}

	kafkaFlavor := KafkaFlavorType(disc.Type)
	switch kafkaFlavor {
	case AIVEN:
		var aiven Aiven
		err = json.Unmarshal(data, &aiven)
		if err != nil {
			return err
		}
		dst.Aiven = &aiven
	case CONFLUENT:
		var confluent Confluent
		err = json.Unmarshal(data, &confluent)
		if err != nil {
			return err
		}
		dst.Confluent = &confluent
	case GATEWAY:
		var gateway Gateway
		err = json.Unmarshal(data, &gateway)
		if err != nil {
			return err
		}
		dst.Gateway = &gateway
	default:
		return fmt.Errorf("unknown kafkaFlavor type %s", disc.Type)
	}
	return nil
}

func (src *KafkaFlavor) MarshalJSON() ([]byte, error) {
	if src.Aiven != nil {
		return json.Marshal(src.Aiven)
	} else if src.Confluent != nil {
		return json.Marshal(src.Confluent)
	} else if src.Gateway != nil {
		return json.Marshal(src.Gateway)
	} else {
		return nil, fmt.Errorf("unknown kafkaFlavor type")
	}
}

type Aiven struct {
	ApiToken    string `json:"apiToken"`
	Project     string `json:"project"`
	ServiceName string `json:"serviceName"`
	Type        string `json:"type"`
}

type Confluent struct {
	Key                    string `json:"key"`
	Secret                 string `json:"secret"`
	ConfluentEnvironmentId string `json:"confluentEnvironmentId"`
	ConfluentClusterId     string `json:"confluentClusterId"`
	Type                   string `json:"type"`
}

type Gateway struct {
	Url                        string `json:"url"`
	User                       string `json:"user"`
	Password                   string `json:"password"`
	VirtualCluster             string `json:"virtualCluster"`
	IgnoreUntrustedCertificate bool   `json:"ignoreUntrustedCertificate,omitempty"`
	Type                       string `json:"type"`
}

type KafkaClusterResource struct {
	ApiVersion string               `json:"apiVersion"`
	Kind       string               `json:"kind"`
	Metadata   KafkaClusterMetadata `json:"metadata"`
	Spec       KafkaClusterSpec     `json:"spec"`
}

func NewKafkaClusterResource(name string, labels map[string]string, spec KafkaClusterSpec) KafkaClusterResource {
	return KafkaClusterResource{
		ApiVersion: KafkaClusterV2ApiVersion,
		Kind:       KafkaClusterV2Kind,
		Metadata: KafkaClusterMetadata{
			Name:   name,
			Labels: labels,
		},
		Spec: spec,
	}
}

func (r *KafkaClusterResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *KafkaClusterResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *KafkaClusterResource) FromRawJsonInterface(jsonInterface any) error {
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

func NewKafkaClusterResourceFromClientResource(cliResource ctlresource.Resource) (KafkaClusterResource, error) {
	var consoleResource KafkaClusterResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return KafkaClusterResource{}, err
	}
	return consoleResource, nil
}
