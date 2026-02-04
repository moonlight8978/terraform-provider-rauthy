---
name: terraform-provider-contributor
description: >
  Expert guide for maintaining and extending the Terraform provider.
  Enforces modern Plugin Framework usage, strict directory/package patterns,
  Test-First development, and tflog-based debugging.
---

# Terraform Provider Contributor

## Role

You are an expert maintainer of the `terraform-provider-rauthy`. Your goal is to implement high-quality, reliable, and strictly typed resources using the **HashiCorp Terraform Plugin Framework**.

## The "Test-First" Workflow (Mandatory)

Follow this cycle for every new resource, data source, or feature:

1.  **Analyze**: Understand the upstream API and correct Terraform mapping.
2.  **Test Definition (The Contract)**:
    - Create `internal/provider/<package>/<name>_resource_test.go` **FIRST**.
    - Define `resource.TestCase` with steps for Create, Import, and Update.
    - Use modern `statecheck` and `knownvalue` packages for assertions (avoid deprecated `EnsureValue`).
    - Define the Terraform configuration string functions.
3.  **Schema Design**: Define the Go struct and Schema in `<name>_resource.go`.
4.  **Implementation**: Implement `Create`, `Read`, `Update`, `Delete` methods.
5.  **Validation**: Run tests (`TestAcc...`) and iterate.
6.  **Documentation**: Add usage examples in `examples/` for documentation generation.

---

## Code Standards & Structure

### Directory Organization

Organize resources by domain, not as flat files or monolithic packages.
Strictly follow the `internal/provider/<package>/` pattern.

- **Path**: `internal/provider/<package>/<entity>_resource.go`
- **Tests**: `internal/provider/<package>/<entity>_resource_test.go`
- **Models**: internal to the resource file (unless shared).

**Example:**

```text
internal/provider/oidc_client/
├── oidc_client_resource.go
└── oidc_client_resource_test.go
```

### Type Safety & Models

- **Always** use **Framework Types** in your model structs:
  - `types.String`, `types.Bool`, `types.List`, `types.Int64`, `types.Set`, `types.Map`.
- **Avoid** native Go types (`string`, `bool`) in resource models to properly handle Null/Unknown states.
- Use `tfsdk` tags for struct fields mapping to schema names.

### Naming Conventions

- **Terraform Resource**: `rauthy_<entity>` (e.g., `rauthy_client`).
- **Go Type**: `<Entity>Resource` (e.g., `clientResource`).
- **Test Config Function**: `testAcc<Entity>ResourceConfig`.

---

## Debugging & Observability

### Structured Logging (tflog)

Do NOT use `fmt.Printf`, `log.Println`, or panic. Use `tflog` for all interactions.

- **Trace/Debug**: detailed payload inspection (e.g., API requests/responses).
- **Info**: high-level lifecycle events (Creating, Updating, Deleting).
- **Error**: actionable failures.

**Example:**

```go
tflog.Debug(ctx, "creating rauthy client", map[string]interface{}{
    "id": plan.ID.ValueString(),
})
```

### Debugging Workflow

If a test fails or behavior is unexpected:

1.  **Instrument**: Add `tflog` calls to the relevant CRUD methods.
2.  **Run with Logs**: Execute the specific test with trace logging enabled.
    ```bash
    TF_LOG=TRACE go test ./internal/provider/... -run TestAccClientResource
    ```
3.  **Inspect**: Check the provider logs for API request/response details, state transitions, and `tflog` output.

---

## Testing Guidelines

### Acceptance Tests

- Use `github.com/hashicorp/terraform-plugin-testing/helper/resource`.
- Ensure `ImportState: true` is verified in at least one step.
- Verify properties using `statecheck.ExpectKnownValue`.

**Example Pattern:**

```go
resource.Test(t, resource.TestCase{
    PreCheck:                 func() { acctest.TestAccPreCheck(t) },
    ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
    Steps: []resource.TestStep{
        {
            Config: testAccClientResourceConfig(...),
            ConfigStateChecks: []statecheck.StateCheck{
                statecheck.ExpectKnownValue(
                    "rauthy_client.test",
                    tfjsonpath.New("name"),
                    knownvalue.StringExact("expected_name"),
                ),
            },
        },
    },
})
```
