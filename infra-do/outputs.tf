output "db_uri" {
  description = "The full URI for connecting to the database cluster"
  value       = digitalocean_database_cluster.postgres.uri
  sensitive   = true
}

output "default_ingress" {
  description = "The default URL to access the app"
  value       = digitalocean_app.app.default_ingress
}

output "live_url" {
  description = "The live URL for the app"
  value       = digitalocean_app.app.live_url
}

output "live_domain" {
  description = "The live domain for the app"
  value       = digitalocean_app.app.live_domain
}

# output "droplet_ip" {
#  value = digitalocean_droplet.influxdb.ipv4_address
#}

output "volume_id" {
  value = digitalocean_volume.influxdb_volume.id
}

output "droplet_executor_ip" {
  value = digitalocean_droplet.executor.ipv4_address
}

output "droplet_influxdb_ip" {
  value = digitalocean_droplet.influxdb.ipv4_address
}
