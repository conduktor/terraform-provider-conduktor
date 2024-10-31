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
