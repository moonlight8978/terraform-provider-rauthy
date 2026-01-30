# Agents Configuration

This file contains configuration and context for AI agents working on the `terraform-provider-rauthy` project.

## Project Overview

This project is a **Terraform Provider for Rauthy** - a lightweight Identity Provider (IdP) written in Rust.

- **Provider Framework**: [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)
- **Target System**: [Rauthy](https://github.com/sebadob/rauthy)
- **Language**: Go

## Directory Structure

### `pkg/rauthy` (SDK/Client Layer)

Contains the Go client interacting with the Rauthy API.

- **Role**: Validates inputs, marshals/unmarshals JSON, handles HTTP requests to Rauthy.
- **Testing**: Unit tests are co-located (e.g., `role.go` tests are in `role_test.go`).
- **Convention**: Pure Go logic, decoupled from Terraform internals.

### `pkg/tfutils` (Terraform Framework Utilities)

Contains common utilities for the Terraform Plugin Framework.

### `internal/provider` (Terraform Layer)

Contains the Terraform provider implementation.

- **Organization**: Grouped by logical resource/component (e.g., `internal/provider/oidc_client/`).
- **Key Files**:
  - `resource.go` / `data_source.go`: The main TF Plugin Framework implementation.
  - `*_test.go`: Acceptance tests using `resource.Test`.
- **`provider.go`**: The root provider definition.
- **`testacc/`**: Acceptance tests helpers
- **`utils/`**: Common utilities for the project

### `examples/` (Example Usage)

Contains example Terraform configurations for using the provider, use to generate docs on Terraform registry

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

## Coding Standards

1.  **Separation of Concerns**: Keep Terraform logic (Schema, State handling) in `internal/provider` and API logic (structs, requests) in `pkg/rauthy`.
2.  **Testing**:
    - Write unit tests for complex logic in `pkg/rauthy`.
    - Write acceptance tests for all Resources and Data Sources in `internal/provider`.
3.  **Naming**: Follow Go and Terraform naming conventions (Snake case for TF attributes, Camel case for Go structs).
4.  **Comments**: Only add comments when necessary or the logic is very complex, or it's tricky for legacy support.
5.  **Documentation**: Document, write examples follows the Terraform Plugin Framework docs

## Common Tasks for Agents

- **Adding a new Resource**:
  1.  Define the API methods/structs in `pkg/rauthy`.
  2.  Add unit tests in `pkg/rauthy`.
  3.  Create a new directory `internal/provider/<name>`.
  4.  Implement the `resource.Resource` interface.
  5.  Register the resource in `internal/provider/provider.go`.
  6.  Add acceptance tests.
