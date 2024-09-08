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

# Provider-specific variables
variable "project_id" {
  description = "GCP Project ID"
  type        = string
  default     = "bee-ci"
}

# Asking for a project_id every time is annyoing. But I had some problems with the below.
# resource "random_string" "project_id" {
#   length  = 8
#   special = false
#   upper   = false
# }

variable "credentials_file" {
  description = "Path to the GCP credentials file"
  type        = string
}

variable "billing_account_id" {
  description = "Can be found with `gcloud billing accounts list`"
  type        = string
}

variable "region" {
  default = "us-central1"
}

variable "zone" {
  default = "us-central1-c"
}
