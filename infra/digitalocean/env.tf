locals {
  env_vars = [
    {
      // This is a bindable variable that will be replaced with the actual URL. See:
      // https://docs.digitalocean.com/products/app-platform/how-to/use-environment-variables
      key   = "SERVER_URL",
      value = format("$%s", "{APP_URL}"),
      scope = "RUN_TIME"
    },
    {
      key   = "MAIN_DOMAIN",
      value = ".pacia.tech",
      scope = "RUN_TIME",
    },
    {
      key   = "FRONTEND_URL",
      value = "https://bee-ci.pacia.tech",
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
      value = digitalocean_database_cluster.postgres.host
      scope = "RUN_TIME"
    },
    {
      key   = "DB_PORT"
      value = digitalocean_database_cluster.postgres.port
      scope = "RUN_TIME"
    },
    {
      key   = "DB_USER"
      value = digitalocean_database_cluster.postgres.user
      scope = "RUN_TIME"
    },
    {
      key   = "DB_PASSWORD"
      value = digitalocean_database_cluster.postgres.password
      scope = "RUN_TIME"
      type  = "SECRET"
    },
    {
      key   = "DB_NAME"
      value = digitalocean_database_db.postgres.name
      scope = "RUN_TIME"
    },
    {
      key   = "DB_OPTS"
      value = "sslmode=require"
      scope = "RUN_TIME"
    },
    {
      key   = "INFLUXDB_URL"
      value = format("http://%s:8086", digitalocean_droplet.influxdb.ipv4_address),
      scope = "RUN_TIME"
    },
    {
      key   = "INFLUXDB_TOKEN"
      value = var.influxdb_token
      scope = "RUN_TIME"
      type  = "SECRET"
    },
    {
      key   = "INFLUXDB_ORG"
      value = var.influxdb_org
      scope = "RUN_TIME"
    },
    {
      key   = "INFLUXDB_BUCKET"
      value = var.influxdb_bucket
      scope = "RUN_TIME"
    },
    {
      key   = "REDIS_ADDRESS"
      value = format("%s:%s", digitalocean_database_cluster.redis.host, digitalocean_database_cluster.redis.port)
      scope = "RUN_TIME"
    },
    {
      key   = "REDIS_PASSWORD"
      value = digitalocean_database_cluster.redis.password
      scope = "RUN_TIME"
      type  = "SECRET"
    },
    {
      key   = "REDIS_USE_TLS",
      value = "true"
      scope = "RUN_TIME"
    }
  ]
}
