resource "digitalocean_database_redis_config" "default" {
  cluster_id             = digitalocean_database_cluster.redis.id
  maxmemory_policy       = "allkeys-lru"
  notify_keyspace_events = "KEA"
  timeout                = 60
}

resource "digitalocean_database_cluster" "redis" {
  name       = "bee-db-cluster-redis"
  engine     = "redis"
  version    = "7"
  size       = "db-s-1vcpu-1gb"
  region     = "sfo2"
  node_count = 1
}
