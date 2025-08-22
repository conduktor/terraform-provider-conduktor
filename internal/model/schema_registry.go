package model

import (
	"encoding/json"
	"fmt"
)

type SchemaRegistryType string

const (
	CONFLUENT SchemaRegistryType = "ConfluentLike"
	GLUE      SchemaRegistryType = "Glue"
)

type SchemaRegistry struct {
	ConfluentLike *ConfluentLike
	Glue          *Glue
}

func (s *SchemaRegistry) UnmarshalJSON(data []byte) error {
	var disc Discriminable
	err := json.Unmarshal(data, &disc)
	if err != nil {
		return err
	}

	schemaRegistryType := SchemaRegistryType(disc.Type)
	switch schemaRegistryType {
	case CONFLUENT:
		var confluentLike ConfluentLike
		err = json.Unmarshal(data, &confluentLike)
		if err != nil {
			return err
		}
		s.ConfluentLike = &confluentLike
	case GLUE:
		var glue Glue
		err = json.Unmarshal(data, &glue)
		if err != nil {
			return err
		}
		s.Glue = &glue
	default:
		return fmt.Errorf("unknown schemaRegistry type %s", disc.Type)
	}
	return nil
}

func (s SchemaRegistry) MarshalJSON() ([]byte, error) {
	if s.ConfluentLike != nil {
		return json.Marshal(s.ConfluentLike)
	} else if s.Glue != nil {
		return json.Marshal(s.Glue)
	} else {
		return nil, fmt.Errorf("unknown schemaRegistry type")
	}
}

type ConfluentLike struct {
	Type                       string                              `json:"type"`
	Url                        string                              `json:"url"`
	Security                   ConfluentLikeSchemaRegistrySecurity `json:"security"`
	Properties                 string                              `json:"properties,omitempty"`
	IgnoreUntrustedCertificate bool                                `json:"ignoreUntrustedCertificate"`
}

type ConfluentSecurityType string

const (
	BASIC_AUTH   ConfluentSecurityType = "BasicAuth"
	BEARER_TOKEN ConfluentSecurityType = "BearerToken"
	NO_SECURITY  ConfluentSecurityType = "NoSecurity"
	SSL_AUTH     ConfluentSecurityType = "SSLAuth"
)

type ConfluentLikeSchemaRegistrySecurity struct {
	BasicAuth   *BasicAuth
	BearerToken *BearerToken
	NoSecurity  *NoSecurity
	SSLAuth     *SSLAuth
}

func (s *ConfluentLikeSchemaRegistrySecurity) UnmarshalJSON(bytes []byte) error {
	var disc Discriminable
	err := json.Unmarshal(bytes, &disc)
	if err != nil {
		return err
	}

	confluentSecurityType := ConfluentSecurityType(disc.Type)
	switch confluentSecurityType {
	case BASIC_AUTH:
		var basic BasicAuth
		err = json.Unmarshal(bytes, &basic)
		if err != nil {
			return err
		}
		s.BasicAuth = &basic
	case BEARER_TOKEN:
		var bearertoken BearerToken
		err = json.Unmarshal(bytes, &bearertoken)
		if err != nil {
			return err
		}
		s.BearerToken = &bearertoken
	case NO_SECURITY:
		var nosecurity NoSecurity
		err = json.Unmarshal(bytes, &nosecurity)
		if err != nil {
			return err
		}
		s.NoSecurity = &nosecurity
	case SSL_AUTH:
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

func (s ConfluentLikeSchemaRegistrySecurity) MarshalJSON() ([]byte, error) {
	if s.BasicAuth != nil {
		return json.Marshal(s.BasicAuth)
	} else if s.BearerToken != nil {
		return json.Marshal(s.BearerToken)
	} else if s.NoSecurity != nil {
		return json.Marshal(s.NoSecurity)
	} else if s.SSLAuth != nil {
		return json.Marshal(s.SSLAuth)
	} else {
		return nil, fmt.Errorf("unknown confluentLikeSchemaRegistrySecurity type")
	}
}

type NoSecurity struct {
	Type string `json:"type"`
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
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Glue struct {
	Type         string         `json:"type"`
	RegistryName string         `json:"registryName,omitempty"`
	Region       string         `json:"region"`
	Security     AmazonSecurity `json:"security"`
}

type AmazonSecurityType string

const (
	CREDENTIALS  AmazonSecurityType = "Credentials"
	FROM_CONTEXT AmazonSecurityType = "FromContext"
	FROM_ROLE    AmazonSecurityType = "FromRole"
	IAM_ANYWHERE AmazonSecurityType = "IAMAnywhere"
)

type AmazonSecurity struct {
	Credentials *Credentials
	FromContext *FromContext
	FromRole    *FromRole
	IAMAnywhere *IAMAnywhere
}

func (s *AmazonSecurity) UnmarshalJSON(bytes []byte) error {
	var disc Discriminable
	err := json.Unmarshal(bytes, &disc)
	if err != nil {
		return err
	}

	amazonSecurityType := AmazonSecurityType(disc.Type)
	switch amazonSecurityType {
	case CREDENTIALS:
		var creds Credentials
		err = json.Unmarshal(bytes, &creds)
		if err != nil {
			return err
		}
		s.Credentials = &creds
	case FROM_CONTEXT:
		var fromcontext FromContext
		err = json.Unmarshal(bytes, &fromcontext)
		if err != nil {
			return err
		}
		s.FromContext = &fromcontext
	case FROM_ROLE:
		var fromrole FromRole
		err = json.Unmarshal(bytes, &fromrole)
		if err != nil {
			return err
		}
		s.FromRole = &fromrole
	case IAM_ANYWHERE:
		var iamanywhere IAMAnywhere
		err = json.Unmarshal(bytes, &iamanywhere)
		if err != nil {
			return err
		}
		s.IAMAnywhere = &iamanywhere
	default:
		return fmt.Errorf("unknown amazonSecurity type %s", disc.Type)
	}
	return nil
}

func (s AmazonSecurity) MarshalJSON() ([]byte, error) {
	if s.Credentials != nil {
		return json.Marshal(s.Credentials)
	} else if s.FromContext != nil {
		return json.Marshal(s.FromContext)
	} else if s.FromRole != nil {
		return json.Marshal(s.FromRole)
	} else if s.IAMAnywhere != nil {
		return json.Marshal(s.IAMAnywhere)
	} else {
		return nil, fmt.Errorf("unknown amazonSecurity type")
	}
}

type Credentials struct {
	AccessKeyId string `json:"accessKeyId"`
	SecretKey   string `json:"secretKey"`
	Type        string `json:"type"`
}

type FromContext struct {
	Profile string `json:"profile,omitempty"`
	Type    string `json:"type"`
}

type FromRole struct {
	Role string `json:"role"`
	Type string `json:"type"`
}

type IAMAnywhere struct {
	TrustAnchorArn string `json:"trustAnchorArn"`
	ProfileArn     string `json:"profileArn"`
	RoleArn        string `json:"roleArn"`
	Certificate    string `json:"certificate"`
	PrivateKey     string `json:"privateKey"`
	Type           string `json:"type"`
}
