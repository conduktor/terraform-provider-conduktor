package mapper

import "strings"

// Prefix of labels managed by Conduktor Console.
const managedLabelPrefix = "conduktor.io/"

// SplitLabels separates user-defined labels from managed labels.
func SplitLabels(labels map[string]string) (map[string]string, map[string]string) {
	userLabels := make(map[string]string)
	managedLabels := make(map[string]string)

	for key, value := range labels {
		if strings.HasPrefix(key, managedLabelPrefix) {
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
