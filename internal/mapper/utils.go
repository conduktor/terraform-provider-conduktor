package mapper

import "strings"
import schema "github.com/conduktor/terraform-provider-conduktor/internal/schema/validation"

// SplitLabels separates user-defined labels from managed labels.
func SplitLabels(labels map[string]string) (map[string]string, map[string]string) {
	if labels == nil {
		return nil, nil
	}

	userLabels := make(map[string]string)
	managedLabels := make(map[string]string)

	for key, value := range labels {
		if strings.HasPrefix(key, schema.ManagedLabelsPrefix) {
			managedLabels[key] = value
		} else {
			userLabels[key] = value
		}
	}

	return userLabels, managedLabels
}

// MergeLabels combines user-defined labels with managed labels.
func MergeLabels(userLabels, managedLabels map[string]string) map[string]string {
	mergedLabels := make(map[string]string)

	for key, value := range userLabels {
		mergedLabels[key] = value
	}

	for key, value := range managedLabels {
		mergedLabels[key] = value
	}

	return mergedLabels
}
