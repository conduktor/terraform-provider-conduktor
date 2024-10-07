
resource "conduktor_user_v2" "test" {
  name = "pam.beesly@dunder.mifflin.com"
  spec {
    firstname = "Pam"
    lastname  = "Beesly"
    permissions = [
      {
        resource_type = "TOPIC"
        permissions   = ["topicViewConfig", "topicConsume", "topicProduce"]
        name          = "team1.test-topic"
        pattern_type  = "LITERAL"
        cluster       = "*"
      }
    ]
  }
}
