output "default_ingress" {
  description = "The default URL to access the app"
  value       = digitalocean_app.app.default_ingress
}

output "live_url" {
  description = "The live URL for the app"
  value       = digitalocean_app.app.live_url
}

output "db_postgres_uri" {
  description = "Connection string for Postgres database cluster"
  value       = digitalocean_database_cluster.postgres.uri
  sensitive   = true
}

output "db_redis_uri" {
  description = "Connection string for Redis database cluster"
  value       = digitalocean_database_cluster.redis.uri
  sensitive   = true
}

output "droplet_executor_ip" {
  value = digitalocean_droplet.executor.ipv4_address
}

output "droplet_influxdb_ip" {
  value = digitalocean_droplet.influxdb.ipv4_address
}
