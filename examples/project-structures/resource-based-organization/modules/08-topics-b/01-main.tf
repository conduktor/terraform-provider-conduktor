resource "conduktor_console_topic_v2" "website-analytics-admin-json" {
  name    = "website-analytics.admin.notJsonOrAvro"
  cluster = "my-cluster"
  labels = {
    data-criticality                    = "",
    replication                         = "none",
    url                                 = "http:stu.com",
    "conduktor.io/application"          = "website-analytics"
    "conduktor.io/application-instance" = "website-analytics-dev"
  }
  description = "# Admin logs for website analytics team\nInformation about admin activities in the website-analytics domain.\nThis description is not editable from the UI."
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
