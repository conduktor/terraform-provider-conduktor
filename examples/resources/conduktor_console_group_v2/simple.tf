resource "conduktor_console_group_v2" "example" {
  name = "simple-group"
  spec {
    display_name = "Simple Group"
    description  = "Simple group description"
  }
}
