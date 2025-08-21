package console

import (
	"encoding/json"
	"fmt"
	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const KsqlDBClusterV2Kind = "KsqlDBCluster"
const KsqlDBClusterV2ApiVersion = "v2"

type KsqlDBClusterMetadata struct {
	Name    string `json:"name"`
	Cluster string `json:"cluster"`
}

func (r KsqlDBClusterMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type KsqlDBClusterSpec struct {
	Url                        string                 `json:"url"`
	DisplayName                string                 `json:"displayName"`
	IgnoreUntrustedCertificate bool                   `json:"ignoreUntrustedCertificate"`
	Headers                    map[string]string      `json:"headers,omitempty"`
	Security                   *KsqlDBClusterSecurity `json:"security,omitempty"`
}

type KsqlDBClusterSecurityType string

const (
	KSQLDB_BASIC_AUTH   KsqlDBClusterSecurityType = "BasicAuth"
	KSQLDB_BEARER_TOKEN KsqlDBClusterSecurityType = "BearerToken"
	KSQLDB_SSL_AUTH     KsqlDBClusterSecurityType = "SSLAuth"
)

type KsqlDBClusterSecurity struct {
	BasicAuth   *KsqlDBClusterBasicAuth
	BearerToken *KsqlDBClusterBearerToken
	SSLAuth     *KsqlDBClusterSSLAuth
}

func (s *KsqlDBClusterSecurity) UnmarshalJSON(bytes []byte) error {
	var disc model.Discriminable
	err := json.Unmarshal(bytes, &disc)
	if err != nil {
		return err
	}

	securityType := KsqlDBClusterSecurityType(disc.Type)
	switch securityType {
	case KSQLDB_BASIC_AUTH:
		var basic KsqlDBClusterBasicAuth
		err = json.Unmarshal(bytes, &basic)
		if err != nil {
			return err
		}
		s.BasicAuth = &basic
	case KSQLDB_BEARER_TOKEN:
		var bearertoken KsqlDBClusterBearerToken
		err = json.Unmarshal(bytes, &bearertoken)
		if err != nil {
			return err
		}
		s.BearerToken = &bearertoken
	case KSQLDB_SSL_AUTH:
		var sslauth KsqlDBClusterSSLAuth
		err = json.Unmarshal(bytes, &sslauth)
		if err != nil {
			return err
		}
		s.SSLAuth = &sslauth
	default:
		return nil
	}
	return nil
}

func (s KsqlDBClusterSecurity) MarshalJSON() ([]byte, error) {
	if s.BasicAuth != nil {
		return json.Marshal(s.BasicAuth)
	} else if s.BearerToken != nil {
		return json.Marshal(s.BearerToken)
	} else if s.SSLAuth != nil {
		return json.Marshal(s.SSLAuth)
	} else {
		return nil, nil
	}
}

type KsqlDBClusterSSLAuth struct {
	Key              string `json:"key"`
	CertificateChain string `json:"certificateChain"`
	Type             string `json:"type"`
}

type KsqlDBClusterBearerToken struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

type KsqlDBClusterBasicAuth struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type KsqlDBClusterResource struct {
	ApiVersion string                `json:"apiVersion"`
	Kind       string                `json:"kind"`
	Metadata   KsqlDBClusterMetadata `json:"metadata"`
	Spec       KsqlDBClusterSpec     `json:"spec"`
}

func NewKsqlDBClusterResource(name string, cluster string, spec KsqlDBClusterSpec) KsqlDBClusterResource {
	return KsqlDBClusterResource{
		ApiVersion: KsqlDBClusterV2ApiVersion,
		Kind:       KsqlDBClusterV2Kind,
		Metadata: KsqlDBClusterMetadata{
			Name:    name,
			Cluster: cluster,
		},
		Spec: spec,
	}
}

func (r *KsqlDBClusterResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *KsqlDBClusterResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *KsqlDBClusterResource) FromRawJsonInterface(jsonInterface any) error {
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

func NewKsqlDBClusterResourceFromClientResource(cliResource ctlresource.Resource) (KsqlDBClusterResource, error) {
	var consoleResource KsqlDBClusterResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return KsqlDBClusterResource{}, err
	}
	return consoleResource, nil
}
