package model

type Permission struct {
	ResourceType string   `json:"resourceType"`
	Permissions  []string `json:"permissions"`
	Name         string   `json:"name,omitempty"`
	PatternType  string   `json:"patternType,omitempty"`
	Cluster      string   `json:"cluster,omitempty"`
	KafkaConnect string   `json:"kafkaConnect,omitempty"`
}

type ApplicationGroupPermission struct {
	AppInstance    string   `json:"app_instance"`
	PatternType    string   `json:"pattern_type"`
	ConnectCluster string   `json:"connect_cluster,omitempty"`
	ResourceType   string   `json:"resource_type"`
	Name           string   `json:"name"`
	Permissions    []string `json:"permissions"`
}
