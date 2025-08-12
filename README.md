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
> - Let us know if you have [feedback](https://support.conduktor.io/hc/en-gb/requests/new?ticket_form_id=17438365654417) or wish to be a design partner.

**Table of Contents**
- [Supported resources](#supported-resources)
- [Install](#install)
- [Usage/Examples](#usageexamples)
  - [Provider authentication](#provider-authentication)
    - [Conduktor Console](#conduktor-console)
        - [API key](#api-key)
        - [Admin credentials](#admin-credentials)
    - [Conduktor Gateway](#conduktor-gateway)
        - [Admin credentials](#admin-credentials-2)
  - [Multi client configuration](#multi-client-configuration)
- [Development](#development)
  - [Requirements](#requirements)
  - [Install git hooks](#install-git-hooks)
  - [Building The Provider](#building-the-provider)
    - [Build and install provider in local Terraform registry](#build-and-install-provider-in-local-terraform-registry)
  - [Adding Dependencies](#adding-dependencies)
  - [Codegen](#codegen)
  - [Run acceptance tests](#run-acceptance-tests)
  - [Misc](#misc)
- [Resources Roadmap](#resources-roadmap)
- [Contributing](#contributing)
- [License](#license)

## Supported resources

### Console resources

- [Console user](./docs/resources/console_user_v2.md)
- [Console group](./docs/resources/console_group_v2.md)
- [Console partner zone](./docs/resources/console_partner_zone_v2.md)
- [Kafka clusters with schema registry](./docs/resources/console_kafka_cluster_v2.md)
- [Kafka connect server](./docs/resources/console_kafka_connect_v2.md)
- [KsqlDB cluster](./docs/resources/console_ksqldb_cluster_v2.md)

### Gateway resources

- [Gateway service account](./docs/resources/gateway_service_account_v2.md)
- [Gateway interceptor](./docs/resources/gateway_interceptor_v2.md)

### Self-service resources

- [Console application](./docs/resources/console_application_v1.md)
- [Console application group](./docs/resources/console_application_group_v1.md)
- [Console application instance](./docs/resources/console_application_instance_v1.md)
- [Console application instance permission](./docs/resources/console_application_instance_permission_v1.md)
- [Console service account](./docs/resources/console_service_account_v1.md)
- [Console resource policy](./docs/resources/console_resource_policy_v1.md)
- [Console topic policy](./docs/resources/console_topic_policy_v1.md)

### Kafka resources
- [Topic](./docs/resources/console_topic_v2.md)
- [Connector](./docs/resources/console_connector_v2.md)

### Generic resource

- [Generic](./docs/resources/generic.md) :warning: This resource is experimental and should be used with care.

## Install

Provider should be installed automatically with `terraform init`, but it's recommended to pin a specific version or range of versions using the following [`required_providers` configuration](https://developer.hashicorp.com/terraform/language/providers/requirements) :

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
  mode        = "console"
  base_url    = "http://localhost:8080"
  api_token   = "your-api-key" # can also use admin email/password to authenticate.
}

# register an external user bob with PLATFORM.userView permission
resource "conduktor_console_user_v2" "bob" {
  name = "bob@mycompany.io"
  spec = {
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
resource "conduktor_console_group_v2" "qa" {
  name = "qa"
  spec = {
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

> [!IMPORTANT]
> It is required to specify the provider `mode` to use, as it will determine the authentication method.

The provider can be used in two modes: `console` and `gateway`.

Example using console mode:
```hcl
provider "conduktor" {
  mode = "console"
  # ...
}
```

#### Conduktor Console

To use Conduktor Console API, the Terraform provider needs to authenticate against it.

For that we offer two possibilities:

##### API key

Use an already manually forged API key. See [documentation](https://docs.conduktor.io/platform/reference/api-reference/#generate-an-api-key) to create one.

Using HCL `api_token` attribute
```hcl
provider "conduktor" {
  mode      = "console"
  api_token = "your-api-key"
}
```
Using environment variables `CDK_API_TOKEN` or `CDK_API_KEY`.

##### Admin credentials
Use local user (usually admin) credentials pair. This will login against the API and use an ephemeral access token to make API calls.


Using HCL `admin_user`/`admin_password` attributes
```hcl
provider "conduktor" {
  mode           = "console"
  admin_user     = "admin@my-org.com"
  admin_password = "admin-password"
}
```
Using environment variables `CDK_ADMIN_EMAIL` or `CDK_ADMIN_PASSWORD`.

Either way be aware that API Key and admin credentials are sensitive data and should be stored and provided to Terraform [properly](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables).

#### Conduktor Gateway

To use Conduktor Gateway API, the Terraform provider needs to authenticate against it.

##### Admin credentials
Use local user (usually admin) credentials pair. Those will be used in the authentication header for the HTTP requests against the API.

Using HCL `admin_user`/`admin_password` attributes
```hcl
provider "conduktor" {
  mode           = "gateway"
  admin_user     = "admin@my-org.com"
  admin_password = "admin-password"
}
```
Using environment variables `CDK_ADMIN_EMAIL` or `CDK_ADMIN_PASSWORD`.

### Multi client configuration

Conduktor provider can also be configured to use multiple clients, each with its own authentication method.

For this we will make use of the `alias` attribute in the provider definition. Further information can be found on the official [Terraform Documentation](https://developer.hashicorp.com/terraform/language/providers/configuration#alias-multiple-provider-configurations).

```hcl
provider "conduktor" {
  alias    = "console"
  mode     = "console"
  # ...
}

provider "conduktor" {
  alias    = "gateway"
  mode     = "gateway"
  # ...
}
```

You will also need to specify the provider alias when defining resources.
``` hcl
resource "conduktor_console_user_v2" "user" {
  provider = conduktor.console
  # ...
}

resource "conduktor_gateway_service_account_v2" "gateway_sa" {
  provider = conduktor.gateway
  # ...
}
```

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

**Documentation** in [`docs`](./docs/) folder is generated using [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs).

**Terraform schema definition** in [`schema`](./internal/schema/) module are generated using [tfplugingen-framework](https://github.com/hashicorp/terraform-plugin-codegen-framework) from currently manually maintained code spec [json file](./provider_code_spec.json).

### Run acceptance tests

```shell
# Optional
export CDK_LICENSE="your_license_here"
make testacc
```
This action will start a testing environment using [Docker Compose](./docker-compose.yaml) and run all acceptance tests against it. Test environment is destroyed at the end.

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

Need a resource to unblock a use case? [Feedback](https://support.conduktor.io/hc/en-gb/requests/new?ticket_form_id=17438365654417) to the Product team directly.

We are currently reviewing the following resources:

1. Console resources:
  - [Alert](https://docs.conduktor.io/platform/reference/resource-reference/console/#alert)
2. Kafka resources:
  - [Subject](https://docs.conduktor.io/platform/reference/resource-reference/kafka/#subject)
3. Gateway resources:
  - [GatewayGroup](https://docs.conduktor.io/gateway/reference/resources-reference/#gatewaygroup)
  - [ConcentrationRule](https://docs.conduktor.io/gateway/reference/resources-reference/#concentrationrule)
  - [VirtualCluster](https://docs.conduktor.io/gateway/reference/resources-reference/#virtualcluster)
  - [AliasTopic](https://docs.conduktor.io/gateway/reference/resources-reference/#aliastopic)

> [!NOTE]
>
> This list is not exhaustive and can change depending on requests and needs.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
