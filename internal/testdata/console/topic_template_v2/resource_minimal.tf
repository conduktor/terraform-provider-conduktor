resource "conduktor_console_topic_template_v2" "minimal" {
  name = "minimal-template"
  spec = {
    display_name = "Minimal Template"
    defaults     = {}
  }
}
