package customtypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectSchemaFormat(t *testing.T) {
	tests := []struct {
		name     string
		schema   string
		expected string
	}{
		// Protobuf tests
		{
			name:     "protobuf with syntax declaration",
			schema:   `syntax = "proto3";`,
			expected: "PROTOBUF",
		},
		{
			name: "protobuf with message",
			schema: `syntax = "proto3";
message Person {
  string name = 1;
  int32 id = 2;
}`,
			expected: "PROTOBUF",
		},
		{
			name: "protobuf with service",
			schema: `syntax = "proto3";
service MyService {
  rpc GetUser(UserRequest) returns (UserResponse);
}`,
			expected: "PROTOBUF",
		},
		{
			name: "protobuf with enum",
			schema: `syntax = "proto3";
enum Status {
  PENDING = 0;
  APPROVED = 1;
}`,
			expected: "PROTOBUF",
		},

		// AVRO tests - primitives
		{
			name:     "avro null primitive",
			schema:   `"null"`,
			expected: "AVRO",
		},
		{
			name:     "avro string primitive",
			schema:   `"string"`,
			expected: "AVRO",
		},
		{
			name:     "avro int primitive",
			schema:   `"int"`,
			expected: "AVRO",
		},

		// AVRO tests - complex types
		{
			name: "avro record with fields",
			schema: `{
				"type": "record",
				"name": "User",
				"fields": [
					{"name": "name", "type": "string"},
					{"name": "age", "type": "int"}
				]
			}`,
			expected: "AVRO",
		},
		{
			name: "avro enum with symbols",
			schema: `{
				"type": "enum",
				"name": "Status",
				"symbols": ["PENDING", "APPROVED", "REJECTED"]
			}`,
			expected: "AVRO",
		},

		// JSON Schema tests
		{
			name: "json schema with $schema property",
			schema: `{
				"$schema": "http://json-schema.org/draft-07/schema#",
				"type": "object",
				"properties": {
					"name": {"type": "string"}
				}
			}`,
			expected: "JSON",
		},
		{
			name: "json schema with properties",
			schema: `{
				"type": "object",
				"properties": {
					"name": {"type": "string"},
					"age": {"type": "integer"}
				},
				"required": ["name"]
			}`,
			expected: "JSON",
		},

		// Edge cases and unknown
		{
			name:     "empty string",
			schema:   "",
			expected: "UNKNOWN",
		},
		{
			name:     "whitespace only",
			schema:   "   \n\t  ",
			expected: "UNKNOWN",
		},
		{
			name:     "invalid json",
			schema:   `{"invalid": json}`,
			expected: "UNKNOWN",
		},
		{
			name:     "plain text",
			schema:   "this is just plain text",
			expected: "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectSchemaFormat(tt.schema)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeAvroSchema(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		normalized  string
		expectError bool
	}{
		{
			name:       "primitive string",
			input:      `"string"`,
			normalized: `"string"`,
		},
		{
			name:       "primitive string",
			input:      `  { "type": "string" }  `,
			normalized: `"string"`,
		},
		{
			name: "simple record schema",
			input: `{
				"type": "record",
				"name": "User", 
				"fields": [
					{"name": "id", "type": "int"},
					{"name": "name", "type": "string"}
				]
			}`,
			normalized: `{"name":"User","type":"record","fields":[{"name":"id","type":"int"},{"name":"name","type":"string"}]}`,
		},
		{
			name: "enum schema",
			input: `{
				"type": "enum",
				"name": "Status",
				"symbols": ["A", "B", "C"]
			}`,
			normalized: `{"name":"Status","type":"enum","symbols":["A","B","C"]}`,
		},
		{
			name: "array schema",
			input: `{
				"type": "array",
				"items": "string"
			}`,
			normalized: `{"type":"array","items":"string"}`,
		},
		{
			name:        "invalid json",
			input:       `{"invalid": json}`,
			expectError: true,
		},
		{
			name:        "invalid avro schema",
			input:       `{"type": "invalid_type"}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizeAvroSchema(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, result)

			// Verify the result is valid AVRO by trying to parse it again
			_, err = normalizeAvroSchema(result)
			assert.NoError(t, err)
		})
	}
}

func TestNormalizeProtobufSchema(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		normalized  string
		expectError bool
	}{
		{
			name: "simple message",
			input: `syntax = "proto3";

message   Person   {

  string  name = 1;

  int32 id=2;
}`,
			normalized: `syntax = "proto3";
message Person {
  string name = 1;
  int32 id = 2;
}`,
		},
		{
			name: "proto2 syntax",
			input: `syntax = "proto2";

message SearchRequest {
  optional string query = 1;
  optional int32 page_number = 2;
  optional int32 results_per_page = 3;
}`,
			normalized: `syntax = "proto2";
message SearchRequest {
  optional string query = 1;
  optional int32 page_number = 2;
  optional int32 results_per_page = 3;
}`,
		},
		{
			name: "edition proto",
			input: `edition = "2023";

message SearchRequest {
  string query = 1;

  int32 page_number = 2;

  int32 results_per_page = 3;
}`,
			normalized: `edition = "2023";
message SearchRequest {
  string query = 1;
  int32 page_number = 2;
  int32 results_per_page = 3;
}`,
		},
		{
			name: "with imports",
			input: `syntax = "proto3";
import "google/protobuf/timestamp.proto";
import "common/types.proto";
message Test {
  string value = 1;
}`,
			normalized: `syntax = "proto3";
import "google/protobuf/timestamp.proto";
import "common/types.proto";
message Test {
  string value = 1;
}`,
		},
		{
			name: "with comments",
			input: `// This is a comment
syntax = "proto3"; // syntax declaration
// Message definition
message Person {
  string name = 1; // Person's name
  int32 id = 2; // Person's ID
  // Nested message
  message Address {
	string street = 1;
	string city = 2;
  }
  Address address = 3; // Person's address
}`,
			normalized: `syntax = "proto3";
message Person {
  string name = 1;
  int32 id = 2;
  message Address {
	string street = 1;
	string city = 2;
  }
  Address address = 3;
}`,
		},
		{
			name: "service definition",
			input: `syntax = "proto3";
service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
}`,
			normalized: `syntax = "proto3";
service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
}`,
		},
		{
			name: "enum definition",
			input: `syntax = "proto3";
enum Status {
  UNKNOWN = 0;
  PENDING = 1;
  APPROVED = 2;
}`,
			normalized: `syntax = "proto3";
enum Status {
  UNKNOWN = 0;
  PENDING = 1;
  APPROVED = 2;
}`,
		},
		{
			name:        "invalid protobuf",
			input:       `invalid protobuf syntax`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizeProtobufSchema(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, result)
		})
	}
}

func TestNormalizeJSONSchema(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		normalized  string
		expectError bool
	}{
		{
			name: "simple object schema",
			input: `{
				"type": "object",
				"properties": {
					"name": {"type": "string"},
					"age": {"type": "integer"}
				}
			}`,
			normalized: `{"properties":{"age":{"type":"integer"},"name":{"type":"string"}},"type":"object"}`,
		},
		{
			name: "array schema",
			input: `{
				"type": "array",
				"items": {
					"type": "string"
				}
			}`,
			normalized: `{"items":{"type":"string"},"type":"array"}`,
		},
		{
			name: "schema with $schema property",
			input: `{
				"$schema": "http://json-schema.org/draft-07/schema#",
				"type": "object",
				"properties": {
					"user": {
						"type": "object",
						"properties": {
							"name": {"type": "string"}
						}
					}
				}
			}`,
			normalized: `{"$schema":"http://json-schema.org/draft-07/schema#","properties":{"user":{"properties":{"name":{"type":"string"}},"type":"object"}},"type":"object"}`,
		},
		{
			name:        "invalid json",
			input:       `{"invalid": json}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizeJSONSchema(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, result)

			// Verify the result is valid JSON
			var jsonData interface{}
			err = json.Unmarshal([]byte(result), &jsonData)
			assert.NoError(t, err)
		})
	}
}

func TestIsProtobufSchema(t *testing.T) {
	tests := []struct {
		name     string
		schema   string
		expected bool
	}{
		{
			name:     "valid protobuf with syntax",
			schema:   `syntax = "proto3";`,
			expected: true,
		},
		{
			name: "valid protobuf with message",
			schema: `syntax = "proto3";
message Person {
  string name = 1;
}`,
			expected: true,
		},
		{
			name: "valid syntax edition",
			schema: `edition = "2023";

message SearchRequest {
  string query = 1;
  int32 page_number = 2;
  int32 results_per_page = 3;
}`,
			expected: true,
		},
		{
			name: "valid syntax proto2",
			schema: `syntax = "proto2";

message SearchRequest {
  optional string query = 1;
  optional int32 page_number = 2;
  optional int32 results_per_page = 3;
}`,
			expected: true,
		},
		{
			name:     "invalid protobuf",
			schema:   `{"type": "object"}`,
			expected: false,
		},
		{
			name:     "plain text",
			schema:   `just some text`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isProtobufSchema(tt.schema)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsAvroSchema(t *testing.T) {
	tests := []struct {
		name     string
		schema   string
		expected bool
	}{
		{
			name:     "primitive string",
			schema:   `"string"`,
			expected: true,
		},
		{
			name: "record type",
			schema: `{
				"type": "record",
				"name": "User",
				"fields": [
					{"name": "name", "type": "string"}
				]
			}`,
			expected: true,
		},
		{
			name:     "invalid avro",
			schema:   `{"type": "invalid_type"}`,
			expected: false,
		},
		{
			name:     "plain text",
			schema:   `just text`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAvroSchema(tt.schema)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsJSONSchema(t *testing.T) {
	tests := []struct {
		name     string
		schema   string
		expected bool
	}{
		{
			name: "valid json schema with $schema",
			schema: `{
				"$schema": "http://json-schema.org/draft-07/schema#",
				"type": "object"
			}`,
			expected: true,
		},
		{
			name: "valid json schema with properties",
			schema: `{
				"type": "object",
				"properties": {
					"name": {"type": "string"}
				}
			}`,
			expected: true,
		},
		{
			name:     "plain json (not schema)",
			schema:   `{"data": "value"}`,
			expected: false,
		},
		{
			name:     "invalid json",
			schema:   `{invalid json}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isJSONSchema(tt.schema)
			assert.Equal(t, tt.expected, result)
		})
	}
}
