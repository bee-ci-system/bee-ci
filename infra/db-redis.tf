resource "digitalocean_database_valkey_config" "default" {
  cluster_id              = digitalocean_database_cluster.redis.id
  valkey_maxmemory_policy = "allkeys-lru"
  notify_keyspace_events  = "KEA"
  timeout                 = 60
}

# Actually "valkey" but we keep using the name "redis".
resource "digitalocean_database_cluster" "redis" {
  name       = "bee-db-cluster-redis"
  engine     = "valkey"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "sfo3"
  node_count = 1
}
