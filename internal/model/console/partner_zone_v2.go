package console

import (
	"encoding/json"
	"fmt"

	ctlresource "github.com/conduktor/ctl/resource"
	model "github.com/conduktor/terraform-provider-conduktor/internal/model"
	jsoniter "github.com/json-iterator/go"
)

const PartnerZoneV2Kind = "PartnerZone"
const PartnerZoneV2ApiVersion = "v2"

type PartnerZoneConsoleMetadata struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
}

func (r PartnerZoneConsoleMetadata) String() string {
	return fmt.Sprintf(`name: %s`, r.Name)
}

type PartnerZoneAuthenticationMode struct {
	ServiceAccount string `json:"serviceAccount"`
	Type           string `json:"type"`
}

type PartnerZoneTopic struct {
	Name         string `json:"name"`
	BackingTopic string `json:"backingTopic"`
	Permission   string `json:"permission"`
}

type PartnerZonePartner struct {
	Name  string `json:"name"`
	Role  string `json:"role,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type PartnerZoneTrafficControlPolicies struct {
	MaxProduceRate    int64 `json:"maxProduceRate,omitempty"`
	MaxConsumeRate    int64 `json:"maxConsumeRate,omitempty"`
	LimitCommitOffset int64 `json:"limitCommitOffset,omitempty"`
}

type PartnerZoneToAdd struct {
	Key              string `json:"key"`
	Value            string `json:"value"`
	OverrideIfExists bool   `json:"overrideIfExists"`
}

type PartnerZoneToRemove struct {
	KeyRegex string `json:"keyRegex"`
}

type PartnerZoneHeaders struct {
	AddOnProduce    []PartnerZoneToAdd    `json:"addOnProduce,omitempty"`
	RemoveOnConsume []PartnerZoneToRemove `json:"removeOnConsume,omitempty"`
}

type PartnerZoneConsoleSpec struct {
	Cluster                string                            `json:"cluster"`
	DisplayName            string                            `json:"displayName,omitempty"`
	Description            string                            `json:"description,omitempty"`
	Url                    string                            `json:"url,omitempty"`
	AuthenticationMode     PartnerZoneAuthenticationMode     `json:"authenticationMode"`
	Topics                 []PartnerZoneTopic                `json:"topics,omitempty"`
	Partner                PartnerZonePartner                `json:"partner,omitempty"`
	TrafficControlPolicies PartnerZoneTrafficControlPolicies `json:"trafficControlPolicies,omitempty"`
	Headers                PartnerZoneHeaders                `json:"headers,omitempty"`
}

type PartnerZoneConsoleResource struct {
	Kind       string                     `json:"kind"`
	ApiVersion string                     `json:"apiVersion"`
	Metadata   PartnerZoneConsoleMetadata `json:"metadata"`
	Spec       PartnerZoneConsoleSpec     `json:"spec"`
}

func NewPartnerZoneConsoleResource(meta PartnerZoneConsoleMetadata, spec PartnerZoneConsoleSpec) PartnerZoneConsoleResource {
	return PartnerZoneConsoleResource{
		Kind:       PartnerZoneV2Kind,
		ApiVersion: PartnerZoneV2ApiVersion,
		Metadata:   meta,
		Spec:       spec,
	}
}

func (r *PartnerZoneConsoleResource) ToClientResource() (ctlresource.Resource, error) {
	return model.ToClientResource(r)
}

func (r *PartnerZoneConsoleResource) FromClientResource(cliResource ctlresource.Resource) error {
	err := jsoniter.Unmarshal(cliResource.Json, r)
	if err != nil {
		return err
	}
	return nil
}

func (r *PartnerZoneConsoleResource) FromRawJsonInterface(jsonInterface any) error {
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

func NewPartnerZoneConsoleResourceFromClientResource(cliResource ctlresource.Resource) (PartnerZoneConsoleResource, error) {
	var consoleResource PartnerZoneConsoleResource
	err := consoleResource.FromClientResource(cliResource)
	if err != nil {
		return PartnerZoneConsoleResource{}, err
	}
	return consoleResource, nil
}
