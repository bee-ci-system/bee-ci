locals {
  env_vars = [
    {
      key   = "PORT"
      value = "8080"
      scope = "RUN_TIME"
    },
    {
      key   = "GITHUB_APP_ID"
      value = var.github_app_id
      scope = "RUN_TIME"
    },
    {
      key   = "GITHUB_APP_CLIENT_ID"
      value = var.github_app_client_id
      scope = "RUN_TIME"
    },
    {
      key   = "GITHUB_APP_WEBHOOK_SECRET"
      value = var.github_app_webhook_secret
      scope = "RUN_TIME"
      type  = "SECRET"
    },
    {
      key   = "GITHUB_APP_PRIVATE_KEY_BASE64"
      value = var.github_app_private_key_base64
      scope = "RUN_TIME"
      type  = "SECRET"
    },
    {
      key   = "GITHUB_APP_CLIENT_SECRET"
      value = var.github_app_client_secret
      scope = "RUN_TIME"
      type  = "SECRET"
    },
    {
      key   = "DB_HOST"
      value = digitalocean_database_cluster.main-db-cluster.host
      scope = "RUN_TIME"
    }
    , {
      key   = "DB_PORT"
      value = digitalocean_database_cluster.main-db-cluster.port
      scope = "RUN_TIME"
    },
    {
      key   = "DB_USER"
      value = digitalocean_database_cluster.main-db-cluster.user
      scope = "RUN_TIME"
    },
    {
      key   = "DB_PASSWORD"
      value = digitalocean_database_cluster.main-db-cluster.password
      scope = "RUN_TIME"
      type  = "SECRET"
    },
    {
      key   = "DB_NAME"
      value = digitalocean_database_db.main-db.name
      scope = "RUN_TIME"
    },
    {
      key   = "DB_OPTS"
      value = "sslmode=require"
      scope = "RUN_TIME"
    }
  ]
}