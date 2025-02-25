package model

type ResourceWithOwnership struct {
	Type           string `json:"type"`
	Name           string `json:"name"`
	PatternType    string `json:"patternType"`
	ConnectCluster string `json:"connectCluster,omitempty"`
	OwnershipMode  string `json:"ownershipMode,omitempty"`
}
