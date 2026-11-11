resource "conduktor_console_topic_template_v2" "simple" {
  name = "simple-template"
  spec = {
    display_name = "Simple Template"
    defaults = {
      spec = {
        partitions         = 3
        replication_factor = 1
        configs = {
          "retention.ms" = "604800000"
        }
      }
    }
  }
}
