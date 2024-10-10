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

variable "influxdb_org" {
  type = string
}

variable "influxdb_bucket" {
  type = string
}

variable "influxdb_password" {
  type      = string
  sensitive = true
}

variable "influxdb_token" {
  type      = string
  sensitive = true
}

# InfluxDB-specific variables

variable "INFLUXDB_PASSWORD" {
  description = "InfluxDB password"
  type        = string
}



# Provider-specific variables

variable "do_token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}
