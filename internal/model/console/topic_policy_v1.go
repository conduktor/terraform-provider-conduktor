package console

import (
	"encoding/json"
	"fmt"
	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const TopicPolicyV1Kind = "TopicPolicy"
const TopicPolicyV1ApiVersion = "v1"

type TopicPolicyMetadata struct {
	Name string `json:"name"`
}

func (r TopicPolicyMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type Constraint struct {
	AllowedKeys *AllowedKeys
	Match       *Match
	NoneOf      *NoneOf
	OneOf       *OneOf
	Range       *Range
}

func (dst *Constraint) UnmarshalJSON(data []byte) error {
	var disc model.Discriminable
	err := json.Unmarshal(data, &disc)
	if err != nil {
		return err
	}

	switch disc.Constraint {
	case "AllowedKeys":
		var allowedKeys AllowedKeys
		err = json.Unmarshal(data, &allowedKeys)
		if err != nil {
			return err
		}
		dst.AllowedKeys = &allowedKeys
	case "Match":
		var match Match
		err = json.Unmarshal(data, &match)
		if err != nil {
			return err
		}
		dst.Match = &match
	case "NoneOf":
		var noneOf NoneOf
		err = json.Unmarshal(data, &noneOf)
		if err != nil {
			return err
		}
		dst.NoneOf = &noneOf
	case "OneOf":
		var oneOf OneOf
		err = json.Unmarshal(data, &oneOf)
		if err != nil {
			return err
		}
		dst.OneOf = &oneOf
	case "Range":
		var rangeVal Range
		err = json.Unmarshal(data, &rangeVal)
		if err != nil {
			return err
		}
		dst.Range = &rangeVal
	default:
		return fmt.Errorf("unknown constraint type %s", disc.Constraint)
	}
	return nil
}

func (src *Constraint) MarshalJSON() ([]byte, error) {
	if src.AllowedKeys != nil {
		return json.Marshal(src.AllowedKeys)
	} else if src.Match != nil {
		return json.Marshal(src.Match)
	} else if src.NoneOf != nil {
		return json.Marshal(src.NoneOf)
	} else if src.OneOf != nil {
		return json.Marshal(src.OneOf)
	} else if src.Range != nil {
		return json.Marshal(src.Range)
	} else {
		return nil, fmt.Errorf("unknown constraint type")
	}
}

type AllowedKeys struct {
	Constraint string   `json:"constraint"`
	Optional   bool     `json:"optional,omitempty"`
	Keys       []string `json:"keys"`
}

type Match struct {
	Constraint string `json:"constraint"`
	Optional   bool   `json:"optional,omitempty"`
	Pattern    string `json:"pattern"`
}

type NoneOf struct {
	Constraint string   `json:"constraint"`
	Optional   bool     `json:"optional,omitempty"`
	Values     []string `json:"values"`
}

type OneOf struct {
	Constraint string   `json:"constraint"`
	Optional   bool     `json:"optional,omitempty"`
	Values     []string `json:"values"`
}

type Range struct {
	Constraint string `json:"constraint"`
	Optional   bool   `json:"optional,omitempty"`
	Min        int64  `json:"min"`
	Max        int64  `json:"max"`
}

type TopicPolicySpec struct {
	Policies map[string]*Constraint `json:"policies"`
}

type TopicPolicyResource struct {
	ApiVersion string              `json:"apiVersion"`
	Kind       string              `json:"kind"`
	Metadata   TopicPolicyMetadata `json:"metadata"`
	Spec       TopicPolicySpec     `json:"spec"`
}

func NewTopicPolicyResource(name string, spec TopicPolicySpec) TopicPolicyResource {
	return TopicPolicyResource{
		ApiVersion: TopicPolicyV1ApiVersion,
		Kind:       TopicPolicyV1Kind,
		Metadata: TopicPolicyMetadata{
			Name: name,
		},
		Spec: spec,
	}
}

func (r *TopicPolicyResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *TopicPolicyResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *TopicPolicyResource) FromRawJsonInterface(jsonInterface interface{}) error {
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

func NewTopicPolicyResourceFromClientResource(cliResource ctlresource.Resource) (TopicPolicyResource, error) {
	var consoleResource TopicPolicyResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return TopicPolicyResource{}, err
	}
	return consoleResource, nil
}
