# Agents Configuration

This file contains configuration and context for AI agents working on the `terraform-provider-rauthy` project.

## Project Overview

This project is a **Terraform Provider for Rauthy** - a lightweight Identity Provider (IdP) written in Rust.

- **Provider Framework**: [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)
- **Target System**: [Rauthy](https://github.com/sebadob/rauthy)
- **Language**: Go
- **Current Resources**: `rauthy_role`, `rauthy_oidc_client`, `rauthy_oidc_client_secret`, `rauthy_auth_provider`, `rauthy_password_policy`

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Terraform Core                          │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              internal/provider (TF Layer)                   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ provider.go - Provider configuration & registration  │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ role/, oidc_client/, etc. - Resource implementations │   │
│  │ - Schema definitions (TF attributes)                 │   │
│  │ - CRUD operations (Create/Read/Update/Delete)        │   │
│  │ - State management                                   │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              pkg/rauthy (API Client Layer)                  │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ client.go - HTTP client & authentication             │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ role.go, oidc_client.go, etc. - API methods          │   │
│  │ - Request/Response structs                           │   │
│  │ - HTTP method calls                                  │   │
│  │ - JSON marshaling/unmarshaling                       │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
                   ┌────────────────┐
                   │  Rauthy API    │
                   └────────────────┘
```

## Directory Structure

### `pkg/rauthy` (SDK/Client Layer)

Contains the Go client interacting with the Rauthy API.

- **Role**: Validates inputs, marshals/unmarshals JSON, handles HTTP requests to Rauthy.
- **Testing**: Unit tests are co-located (e.g., `role.go` tests are in `role_test.go`).
- **Convention**: Pure Go logic, decoupled from Terraform internals.
- **Key Files**:
  - `client.go`: HTTP client, request wrapper, authentication
  - `authenticator.go`: Authentication interface
  - `<resource>.go`: API methods for each resource type (CRUD operations)
  - `<resource>_test.go`: Unit tests using mock HTTP servers

### `pkg/tfutils` (Terraform Framework Utilities)

Contains common utilities for the Terraform Plugin Framework.

### `internal/provider` (Terraform Layer)

Contains the Terraform provider implementation.

- **Organization**: Grouped by logical resource/component (e.g., `internal/provider/role/`).
- **Key Files**:
  - `<resource>_resource.go`: The main TF Plugin Framework implementation.
  - `<resource>_resource_test.go`: Acceptance tests using `resource.Test`.
- **`provider.go`**: The root provider definition, resource registration.
- **`acctest/`**: Acceptance test helpers (provider factories, pre-check functions).
- **`utils/`**: Common utilities for the provider layer.

### `examples/` (Example Usage)

Contains example Terraform configurations for using the provider. These are used to generate documentation on the Terraform registry.

## Coding Patterns

### 1. Adding a New Resource - Complete Workflow

#### Step 1: Define API Client in `pkg/rauthy`

Create `pkg/rauthy/<resource>.go`:

```go
package rauthy

import (
    "context"
    "fmt"
    "net/http"
)

// Define the API response struct
type MyResource struct {
    Id   string `json:"id"`
    Name string `json:"name"`
    // Add other fields as needed
}

// Define the API request struct
type MyResourceRequest struct {
    Name string `json:"name"`
    // Add other fields as needed
}

// CRUD operations
func (c *Client) CreateMyResource(ctx context.Context, req *MyResourceRequest) (MyResource, error) {
    var created MyResource
    if _, err := c.Request(ctx, http.MethodPost, "/my-resources", req, &created); err != nil {
        return created, err
    }
    return created, nil
}

func (c *Client) GetMyResource(ctx context.Context, id string) (*MyResource, error) {
    var resource MyResource
    if _, err := c.Request(ctx, http.MethodGet, fmt.Sprintf("/my-resources/%s", id), nil, &resource); err != nil {
        return nil, err
    }
    return &resource, nil
}

func (c *Client) UpdateMyResource(ctx context.Context, id string, req *MyResourceRequest) (*MyResource, error) {
    var updated MyResource
    if _, err := c.Request(ctx, http.MethodPut, fmt.Sprintf("/my-resources/%s", id), req, &updated); err != nil {
        return nil, err
    }
    return &updated, nil
}

func (c *Client) DeleteMyResource(ctx context.Context, id string) error {
    if _, err := c.Request(ctx, http.MethodDelete, fmt.Sprintf("/my-resources/%s", id), nil, nil); err != nil {
        return err
    }
    return nil
}
```

#### Step 2: Write Unit Tests in `pkg/rauthy`

Create `pkg/rauthy/<resource>_test.go`:

```go
package rauthy_test

import (
    "context"
    "net/http"
    "testing"

    "github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
    "github.com/stretchr/testify/assert"
)

var myResourceResponse = `{
    "id": "resource-1",
    "name": "Test Resource"
}`

func TestCreateMyResource(t *testing.T) {
    ts := CreateServer(myResourceResponse, http.StatusOK)
    defer ts.Close()

    client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

    created, err := client.CreateMyResource(context.Background(), &rauthy.MyResourceRequest{Name: "Test Resource"})
    assert.NoError(t, err)
    assert.Equal(t, "resource-1", created.Id)
    assert.Equal(t, "Test Resource", created.Name)
}

// Add tests for Get, Update, Delete operations
```

#### Step 3: Implement Terraform Resource

Create `internal/provider/<resource>/<resource>_resource.go`:

```go
package myresource

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
)

var _ resource.Resource = &MyResourceResource{}
var _ resource.ResourceWithImportState = &MyResourceResource{}

func NewMyResourceResource() resource.Resource {
    return &MyResourceResource{}
}

type MyResourceResource struct {
    client *rauthy.Client
}

type MyResourceResourceModel struct {
    Id   types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
}

func (r *MyResourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_my_resource"
}

func (r *MyResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "My Resource description",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Resource ID",
                Computed:            true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Resource name",
                Required:            true,
            },
        },
    }
}

func (r *MyResourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*rauthy.Client)
    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *rauthy.Client, got: %T", req.ProviderData),
        )
        return
    }

    r.client = client
}

func (r *MyResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data MyResourceResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    created, err := r.client.CreateMyResource(ctx, data.ToRequest())
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create resource: %s", err))
        return
    }

    data.Id = types.StringValue(created.Id)
    data.Name = types.StringValue(created.Name)

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MyResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data MyResourceResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    resource, err := r.client.GetMyResource(ctx, data.Id.ValueString())
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read resource: %s", err))
        return
    }

    data.Id = types.StringValue(resource.Id)
    data.Name = types.StringValue(resource.Name)

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MyResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data MyResourceResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    updated, err := r.client.UpdateMyResource(ctx, data.Id.ValueString(), data.ToRequest())
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update resource: %s", err))
        return
    }

    data.Name = types.StringValue(updated.Name)

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MyResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data MyResourceResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    if err := r.client.DeleteMyResource(ctx, data.Id.ValueString()); err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete resource: %s", err))
        return
    }
}

func (r *MyResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource, err := r.client.GetMyResource(ctx, req.ID)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import resource: %s", err))
        return
    }

    model := MyResourceResourceModel{
        Id:   types.StringValue(resource.Id),
        Name: types.StringValue(resource.Name),
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r MyResourceResourceModel) ToRequest() *rauthy.MyResourceRequest {
    return &rauthy.MyResourceRequest{
        Name: r.Name.ValueString(),
    }
}
```

#### Step 4: Register Resource in Provider

Edit `internal/provider/provider.go`:

```go
import (
    // ... other imports
    "github.com/moonlight8978/terraform-provider-rauthy/internal/provider/myresource"
)

func (p *RauthyProvider) Resources(ctx context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        // ... existing resources
        myresource.NewMyResourceResource,
    }
}
```

#### Step 5: Write Acceptance Tests

Create `internal/provider/<resource>/<resource>_resource_test.go`:

```go
package myresource_test

import (
    "fmt"
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
    "github.com/moonlight8978/terraform-provider-rauthy/internal/provider/acctest"
)

func TestAccMyResourceResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { acctest.TestAccPreCheck(t) },
        ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccMyResourceResourceConfig("test"),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("rauthy_my_resource.test", "name", "test"),
                    resource.TestCheckResourceAttrSet("rauthy_my_resource.test", "id"),
                ),
            },
            {
                ResourceName:      "rauthy_my_resource.test",
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}

func testAccMyResourceResourceConfig(name string) string {
    return fmt.Sprintf(`
resource "rauthy_my_resource" "test" {
    name = %[1]q
}
`, name)
}
```

### 2. Error Handling Pattern

**Consistent error handling across the codebase:**

```go
// In pkg/rauthy - API client errors
if err != nil {
    return nil, fmt.Errorf("failed to <operation>: %w", err)
}

// In internal/provider - Terraform diagnostics
if err != nil {
    resp.Diagnostics.AddError(
        "Client Error",
        fmt.Sprintf("Unable to <operation> <resource>: %s", err),
    )
    return
}
```

### 3. State Management Pattern

**Always follow this pattern for CRUD operations:**

1. **Get data from request** (Plan for Create/Update, State for Read/Delete)
2. **Check for diagnostics errors** before proceeding
3. **Call API client method**
4. **Handle errors** with diagnostics
5. **Update model** with response data
6. **Set state** with updated model

## Testing Best Practices

### Unit Testing (pkg/rauthy)

- **Location**: Co-located with source files (`<resource>_test.go`)
- **Purpose**: Test API client logic, request/response handling
- **Pattern**: Use mock HTTP servers via `CreateServer` helper
- **Run**: `task test:pkg` or `task test:pkg -- pkg/rauthy/role_test.go`

**Example test structure:**

```go
func TestCreateResource(t *testing.T) {
    // 1. Create mock server with expected response
    ts := CreateServer(expectedResponse, http.StatusOK)
    defer ts.Close()

    // 2. Create client pointing to mock server
    client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

    // 3. Call method under test
    result, err := client.CreateResource(context.Background(), &rauthy.ResourceRequest{...})

    // 4. Assert expectations
    assert.NoError(t, err)
    assert.Equal(t, "expected-id", result.Id)
}
```

### Acceptance Testing (internal/provider)

- **Location**: `internal/provider/<resource>/<resource>_resource_test.go`
- **Purpose**: Test Terraform resource behavior end-to-end
- **Requirements**: Running Rauthy instance (set via environment variables)
- **Run**: `task test:provider` or `task test:provider -- internal/provider/role/`

**Required environment variables:**

```bash
export TF_ACC=1
export RAUTHY_ENDPOINT=http://localhost:8080
export RAUTHY_API_KEY=your-api-key
export RAUTHY_INSECURE=true  # for local testing
```

**Example acceptance test structure:**

```go
func TestAccResourceName(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { acctest.TestAccPreCheck(t) },
        ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Step 1: Create resource
            {
                Config: testAccResourceConfig("initial-value"),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("rauthy_resource.test", "attribute", "initial-value"),
                    resource.TestCheckResourceAttrSet("rauthy_resource.test", "id"),
                ),
            },
            // Step 2: Update resource
            {
                Config: testAccResourceConfig("updated-value"),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("rauthy_resource.test", "attribute", "updated-value"),
                ),
            },
            // Step 3: Import resource
            {
                ResourceName:      "rauthy_resource.test",
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}
```

## Development Workflows

We use `task` (Taskfile) for common operations.

### Testing

- **Run SDK/Client Tests**:
  ```bash
  task test:pkg
  # OR for specific files
  task test:pkg -- pkg/rauthy/role_test.go
  ```
- **Run Provider Acceptance Tests**:
  ```bash
  task test:provider
  # OR for specific files
  task test:provider -- internal/provider/role/
  ```
  _Note: Acceptance tests require a running Rauthy instance._

### Dependencies

- **Install/Update**:
  ```bash
  task prepare
  ```

### Local Development

- **Start Dev Container**:
  ```bash
  task up
  ```

## Troubleshooting Guide

### Common Issues

#### 1. Acceptance Tests Fail with "connection refused"

**Problem**: Tests can't connect to Rauthy instance.

**Solution**:

- Ensure Rauthy is running: `docker run -p 8080:8080 ghcr.io/sebadob/rauthy:latest`
- Verify `RAUTHY_ENDPOINT` environment variable is set correctly
- Check `RAUTHY_API_KEY` is valid

#### 2. "Unexpected Resource Configure Type" Error

**Problem**: Resource's `Configure` method receives wrong type.

**Solution**:

- Ensure `provider.go` sets both `resp.ResourceData` and `resp.DataSourceData` to the client
- Check resource's `Configure` method type assertion matches `*rauthy.Client`

#### 3. Import State Not Working

**Problem**: `terraform import` fails or doesn't populate attributes.

**Solution**:

- Implement `resource.ResourceWithImportState` interface
- In `ImportState` method, fetch resource from API and populate all attributes
- Ensure all computed and required attributes are set

#### 4. Unit Tests Fail with JSON Unmarshaling Errors

**Problem**: Mock server response doesn't match expected struct.

**Solution**:

- Verify JSON tags on struct fields match API response
- Check mock response JSON is valid
- Use actual API response as reference for mock data

### Debugging Tips

1. **Enable Terraform Logging**:

   ```bash
   export TF_LOG=DEBUG
   export TF_LOG_PATH=terraform.log
   ```

2. **Test API Client Separately**:
   - Write unit tests for API client methods first
   - Use tools like `curl` or Postman to verify API responses
   - Compare actual API responses with your structs

3. **Use Acceptance Test Debugging**:
   ```bash
   TF_ACC=1 go test -v -run TestAccSpecificTest ./internal/provider/resource/
   ```

## Coding Standards

1. **Separation of Concerns**: Keep Terraform logic (Schema, State handling) in `internal/provider` and API logic (structs, requests) in `pkg/rauthy`.
2. **Testing**:
   - Write unit tests for complex logic in `pkg/rauthy`.
   - Write acceptance tests for all Resources and Data Sources in `internal/provider`.
3. **Naming**: Follow Go and Terraform naming conventions (Snake case for TF attributes, Camel case for Go structs).
4. **Comments**: Only add comments when necessary or the logic is very complex, or it's tricky for legacy support.
5. **Documentation**: Document and write examples following the Terraform Plugin Framework docs.
6. **Error Messages**: Be descriptive and actionable in error messages.

## Quick Reference

### File Location Cheatsheet

| Task                | Location                        | File Pattern                       |
| ------------------- | ------------------------------- | ---------------------------------- |
| Add API method      | `pkg/rauthy/`                   | `<resource>.go`                    |
| Add unit test       | `pkg/rauthy/`                   | `<resource>_test.go`               |
| Add TF resource     | `internal/provider/<resource>/` | `<resource>_resource.go`           |
| Add acceptance test | `internal/provider/<resource>/` | `<resource>_resource_test.go`      |
| Register resource   | `internal/provider/`            | `provider.go`                      |
| Add example         | `examples/`                     | `resources/<resource>/resource.tf` |

### Resource Implementation Checklist

- [ ] Define API structs in `pkg/rauthy/<resource>.go`
- [ ] Implement CRUD methods in API client
- [ ] Write unit tests in `pkg/rauthy/<resource>_test.go`
- [ ] Run unit tests: `task test:pkg -- pkg/rauthy/<resource>_test.go`
- [ ] Create resource directory: `internal/provider/<resource>/`
- [ ] Implement `<resource>_resource.go` with all interfaces
- [ ] Register resource in `internal/provider/provider.go`
- [ ] Write acceptance tests in `<resource>_resource_test.go`
- [ ] Run acceptance tests: `task test:provider -- internal/provider/<resource>/`
- [ ] Add example in `examples/resources/<resource>/resource.tf`
- [ ] Verify import functionality works

### Common Commands

```bash
# Install dependencies
task prepare

# Run all unit tests
task test:pkg

# Run specific unit test
task test:pkg -- pkg/rauthy/role_test.go

# Run all acceptance tests (requires Rauthy instance)
TF_ACC=1 task test:provider

# Run specific acceptance test
TF_ACC=1 task test:provider -- internal/provider/role/

# Start dev container
task up

# Format code
go fmt ./...

# Run linter
golangci-lint run
```

## Chore

- Always create temporary files or planning files in `tmp` directory, name it meaningfully

## Additional Resources

- [Terraform Plugin Framework Documentation](https://developer.hashicorp.com/terraform/plugin/framework)
- [Rauthy API Documentation](https://github.com/sebadob/rauthy)
- [Terraform Plugin Testing](https://developer.hashicorp.com/terraform/plugin/testing)
