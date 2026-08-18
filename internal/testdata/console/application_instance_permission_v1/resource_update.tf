
resource "conduktor_console_application_instance_permission_v1" "test" {
  name         = "appinstance-permission"
  application  = "myapp"
  app_instance = "my-app-instance"
  spec = {
    resource = {
      type         = "TOPIC"
      name         = "my-topic"
      pattern_type = "LITERAL"
    }
    user_permission            = "READ"
    service_account_permission = "NONE"
    granted_to                 = "my-app-instance"
  }
}
