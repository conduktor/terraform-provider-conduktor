package console_kafka_subject_v2

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func normalizeAvroSchema(schemaStr, format string) (string, error) {
	if strings.ToUpper(format) != "AVRO" {
		return schemaStr, nil // Only normalize AVRO schemas
	}

	var schema interface{}
	if err := json.Unmarshal([]byte(schemaStr), &schema); err != nil {
		return schemaStr, fmt.Errorf("failed to parse schema JSON: %w", err)
	}

	normalized := normalizeAvroObject(schema)

	normalizedBytes, err := json.Marshal(normalized)
	if err != nil {
		return schemaStr, fmt.Errorf("failed to marshal normalized schema: %w", err)
	}

	return string(normalizedBytes), nil
}

func normalizeAvroObject(obj interface{}) interface{} {
	switch v := obj.(type) {
	case map[string]interface{}:
		return normalizeAvroMap(v)
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = normalizeAvroObject(item)
		}
		return result
	case string:
		return v
	default:
		return v
	}
}

func normalizeAvroMap(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Normalize field order according to Schema Registry canonical form
	fieldOrder := []string{"name", "type", "fields", "symbols", "items", "values", "size"}

	// Add ordered fields first
	for _, field := range fieldOrder {
		if val, exists := m[field]; exists {
			result[field] = normalizeAvroObject(val)
		}
	}

	// Add remaining fields in alphabetical order
	var remainingFields []string
	for key := range m {
		found := false
		for _, orderedField := range fieldOrder {
			if key == orderedField {
				found = true
				break
			}
		}
		if !found {
			remainingFields = append(remainingFields, key)
		}
	}
	sort.Strings(remainingFields)

	for _, field := range remainingFields {
		result[field] = normalizeAvroObject(m[field])
	}

	return result
}
