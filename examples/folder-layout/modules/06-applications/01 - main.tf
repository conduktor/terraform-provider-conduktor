# web analytics
resource "conduktor_console_application_v1" "website-analytics" {
  # name = "website-analytics"
  name = var.application_name
  spec = {
    title       = "Website Analytics"
    description = "Application for streaming web analytics"
    owner       = "website-analytics-team"
  }
}