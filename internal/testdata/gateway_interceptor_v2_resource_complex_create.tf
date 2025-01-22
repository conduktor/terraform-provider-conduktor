#
# resource "conduktor_gateway_interceptor_v2" "test-encryption" {
#   name = "myEncryptPlugin"
#   scope {
#     vcluster = "passthrough"
#   }
#   spec {
#     plugin_class = "io.conduktor.gateway.interceptor.EncryptSchemaBasedPlugin"
#     priority     = 1
#     config = jsonencode({
#       "schemaDataMode" = "convert_json"
#       "kmsConfig" = {
#       }
#       "tags" = ["PII", "ENCRYPTION"]
#       "defaultAlgorithm" = "AES128_EAX"
#       "defaultKeySecretId" = "myDefaultKeySecret"
#       "namespace" = "conduktor."
#     })
#   }
# }

resource "conduktor_gateway_interceptor_v2" "test-encryption" {
  name = "myEncryptPlugin"
  spec {
    plugin_class = "io.conduktor.gateway.interceptor.EncryptPlugin"
    priority     = 1
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
