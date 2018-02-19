variable "name" {
  default = "strava-commute"
}

variable "http_methods" {
  type    = "list"
  default = ["GET", "POST"]
}

variable "strava_api_token" {}
