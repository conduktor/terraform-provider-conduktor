---
page_title: "Condutkor : conduktor_kafka_cluster_v2 "
subcategory: "console/v2"
description: |-
    Resource for managing Conduktor Kafka cluster and Schema registry.
    This resource allows you to create, read, update and delete Kafka clusters and Schema registries in Conduktor.
---

# conduktor_kafka_cluster_v2

Resource for managing Conduktor Kafka cluster and Schema registry.
This resource allows you to create, read, update and delete Kafka clusters and Schema registries in Conduktor.

## Example Usage

### Simple Kafka cluster without Schema registry
This example creates a simple Kafka cluster without authentication resource and without Schema Registry.
```terraform
resource "conduktor_kafka_cluster_v2" "simple" {
  name = "simple-cluster"
  spec {
    display_name                 = "Simple kafka Cluster"
    icon                         = "kafka"
    color                        = "#000000"
    bootstrap_servers            = "localhost:9092"
    ignore_untrusted_certificate = true
  }
}
```

### Confluent Kafka cluster with Schema registry
This example creates a Confluent Kafka cluster and Schema Registry resource.
The Schema Registry authentication use mTLS.
```terraform
resource "conduktor_kafka_cluster_v2" "confluent" {
  name = "confluent-cluster"
  labels = {
    "env" = "staging"
  }
  spec {
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
      type                     = "Confluent"
      key                      = "yourApiKey123456"
      secret                   = "yourApiSecret123456"
      confluent_environment_id = "env-12345"
      confluent_cluster_id     = "lkc-67890"
    }
    schema_registry = {
      type                         = "ConfluentLike"
      url                          = "https://bbb-bbbb.us-west4.gcp.confluent.cloud:8081"
      ignore_untrusted_certificate = false
      security = {
        type              = "SSLAuth"
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

### Aiven Kafka cluster with Schema registry
This example creates an Aiven Kafka cluster and Schema Registry resource.
The Schema Registry authentication use basic auth.
```terraform
resource "conduktor_kafka_cluster_v2" "aiven" {
  name = "aiven-cluster"
  labels = {
    "env" = "test"
  }
  spec {
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
      type         = "Aiven"
      api_token    = "a1b2c3d4e5f6g7h8i9j0"
      project      = "my-kafka-project"
      service_name = "my-kafka-service"
    }
    schema_registry = {
      type                         = "ConfluentLike"
      url                          = "https://sr.aiven.io:8081"
      ignore_untrusted_certificate = false
      security = {
        type     = "BasicAuth"
        username = "uuuuuuu"
        password = "ppppppp"
      }
    }
  }
}
```

### AWS MSK with Glue Schema registry
This example creates an AWS MSK Kafka Cluster resource and a Glue Schema Registry.
```terraform
resource "conduktor_kafka_cluster_v2" "aws_msk" {
  name = "aws-cluster"
  labels = {
    "env" = "prod"
  }
  spec {
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
      type          = "Glue"
      region        = "eu-west-1"
      registry_name = "default"
      security = {
        type          = "Credentials"
        access_key_id = "accessKey"
        secret_key    = "secretKey"
      }
    }
  }
}
```

### Conduktor Gateway Kafka cluster with Schema registry
This example creates a Conduktor Gateway Kafka Cluster resource and a Schema Registry.
The Schema Registry authentication use bearer token.
```terraform
resource "conduktor_kafka_cluster_v2" "gateway" {
  name = "gateway-cluster"
  labels = {
    "env" = "prod"
  }
  spec {
    display_name      = "Gateway Cluster"
    bootstrap_servers = "gateway:9092"
    properties = {
      "sasl.jaas.config"  = "org.apache.kafka.common.security.plain.PlainLoginModule required username='admin' password='admin-secret';"
      "security.protocol" = "SASL_SSL"
      "sasl.mechanism"    = "PLAIN"
    }
    icon                         = "shield-blank"
    ignore_untrusted_certificate = true
    kafka_flavor = {
      type                         = "Gateway"
      url                          = "http://gateway:8088"
      user                         = "admin"
      password                     = "admin"
      virtual_cluster              = "vc1"
      ignore_untrusted_certificate = true
    }
    schema_registry = {
      type                         = "ConfluentLike"
      url                          = "http://localhost:8081"
      ignore_untrusted_certificate = true
      security = {
        type  = "BearerToken"
        token = "auth-token"
      }
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Kafka cluster name, must be unique, act as ID for import

### Optional

- `labels` (Map of String) Kafka cluster labels
- `spec` (Block, Optional) (see [below for nested schema](#nestedblock--spec))

<a id="nestedblock--spec"></a>
### Nested Schema for `spec`

Required:

- `bootstrap_servers` (String) List of bootstrap servers for the Kafka cluster separated by comma
- `display_name` (String) Kafka cluster display name

Optional:

- `color` (String) Kafka cluster icon color in hexadecimal format like `#FF0000`
- `icon` (String) Kafka cluster icon. List of available icons can be found [here](https://docs.conduktor.io/platform/reference/resource-reference/console/#icon-sets)
- `ignore_untrusted_certificate` (Boolean) Ignore untrusted certificate for Kafka cluster
- `kafka_flavor` (Attributes) Schema registry configuration (see [below for nested schema](#nestedatt--spec--kafka_flavor))
- `properties` (Map of String) Kafka cluster properties
- `schema_registry` (Attributes) Schema registry configuration (see [below for nested schema](#nestedatt--spec--schema_registry))

<a id="nestedatt--spec--kafka_flavor"></a>
### Nested Schema for `spec.kafka_flavor`

Required:

- `type` (String) Kafka provider type : `Confluent`, `Aiven`, `Gateway`. More detail on our [documentation](https://docs.conduktor.io/platform/reference/resource-reference/console/#kafka-provider)

Optional:

- `api_token` (String, Sensitive) Aiven API token. Required if type is `Aiven`
- `confluent_cluster_id` (String) Confluent cluster identifier. Required if type is `Confluent`
- `confluent_environment_id` (String) Confluent environment identifier. Required if type is `Confluent`
- `ignore_untrusted_certificate` (Boolean) Ignore untrusted certificate for Gateway Admin API. Only used if type is `Gateway`
- `key` (String, Sensitive) Confluent API key. Required if type is `Confluent`
- `password` (String, Sensitive) Conduktor Gateway Admin password. Required if type is `Gateway`
- `project` (String) Aiven project name. Required if type is `Aiven`
- `secret` (String, Sensitive) Confluent API secret. Required if type is `Confluent`
- `service_name` (String) Aiven service name. Required if type is `Aiven`
- `url` (String) Conduktor Gateway Admin API URL. Required if type is `Gateway`
- `user` (String) Conduktor Gateway Admin user. Required if type is `Gateway`
- `virtual_cluster` (String) Conduktor Gateway Virtual cluster name. Only used if type is `Gateway`


<a id="nestedatt--spec--schema_registry"></a>
### Nested Schema for `spec.schema_registry`

Required:

- `security` (Attributes) Schema registry configuration. Required if type is `ConfluentLike` or `Glue` (see [below for nested schema](#nestedatt--spec--schema_registry--security))
- `type` (String) Schema registry type valid values are: `ConfluentLike`, `Glue`

More detail on our [documentation](https://docs.conduktor.io/platform/reference/resource-reference/console/#schema-registry)

Optional:

- `ignore_untrusted_certificate` (Boolean) Ignore untrusted certificate for schema registry. Only used if type is `ConfluentLike`
- `properties` (String) Schema registry properties. Only used if type is `ConfluentLike`
- `region` (String) Glue Schema registry AWS region. Required if type is `Glue`
- `registry_name` (String) Glue Schema registry name. Only used if type is `Glue`
- `url` (String) Schema registry URL. Required if type is `ConfluentLike`

<a id="nestedatt--spec--schema_registry--security"></a>
### Nested Schema for `spec.schema_registry.security`

Required:

- `type` (String) Schema registry security type. Required if type is `ConfluentLike` or `Glue`.

Valid values are:

- For **ConfluentLike** : `NoSecurity`, `BasicAuth`, `BearerToken`, `SSLAuth` 

- For **Glue** : `Credentials`, `FromContext`, `FromRole`, `IAMAnywhere`

 More detail on our [documentation](https://docs.conduktor.io/platform/reference/resource-reference/console/#schema-registry)

Optional:

- `access_key_id` (String, Sensitive) Glue Schema registry AWS access key ID. Required if type is Glue with security `Credentials`
- `certificate` (String) Glue Schema registry AWS certificate. Required if type is Glue with security `IAMAnywhere`
- `certificate_chain` (String) Schema registry SSL auth certificate chain PEM. Required if security type is `SSLAuth`
- `key` (String, Sensitive) Schema registry SSL auth private key PEM. Required if security type is `SSLAuth`
- `password` (String, Sensitive) Schema registry basic auth password. Required if security type is `BasicAuth`
- `private_key` (String) Glue Schema registry AWS private key. Required if type is Glue with security `IAMAnywhere`
- `profile` (String) Glue Schema registry AWS profile name. Required if type is Glue with security `FromContext`
- `profile_arn` (String) Glue Schema registry AWS profile ARN. Required if type is Glue with security `IAMAnywhere`
- `role` (String) Glue Schema registry AWS role ARN. Required if type is Glue with security `FromRole`
- `role_arn` (String) Glue Schema registry AWS role ARN. Required if type is Glue with security `IAMAnywhere`
- `secret_key` (String, Sensitive) Glue Schema registry AWS secret key. Required if type is Glue with security `Credentials`
- `token` (String, Sensitive) Schema registry bearer token. Required if security type is `BearerToken`
- `trust_anchor_arn` (String) Glue Schema registry AWS trust anchor ARN. Required if type is Glue with security `IAMAnywhere`
- `username` (String) Schema registry basic auth username. Required if security type is `BasicAuth`





