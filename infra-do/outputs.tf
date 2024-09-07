output "db_uri" {
  description = "The full URI for connecting to the database cluster"
  value       = digitalocean_database_cluster.main-db-cluster.uri
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
