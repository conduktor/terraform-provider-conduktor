terraform {
  required_version = ">= v1.10.4"
  required_providers {
    conduktor = {
      # source  = "terraform.local/conduktor/conduktor" # local provider
      source  = "conduktor/conduktor"
      version = ">= 0.4.0"
    }
  }
}

provider "conduktor" {
  alias          = "console"
  mode           = "console"
  base_url       = "http://localhost:8080" # or env vars CDK_CONSOLE_URL or CDK_BASE_URL
  # admin_user     = "admin@conduktor.io"    # or env var CDK_ADMIN_EMAIL
  # admin_password = "testP4ss!"             # or env var CDK_ADMIN_PASSWORD
  api_token = "D7AIHHW2cr4=.Nogw352hIuhiGtPy5laOG5s1gdUOxGX52tcY+qtyS++Y97t+yMK2j7A2VE9pwQxk+SVFBFWrMJoZiqFVrEI5mvLXBovAWXMKtjCWQVohkevQDAtZNxc947sUcSrWUHeONJCX9UIKjQ9tdSES7LU+DPg/+DUTgxvqU+meCNnhCJLXPN6zBNCsw5kIHct89tSrFctHGExZ9cBkvNN9BozsaZVWuJGjjHZ6TO/o0ZwsyTWMO3kgqMMyYz9U2+rvHdeb4Uazg0kObdl4AcVPbQCvQq2ds6O6/hd+lv57XTJgIuAGC/LRYifZigJhlHoAVg6q3sSWm9BMUqYhZOIWKYaWRg==" # admin
  # api_token = "B6J/NscIEHA=.9tFH9oCT2fTHDr4RsPuuXR5zq40JCRb2ss3zCcK/0+Ghbi3jeyDuxcveHa7QK/JtD+9WyCNMqZVTluDVnLS/e56uPNVLt/kQvA1RNpSSmQMKOAF01RE7ylJU1/r/gOVIHFZih971UlhoLkkG6+s28lh8JqpRgP6u1pYGfu4WMcCTOzaiGvO4E4qcMuLrFvoaNfeXdnPLmuom8lsnpYfGOA43nNObWNfP7hmevHemaAjvJ9AU5Os7Iv61G+7Lv0njpt24cAqlzPXWIR5wLJ3KymeDmQIOeoUtlLjmiSZgbrPP/rpHkY36WNAPOi42eM6SWwkIDQ9oaWImVyCi6Cd1jQ==" # web-dev

  insecure = true # or env var CDK_INSECURE
}

provider "conduktor" {
  alias          = "gateway"
  mode           = "gateway"
  base_url       = "http://localhost:8888"
  admin_user     = "admin"
  admin_password = "conduktor"

  insecure = true # or env var CDK_INSECURE
}
