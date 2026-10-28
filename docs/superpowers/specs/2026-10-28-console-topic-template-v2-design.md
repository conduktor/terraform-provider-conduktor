# console_topic_template_v2 Design

Date: 2026-10-28

## Summary

Add `conduktor_console_topic_template_v2` as a first-class Terraform resource, backed by the Conduktor Console `/public/console/v2/topic-template` API. Topic Templates let platform teams define reusable topic blueprints (partitions, replication factor, configs, naming patterns with Handlebars variables) that Console users can apply when creating topics.

## API

- **List:** `GET /public/console/v2/topic-template`
- **Get:** `GET /public/console/v2/topic-template/{name}`
- **Create/Update:** `PUT /public/console/v2/topic-template` (upsert — single endpoint for both)
- **Delete:** `DELETE /public/console/v2/topic-template/{name}`

Response wraps the resource in `{ "resource": {...}, "upsertResult": "Created|Updated|NotChanged" }`. Read the `resource` field from `apply.Resource` after Apply, exactly as other console resources do.

## Version and License Gating

- Minimum Console version: `v1.26.0`
- Enterprise plan required: gate with `checkEnterprisePlanRequirement(ctx, client, consoleVersion, topicTemplateMinimumVersion, &resp.Diagnostics)` — same pattern as `console_topic_policy_v1` (which used `v1.43.0` as the enterprise gate; here we use the same min version `v1.26.0` for both the minimum and enterprise gate).

## TF Schema (approved)

```
conduktor_console_topic_template_v2
├── name            required string, RequiresReplace, pattern ^[0-9a-z_\-.]+$
├── labels          optional map[string]string (metadata labels), validation.Labels()
└── spec            required single_nested
    ├── display_name  required string
    ├── description   optional string
    └── defaults      required single_nested
        ├── metadata  optional single_nested
        │   ├── name    optional string  (handlebars pattern validator)
        │   └── labels  optional map[string]string
        └── spec      optional single_nested
            ├── partitions          optional int64, Between(1, 2147483647)
            ├── replication_factor  optional int64, Between(1, 2147483647)
            └── configs             optional map[string]string
```

No oneOf/ADT — this is a plain nested struct. No `RequiresReplace` on any spec field.

## Internal Model (`internal/model/console/topic_template_v2.go`)

```go
const TopicTemplateV2Kind       = "TopicTemplate"
const TopicTemplateV2ApiVersion = "v2"

type TopicTemplateMetadata struct {
    Name   string            `json:"name"`
    Labels map[string]string `json:"labels,omitempty"`
}

type TopicTemplateDefaultsMetadata struct {
    Name   string            `json:"name,omitempty"`
    Labels map[string]string `json:"labels,omitempty"`
}

type TopicTemplateDefaultsSpec struct {
    Partitions        *int64            `json:"partitions,omitempty"`
    ReplicationFactor *int64            `json:"replicationFactor,omitempty"`
    Configs           map[string]string `json:"configs,omitempty"`
}

type TopicTemplateDefaults struct {
    Metadata TopicTemplateDefaultsMetadata `json:"metadata"`
    Spec     TopicTemplateDefaultsSpec     `json:"spec"`
}

type TopicTemplateSpec struct {
    DisplayName string                `json:"displayName"`
    Description string                `json:"description,omitempty"`
    Defaults    TopicTemplateDefaults `json:"defaults"`
}

type TopicTemplateResource struct {
    ApiVersion string                 `json:"apiVersion"`
    Kind       string                 `json:"kind"`
    Metadata   TopicTemplateMetadata  `json:"metadata"`
    Spec       TopicTemplateSpec      `json:"spec"`
}
```

No custom `(Un)MarshalJSON` needed — no discriminated unions.

`Partitions` and `ReplicationFactor` are `*int64` (pointer) so they can be omitted from the JSON payload when not set, avoiding the API interpreting `0` as a value.

## Mapper

`internal/mapper/console_topic_template_v2/`

- `TFToInternalModel`: straightforward field-by-field mapping. Use `schemaUtils.AttrIsSet` guards on the optional `defaults.metadata` and `defaults.spec` nested objects. Use `schemaUtils.MapValueToStringMap` for all `map[string]string` fields.
- `InternalModelToTerraform`: reverse mapping. Populate nullable nested objects using the same `ValueMapFromTypes` / `types.ObjectValue` pattern as existing mappers.
- Mapper unit test: `mapper_test.go` with round-trip from `api.json` fixture.

## Resource (`internal/provider/console_topic_template_v2_resource.go`)

- API path constants:
  ```go
  const topicTemplateV2ApiPath       = "/public/console/v2/topic-template"
  const topicTemplateV2MinimumVersion = "v1.26.0"
  ```
- `Configure`: assert Console mode, check `>= v1.26.0`, call `checkEnterprisePlanRequirement` with the same version.
- `Create`/`Update`: `r.apiClient.Apply(ctx, topicTemplateV2ApiPath, consoleResource)` → read `apply.Resource` via `FromRawJsonInterface`.
- `Read`: `r.apiClient.Describe(ctx, topicTemplateV2ApiPath+"/"+name)` → unmarshal → mapper.
- `Delete`: `r.apiClient.Delete(ctx, client.CONSOLE, topicTemplateV2ApiPath+"/"+name, nil)`.
- `ImportState`: passthrough on `name`.
- Register in `provider.go` `Resources()` slice.

## Test Data (`internal/testdata/console/topic_template_v2/`)

- `api.json` — full API response with all fields populated
- `resource_create.tf` — create with all fields
- `resource_update.tf` — update display_name + configs
- `resource_minimal.tf` — only required fields (name + display_name + empty defaults)

Acceptance test gated behind `test.CheckMinimumVersionRequirement(t, "v1.26.0")` and `test.CheckEnterpriseEnabled(t)`.

## Docs Template (`templates/resources/console_topic_template_v2.md.tmpl`)

Standard template with `simple.tf` and `complex.tf` examples under `examples/resources/conduktor_console_topic_template_v2/`.

## File Checklist

1. `provider_code_spec.json` — add resource entry, `make generate`
2. `internal/model/console/topic_template_v2.go`
3. `internal/mapper/console_topic_template_v2/tf_to_internal_mapper.go`
4. `internal/mapper/console_topic_template_v2/internal_to_tf_mapper.go`
5. `internal/mapper/console_topic_template_v2/mapper_test.go`
6. `internal/testdata/console/topic_template_v2/api.json`
7. `internal/testdata/console/topic_template_v2/resource_create.tf`
8. `internal/testdata/console/topic_template_v2/resource_update.tf`
9. `internal/testdata/console/topic_template_v2/resource_minimal.tf`
10. `internal/provider/console_topic_template_v2_resource.go`
11. `internal/provider/console_topic_template_v2_resource_test.go`
12. `examples/resources/conduktor_console_topic_template_v2/simple.tf`
13. `examples/resources/conduktor_console_topic_template_v2/complex.tf`
14. `templates/resources/console_topic_template_v2.md.tmpl`
15. `internal/provider/provider.go` — register resource
