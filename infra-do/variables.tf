variable "do_token" {
  description = "DigitalOcean API token"
  sensitive   = true
}

variable "github_app_id" {
  description = "GitHub App ID"
  type        = string
}

variable "github_app_webhook_secret" {
  description = "GitHub App Webhook Secret"
  type        = string
}

variable "github_app_private_key_base64" {
  description = "GitHub App Private Key Base64"
  type        = string
}

variable "github_app_client_id" {
  description = "GitHub App Client ID"
  type        = string
}

variable "github_app_client_secret" {
  description = "GitHub App Client Secret"
  type        = string
}
