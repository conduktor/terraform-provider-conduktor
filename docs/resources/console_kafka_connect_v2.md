---
page_title: "Conduktor : conduktor_console_kafka_connect_v2 "
subcategory: "console/v2"
description: |-
    Resource for managing Conduktor Kafka Connect servers definition linked to an existing Kafka cluster definition inside Conduktor Console.
    This resource allows you to create, read, update and delete Kafka Connect servers connections from Conduktor Console.
---

# conduktor_console_kafka_connect_v2

Resource for managing Conduktor Kafka Connect servers definition linked to an existing Kafka cluster definition inside Conduktor Console.
This resource allows you to create, read, update and delete Kafka Connect servers connections from Conduktor Console.

## Example Usage

### Simple Kafka Connect server
This example creates a simple Kafka Connect server connection without any authentication.
```terraform
resource "conduktor_console_kafka_cluster_v2" "minimal" {
  name = "mini-cluster"
  spec = {
    display_name      = "Minimal Cluster"
    bootstrap_servers = "localhost:9092"
  }
}

resource "conduktor_console_kafka_connect_v2" "simple" {
  name    = "simple-connect"
  cluster = conduktor_console_kafka_cluster_v2.minimal.name
  spec = {
    display_name = "Simple Connect Server"
    urls         = "http://localhost:8083"
  }
}
```

### Basic Kafka Connect server
This example creates a complex Kafka Connect server connection with basic authentication.
```terraform
resource "conduktor_console_kafka_cluster_v2" "minimal" {
  name = "mini-cluster"
  spec = {
    display_name      = "Minimal Cluster"
    bootstrap_servers = "localhost:9092"
  }
}

resource "conduktor_console_kafka_connect_v2" "basic" {
  name    = "basic-connect"
  cluster = conduktor_console_kafka_cluster_v2.minimal.name
  labels = {
    description   = "This is a complex connect using basic authentication"
    documentation = "https://docs.mycompany.com/complex-connect"
    env           = "dev"
  }
  spec = {
    display_name = "Basic Connect server"
    urls         = "http://localhost:8083"
    headers = {
      X-PROJECT-HEADER = "value"
      Cache-Control    = "no-cache"
    }
    ignore_untrusted_certificate = false
    security = {
      basic_auth = {
        username = "user"
        password = "password"
      }
    }
  }
}
```

### Bearer token Kafka Connect server
This example creates a complex Kafka Connect server connection with bearer token authentication.
```terraform
resource "conduktor_console_kafka_cluster_v2" "minimal" {
  name = "mini-cluster"
  spec = {
    display_name      = "Minimal Cluster"
    bootstrap_servers = "localhost:9092"
  }
}

resource "conduktor_console_kafka_connect_v2" "bearer" {
  name    = "bearer-connect"
  cluster = conduktor_console_kafka_cluster_v2.minimal.name
  labels = {
    description   = "This is a complex connect using bearer token authentication"
    documentation = "https://docs.mycompany.com/complex-connect"
    env           = "dev"
  }
  spec = {
    display_name = "Bearer Connect server"
    urls         = "http://localhost:8083"
    headers = {
      X-PROJECT-HEADER = "value"
      Cache-Control    = "no-cache"
    }
    ignore_untrusted_certificate = false
    security = {
      bearer_token = {
        token = "token"
      }
    }
  }
}
```

### mTLS Kafka Connect server
This example creates a complex Kafka Connect server connection with mTLS authentication.
```terraform
resource "conduktor_console_kafka_cluster_v2" "minimal" {
  name = "mini-cluster"
  spec = {
    display_name      = "Minimal Cluster"
    bootstrap_servers = "localhost:9092"
  }
}

resource "conduktor_console_kafka_connect_v2" "mtls" {
  name    = "mtls-connect"
  cluster = conduktor_console_kafka_cluster_v2.minimal.name
  labels = {
    description   = "This is a complex connect using mTLS authentication"
    documentation = "https://docs.mycompany.com/complex-connect"
    env           = "dev"
  }
  spec = {
    display_name = "mTLS Connect server"
    urls         = "http://localhost:8083"
    headers = {
      X-PROJECT-HEADER = "value"
      Cache-Control    = "no-cache"
    }
    ignore_untrusted_certificate = false
    security = {
      ssl_auth = {
        key               = <<EOT
-----BEGIN PRIVATE KEY-----
MIIOXzCCDUegAwIBAgIRAPRytMVYJNUgCbhnA+eYumgwDQYJKoZIhvcNAQELBQAw
...
IFyCs+xkcgvHFtBjjel4pnIET0agtbGJbGDEQBNxX+i4MDA=
-----END PRIVATE KEY-----
EOT
        certificate_chain = <<EOT
-----BEGIN CERTIFICATE-----
MIIOXzCCDUegAwIBAgIRAPRytMVYJNUgCbhnA+eYumgwDQYJKoZIhvcNAQELBQAw
...
IFyCs+xkcgvHFtBjjel4pnIET0agtbGJbGDEQBNxX+i4MDA=
-----END CERTIFICATE-----
EOT
      }
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster` (String) Kafka cluster name linked with the current Kafka connect server. Must already exist in Conduktor Console
- `name` (String) Kafka connect server name, must be unique, acts as an ID for import
- `spec` (Attributes) Kafka connect server specification (see [below for nested schema](#nestedatt--spec))

### Optional

- `labels` (Map of String) Kafka connect server labels

<a id="nestedatt--spec"></a>
### Nested Schema for `spec`

Required:

- `display_name` (String) Kafka connect server display name
- `urls` (String) URL of a Kafka Connect cluster. **Multiple URLs are not supported for now**

Optional:

- `headers` (Map of String) Key-Value HTTP headers to add to requests
- `ignore_untrusted_certificate` (Boolean) Ignore untrusted certificate for Kafka connect server requests
- `security` (Attributes) Kafka connect server security configuration. One of `basic_auth`, `bearer_token`, `ssl_auth` (see [below for nested schema](#nestedatt--spec--security))

<a id="nestedatt--spec--security"></a>
### Nested Schema for `spec.security`

Optional:

- `basic_auth` (Attributes) Basic auth for Kafka connect server security configuration. (see [below for nested schema](#nestedatt--spec--security--basic_auth))
- `bearer_token` (Attributes) Bearer token for Kafka connect server security configuration. (see [below for nested schema](#nestedatt--spec--security--bearer_token))
- `ssl_auth` (Attributes) SSL auth (mTLS) for Kafka connect server security configuration. (see [below for nested schema](#nestedatt--spec--security--ssl_auth))

<a id="nestedatt--spec--security--basic_auth"></a>
### Nested Schema for `spec.security.basic_auth`

Required:

- `password` (String, Sensitive) Kafka connect server basic auth password.
- `username` (String) Kafka connect server basic auth username.


<a id="nestedatt--spec--security--bearer_token"></a>
### Nested Schema for `spec.security.bearer_token`

Required:

- `token` (String, Sensitive) Kafka connect server bearer token.


<a id="nestedatt--spec--security--ssl_auth"></a>
### Nested Schema for `spec.security.ssl_auth`

Required:

- `certificate_chain` (String) Kafka connect server mTLS auth certificate chain PEM.
- `key` (String, Sensitive) Kafka connect server mTLS auth private key PEM.






## Import

In order to import a Kafka Connect server connection into Conduktor, you need to know the Kafka cluster ID and the Kafka Connect server ID.

The import ID is constructed as follows: `< cluster_id >/< connect_id >`.

For example, using an [`import` block](https://developer.hashicorp.com/terraform/language/import) :
```terraform
import {
  to = conduktor_console_kafka_connect_v2.example
  id = "mini-cluster/import-connect" # Import "import-connect" Connect server for "mini-cluster" Kafka cluster
}
```

Using the `terraform import` command:
```shell
terraform import conduktor_console_kafka_connect_v2.example mini-cluster/import-connect
```
