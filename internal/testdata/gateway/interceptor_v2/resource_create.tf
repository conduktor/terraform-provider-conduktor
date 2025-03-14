
resource "conduktor_gateway_interceptor_v2" "topic-policy-default" {
  name = "enforce-partition-limit-test"
  scope = {
    vcluster = "passthrough"
  }
  spec = {
    plugin_class = "io.conduktor.gateway.interceptor.safeguard.CreateTopicPolicyPlugin"
    priority     = 1
    config = jsonencode({
      topic = "myprefix-.*"
      numPartition = {
        min    = 5
        max    = 5
        action = "INFO"
      }
    })
  }
}

resource "conduktor_gateway_interceptor_v2" "topic-policy-vcluster_sa" {
  name = "enforce-partition-limit-test"
  scope = {
    vcluster = "vcluster_sa"
    username = "my.user"
  }
  spec = {
    plugin_class = "io.conduktor.gateway.interceptor.safeguard.CreateTopicPolicyPlugin"
    priority     = 4
    config = jsonencode({
      topic = "other-.*"
      numPartition = {
        min    = 3
        max    = 6
        action = "INFO"
      }
    })
  }
  # force dependency on the default topic policy to force sequential creation of the interceptors with same name
  # due to a limitaion in current Gateway API
  depends_on = [conduktor_gateway_interceptor_v2.topic-policy-default]
}

resource "conduktor_gateway_interceptor_v2" "schema-encryption" {
  name = "schema-encryption"
  scope = {
    vcluster = "vcluster_sa" # should already be setup by init
    group    = "group-a"     # should already be setup by init
  }
  spec = {
    plugin_class = "io.conduktor.gateway.interceptor.EncryptSchemaBasedPlugin"
    priority     = 2
    config = jsonencode({
      "schemaDataMode" = "convert_json"
      "kmsConfig" = {
      }
      "tags"               = ["PII", "ENCRYPTION"]
      "defaultAlgorithm"   = "AES128_EAX"
      "defaultKeySecretId" = "in-memory-kms://myDefaultKeySecret"
      "namespace"          = "conduktor."
    })
  }
}

resource "conduktor_gateway_interceptor_v2" "full-encryption" {
  name = "full-encryption"
  spec = {
    plugin_class = "io.conduktor.gateway.interceptor.EncryptPlugin"
    priority     = 3
    config = jsonencode({
      "topic" = "full-encrypt.*"
      "kmsConfig" = {
        "aws" = {
          "basicCredentials" = {
            "accessKey" = "test"
            "secretKey" = "test"
          }
        }
      }

      "recordValue" = {
        "payload" = {
          "keySecretId" = "aws-kms://test-arn"
          "algorithm"   = "AES128_GCM"
        }
      }
    })
  }
}

resource "conduktor_gateway_interceptor_v2" "datamasking" {
  name = "mask-sensitive-fields"
  spec = {
    plugin_class = "io.conduktor.gateway.interceptor.FieldLevelDataMaskingPlugin"
    priority     = 100
    config = jsonencode({
      "topic" = "^[A-Za-z]*_masked$"
      "policies" = [
        {
          "name" = "Mask credit card"
          "rule" = {
            "type" = "MASK_ALL"
          },
          "fields" = ["profile.creditCardNumber", "contact.email"]
        },
        {
          "name" = "Partial mask phone"
          "rule" = {
            "type"          = "MASK_FIRST_N"
            "maskingChar"   = "*"
            "numberOfChars" = 9
          },
          "fields" = ["contact.phone"]
        }
      ]
    })
  }
}

