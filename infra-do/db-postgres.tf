resource "digitalocean_database_db" "postgres" {
  cluster_id = digitalocean_database_cluster.postgres.id
  name       = "bee"

  provisioner "local-exec" {
    command = <<EOF
      psql -f ../backend/sql-scripts/1-schema.sql
      psql -f ../backend/sql-scripts/2-triggers.sql
      psql -f ../backend/sql-scripts/3-views.sql
      psql -f ../backend/sql-scripts/100-seed.sql
    EOF

    environment = {
      "PGUSER" = digitalocean_database_cluster.postgres.user # -U
      "PGHOST" = digitalocean_database_cluster.postgres.host # -h
      "PGPORT" = digitalocean_database_cluster.postgres.port # -p
      "PGDATABASE" = digitalocean_database_db.postgres.name  # -d
      "PGSSLMODE" = "require"                                # -c sslmode=require
      "PGPASSWORD" = digitalocean_database_cluster.postgres.password
    }
  }
}

resource "digitalocean_database_cluster" "postgres" {
  name       = "bee-db-cluster-postgres"
  engine     = "pg"
  version    = "16"
  size       = "db-s-1vcpu-1gb"
  region     = "sfo2"
  node_count = 1
}
