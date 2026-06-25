# AGENTS.md

This file provides guidance to coding agents when working with code in this repository.

## Overview

Terraform provider for Conduktor Console and Conduktor Gateway. Built with the [terraform-plugin-framework](https://developer.hashicorp.com/terraform/plugin/framework). Resources are split into two `mode`s — `console` and `gateway` — selected at provider configure time; a single Terraform configuration can use both via provider `alias`.

## External references

When adding or changing a resource, consult the upstream Conduktor docs to confirm field names, required vs optional attributes, discriminator values, and version-gated behaviours. The provider mirrors these APIs — don't infer schema shape from sibling TF resources alone.

- Console product/resource reference: <https://docs.conduktor.io/guide/reference/console-reference>
- Gateway product/resource reference: <https://docs.conduktor.io/guide/reference/gateway-reference>
- Console OpenAPI (browseable): <https://developers.conduktor.io/?product=console>
- Gateway v2 OpenAPI (browseable): <https://developers.conduktor.io/?product=gateway&gatewayApiVersion=v2>
- All Console + Gateway OpenAPI versions (raw, fetchable): <https://github.com/conduktor/developers.conduktor.io/tree/main/openapi> — start from [`openapi/manifest.json`](https://github.com/conduktor/developers.conduktor.io/blob/main/openapi/manifest.json) to discover available product/version pairs and their spec file paths. Prefer this source when an agent needs to grep or diff the spec.

## Common commands

```sh
make build                  # go build the provider binary
make go-fmt                 # gofmt
make tf-fmt                 # terraform fmt -recursive on examples/
make go-lint                # golangci-lint (version pinned via GO_LINT_VERSION)
make generate               # runs `go generate ./...` — regenerates schema/, docs/, and tf-fmts examples/
make testacc                # full acceptance test cycle: start docker env, run tests, tear down
make start_test_env         # bring up docker-compose env + seed it via ctl
make test                   # run acceptance tests against an already-running env (TF_ACC=1)
make clean                  # tear down docker-compose env
make deploy-locally         # build + install into ~/.terraform.d/plugins as terraform.local/conduktor/conduktor (VERSION=0.0.1 by default)
```

Run a single acceptance test:
```sh
make start_test_env
TESTARGS='-run TestAccUserV2Resource' make test
make clean
```

Acceptance tests REQUIRE `TF_ACC=1` (set by `make test`) and a live Console+Gateway running on the URLs in `.env`. Some tests are gated behind `CDK_LICENSE` (skipped when absent) and version checks via `internal/test.CheckMinimumVersionRequirement`.

## Architecture

The provider is structured as a pipeline per resource: **Terraform schema model ⇄ internal Go model ⇄ Conduktor API**.

For every resource (e.g. `console_user_v2`), the following parallel files exist:

| Layer | Location | Purpose |
|---|---|---|
| Schema (generated) | `internal/schema/resource_<name>/<name>_resource_gen.go` | Terraform attribute schema + TF model struct. Generated from `provider_code_spec.json` by `tfplugingen-framework`. **Do not edit by hand.** |
| Internal model | `internal/model/{console,gateway}/<name>.go` | Hand-written Go struct mirroring the Conduktor API JSON/YAML payload. |
| Mapper | `internal/mapper/<name>/<name>_resource_mapper.go` | `TFToInternalModel` / `InternalModelToTerraform` conversion. Mapper unit tests live next to the mapper. |
| Resource | `internal/provider/<name>_resource.go` | Terraform CRUD glue: reads plan/state, calls mapper, hits API via `internal/client`. |
| Resource acceptance test | `internal/provider/<name>_resource_test.go` | Uses `testdata/` HCL fixtures and `examples/` to verify create/update/import. |
| Docs template | `templates/resources/<name>.md.tmpl` → renders to `docs/resources/<name>.md` via `tfplugindocs` (`make generate`). |
| Example HCL | `examples/resources/<name>/` — also used by docs generation. |

**Adding/changing a resource attribute** typically requires editing `provider_code_spec.json`, running `make generate`, then updating the mapper + internal model + tests.

### Provider entry points

- `main.go` registers the provider server at `registry.terraform.io/conduktor/conduktor`.
- `internal/provider/provider.go` — `ConduktorProvider` configures the HTTP client based on `mode` and exposes the full resource list in `Resources()`. `ProviderData{Mode, Client}` is passed to each resource via `Configure`.
- Every resource's `Configure` asserts the provider mode (console vs gateway). Resources fail loudly if attached to the wrong mode.
- `internal/provider/provider.go:checkEnterprisePlanRequirement` is the shared helper for resources that are enterprise-only on Console ≥ a given version.
- A package-level `resourceMutex sync.Mutex` exists to serialize resource ops when needed.

### Client

`internal/client/` wraps `resty`. `Make()` picks auth based on `Mode`:
- **Console**: `api_token` if provided, otherwise login with `admin_user`/`admin_password` to mint a bearer.
- **Gateway**: HTTP Basic with `admin_user`/`admin_password`, validated against `/metrics` at startup.

Config resolution is layered: HCL attribute → typed env vars (`CDK_API_TOKEN`, `CDK_BASE_URL`, `CDK_ADMIN_EMAIL`/`PASSWORD`, `CDK_GATEWAY_*`, etc.). See `internal/client/client_config.go`.

### Generic resource

`conduktor_generic` (in `internal/provider/generic_resource.go`) accepts raw YAML manifests and posts them through the Conduktor `ctl` library. It's the escape hatch for kinds without a dedicated typed resource — useful when a new API resource hasn't been promoted to a first-class TF resource yet.

### Custom types and shared validators

- `internal/customtypes/` — custom Terraform types, notably `NormalizedYaml` (used by `conduktor_generic`) which normalizes YAML for plan stability.
- `internal/schema/validation/` — shared validators (labels, non-empty string, permission resource types).
- `internal/schema/default/` — shared schema defaults.
- `internal/planmodifiers/always_use_state.go` — plan modifier that suppresses diffs for computed-only attributes.

## Test environment

`make start_test_env` runs docker-compose (Console + Cortex + Postgres + Redpanda + kafka-connect + Gateway), waits for health, then runs `scripts/setup_test_env.sh` which uses `github.com/conduktor/ctl` (version pinned in `go.mod`) to seed `testdata/init/*.yaml`. Enterprise seed and Console-1.27+ seed are conditionally applied based on `CDK_LICENSE` and `CONDUKTOR_CONSOLE_IMAGE`.

Image tags and credentials are in `.env`. Acceptance tests read `CDK_BASE_URL`, `CDK_ADMIN_EMAIL`, `CDK_ADMIN_PASSWORD` (enforced by `internal/test.TestAccPreCheck`).

Test fixtures live in `internal/testdata/{console,gateway}/<resource>/` — `resource_create.tf`, `resource_update.tf`, `resource_minimal.tf`, `api.json`. Helpers in `internal/test/test_helpers.go` load them via `TestAccTestdata` / `TestAccExample`.

## Codegen notes

- `provider_code_spec.json` is **manually maintained** (despite its codegen-spec format) — edit it, then `make generate`.
- `go generate ./...` runs three steps from `main.go`: schema gen → `terraform fmt` on examples → `tfplugindocs`.
- Generated files under `internal/schema/resource_*/` are excluded from lint exclusions but should not be hand-edited.

### Forked codegen toolchain

The schema generator and its spec format are **forks** of the upstream HashiCorp tools, pulled in via `replace` directives in `go.mod`:

- `github.com/hashicorp/terraform-plugin-codegen-framework` → [`github.com/conduktor/terraform-plugin-codegen-framework@cdk-fix-model-name-conflicts`](https://github.com/conduktor/terraform-plugin-codegen-framework/tree/cdk-fix-model-name-conflicts) — the generator itself.
- `github.com/hashicorp/terraform-plugin-codegen-spec` → [`github.com/conduktor/terraform-plugin-codegen-spec@cdk-add-single_nested-custom_type_name-override`](https://github.com/conduktor/terraform-plugin-codegen-spec/commits/cdk-add-single_nested-custom_type_name-override/) — the JSON schema that `provider_code_spec.json` conforms to.

The key fork-only feature is **`custom_type_name`** on `single_nested` attributes. It lets the same TF attribute name appear in multiple places with different generated Go type names, e.g. a two-level nested `security` object on schema-registry variants where each level/variant would otherwise collide on the generated `SecurityValue` type. If you see `custom_type_name` in `provider_code_spec.json`, that's the fork — upstream codegen would reject or ignore it.

To bump the fork: `go get github.com/conduktor/terraform-plugin-codegen-framework@cdk-fix-model-name-conflicts` (and equivalent for the spec module), then `go mod tidy`.

## Conventions

- Pre-commit hooks (`make setup-hooks`) run `go-fmt`, `tf-fmt`, `go-lint`, plus secret/large-file checks. Linter config in `.golangci.yml` (golangci-lint v2 config style).
- Commit messages follow Conventional Commits (enforced by release-drafter).
- Resource versions (`_v1`, `_v2`) reflect the Conduktor API version, not the TF resource version. New resources should track the latest stable Conduktor API.

## Representing API `oneOf` / `anyOf` (ADT) in the provider

Conduktor APIs frequently use discriminated unions (OpenAPI `oneOf` with a `type` field). The TF type system [doesn't support oneOf](https://github.com/hashicorp/terraform-plugin-codegen-openapi/blob/main/DESIGN.md#known-limitations) — every attribute must have a single static type. **Pick the right pattern before adding a new field:**

### Vocabulary

Three patterns were considered. Only two are in use; the third is documented so it isn't reinvented.

| Pattern | Status | Shape | Use when |
|---|---|---|---|
| **Variant wrapper** | **Default** | `field = { variant_a = {...} }` — `single_nested` wrapper holding one optional sub-attribute per variant. Discriminator `type` is hidden from TF schema; implied by which sub-object is set. | A oneOf appears as a **singleton field** (auth config, flavor, schema registry, scope variant). |
| **Discriminated flat record** | Kept only for lists | `field = { type = "...", optional_field_a = ..., optional_field_b = ... }` — single flat object with all variant fields optional and a required `type` discriminator. | `oneOf` / `anyOf` appears as **elements of a list/set** (e.g. permissions on users/groups), where per-element nesting would explode HCL noise. |
| **Suffix-named siblings** | Rejected | `field_variant_a = {...}`, `field_variant_b = {...}` at the same nesting level. | Never. Pollutes the parent namespace and blocks common fields on a wrapper. |

### Rationale (settled, not up for re-litigation per attribute)

- The **variant wrapper** is the default because a flat record at top level forces every variant-specific field to be optional in the schema (API becomes the only validator), risks field-name collisions across variants, and is hard to document. The wrapper lets each variant keep required fields in the schema, avoids collisions, reads cleanly in HCL, and lets common fields live on the wrapper — at the cost of one "virtual" nesting level that doesn't exist in the API JSON.
- The **discriminated flat record** survives only for unions inside a list: wrapping each list element in a per-variant nested object would force HCL like `permissions[*].topic = {...}` next to `permissions[*].platform = {...}`, which is awful for what is conceptually a flat list of mixed-shape records.
- If a oneOf field becomes a list of oneOfs, switch from variant wrapper to discriminated flat record (and vice versa).

### Variant wrapper — anatomy

Canonical example: `console_kafka_cluster_v2`'s `kafka_flavor` (variants: `aiven`, `confluent`, `gateway`). Also see `internal/model/schema_registry.go` (`ConfluentLike` / `Glue`, with a nested second-level oneOf for `security`).

- **Internal model** (`internal/model/console/kafka_cluster_v2.go:34-94`):
  - Wrapper struct holds one **pointer per variant**: `type KafkaFlavor struct { Aiven *Aiven; Confluent *Confluent; Gateway *Gateway }`. Exactly one is non-nil by invariant.
  - Variant constants typed as `KafkaFlavorType` string enum.
  - Each variant struct has a literal `Type string` JSON field set to the discriminator value.
  - Custom `UnmarshalJSON` reads `type` via `model.Discriminable` (`internal/model/utils.go:9`) and switches to the matching variant.
  - Custom `MarshalJSON` walks the non-nil pointer and delegates `json.Marshal` to that variant struct (which carries the `type` field, so the discriminator re-emerges).
- **TF schema** (`internal/schema/resource_console_kafka_cluster_v2/...gen.go`): wrapper is `SingleNestedAttribute`; each variant is an optional `SingleNestedAttribute` with no `type` attribute.
- **Mapper** (`internal/mapper/console_kafka_cluster_v2/tf_to_internal_mapper.go:60-119` and `internal_to_tf_mapper.go`):
  - TF → internal: guard with `if r.IsNull() { return nil, nil }`, then for each variant `if schemaUtils.AttrIsSet(wrapper.Variant) { build &console.Variant{Type: string(console.VARIANT), ...} }`, return the wrapper with the populated pointers.
  - Internal → TF: mirror — bail out if all pointers are nil, otherwise populate exactly one entry in `valuesMap`.
- **Helper**: `schemaUtils.AttrIsSet(attr)` (`internal/schema/schema_utils.go:261`) = `!IsNull && !IsUnknown`. Use this everywhere instead of inlining the check.

### Discriminated flat record — anatomy

Canonical example: `console_user_v2` / `console_group_v2` permissions.

- **Internal model** (`internal/model/permissions.go:5-13`): single flat `Permission` struct with all variant fields as `omitempty` strings, plus a required `ResourceType` discriminator. **No custom (Un)MarshalJSON** — discriminator is just another field.
- **TF schema**: each variant field is `Optional`; `resource_type` is required.
- **Validator**: `internal/schema/validation/permission_resource_type.go` enforces "field X may only be set when `resource_type` is in {…}" via `permissionFieldsByResourceType`. Wire it into the schema with `validators.Set: []validator.Set{ PermissionResourceType() }`.
- **Mapper**: shared helpers `schemaUtils.SetValueToPermissionArray` / `PermissionArrayToSetValue` (`internal/schema/schema_utils.go:115-232`) — generalised across `users`/`groups` via a `Resource` enum, so per-resource mappers stay short (`internal/mapper/console_user_v2/user_v2_resource_mapper.go:15-29`).
- **API-strips-fields gotcha** (`internal/model/permissions.go:15-89`): the Console API drops optional fields it deems irrelevant for a given `resource_type` (e.g. it won't echo back `kafka_connect` on a `TOPIC` permission). Naïvely using the response as the new TF state triggers *"Provider produced inconsistent result after apply"*. Discriminated-flat-record resources MUST run the API response through `model.MergeWithPlannedPermissions(planned, response)` before saving state — it re-attaches planned-but-stripped fields by matching on `(resourceType, permissions, non-empty fields)`.

### Adding a new variant to an existing variant wrapper

1. Add the variant struct to `internal/model/<domain>/<resource>.go` with a `Type string` field and a typed constant.
2. Extend the wrapper struct with a `*NewVariant` pointer.
3. Extend `UnmarshalJSON`'s `switch disc.Type` and `MarshalJSON`'s pointer chain.
4. Edit `provider_code_spec.json` to add a new optional `single_nested` sub-attribute under the wrapper, then `make generate`.
5. Add a TF→internal branch in the mapper (`AttrIsSet(wrapper.NewVariant)` → build struct with `Type: string(NEW_VARIANT)`) and an internal→TF branch (`if r.NewVariant != nil { ... }`).
6. Add fixtures under `internal/testdata/<domain>/<resource>/` and a mapper unit test in `internal/mapper/<resource>/mapper_test.go`.

### Notes / pitfalls

- Never expose the discriminator (`type`) as a TF attribute in a variant wrapper — it would create two sources of truth and break the wrapper invariant.
- In a variant wrapper's `MarshalJSON`, returning an error on "no pointer set" is the convention — silent empty marshals will produce malformed API payloads.
- `gateway_interceptor_v2.GatewayInterceptorScope` (`internal/model/gateway/interceptor_v2.go:15-19`) is **not** a oneOf — it's a plain object with mutually contextual optional fields. Don't model new oneOfs on it.

## Adding a new resource (or data source)

The provider's per-resource layout (schema ⇄ model ⇄ mapper ⇄ resource) is described above; this is the end-to-end checklist for promoting a new Conduktor API kind to a first-class TF resource. Follow the order — later steps rely on artifacts from earlier ones.

1. **Define the TF schema spec.** Add a new resource entry to `provider_code_spec.json`, modelled on the Conduktor OpenAPI definition for the kind. Map attribute types (string/int/bool/list/map/single_nested), mark `Required` vs `Optional` vs `Computed`, pick the right oneOf pattern (see *Representing API `oneOf` / `anyOf`* above). Cross-reference an existing resource of similar shape — don't invent attribute layouts.
   1. Run `make generate` to regenerate `internal/schema/resource_<name>/<name>_resource_gen.go`. This file is generated; never hand-edit it.
2. **Add the internal API model.** Create `internal/model/{console,gateway}/<name>.go` mirroring the Conduktor API JSON/YAML payload (`Kind`, `ApiVersion`, `Metadata`, `Spec`). Add `NewXxxResource(...)`, `ToClientResource`, `FromClientResource`, `FromRawJsonInterface`, `NewXxxResourceFromClientResource` helpers — copy the shape from a sibling model. Custom `(Un)MarshalJSON` is only required for variant-wrapper oneOfs (see ADT section).
3. **Write the mapper.** Create `internal/mapper/<name>/<name>_resource_mapper.go` with `TFToInternalModel(ctx, *schema.<Name>Model) (model.<Name>Resource, error)` and `InternalModelToTerraform(ctx, *model.<Name>Resource) (schema.<Name>Model, error)`. Use `schemaUtils.AttrIsSet`, `MapValueToStringMap`, `SetValueToStringArray` and friends from `internal/schema/schema_utils.go` rather than re-implementing them.
   1. Add a mapper unit test `<name>_resource_mapper_test.go` next to the mapper, asserting **full round-trip**: load `api.json` fixture → ctlresource → internal model → TF model → internal model → JSON, and assert byte-equivalence (or `cmp.Equal` with the original). The fixture lives under `internal/testdata/<domain>/<name>/api.json`. Model `TestUserV2ModelMapping` in `internal/mapper/console_user_v2/user_v2_resource_mapper_test.go` for the canonical pattern.
4. **Implement the resource.** Create `internal/provider/<name>_resource.go`. Copy the skeleton from a sibling resource of the same mode (console/gateway); the boilerplate is:
   - `Metadata` → resource type name `conduktor_<name>`.
   - `Schema` → call the generated `<Name>ResourceSchema(ctx)`.
   - `Configure` → assert provider mode and stash `*client.Client`. Fail loudly if attached to the wrong mode.
   - `Create` / `Read` / `Update` / `Delete` → use the mapper, then `r.apiClient.Apply` / `Describe` / `Delete` against the API endpoint paths (define `xxxApiPutPath` / `xxxApiGetPath` constants at the top of the file).
   - `ImportState` → key by `name` (or composite key for nested resources — see `console_topic_v2_resource.go` for `cluster/topic_name` import).
   - Register the constructor in `internal/provider/provider.go`'s `Resources()` slice.
   - If the resource is enterprise-only on Console ≥ a version, gate it with `checkEnterprisePlanRequirement`.
5. **Write the acceptance test.** Create `internal/provider/<name>_resource_test.go` covering: Create + Read (assert key attrs), Import (`ImportStateVerify: true`, identifier set), Update + Read (assert changed attrs), Delete (automatic). Fixtures live in `internal/testdata/<domain>/<name>/` as `resource_create.tf`, `resource_update.tf`, and `resource_minimal.tf` — loaded via `test.TestAccTestdata`. Gate tests behind `test.CheckEnterpriseEnabled` or `test.CheckMinimumVersionRequirement` if needed.
6. **Add user-facing examples.** Under `examples/resources/conduktor_<name>/`, add at least `simple.tf` and ideally `complex.tf` (or domain-specific snippets). These files are used by both the docs renderer (step 7) and the acceptance tests when convenient — load them with `test.TestAccExample(...)`.
7. **Write the docs template.** Add `templates/resources/<name>.md.tmpl` with frontmatter (`page_title`, `subcategory`, `description`), a short prose intro, one or more `{{tffile "examples/resources/conduktor_<name>/<snippet>.tf"}}` directives per example, and `{{ .SchemaMarkdown | trimspace }}` at the bottom to splice in the generated attribute reference.
   1. Run `make generate` to regenerate `docs/resources/<name>.md` via `tfplugindocs`. Don't hand-edit the generated file.
8. **Update the README.** If `README.md` keeps an inline list of supported resources, add the new entry; otherwise no change. The canonical user-facing list lives on docs.conduktor.io and in the Terraform Registry, not in this repo.
9. **Pre-flight before opening the PR.** Run, in order: `make go-fmt`, `make tf-fmt`, `make go-lint`, the mapper unit tests (`go test ./internal/mapper/<name>/...`), and a full `make testacc` (or `make start_test_env` + targeted `TESTARGS='-run TestAcc<Name>Resource' make test` + `make clean`). The pre-commit hook (`make setup-hooks`) will run fmt + lint + secret/large-file checks on commit — fix any failures rather than bypassing.
