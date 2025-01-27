resource "conduktor_gateway_interceptor_v2" "field-encryption" {
  name = "field-encryption"
  spec {
    plugin_class = "io.conduktor.gateway.interceptor.EncryptPlugin"
    priority     = 1
    config = jsonencode({
      "topic" = "encrypt.*"
      "kmsConfig" = {
        "vault" = {
          "uri"     = "http://vault:8200"
          "token"   = "test"
          "version" = 1
        }
      }
      "recordValue" = {
        "fields" = [
          {
            "fieldName"   = "password"
            "keySecretId" = "vault-kms://vault:8200/transit/keys/password-secret"
            "algorithm"   = "AES128_GCM"
          },
          {
            "fieldName"   = "visa"
            "keySecretId" = "vault-kms://vault:8200/transit/keys/{{record.header.test-header}}-visa-secret-{{record.key}}-{{record.value.username}}-{{record.value.education.account.accountId}}"
            "algorithm"   = "AES128_GCM"
          },
          {
            "fieldName"   = "education.account.username"
            "keySecretId" = "in-memory-kms://myDefaultKeySecret"
            "algorithm"   = "AES128_GCM"
          }
        ]
      }
    })
  }
}
