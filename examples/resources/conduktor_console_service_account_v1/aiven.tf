

resource "conduktor_console_service_account_v1" "aiven_sa" {
  name    = "aiven-service-account"
  cluster = "aiven-cluster"
  spec = {
    authorization = {
      aiven = {
        acls = [
          {
            resource_type = "TOPIC"
            name          = "click.event-stream.avro"
            permission    = "readwrite"
          },
          {
            resource_type = "TOPIC"
            name          = "public*"
            permission    = "read"
          },
          {
            resource_type = "SCHEMA"
            name          = "Subject:click.event-stream.avro"
            permission    = "schema_registry_write"
          },
        ]
      }
    }
  }
}
