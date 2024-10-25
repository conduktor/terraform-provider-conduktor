package model

import (
	"encoding/json"
	"fmt"
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

	switch disc.Type {
	case "ConfluentLike":
		var confluentLike ConfluentLike
		err = json.Unmarshal(data, &confluentLike)
		if err != nil {
			return err
		}
		s.ConfluentLike = &confluentLike
	case "Glue":
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
	case "NoSecurity":
		var nosecurity NoSecurity
		err = json.Unmarshal(bytes, &nosecurity)
		if err != nil {
			return err
		}
		s.NoSecurity = &nosecurity
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
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
}

type Glue struct {
	Type         string         `json:"type"`
	RegistryName string         `json:"registryName,omitempty"`
	Region       string         `json:"region"`
	Security     AmazonSecurity `json:"security"`
}

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
	switch disc.Type {
	case "Credentials":
		var creds Credentials
		err = json.Unmarshal(bytes, &creds)
		if err != nil {
			return err
		}
		s.Credentials = &creds
	case "FromContext":
		var fromcontext FromContext
		err = json.Unmarshal(bytes, &fromcontext)
		if err != nil {
			return err
		}
		s.FromContext = &fromcontext
	case "FromRole":
		var fromrole FromRole
		err = json.Unmarshal(bytes, &fromrole)
		if err != nil {
			return err
		}
		s.FromRole = &fromrole
	case "IAMAnywhere":
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
