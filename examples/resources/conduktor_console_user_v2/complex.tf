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
      }
    ]
  }
}
