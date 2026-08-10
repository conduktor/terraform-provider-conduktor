# connector_v2: Fix `config` Required + Add `initial_state` Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix `spec.config` to be `required` (matching the API) and add the new `spec.initial_state` optional string field (added in Console 1.46.0, values: `RUNNING`, `PAUSED`, `STOPPED`).

**Architecture:** Changes flow through three layers in order: internal model → `provider_code_spec.json` + `make generate` (regenerates schema) → mapper → tests. The schema gen file is never hand-edited.

**Tech Stack:** Go, Terraform Plugin Framework, `provider_code_spec.json` (codegen spec), `make generate` (runs `go generate ./...`), `go test`

---

## File Map

| File | Action | Responsibility |
|---|---|---|
| `provider_code_spec.json` | Modify | `config` → `required`; add `initial_state` optional string with OneOf validator |
| `internal/schema/resource_console_connector_v2/console_connector_v2_resource_gen.go` | Regenerated | Never hand-edit; updated by `make generate` |
| `internal/model/console/connector_v2.go` | Modify | Add `InitialState string` with `omitempty` to `ConnectorConsoleSpec` |
| `internal/mapper/console_connector_v2/console_connector_v2_resource_mapper.go` | Modify | Map `initial_state` ↔ `InitialState` in both directions; update `NewSpecValue` calls |
| `internal/mapper/console_connector_v2/console_connector_v2_resource_mapper_test.go` | Modify | Assert `initial_state` round-trip |
| `internal/testdata/console/connector_v2/api.json` | Modify | Add `"initialState": "RUNNING"` to spec |
| `internal/testdata/console/connector_v2/resource_create.tf` | Modify | Add `initial_state = "RUNNING"` |
| `internal/testdata/console/connector_v2/resource_update.tf` | Modify | Add `initial_state = "PAUSED"` |
| `internal/provider/console_connector_v2_resource_test.go` | Modify | Assert `spec.initial_state` in create/update steps |

---

## Task 1: Update internal model to add `InitialState`

**Files:**
- Modify: `internal/model/console/connector_v2.go`

- [ ] **Step 1: Add `InitialState` to `ConnectorConsoleSpec`**

In `internal/model/console/connector_v2.go`, change `ConnectorConsoleSpec` from:

```go
type ConnectorConsoleSpec struct {
	Config map[string]string `json:"config"`
}
```

to:

```go
type ConnectorConsoleSpec struct {
	Config       map[string]string `json:"config"`
	InitialState string            `json:"initialState,omitempty"`
}
```

- [ ] **Step 2: Verify it compiles**

```bash
cd /Users/ben/conduktor/terraform-provider-conduktor && go build ./...
```

Expected: no output (clean build).

- [ ] **Step 3: Commit**

```bash
git add internal/model/console/connector_v2.go
git commit -m "feat(connector-v2): add InitialState field to ConnectorConsoleSpec"
```

---

## Task 2: Update `provider_code_spec.json` and regenerate schema

**Files:**
- Modify: `provider_code_spec.json`
- Regenerated: `internal/schema/resource_console_connector_v2/console_connector_v2_resource_gen.go`

- [ ] **Step 1: Change `config` from `computed_optional` to `required`**

Find the `spec` single_nested block for `console_connector_v2` in `provider_code_spec.json` (around line 1191). Change the `config` map attribute's `computed_optional_required` field from `"computed_optional"` to `"required"`:

```json
{
  "name": "config",
  "map": {
    "description": "Must be valid Kafka Connect Connector configs",
    "computed_optional_required": "required",
    "element_type": {
      "string": {}
    }
  }
}
```

- [ ] **Step 2: Add `initial_state` attribute inside `spec.attributes`**

In the same `spec` single_nested block (around line 1195), add `initial_state` after the `config` attribute so the `attributes` array looks like:

```json
"attributes": [
  {
    "name": "config",
    "map": {
      "description": "Must be valid Kafka Connect Connector configs",
      "computed_optional_required": "required",
      "element_type": {
        "string": {}
      }
    }
  },
  {
    "name": "initial_state",
    "string": {
      "description": "Initial state of the connector after creation or update. Valid values are RUNNING, PAUSED, STOPPED. NOTE: this field has been introduced with Console 1.46.0 and it will not work with previous versions",
      "computed_optional_required": "optional",
      "validators": [
        {
          "custom": {
            "imports": [
              {
                "path": "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
              }
            ],
            "schema_definition": "stringvalidator.OneOf(\"RUNNING\", \"PAUSED\", \"STOPPED\")"
          }
        }
      ]
    }
  }
]
```

- [ ] **Step 3: Run `make generate` to regenerate the schema**

```bash
cd /Users/ben/conduktor/terraform-provider-conduktor && make generate
```

Expected: runs without error. The file `internal/schema/resource_console_connector_v2/console_connector_v2_resource_gen.go` is updated.

- [ ] **Step 4: Verify the generated schema has the changes**

```bash
grep -n "initial_state\|InitialState\|required\|Required" internal/schema/resource_console_connector_v2/console_connector_v2_resource_gen.go | head -30
```

Expected: you see `initial_state` as an optional string attribute, `config` marked as `Required: true`, and an `InitialState` field on `SpecValue`.

- [ ] **Step 5: Verify it compiles**

```bash
go build ./...
```

Expected: no output.

- [ ] **Step 6: Commit**

```bash
git add provider_code_spec.json internal/schema/resource_console_connector_v2/console_connector_v2_resource_gen.go
git commit -m "feat(connector-v2): make config required and add initial_state to spec schema"
```

---

## Task 3: Update the mapper

**Files:**
- Modify: `internal/mapper/console_connector_v2/console_connector_v2_resource_mapper.go`

The mapper currently builds `SpecValue` with only `config`. After `make generate`, `SpecValue` will also have an `InitialState basetypes.StringValue` field and the `NewSpecValue` constructor will require it.

- [ ] **Step 1: Update `TFToInternalModel` to map `initial_state`**

In `TFToInternalModel`, update the returned `ConnectorConsoleSpec` to include `InitialState`:

```go
return console.NewConnectorConsoleResource(
    console.ConnectorConsoleMetadata{
        Name:           r.Name.ValueString(),
        Cluster:        r.Cluster.ValueString(),
        ConnectCluster: r.ConnectCluster.ValueString(),
        Labels:         mapper.MergeLabels(managedLabels, userLabels),
        Description:    r.Description.ValueString(),
        AutoRestart:    autoRestart,
    },
    console.ConnectorConsoleSpec{
        Config:       config,
        InitialState: r.Spec.InitialState.ValueString(),
    },
), nil
```

Note: `ValueString()` on a null `StringValue` returns `""`, which is correct — the `omitempty` tag on `InitialState` will then omit it from the JSON payload.

- [ ] **Step 2: Update `InternalModelToTerraform` to include `initial_state` in `NewSpecValue`**

The `NewSpecValue` constructor will now require `initial_state` in both the `attributeTypes` and `attributes` maps. Update the call:

```go
var initialState types.String
if r.Spec.InitialState == "" {
    initialState = types.StringNull()
} else {
    initialState = schema.NewStringValue(r.Spec.InitialState)
}

specValue, diag := connector.NewSpecValue(
    map[string]attr.Type{
        "config":        config.Type(ctx),
        "initial_state": basetypes.StringType{},
    },
    map[string]attr.Value{
        "config":        config,
        "initial_state": initialState,
    },
)
```

Also add the `basetypes` import if it is not already present (it is already imported in this file for `autoRestartInternalToTerraform`).

- [ ] **Step 3: Verify it compiles**

```bash
cd /Users/ben/conduktor/terraform-provider-conduktor && go build ./...
```

Expected: no output.

- [ ] **Step 4: Commit**

```bash
git add internal/mapper/console_connector_v2/console_connector_v2_resource_mapper.go
git commit -m "feat(connector-v2): map initial_state in TFToInternalModel and InternalModelToTerraform"
```

---

## Task 4: Update mapper unit test fixture and assertions

**Files:**
- Modify: `internal/testdata/console/connector_v2/api.json`
- Modify: `internal/mapper/console_connector_v2/console_connector_v2_resource_mapper_test.go`

- [ ] **Step 1: Add `initialState` to the test fixture**

In `internal/testdata/console/connector_v2/api.json`, add `"initialState": "RUNNING"` to the `spec` object:

```json
{
  "kind": "Connector",
  "apiVersion": "v2",
  "metadata": {
    "name": "connector",
    "cluster": "cluster",
    "connectCluster": "connect",
    "labels": {
      "conduktor.io/application": "test-app",
      "conduktor.io/application-instance": "test-app-instance",
      "kind": "connector",
      "data-criticality": "C0",
      "environment": "prod",
      "team": "analytics"
    },
    "description": "This is a connector",
    "autoRestart": {
      "enabled": true,
      "frequencySeconds": 500
    }
  },
  "spec": {
    "config": {
      "connector.class": "io.connect.jdbc.JdbcSourceConnector",
      "tasks.max": "1",
      "topic": "click.pageviews",
      "connection.url": "jdbc:mysql://127.0.0.1:3306/sample?verifyServerCertificate=false&useSSL=true&requireSSL=true",
      "consumer.override.sasl.jaas.config": "o.a.k.s.s.ScramLoginModule required username='<user>' password='<password>';"
    },
    "initialState": "RUNNING"
  }
}
```

- [ ] **Step 2: Add `initialState` assertions to the mapper test**

In `TestConnectorV2ModelMapping` in `internal/mapper/console_connector_v2/console_connector_v2_resource_mapper_test.go`:

After the existing `assert.Equal(t, map[string]any{"config": config}, ctlResource.Spec)` line, update the spec assertion to include `initialState`:

```go
assert.Equal(t, map[string]any{"config": config, "initialState": "RUNNING"}, ctlResource.Spec)
```

After `assert.Equal(t, "io.connect.jdbc.JdbcSourceConnector", internal.Spec.Config["connector.class"])`, add:

```go
assert.Equal(t, "RUNNING", internal.Spec.InitialState)
```

After `assert.Equal(t, types.BoolValue(true), tfModel.AutoRestart.Enabled)`, add:

```go
assert.Equal(t, types.StringValue("RUNNING"), tfModel.Spec.InitialState)
```

After `assert.Equal(t, "io.connect.jdbc.JdbcSourceConnector", internal2.Spec.Config["connector.class"])`, add:

```go
assert.Equal(t, "RUNNING", internal2.Spec.InitialState)
```

- [ ] **Step 3: Run the mapper unit test and verify it passes**

```bash
cd /Users/ben/conduktor/terraform-provider-conduktor && go test ./internal/mapper/console_connector_v2/... -v -run TestConnectorV2ModelMapping
```

Expected: `PASS` with all assertions green.

- [ ] **Step 4: Commit**

```bash
git add internal/testdata/console/connector_v2/api.json \
        internal/mapper/console_connector_v2/console_connector_v2_resource_mapper_test.go
git commit -m "test(connector-v2): add initial_state to mapper fixture and unit test assertions"
```

---

## Task 5: Update acceptance test fixtures and test assertions

**Files:**
- Modify: `internal/testdata/console/connector_v2/resource_create.tf`
- Modify: `internal/testdata/console/connector_v2/resource_update.tf`
- Modify: `internal/provider/console_connector_v2_resource_test.go`

- [ ] **Step 1: Add `initial_state` to `resource_create.tf`**

Replace the contents of `internal/testdata/console/connector_v2/resource_create.tf`:

```hcl
resource "conduktor_console_connector_v2" "test" {
  name            = "connector-test"
  cluster         = "kafka-cluster"
  connect_cluster = "kafka-connect"
  labels = {
    key1 = "value1"
  }
  description = "description"
  auto_restart = {
    enabled           = true
    frequency_seconds = 800
  }
  spec = {
    config = {
      "connector.class" = "org.apache.kafka.connect.tools.MockSourceConnector"
      "tasks.max"       = "1"
      "topic"           = "click.pageviews"
      "file"            = "/etc/kafka/consumer.properties"
    }
    initial_state = "RUNNING"
  }
}
```

- [ ] **Step 2: Add `initial_state` to `resource_update.tf`**

Replace the contents of `internal/testdata/console/connector_v2/resource_update.tf`:

```hcl
resource "conduktor_console_connector_v2" "test" {
  name            = "connector-test"
  cluster         = "kafka-cluster"
  connect_cluster = "kafka-connect"
  labels = {
    "env" = "test"
    "sec" = "C1"
  }
  description = "description update"
  auto_restart = {
    enabled = false
  }
  spec = {
    config = {
      "connector.class" = "org.apache.kafka.connect.tools.MockSourceConnector"
      "tasks.max"       = "2"
      "topic"           = "click.pageviews.new"
      "file"            = "/etc/kafka/producer.properties"
    }
    initial_state = "PAUSED"
  }
}
```

- [ ] **Step 3: Add `initial_state` assertions to the acceptance test**

In `internal/provider/console_connector_v2_resource_test.go`, in `TestAccConnectorV2Resource`:

In the Create step checks (after the existing `spec.config.*` assertions), add:
```go
resource.TestCheckResourceAttr(resourceRef, "spec.initial_state", "RUNNING"),
```

In the Update step checks (after the existing `spec.config.*` assertions), add:
```go
resource.TestCheckResourceAttr(resourceRef, "spec.initial_state", "PAUSED"),
```

Also add a `test.CheckMinimumVersionRequirement` guard at the top of `TestAccConnectorV2Resource` for the `initial_state` field — or leave it gated only by the existing `connectorMininumRecommendedVersion` check since this is additive. Since `initial_state` is optional and doesn't cause errors on older Console (the API just ignores the field if unsupported), no extra version gate is needed. The description on the attribute documents the 1.46.0 minimum.

- [ ] **Step 4: Run `go vet` to catch any issues before the full acc test cycle**

```bash
cd /Users/ben/conduktor/terraform-provider-conduktor && go vet ./...
```

Expected: no output.

- [ ] **Step 5: Commit**

```bash
git add internal/testdata/console/connector_v2/resource_create.tf \
        internal/testdata/console/connector_v2/resource_update.tf \
        internal/provider/console_connector_v2_resource_test.go
git commit -m "test(connector-v2): add initial_state to acceptance test fixtures and assertions"
```

---

## Task 6: Lint, fmt, and final verification

**Files:** none (no code changes)

- [ ] **Step 1: Run `make go-fmt`**

```bash
cd /Users/ben/conduktor/terraform-provider-conduktor && make go-fmt
```

Expected: no output (nothing to reformat).

- [ ] **Step 2: Run `make tf-fmt`**

```bash
make tf-fmt
```

Expected: no output.

- [ ] **Step 3: Run `make go-lint`**

```bash
make go-lint
```

Expected: no lint errors. If there are errors, fix them (they are likely unused imports or similar).

- [ ] **Step 4: Run the full mapper unit test suite**

```bash
go test ./internal/mapper/... -v
```

Expected: all tests pass.

- [ ] **Step 5: (Optional, requires live env) Run the targeted acceptance test**

```bash
make start_test_env
TESTARGS='-run TestAccConnectorV2Resource' make test
make clean
```

Expected: all steps pass, `initial_state` attribute is correctly applied and read back.

- [ ] **Step 6: Commit any lint fixes (if needed)**

```bash
git add -p
git commit -m "fix(connector-v2): lint and format fixes"
```

---

## Self-Review Notes

- `config` is changed to `required` in the spec JSON — this matches the API constraint and is consistent with how the field already worked in practice (applying without it would fail at the API).
- `initial_state` is `omitempty` in the internal model so when a user doesn't set it, the field is absent from the API payload entirely (the API's default behaviour applies).
- The `NewSpecValue` call in the mapper must pass **both** `config` and `initial_state` after regeneration — the generated constructor validates exact key presence.
- The `SpecValue` struct in the generated file will have `InitialState basetypes.StringValue` — the mapper test assertion `tfModel.Spec.InitialState` references this field directly.
- No `MarshalJSON`/`UnmarshalJSON` customisation needed — `initialState` is a plain scalar field, not a oneOf.
- `resource_minimal.tf` is intentionally left without `initial_state` to exercise the optional path.
