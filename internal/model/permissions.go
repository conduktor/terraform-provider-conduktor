package model

type Permission struct {
	ResourceType string   `json:"resourceType"`
	Permissions  []string `json:"permissions"`
	Name         string   `json:"name,omitempty"`
	PatternType  string   `json:"patternType,omitempty"`
	Cluster      string   `json:"cluster,omitempty"`
	KafkaConnect string   `json:"kafkaConnect,omitempty"`
	KsqlDB       string   `json:"ksqlDB,omitempty"`
}
