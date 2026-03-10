# Terraform Provider Rauthy

This provider is used to interact with the awesome [Rauthy](https://github.com/sebadob/rauthy) project, but it lacks API support on some endpoints, hopefully it will be fully supported in the future. If you like this provider, make sure to star sebadob's awesome projects as well.

It is currently built to support our [fork](https://github.com/moonlight8978/rauthy), which adds API support at some endpoints alongside some customized behavior.

# Getting started

To install the provider, add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    rauthy = {
      source  = "moonlight8978/rauthy"
      version = "~> 1.0"
    }
  }
}

provider "rauthy" {
  # Configuration options
}
```

Please refer to the [Terraform Registry Documentation](https://registry.terraform.io/providers/moonlight8978/rauthy/latest/docs) for detailed provider documentation and examples.

# Development

To set up the project locally, ensure you have Go and [Task](https://taskfile.dev/) installed. You can also use the integrated Devcontainer for a fully configured environment.

To install dependencies and prepare the project, run:

```bash
task prepare
```

To run the dependencies (rauthy, mailcrab) in Docker, you can start it with:

```bash
task up
```

The provider itself had not supported container yet.

### Running Tests

You can run the different test suites using the configured task runner commands:

- Run module package tests:

  ```bash
  task test:pkg
  ```

- Run provider acceptance tests:
  ```bash
  task test:provider
  ```

# License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
