---
page_title: "Conduktor : conduktor_console_kafka_cluster_v2 "
subcategory: "console/v2"
description: |-
    Resource for managing Conduktor Kafka cluster definition with optional Schema registry.
    This resource allows you to create, read, update and delete Kafka cluster and Schema registry definitions in Conduktor.
---

# conduktor_console_kafka_cluster_v2

Resource for managing Conduktor Kafka cluster and Schema registry definitions.
This resource allows you to create, read, update and delete Kafka clusters and Schema registry definitions in Conduktor.

## Example Usage

### Simple Kafka cluster without Schema registry
This example creates a simple Kafka cluster definition without authentication resource and without Schema Registry.
```terraform
resource "conduktor_console_kafka_cluster_v2" "simple" {
  name = "simple-cluster"
  spec = {
    display_name                 = "Simple kafka Cluster"
    icon                         = "kafka"
    color                        = "#000000"
    bootstrap_servers            = "localhost:9092"
    ignore_untrusted_certificate = true
  }
}
```

### Confluent Kafka cluster with Schema registry
This example creates a Confluent Kafka cluster and Schema Registry definition resource.
The Schema Registry authentication uses mTLS.
```terraform
resource "conduktor_console_kafka_cluster_v2" "confluent" {
  name = "confluent-cluster"
  labels = {
    "env" = "staging"
  }
  spec = {
    display_name      = "Confluent Cluster"
    bootstrap_servers = "aaa-aaaa.us-west4.gcp.confluent.cloud:9092"
    properties = {
      "sasl.jaas.config"  = "org.apache.kafka.common.security.plain.PlainLoginModule required username='admin' password='admin-secret';"
      "security.protocol" = "SASL_PLAINTEXT"
      "sasl.mechanism"    = "PLAIN"
    }
    icon                         = "kafka"
    ignore_untrusted_certificate = false
    kafka_flavor = {
      confluent = {
        key                      = "yourApiKey123456"
        secret                   = "yourApiSecret123456"
        confluent_environment_id = "env-12345"
        confluent_cluster_id     = "lkc-67890"
      }
    }
    schema_registry = {
      confluent_like = {
        url                          = "https://bbb-bbbb.us-west4.gcp.confluent.cloud:8081"
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
  }
}
```

### Aiven Kafka cluster with Schema registry
This example creates an Aiven Kafka cluster and Schema Registry definition resource.
The Schema Registry authentication uses basic auth.
```terraform
resource "conduktor_console_kafka_cluster_v2" "aiven" {
  name = "aiven-cluster"
  labels = {
    "env" = "test"
  }
  spec = {
    display_name      = "Aiven Cluster"
    bootstrap_servers = "cluster.aiven.io:9092"
    properties = {
      "sasl.jaas.config"  = "org.apache.kafka.common.security.plain.PlainLoginModule required username='admin' password='admin-secret';"
      "security.protocol" = "SASL_SSL"
      "sasl.mechanism"    = "PLAIN"
    }
    icon                         = "crab"
    ignore_untrusted_certificate = true
    kafka_flavor = {
      aiven = {
        api_token    = "a1b2c3d4e5f6g7h8i9j0"
        project      = "my-kafka-project"
        service_name = "my-kafka-service"
      }
    }
    schema_registry = {
      confluent_like = {
        url                          = "https://sr.aiven.io:8081"
        ignore_untrusted_certificate = false
        security = {
          basic_auth = {
            username = "uuuuuuu"
            password = "ppppppp"
          }
        }
      }
    }
  }
}
```

### AWS MSK with Glue Schema registry
This example creates an AWS MSK Kafka Cluster and a Glue Schema Registry definition resource.
```terraform
resource "conduktor_console_kafka_cluster_v2" "aws_msk" {
  name = "aws-cluster"
  labels = {
    "env" = "prod"
  }
  spec = {
    display_name      = "AWS MSK Cluster"
    bootstrap_servers = "b-3-public.xxxxx.yyyyy.zz.kafka.eu-west-1.amazonaws.com:9198,b-2-public.xxxxx.yyyyy.zz.kafka.eu-west-1.amazonaws.com:9198,b-1-public.xxxxx.yyyyy.zz.kafka.eu-west-1.amazonaws.com:9198"
    properties = {
      "sasl.jaas.config"                   = "software.amazon.msk.auth.iam.IAMLoginModule required awsRoleArn='arn:aws:iam::123456789123:role/MSK-role';"
      "sasl.client.callback.handler.class" = "software.amazon.msk.auth.iam.IAMClientCallbackHandler"
      "security.protocol"                  = "SASL_SSL"
      "sasl.mechanism"                     = "AWS_MSK_IAM"
    }
    icon                         = "kafka"
    color                        = "#FF0000"
    ignore_untrusted_certificate = true
    schema_registry = {
      glue = {
        region        = "eu-west-1"
        registry_name = "default"
        security = {
          credentials = {
            access_key_id = "accessKey"
            secret_key    = "secretKey"
          }
        }
      }
    }
  }
}
```

### Conduktor Gateway Kafka cluster with Schema registry
This example creates a Conduktor Gateway Kafka Cluster and Schema Registry definition resource.
The Schema Registry authentication uses a bearer token.
```terraform
resource "conduktor_console_kafka_cluster_v2" "gateway" {
  name = "gateway-cluster"
  labels = {
    "env" = "prod"
  }
  spec = {
    display_name      = "Gateway Cluster"
    bootstrap_servers = "gateway:6969"
    properties = {
      "sasl.jaas.config"  = "org.apache.kafka.common.security.plain.PlainLoginModule required username='admin' password='admin-secret';"
      "security.protocol" = "SASL_SSL"
      "sasl.mechanism"    = "PLAIN"
    }
    icon                         = "shield-blank"
    ignore_untrusted_certificate = true
    kafka_flavor = {
      gateway = {
        url                          = "http://gateway:8888"
        user                         = "admin"
        password                     = "admin"
        virtual_cluster              = "passthrough"
        ignore_untrusted_certificate = true
      }
    }
    schema_registry = {
      confluent_like = {
        url                          = "http://localhost:8081"
        ignore_untrusted_certificate = true
        security = {
          bearer_token = {
            token = "auth-token"
          }
        }
      }
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Kafka cluster name, must be unique, acts as an ID for import
- `spec` (Attributes) Kafka cluster specification (see [below for nested schema](#nestedatt--spec))

### Optional

- `labels` (Map of String) Kafka cluster labels

<a id="nestedatt--spec"></a>
### Nested Schema for `spec`

Required:

- `bootstrap_servers` (String) List of bootstrap servers for the Kafka cluster separated by comma
- `display_name` (String) Kafka cluster display name

Optional:

- `color` (String) Kafka cluster icon color in hexadecimal format like `#FF0000`
- `icon` (String) Kafka cluster icon. List of available icons can be found [here](https://docs.conduktor.io/platform/reference/resource-reference/console/#icon-sets)
- `ignore_untrusted_certificate` (Boolean) Ignore untrusted certificate for Kafka cluster
- `kafka_flavor` (Attributes) Kafka flavor configuration. One of `confluent`, `aiven`, `gateway` (see [below for nested schema](#nestedatt--spec--kafka_flavor))
- `properties` (Map of String) Kafka cluster properties
- `schema_registry` (Attributes) Schema registry configuration. One of `confluent_like`, `glue` (see [below for nested schema](#nestedatt--spec--schema_registry))

<a id="nestedatt--spec--kafka_flavor"></a>
### Nested Schema for `spec.kafka_flavor`

Optional:

- `aiven` (Attributes) Aiven Kafka flavor configuration (see [below for nested schema](#nestedatt--spec--kafka_flavor--aiven))
- `confluent` (Attributes) Confluent Kafka flavor configuration (see [below for nested schema](#nestedatt--spec--kafka_flavor--confluent))
- `gateway` (Attributes) Conduktor Gateway Kafka flavor configuration (see [below for nested schema](#nestedatt--spec--kafka_flavor--gateway))

<a id="nestedatt--spec--kafka_flavor--aiven"></a>
### Nested Schema for `spec.kafka_flavor.aiven`

Required:

- `api_token` (String, Sensitive) Aiven API token.
- `project` (String) Aiven project name.
- `service_name` (String) Aiven service name.


<a id="nestedatt--spec--kafka_flavor--confluent"></a>
### Nested Schema for `spec.kafka_flavor.confluent`

Required:

- `confluent_cluster_id` (String) Confluent cluster identifier.
- `confluent_environment_id` (String) Confluent environment identifier.
- `key` (String, Sensitive) Confluent API key.
- `secret` (String, Sensitive) Confluent API secret.


<a id="nestedatt--spec--kafka_flavor--gateway"></a>
### Nested Schema for `spec.kafka_flavor.gateway`

Required:

- `password` (String, Sensitive) Conduktor Gateway Admin password.
- `url` (String) Conduktor Gateway Admin API URL.
- `user` (String) Conduktor Gateway Admin user.

Optional:

- `ignore_untrusted_certificate` (Boolean) Ignore untrusted certificate for Gateway Admin API.
- `virtual_cluster` (String) Conduktor Gateway Virtual cluster name (default `passthrough`).



<a id="nestedatt--spec--schema_registry"></a>
### Nested Schema for `spec.schema_registry`

Optional:

- `confluent_like` (Attributes) Confluent like schema registry configuration (see [below for nested schema](#nestedatt--spec--schema_registry--confluent_like))
- `glue` (Attributes) AWS Glue schema registry configuration (see [below for nested schema](#nestedatt--spec--schema_registry--glue))

<a id="nestedatt--spec--schema_registry--confluent_like"></a>
### Nested Schema for `spec.schema_registry.confluent_like`

Optional:

- `ignore_untrusted_certificate` (Boolean) Ignore untrusted certificate for schema registry. Only used if type is `ConfluentLike`
- `properties` (String) Schema registry properties. Only used if type is `ConfluentLike`
- `security` (Attributes) Confluent Schema registry security configuration. One of `basic_auth`, `bearer_token`, `ssl_auth`. If none provided, no security is used. (see [below for nested schema](#nestedatt--spec--schema_registry--confluent_like--security))
- `url` (String) Schema registry URL. Required if type is `ConfluentLike`

<a id="nestedatt--spec--schema_registry--confluent_like--security"></a>
### Nested Schema for `spec.schema_registry.confluent_like.security`

Optional:

- `basic_auth` (Attributes) Basic auth schema registry security configuration. (see [below for nested schema](#nestedatt--spec--schema_registry--confluent_like--security--basic_auth))
- `bearer_token` (Attributes) Bearer token schema registry security configuration. (see [below for nested schema](#nestedatt--spec--schema_registry--confluent_like--security--bearer_token))
- `ssl_auth` (Attributes) SSL auth (mTLS) schema registry security configuration. (see [below for nested schema](#nestedatt--spec--schema_registry--confluent_like--security--ssl_auth))

<a id="nestedatt--spec--schema_registry--confluent_like--security--basic_auth"></a>
### Nested Schema for `spec.schema_registry.confluent_like.security.basic_auth`

Required:

- `password` (String, Sensitive) Schema registry basic auth password.
- `username` (String) Schema registry basic auth username.


<a id="nestedatt--spec--schema_registry--confluent_like--security--bearer_token"></a>
### Nested Schema for `spec.schema_registry.confluent_like.security.bearer_token`

Required:

- `token` (String, Sensitive) Schema registry bearer token.


<a id="nestedatt--spec--schema_registry--confluent_like--security--ssl_auth"></a>
### Nested Schema for `spec.schema_registry.confluent_like.security.ssl_auth`

Required:

- `certificate_chain` (String) Schema registry SSL auth certificate chain PEM.
- `key` (String, Sensitive) Schema registry SSL auth private key PEM.




<a id="nestedatt--spec--schema_registry--glue"></a>
### Nested Schema for `spec.schema_registry.glue`

Required:

- `security` (Attributes) Schema registry configuration. One of `credentials`, `from_context`, `from_role`, `iam_anywhere` (see [below for nested schema](#nestedatt--spec--schema_registry--glue--security))

Optional:

- `region` (String) Glue Schema registry AWS region
- `registry_name` (String) Glue Schema registry name

<a id="nestedatt--spec--schema_registry--glue--security"></a>
### Nested Schema for `spec.schema_registry.glue.security`

Optional:

- `credentials` (Attributes) AWS credentials GLUE schema registry security configuration. (see [below for nested schema](#nestedatt--spec--schema_registry--glue--security--credentials))
- `from_context` (Attributes) AWS context GLUE schema registry security configuration. (see [below for nested schema](#nestedatt--spec--schema_registry--glue--security--from_context))
- `from_role` (Attributes) AWS role GLUE schema registry security configuration. (see [below for nested schema](#nestedatt--spec--schema_registry--glue--security--from_role))
- `iam_anywhere` (Attributes) AWS IAM Anywhere GLUE schema registry security configuration. (see [below for nested schema](#nestedatt--spec--schema_registry--glue--security--iam_anywhere))

<a id="nestedatt--spec--schema_registry--glue--security--credentials"></a>
### Nested Schema for `spec.schema_registry.glue.security.credentials`

Required:

- `access_key_id` (String, Sensitive) Glue Schema registry AWS access key ID.
- `secret_key` (String, Sensitive) Glue Schema registry AWS secret key.


<a id="nestedatt--spec--schema_registry--glue--security--from_context"></a>
### Nested Schema for `spec.schema_registry.glue.security.from_context`

Required:

- `profile` (String) Glue Schema registry AWS profile name.


<a id="nestedatt--spec--schema_registry--glue--security--from_role"></a>
### Nested Schema for `spec.schema_registry.glue.security.from_role`

Required:

- `role` (String) Glue Schema registry AWS role ARN.


<a id="nestedatt--spec--schema_registry--glue--security--iam_anywhere"></a>
### Nested Schema for `spec.schema_registry.glue.security.iam_anywhere`

Required:

- `certificate` (String) Glue Schema registry AWS certificate.
- `private_key` (String) Glue Schema registry AWS private key.
- `profile_arn` (String) Glue Schema registry AWS profile ARN.
- `role_arn` (String) Glue Schema registry AWS role ARN.
- `trust_anchor_arn` (String) Glue Schema registry AWS trust anchor ARN.







