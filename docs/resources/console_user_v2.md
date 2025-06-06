---
page_title: "Conduktor : conduktor_console_user_v2 "
subcategory: "iam/v2"
description: |-
    Resource for managing Conduktor users.
    This resource allows you to create, read, update and delete users in Conduktor.
---

# conduktor_console_user_v2

Resource for managing Conduktor users.
This resource allows you to create, read, update and delete users in Conduktor.

## Example Usage

### Simple user without permissions
```terraform
resource "conduktor_console_user_v2" "example" {
  name = "bob@company.io"
  spec = {
    firstname = "Bob"
    lastname  = "Smith"
  }
}
```

### Complex user with permissions
```terraform
resource "conduktor_console_user_v2" "example" {
  name = "bob@company.io"
  spec = {
    firstname = "Bob"
    lastname  = "Smith"
    permissions = [
      {
        resource_type = "PLATFORM"
        permissions   = ["userView", "datamaskingView", "auditLogView"]
      },
      {
        resource_type = "TOPIC"
        name          = "test-topic"
        cluster       = "*"
        pattern_type  = "LITERAL"
        permissions   = ["topicViewConfig", "topicConsume", "topicProduce"]
      },
      {
        resource_type = "KSQLDB"
        cluster       = "*"
        ksqldb        = "*"
        permissions   = ["ksqldbAccess"]
      },
    ]
  }
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) User email, must be unique, acts as an ID for import
- `spec` (Attributes) User specification (see [below for nested schema](#nestedatt--spec))

<a id="nestedatt--spec"></a>
### Nested Schema for `spec`

Optional:

- `firstname` (String) User firstname
- `lastname` (String) User lastname
- `permissions` (Attributes Set) Set of all user permissions (see [below for nested schema](#nestedatt--spec--permissions))

<a id="nestedatt--spec--permissions"></a>
### Nested Schema for `spec.permissions`

Required:

- `permissions` (Set of String) Set of all permissions to apply on the resource. See https://docs.conduktor.io/platform/reference/resource-reference/console/#permissions for more details
- `resource_type` (String) Type of the resource to apply permission on valid values are: CLUSTER, CONSUMER_GROUP, KAFKA_CONNECT, KSQLDB, PLATFORM, SUBJECT, TOPIC

Optional:

- `cluster` (String) Name of the cluster to apply permission, only required if resource_type is TOPIC, SUBJECT, CONSUMER_GROUP, KAFKA_CONNECT, KSQLDB
- `kafka_connect` (String) Name of the Kafka Connect to apply permission, only required if resource_type is KAFKA_CONNECT
- `ksqldb` (String) Name of a valid Kafka Connect cluster, only required if resource_type is KSQLDB
- `name` (String) Name of the resource to apply permission to could be a topic, a cluster, a consumer group, etc. depending on resource_type
- `pattern_type` (String) Type of the pattern to apply permission on valid values are: LITERAL, PREFIXED




