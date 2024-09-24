---
page_title: "Provider: Conduktor"
subcategory: ""
description: |-
    The Conduktor provider is used to interact with the resources supported by Conduktor. The provider needs to be configured with the proper credentials before it can be used.
---

# Conduktor Provider

The Conduktor provider is used to interact with the resources supported by Conduktor. The provider needs to be configured with the proper credentials before it can be used.

> [!WARNING]
> - The Conduktor Terraform provider is currently in **Alpha**.
> - It does not support all Console and Gateway resources yet. See our [resources roadmap](https://github.com/conduktor/terraform-provider-conduktor/blob/main/README.md#resources-roadmap).
> - Let us know if you have [feedback](https://product.conduktor.help/c/74-terraform-provider) or wish to be a design partner.

## Example Usage

```terraform
provider "conduktor" {
  # mandatory console URL
  console_url = "http://localhost:8080" # or env vars CDK_CONSOLE_URL or CDK_BASE_URL

  # authentication either with api token or admin credentials
  api_token = "your-api-token" # or env var CDK_API_TOKEN or CDK_API_KEY
  #admin_email    = "admin@my-org.com" # or env var CDK_ADMIN_EMAIL
  #admin_password = "admin-password"   # or env var CDK_ADMIN_PASSWORD

  # optional http client TLS configuration
  cert     = file("path/to/cert.pem") # or env var CDK_CERT
  insecure = true                     # or env var CDK_INSECURE

  # optional authentication via certificate
  key    = file("path/to/key.pem") # or env var CDK_KEY
  cacert = file("path/to/ca.pem")  # or env var CDK_CA_CERT
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `admin_email` (String) The email of the admin user. May be set using environment variable `CDK_ADMIN_EMAIL`. Required if admin_password is set. If not provided, the API token will be used to authenticate.
- `admin_password` (String, Sensitive) The password of the admin user. May be set using environment variable `CDK_ADMIN_PASSWORD`. Required if admin_email is set. If not provided, the API token will be used to authenticater.
- `api_token` (String, Sensitive) The API token to authenticate with the Conduktor API. May be set using environment variable `CDK_API_TOKEN` or `CDK_API_KEY`. If not provided, admin_email and admin_password will be used to authenticate. See [documentation](https://docs.conduktor.io/platform/reference/api-reference/#generate-an-api-key) for more information.
- `cacert` (String) Root CA certificate in PEM format to verify the Conduktor Console certificate. May be set using environment variable `CDK_CACERT`. If not provided, the system's root CA certificates will be used.
- `cert` (String) Cert in PEM format to authenticate using client certificates. May be set using environment variable `CDK_CERT`. Must be used with key. If key is provided, cert is required. Useful when Console behind a reverse proxy with client certificate authentication.
- `console_url` (String) The URL of the Conduktor Console. May be set using environment variable `CDK_BASE_URL` or `CDK_CONSOLE_URL`. Required either here or in the environment.
- `insecure` (Boolean) Skip TLS verification flag. May be set using environment variable `CDK_INSECURE`.
- `key` (String) Key in PEM format to authenticate using client certificates. May be set using environment variable `CDK_KEY`. Must be used with cert. If cert is provided, key is required. Useful when Console behind a reverse proxy with client certificate authentication.

