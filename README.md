<a name="readme-top" id="readme-top"></a>

<p align="center">
  <img src="https://raw.githubusercontent.com/conduktor/conduktor.io-public/main/logo/transparent.png" width="256px" />
</p>
<h1 align="center">
    <strong>Conduktor Terraform Provider</strong>
</h1>

<p align="center">
    <a href="https://docs.conduktor.io/"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/conduktor/terraform-provider-conduktor/issues">Report Bug</a>
    ·
    <a href="https://github.com/conduktor/terraform-provider-conduktor/issues">Request Feature</a>
    ·
    <a href="https://support.conduktor.io/">Contact support</a>
    <br />
    <br />
    <a href="https://github.com/conduktor/terraform-provider-conduktor/releases/latest"><img alt="GitHub Release" src="https://img.shields.io/github/v/release/conduktor/terraform-provider-conduktor?color=BCFE68"></a>
    ·
    <img alt="License" src="https://img.shields.io/github/license/conduktor/terraform-provider-conduktor?color=BCFE68">
    ·
    <a href="https://registry.terraform.io/providers/conduktor/conduktor"><img alt="Terraform registry" src="https://img.shields.io/badge/Registry-Terraform-BCFE68"></a>
    <br />
    <br />
    <a href="https://conduktor.io/"><img src="https://img.shields.io/badge/Website-conduktor.io-192A4E?color=BCFE68" alt="Scale Data Streaming With Security and Control"></a>
    ·
    <a href="https://twitter.com/getconduktor"><img alt="X (formerly Twitter) Follow" src="https://img.shields.io/twitter/follow/getconduktor?color=BCFE68"></a>
    ·
    <a href="https://conduktor.io/slack"><img src="https://img.shields.io/badge/Slack-Join%20Community-BCFE68?logo=slack" alt="Slack"></a>
</p>

This repository contains the Conduktor Terraform provider, which defines Conduktor resources so that they can be deployed using Infrastructure as Code (IaC).

> [!WARNING]
> - The Conduktor Terraform provider is currently in **Alpha**.
> - It does not support all Console and Gateway resources yet. See our [resources roadmap](#resources-roadmap).
> - Let us know if you have [feedback](https://product.conduktor.help/c/74-terraform-provider) or wish to be a design partner.

## Supported resources

- [Console user](./docs/resources/user_v2.md)
- [Console group](./docs/resources/group_v2.md)
- [Generic](./docs/resources/generic.md) :warning: This resource is experimental and should be used with care.

## Install

Provider should be installed automatically with `terraform init`, but it's recommended to pin a specific version or range of version using following [`required_providers` configuration](https://developer.hashicorp.com/terraform/language/providers/requirements) :

```hcl
terraform {
  required_providers {
    conduktor = {
        source = "conduktor/conduktor"
        version = "~> X.Y" # where X.Y is the current major version and minor version
    }
  }
}
```

## Usage/Examples


```hcl
# configure provider
provider "conduktor" {
  console_url = "http://localhost:8080"
  api_token = "your-api-key" # can also use admin email/password to authenticate.
}

# register an external user bob with PLATFORM.userView permission
resource "conduktor_user_v2" "bob" {
  name = "bob@mycompany.io"
  spec {
    firstname = "Bob"
    lastname  = "Smith"
    permissions = [
      {
        permissions = [ "userView" ]
        resource_type = "PLATFORM"
      },
    ]
  }
}

# create a group with Bob as a member
resource "conduktor_group_v2" "qa" {
  name = "qa"
  spec {
    display_name                 = "QA team"
    description                  = "Quality Assurance team"
    members                      = [ conduktor_user_v2.bob.name ]
    permissions = [
      {
        resource_type = "PLATFORM"
        permissions   = ["userView", "clusterConnectionsManage"]
      }
    ]
  }
}
```

You can find more examples in this repository inside [`example`](./examples/) directory.

Examples can also be found in provider reference documentation available either in [`docs`](./docs/) directory or at [registry.terraform.io/conduktor/conduktor](https://registry.terraform.io/conduktor/conduktor/latest/docs)

You can also check out our [documentation](https://docs.conduktor.io/) for resources reference and provider usage.

### Provider authentication

To use Conduktor Console API, the Terraform provider needs to authenticate against it.

For that we offer two possibilities:

#### API key

Use an already manually forged API key. See [documentation](https://docs.conduktor.io/platform/reference/api-reference/#generate-an-api-key) to create one.

Using HCL `api_token` attribute
```hcl
provider "conduktor" {
  api_token = "your-api-key"
}
```
Using environment variables `CDK_API_TOKEN` or `CDK_API_KEY`.

#### Admin credentials
Use local user (usually admin) credentials pair. This will login against the API and use an ephemeral access token to make API calls.


Using HCL `admin_email`/`admin_password` attributes
```hcl
provider "conduktor" {
  admin_email    = "admin@my-org.com"
  admin_password = "admin-password"
}
```
Using environment variables `CDK_ADMIN_EMAIL` or `CDK_ADMIN_PASSWORD`.

Either way be aware that API Key and admin credentials are sensitive data and should be stored and provided to Terraform [properly](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables).

## Development
### Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23
- [Docker](https://docs.docker.com/get-docker/) with compose to run acceptance tests locally
- [Git hooks](#install-git-hooks) to format/lint code before committing


### Install git hooks
Please install the git hooks to ensure that the code is formatted correctly and pass linter check before committing.

Run `make install-githooks` to install the git hooks.

### Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

#### Build and install provider in local Terraform registry

Will build and install terraform provider locally in `~/.terraform.d/plugins` directory.
Local provider version is set on `VERSION` variable of [GNUmakefile](./GNUmakefile)

```shell
VERSION=0.0.1 make deploy-locally
```
It can then be used on terraform recipe like
```hcl
terraform {
  required_providers {
    conduktor = {
      source  = "terraform.local/conduktor/conduktor" # local provider
      version = ">= 0.0.1"                            # latest version found locally in the plugin cache.
    }
  }
}
```

### Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

### Codegen

The project uses different codegen tool to generate source files.

**Documentation** in [`docs`](./docs/) folder is generated using [tfplugindocs](github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs)

**Terraform schema definition** in [`schema`](./internal/schema/) module are generated using [tfplugingen-framework](github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework) from currently manually maintained code spec [json file](./provider_code_spec.json).

### Run acceptance tests

```shell
# Optional
export CDK_LICENSE="your_license_here"
make testacc
```
This action will start a testing environment using [Docker Compose](./docker-compose.yaml) and run all acceptance tests against it. Test environment is destroy at the end.

You can also start/stop environment and run tests in separate actions using `make start_test_env` / `make test` / `make clean`.

### Misc

```shell
make generate   # run go generate
make build      # run build
make go-fmt     # run go fmt on the project
make go-lint    # run golangci-lint linter
```

## Resources Roadmap

Future versions of the Conduktor Terraform provider will evolve to support more resources.

Need a resource to unblock a use case? [Feedback](https://product.conduktor.help/c/74-terraform-provider) to the Product team directly.

Our current order of priority is:

1. Console resources:
  - [Kafka Clusters with Schema Regsitry](https://docs.conduktor.io/platform/reference/resource-reference/console/#kafkacluster)
  - [Kafka Connect Cluster](https://docs.conduktor.io/platform/reference/resource-reference/console/#kafkaconnectcluster)
2. Kafka resources:
  - [Topic](https://docs.conduktor.io/platform/reference/resource-reference/kafka/#topic)
  - [Subject](https://docs.conduktor.io/platform/reference/resource-reference/kafka/#subject)
  - [Connector](https://docs.conduktor.io/platform/reference/resource-reference/kafka/#connector)
3. Self-service resources:
  - [Application](https://docs.conduktor.io/platform/reference/resource-reference/self-service/#application)
  - [ApplicationInstance](https://docs.conduktor.io/platform/reference/resource-reference/self-service/#application-instance)
  - [TopicPolicy](https://docs.conduktor.io/platform/reference/resource-reference/self-service/#topic-policy)
  - [ApplicationInstancePermission](https://docs.conduktor.io/platform/reference/resource-reference/self-service/#application-instance-permissions)
  - [ApplicationGroup](https://docs.conduktor.io/platform/reference/resource-reference/self-service/#application-group)
4. Gateway resources:
  - [Interceptor](https://docs.conduktor.io/gateway/reference/resources-reference/#interceptor)
  - [GatewayServiceAccount](https://docs.conduktor.io/gateway/reference/resources-reference/#gatewayserviceaccount)
  - [GatewayGroup](https://docs.conduktor.io/gateway/reference/resources-reference/#gatewaygroup)
  - [ConcentrationRule](https://docs.conduktor.io/gateway/reference/resources-reference/#concentrationrule)
  - [VirtualCluster](https://docs.conduktor.io/gateway/reference/resources-reference/#virtualcluster)
  - [AliasTopic](https://docs.conduktor.io/gateway/reference/resources-reference/#aliastopic)

> [!NOTE]
>
> This list is not exaustive and can change depending on requests and needs.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.