# application_v1 policy_ref Field Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add the optional `policy_ref` field (set of strings) to `console_application_v1`, mirroring the identical field already present on `console_application_instance_v1`.

**Architecture:** The field maps directly to `spec.policyRef` in the Conduktor API (array of string, unique items, pattern `^[0-9a-z_\-.]+$`, added in Console 1.35.0). It follows the existing provider convention: a `set` attribute in the TF schema → `[]string` in the internal model → `schema.SetValueToStringArray` / `schema.StringArrayToSetValue` in the mapper. No variant wrapper needed — this is a plain optional list.

**Tech Stack:** Go, terraform-plugin-framework, `provider_code_spec.json` (codegen spec), `make generate` (schema + docs regeneration), golangci-lint.

---

### Task 1: Add `PolicyRef` to the internal model

**Files:**
- Modify: `internal/model/console/application_v1.go`

- [ ] **Step 1: Add the field to `ApplicationConsoleSpec`**

In `internal/model/console/application_v1.go`, update `ApplicationConsoleSpec`:

```go
type ApplicationConsoleSpec struct {
	Description string   `json:"description,omitempty"`
	Title       string   `json:"title"`
	Owner       string   `json:"owner"`
	PolicyRef   []string `json:"policyRef,omitempty"`
}
```

- [ ] **Step 2: Verify it compiles**

```bash
make build
```
Expected: exits 0 with no errors.

- [ ] **Step 3: Commit**

```bash
git add internal/model/console/application_v1.go
git commit -m "feat(application_v1): add PolicyRef field to internal model"
```

---

### Task 2: Add `policy_ref` to the TF schema spec and regenerate

**Files:**
- Modify: `provider_code_spec.json`
- Regenerate (do not hand-edit): `internal/schema/resource_console_application_v1/console_application_v1_resource_gen.go`

- [ ] **Step 1: Add the attribute to `provider_code_spec.json`**

Locate the `console_application_v1` entry's `spec` → `single_nested` → `attributes` array. Add the new attribute after `"owner"`:

```json
{
  "name": "policy_ref",
  "set": {
    "description": "References to resource policies to apply to this application. NOTE: this field has been introduced with Console 1.35.0 and will not work with previous versions",
    "computed_optional_required": "optional",
    "element_type": {
      "string": {}
    },
    "validators": [
      {
        "custom": {
          "imports": [
            {
              "path": "regexp"
            },
            {
              "path": "github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
            },
            {
              "path": "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
            },
            {
              "path": "github.com/hashicorp/terraform-plugin-framework/schema/validator"
            }
          ],
          "schema_definition": "setvalidator.ValueStringsAre(stringvalidator.RegexMatches(regexp.MustCompile(\"^[0-9a-z_\\\\-.]+$\"), \"policy name must match ^[0-9a-z_\\\\-.]+$\"))"
        }
      }
    ]
  }
}
```

- [ ] **Step 2: Regenerate schema and docs**

```bash
make generate
```
Expected: exits 0. Check that `internal/schema/resource_console_application_v1/console_application_v1_resource_gen.go` now includes `PolicyRef` in the `SpecValue` struct and `SpecType` attribute map.

- [ ] **Step 3: Verify it compiles**

```bash
make build
```
Expected: exits 0.

- [ ] **Step 4: Commit**

```bash
git add provider_code_spec.json internal/schema/resource_console_application_v1/ examples/ docs/
git commit -m "feat(application_v1): add policy_ref to TF schema spec and regenerate"
```

---

### Task 3: Update the mapper

**Files:**
- Modify: `internal/mapper/console_application_v1/console_application_v1_resource_mapper.go`

- [ ] **Step 1: Update `TFToInternalModel`**

Replace the existing function body with:

```go
func TFToInternalModel(ctx context.Context, r *apps.ConsoleApplicationV1Model) (console.ApplicationConsoleResource, error) {
	policyRef, diag := schema.SetValueToStringArray(ctx, r.Spec.PolicyRef)
	if diag.HasError() {
		return console.ApplicationConsoleResource{}, mapper.WrapDiagError(diag, "policy_ref", mapper.FromTerraform)
	}

	return console.NewApplicationConsoleResource(
		r.Name.ValueString(),
		console.ApplicationConsoleSpec{
			Title:       r.Spec.Title.ValueString(),
			Description: r.Spec.Description.ValueString(),
			Owner:       r.Spec.Owner.ValueString(),
			PolicyRef:   policyRef,
		},
	), nil
}
```

- [ ] **Step 2: Update `InternalModelToTerraform`**

Replace the existing function body with:

```go
func InternalModelToTerraform(ctx context.Context, r *console.ApplicationConsoleResource) (apps.ConsoleApplicationV1Model, error) {
	policyRef := basetypes.NewSetNull(basetypes.StringType{})
	if r.Spec.PolicyRef != nil {
		var diag diag.Diagnostics
		policyRef, diag = schema.StringArrayToSetValue(r.Spec.PolicyRef)
		if diag.HasError() {
			return apps.ConsoleApplicationV1Model{}, mapper.WrapDiagError(diag, "policy_ref", mapper.IntoTerraform)
		}
	}

	specValue, diag := apps.NewSpecValue(
		map[string]attr.Type{
			"title":       basetypes.StringType{},
			"description": basetypes.StringType{},
			"owner":       basetypes.StringType{},
			"policy_ref":  policyRef.Type(ctx),
		},
		map[string]attr.Value{
			"title":       schema.NewStringValue(r.Spec.Title),
			"description": schema.NewStringValue(r.Spec.Description),
			"owner":       schema.NewStringValue(r.Spec.Owner),
			"policy_ref":  policyRef,
		},
	)
	if diag.HasError() {
		return apps.ConsoleApplicationV1Model{}, mapper.WrapDiagError(diag, "spec", mapper.IntoTerraform)
	}

	return apps.ConsoleApplicationV1Model{
		Name: schema.NewStringValue(r.Metadata.Name),
		Spec: specValue,
	}, nil
}
```

Add the required imports at the top of the file:

```go
import (
	"context"

	mapper "github.com/conduktor/terraform-provider-conduktor/internal/mapper"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	schema "github.com/conduktor/terraform-provider-conduktor/internal/schema"
	apps "github.com/conduktor/terraform-provider-conduktor/internal/schema/resource_console_application_v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)
```

- [ ] **Step 3: Verify it compiles**

```bash
make build
```
Expected: exits 0.

- [ ] **Step 4: Commit**

```bash
git add internal/mapper/console_application_v1/console_application_v1_resource_mapper.go
git commit -m "feat(application_v1): map policy_ref between TF and internal model"
```

---

### Task 4: Update mapper unit test and fixtures

**Files:**
- Modify: `internal/testdata/console/application_v1/api.json`
- Modify: `internal/mapper/console_application_v1/console_application_v1_resource_mapper_test.go`

- [ ] **Step 1: Update `api.json` fixture to include `policyRef`**

Replace `internal/testdata/console/application_v1/api.json` with:

```json
{
  "kind": "Application",
  "apiVersion": "v1",
  "metadata": {
    "name": "application"
  },
  "spec": {
    "title": "application title",
    "description": "application description",
    "owner": "application owner",
    "policyRef": ["my-policy"]
  }
}
```

- [ ] **Step 2: Update the mapper test to assert `PolicyRef`**

In `internal/mapper/console_application_v1/console_application_v1_resource_mapper_test.go`, extend `TestApplicationV1ModelMapping`:

After the existing `assert.Equal(t, "application owner", internal.Spec.Owner)` line, add:
```go
assert.Equal(t, []string{"my-policy"}, internal.Spec.PolicyRef)
```

After the existing `assert.Equal(t, types.StringValue("application owner"), tfModel.Spec.Owner)` line, add:
```go
// policy_ref is a set — check it contains the expected element
policyRefElems := tfModel.Spec.PolicyRef.Elements()
assert.Len(t, policyRefElems, 1)
```

After the existing `assert.Equal(t, "application owner", internal2.Spec.Owner)` line, add:
```go
assert.Equal(t, []string{"my-policy"}, internal2.Spec.PolicyRef)
```

- [ ] **Step 3: Run the mapper unit test and verify it passes**

```bash
go test ./internal/mapper/console_application_v1/...
```
Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add internal/testdata/console/application_v1/api.json \
        internal/mapper/console_application_v1/console_application_v1_resource_mapper_test.go
git commit -m "test(application_v1): add policy_ref assertions to mapper unit test"
```

---

### Task 5: Update acceptance test fixtures and example

**Files:**
- Modify: `internal/testdata/console/application_v1/resource_update.tf`
- Modify: `examples/resources/conduktor_console_application_v1/complex.tf`
- Modify: `internal/provider/console_application_v1_resource_test.go`

- [ ] **Step 1: Add `policy_ref` to the update fixture**

Replace `internal/testdata/console/application_v1/resource_update.tf` with:

```hcl
resource "conduktor_console_application_v1" "test" {
  name = "my-application"
  spec = {
    title       = "My Application"
    description = "My Application description"
    owner       = "admin"
    policy_ref  = ["my-policy"]
  }
}
```

- [ ] **Step 2: Update the complex example**

Replace `examples/resources/conduktor_console_application_v1/complex.tf` with:

```hcl
resource "conduktor_console_application_v1" "example" {
  name = "complex-app"
  spec = {
    title       = "Complex Application"
    description = "Complex Application description"
    owner       = "admin"
    policy_ref  = ["my-policy"]
  }
}
```

- [ ] **Step 3: Add acceptance test check for `policy_ref`**

In `internal/provider/console_application_v1_resource_test.go`, in `TestAccApplicationV1Resource`, update the Update step's `Check` block to add:

```go
resource.TestCheckResourceAttr(resourceRef, "spec.policy_ref.#", "1"),
```

Also add to `TestAccApplicationV1ExampleResource`'s check:

```go
resource.TestCheckResourceAttr("conduktor_console_application_v1.example", "spec.policy_ref.#", "1"),
```

- [ ] **Step 4: Verify lint passes**

```bash
make go-lint
make tf-fmt
```
Expected: exits 0.

- [ ] **Step 5: Commit**

```bash
git add internal/testdata/console/application_v1/resource_update.tf \
        examples/resources/conduktor_console_application_v1/complex.tf \
        internal/provider/console_application_v1_resource_test.go
git commit -m "test(application_v1): update fixtures and acc tests for policy_ref"
```

---

### Task 6: Full lint and build verification

- [ ] **Step 1: Run all linters and formatters**

```bash
make go-fmt && make tf-fmt && make go-lint
```
Expected: all exit 0.

- [ ] **Step 2: Run mapper unit tests**

```bash
go test ./internal/mapper/console_application_v1/...
```
Expected: PASS.

- [ ] **Step 3: Build**

```bash
make build
```
Expected: exits 0.
