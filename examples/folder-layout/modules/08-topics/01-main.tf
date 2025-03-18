# resource "conduktor_console_topic_v2" "production_topic" {
#   name    = "production-topic"
#   cluster = "kafka-cluster"
#   spec = {
#     partitions         = 10
#     replication_factor = 1
#   }

#   lifecycle {
#     prevent_destroy = true
#   }
# }

# resource "conduktor_console_topic_v2" "simple" {
#   name    = "simple"
#   cluster = "kafka-cluster"
#   labels = {
#     domain = "clickstream"
#   }
#   description = "# Simple kafka topic"
#   spec = {
#     partitions         = 3
#     replication_factor = 1
#     configs = {
#       "cleanup.policy" = "delete"
#     }
#   }
# }

# resource "conduktor_console_topic_v2" "complex" {
#   name    = "complex"
#   cluster = "kafka-cluster"
#   labels = {
#     domain  = "clickstream"
#     appcode = "clk"
#   }
#   catalog_visibility      = "PRIVATE"
#   description_is_editable = false
#   description             = "# Complex kafka topic"
#   sql_storage = {
#     retention_time_in_second = 60000
#     enabled                  = true
#   }
#   spec = {
#     partitions         = 3
#     replication_factor = 1
#     configs = {
#       "cleanup.policy" = "delete",
#       "retention.ms"   = "60000"
#     }
#   }
# }


# resource "conduktor_console_topic_v2" "stu-stupid-topic" {
#   catalog_visibility      = "PUBLIC"
#   cluster                 = "kafka-cluster"
#   description             = "My desctipoion"
#   description_is_editable = true
#   labels = {
#     key = "valuepair"
#   }
#   name = "stu-stupid-topic"
#   spec = {
#     configs = {
#       "cleanup.policy" = "delete"
#     }
#     partitions         = 5
#     replication_factor = 1
#   }
# }

resource "conduktor_console_topic_v2" "website-analytics-admin-json" {
  name    = "website-analytics.admin.csv"
  cluster = "my-cluster"
  labels = {
    data-criticality                    = "",
    replication                         = "none",
    url                                 = "http:stu.com",
    "conduktor.io/application"          = "website-analytics"
    "conduktor.io/application-instance" = "website-analytics-dev"
  }
  description = "# Admin logs for website analytics team\nInformation about admin activities in the website-analytics domain.\nhis description is not editable from the UI."
  spec = {
    partitions         = 3
    replication_factor = 1
    configs = {
      "retention.ms"   = "604800000",
      "cleanup.policy" = "delete"
    }
  }

  lifecycle {
    prevent_destroy = true
  }
}