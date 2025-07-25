package model

import (
	"bytes"
	"encoding/json"
	ctlresource "github.com/conduktor/ctl/resource"
)

type Discriminable struct {
	Type       string `json:"type,omitempty"`
	Constraint string `json:"constraint,omitempty"`
}

func ToClientResource(o any) (ctlresource.Resource, error) {
	jsonData, err := json.Marshal(o)
	if err != nil {
		return ctlresource.Resource{}, err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, jsonData, "", "\t")
	if err != nil {
		return ctlresource.Resource{}, err
	}

	ctlResource := ctlresource.Resource{}
	err = ctlResource.UnmarshalJSON(prettyJSON.Bytes())
	if err != nil {
		return ctlresource.Resource{}, err
	}
	return ctlResource, nil
}
