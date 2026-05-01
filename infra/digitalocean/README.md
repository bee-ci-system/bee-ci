# DigitalOcean Terraform

This Terraform root is intended to run as its own Spacelift stack.

## Spacelift stack

- Project root: `infra/digitalocean`
- Terraform workflow tool: Terraform
- Managed state: enabled

Do not configure a backend in this root when using Spacelift-managed state.
Spacelift injects the backend configuration during runs.

## Required variables

Set these as stack environment variables or attach them through a Spacelift context:

- `TF_VAR_do_token`
- `TF_VAR_github_app_id`
- `TF_VAR_github_app_webhook_secret`
- `TF_VAR_github_app_private_key_base64`
- `TF_VAR_github_app_client_id`
- `TF_VAR_github_app_client_secret`
- `TF_VAR_influxdb_password`
- `TF_VAR_influxdb_token`

The `setup_credentials` script remains useful for local Terraform runs because it writes `terraform.tfvars` from 1Password.
