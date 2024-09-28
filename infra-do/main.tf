terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "2.41.0"
    }
  }

  cloud {
    organization = "bpacia"

    workspaces {
      project = "bee-ci"
      name    = "bee-ci-workspace"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

resource "digitalocean_project" "project" {
  name        = "bee-ci-proj-tf"
  description = "A simple container-based CI system"
  environment = "Development"
  resources = [
    digitalocean_app.app.urn,
    digitalocean_database_cluster.main-db-cluster.urn,
    digitalocean_domain.default.urn,
  ]
}

resource "digitalocean_app" "app" {

  spec {
    name   = "bee-ci-tf"
    region = "sfo"

    ingress {
      rule {
        component {
          name = "frontend"
        }
        match {
          path {
            prefix = "/"
          }
        }
      }

      rule {
        component {
          name = "server"
        }
        match {
          path {
            prefix = "/backend"
          }
        }
      }
    }

    database {
      name         = digitalocean_database_db.main-db.name
      db_name      = digitalocean_database_db.main-db.name
      cluster_name = digitalocean_database_cluster.main-db-cluster.name
      production   = true
    }

    service {
      name               = "frontend"
      instance_count     = 1
      instance_size_slug = "apps-s-1vcpu-0.5gb" # doctl apps tier instance-size list

      http_port = 3000

      github {
        repo           = "bee-ci-system/bee-ci"
        branch         = "master"
        deploy_on_push = true
      }

      source_dir      = "./frontend"
      dockerfile_path = "./frontend/Dockerfile"

      health_check {
        http_path             = "/"
        initial_delay_seconds = 10
        period_seconds        = 5
        timeout_seconds       = 1
        success_threshold     = 3
        failure_threshold     = 3
      }
    }

    service {
      name               = "server"
      instance_count     = 3
      instance_size_slug = "apps-s-1vcpu-1gb" # doctl apps tier instance-size list

      dynamic "env" {
        for_each = local.env_vars
        content {
          key   = env.value.key
          value = env.value.value
          scope = env.value.scope
          type  = lookup(env.value, "type", null)
        }
      }

      http_port = 8080

      github {
        repo           = "bee-ci-system/bee-ci"
        branch         = "master"
        deploy_on_push = true
      }

      source_dir      = "./backend"
      dockerfile_path = "./backend/server.dockerfile"

      health_check {
        http_path             = "/"
        initial_delay_seconds = 10
        period_seconds        = 5
        timeout_seconds       = 1
        success_threshold     = 3
        failure_threshold     = 3
      }
    }


    worker {
      name               = "gh-updater"
      instance_count     = 1
      instance_size_slug = "apps-s-1vcpu-0.5gb" # doctl apps tier instance-size list

      dynamic "env" {
        for_each = local.env_vars
        content {
          key   = env.value.key
          value = env.value.value
          scope = env.value.scope
          type  = lookup(env.value, "type", null)
        }
      }

      github {
        repo           = "bee-ci-system/bee-ci"
        branch         = "master"
        deploy_on_push = true
      }

      source_dir      = "./backend"
      dockerfile_path = "./backend/gh-updater.dockerfile"
    }
  }
}

resource "digitalocean_database_db" "main-db" {
  cluster_id = digitalocean_database_cluster.main-db-cluster.id
  name       = "bee"

  provisioner "local-exec" {
    command = <<EOF
      psql -f ../backend/sql-scripts/1-schema.sql
      psql -f ../backend/sql-scripts/2-triggers.sql
      psql -f ../backend/sql-scripts/3-views.sql
      psql -f ../backend/sql-scripts/100-seed.sql
    EOF

    environment = {
      "PGUSER"     = digitalocean_database_cluster.main-db-cluster.user # -U
      "PGHOST"     = digitalocean_database_cluster.main-db-cluster.host # -h
      "PGPORT"     = digitalocean_database_cluster.main-db-cluster.port # -p
      "PGDATABASE" = digitalocean_database_db.main-db.name              # -d
      "PGSSLMODE"  = "require"                                          # -c sslmode=require
      "PGPASSWORD" = digitalocean_database_cluster.main-db-cluster.password
    }
  }
}

resource "digitalocean_database_cluster" "main-db-cluster" {
  name       = "bee-postgres-cluster"
  engine     = "pg"
  version    = "16"
  size       = "db-s-1vcpu-1gb"
  region     = "sfo2"
  node_count = 1
}

resource "digitalocean_domain" "default" {
  name = "bee-ci.karolak.cc"
}
