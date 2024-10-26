# GitHub App-specific variables

variable "github_app_id" {
  description = "GitHub App ID"
  type        = string
}

variable "github_app_webhook_secret" {
  description = "GitHub App Webhook Secret"
  type        = string
  sensitive   = true
}

variable "github_app_private_key_base64" {
  description = "GitHub App Private Key Base64"
  type        = string
  sensitive   = true
}

variable "github_app_client_id" {
  description = "GitHub App Client ID"
  type        = string
}

variable "github_app_client_secret" {
  description = "GitHub App Client Secret"
  type        = string
  sensitive   = true
}

# InfluxDB-specific variables

variable "influxdb_user" {
  type    = string
  default = "beeci"
}

variable "influxdb_org" {
  type    = string
  default = "beeci"
}

variable "influxdb_bucket" {
  type    = string
  default = "home"
}

variable "influxdb_password" {
  type      = string
  sensitive = true
}

variable "influxdb_token" {
  type      = string
  sensitive = true
}

# Provider-specific variables

variable "do_token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}
