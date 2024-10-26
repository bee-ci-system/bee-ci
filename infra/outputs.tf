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

output "db_influx_config_cmd" {
  description = "Command to create a configuration for InfluxDB database"
  value       = "influx config create -n bee_influx_prod -u http://${digitalocean_droplet.influxdb.ipv4_address}:8086 -o ${var.influxdb_org} -t ${var.influxdb_token} -a"
  sensitive   = true
}

output "droplet_executor_ip" {
  description = "IP address of the droplet running executor. Usually used to `ssh` into the droplet for debugging."
  value       = digitalocean_droplet.executor.ipv4_address
}

output "droplet_influxdb_ip" {
  description = "IP address of the droplet running InfluxDB. Usually used to `ssh` into the droplet for debugging."
  value       = digitalocean_droplet.influxdb.ipv4_address
}
