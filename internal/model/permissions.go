package model

import "sort"

type Permission struct {
	ResourceType string   `json:"resourceType"`
	Permissions  []string `json:"permissions"`
	Name         string   `json:"name,omitempty"`
	PatternType  string   `json:"patternType,omitempty"`
	Cluster      string   `json:"cluster,omitempty"`
	KafkaConnect string   `json:"kafkaConnect,omitempty"`
	KsqlDB       string   `json:"ksqlDB,omitempty"`
}

// matchesOnReturnedFields checks if two permissions match based on the fields
// that the API actually returns. The API may strip optional fields like
// kafka_connect, ksqldb, name, pattern_type, or cluster depending on the resource_type.
// Two permissions match if all their non-empty response fields are equal.
func (p Permission) matchesOnReturnedFields(other Permission) bool {
	if p.ResourceType != other.ResourceType {
		return false
	}
	if !stringSlicesEqual(p.Permissions, other.Permissions) {
		return false
	}
	// Compare fields that the API might return: if both are non-empty they must match.
	// If the response field is empty (stripped by API), it's not a distinguishing factor.
	if p.Name != "" && other.Name != "" && p.Name != other.Name {
		return false
	}
	if p.PatternType != "" && other.PatternType != "" && p.PatternType != other.PatternType {
		return false
	}
	if p.Cluster != "" && other.Cluster != "" && p.Cluster != other.Cluster {
		return false
	}
	if p.KafkaConnect != "" && other.KafkaConnect != "" && p.KafkaConnect != other.KafkaConnect {
		return false
	}
	if p.KsqlDB != "" && other.KsqlDB != "" && p.KsqlDB != other.KsqlDB {
		return false
	}
	return true
}

// MergeWithPlannedPermissions merges API response permissions with planned permissions.
// The Console API may strip optional fields (kafka_connect, ksqldb, name, pattern_type, cluster)
// from the response depending on the resource_type. This causes Terraform to report
// "Provider produced inconsistent result after apply" because the planned state has those
// fields set (even if null) but the response doesn't include them.
//
// This function finds matching permissions between planned and response, and preserves
// the planned field values for fields that the API response doesn't include.
func MergeWithPlannedPermissions(planned []Permission, response []Permission) []Permission {
	merged := make([]Permission, len(response))
	copy(merged, response)

	usedPlanned := make([]bool, len(planned))

	for i, resp := range merged {
		for j, plan := range planned {
			if usedPlanned[j] {
				continue
			}
			if resp.matchesOnReturnedFields(plan) {
				// Preserve planned values for fields that the API stripped
				if merged[i].Name == "" && plan.Name != "" {
					merged[i].Name = plan.Name
				}
				if merged[i].PatternType == "" && plan.PatternType != "" {
					merged[i].PatternType = plan.PatternType
				}
				if merged[i].Cluster == "" && plan.Cluster != "" {
					merged[i].Cluster = plan.Cluster
				}
				if merged[i].KafkaConnect == "" && plan.KafkaConnect != "" {
					merged[i].KafkaConnect = plan.KafkaConnect
				}
				if merged[i].KsqlDB == "" && plan.KsqlDB != "" {
					merged[i].KsqlDB = plan.KsqlDB
				}
				usedPlanned[j] = true
				break
			}
		}
	}

	return merged
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aSorted := make([]string, len(a))
	bSorted := make([]string, len(b))
	copy(aSorted, a)
	copy(bSorted, b)
	sort.Strings(aSorted)
	sort.Strings(bSorted)
	for i := range aSorted {
		if aSorted[i] != bSorted[i] {
			return false
		}
	}
	return true
}
